package svc

import (
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/sha512"
	"crypto/x509"
	"log"
	"math/big"
	pb "solxen-tx/internal/svc/proto"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/mr-tron/base58"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/pbkdf2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

func derive(key []byte, chainCode []byte, segment uint32) ([]byte, []byte) {
	// Create buffer
	buf := []byte{0}
	buf = append(buf, key...)
	buf = append(buf, big.NewInt(int64(segment)).Bytes()...)

	// Calculate HMAC hash
	h := hmac.New(sha512.New, chainCode)
	h.Write(buf)
	I := h.Sum(nil)

	// Split result
	IL := I[:32]
	IR := I[32:]

	return IL, IR
}

const Hardened uint32 = 0x80000000

func (s *ServiceContext) GenKeyByWord() {

	logx.Infof("len: %v ProgramID: %v ", s.Config.Sol.Num, s.Config.Sol.ProgramId)

	for i := 0; i < s.Config.Sol.Num; i++ {
		// BIP-39
		seed := pbkdf2.Key([]byte(s.Config.Sol.Mnemonic), []byte("mnemonic"), 2048, 64, sha512.New)

		// BIP-32
		h := hmac.New(sha512.New, []byte("ed25519 seed"))
		h.Write(seed)
		sum := h.Sum(nil)

		derivedSeed := sum[:32]
		chain := sum[32:]

		// BIP-44
		// m/44'/501'/index'/0'/0'
		// m/44'/501'/index'/1'/0'
		hdPath := s.Config.Sol.HdPath
		var path []uint32
		switch hdPath {
		case "m/44'/501'":
			path = []uint32{Hardened + uint32(44), Hardened + uint32(501)}
		case "m/44'/501'/0'":
			path = []uint32{Hardened + uint32(44), Hardened + uint32(501), Hardened + uint32(i)}
		default:
			// m/44'/501'/0'/0'
			path = []uint32{Hardened + uint32(44), Hardened + uint32(501), Hardened + uint32(i), Hardened + uint32(0)}

		}
		for _, segment := range path {
			derivedSeed, chain = derive(derivedSeed, chain, segment)
		}
		key := ed25519.NewKeyFromSeed(derivedSeed)
		// Get Solana wallet
		wallet, err := solana.WalletFromPrivateKeyBase58(base58.Encode(key))
		if err != nil {
			panic(err)
		}
		// address := wallet.PublicKey().String()
		logx.Infof("account: %v", wallet.PublicKey())
		s.AddrList = append(s.AddrList, wallet)
	}

	// err := filepath.Walk(s.Config.Sol.WalletPath, func(path string, info os.FileInfo, err error) error {
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if !info.IsDir() && filepath.Ext(path) == ".json" {
	// 		PrivateKey, err := solana.PrivateKeyFromSolanaKeygenFile(path)
	// 		if err != nil {
	// 			logx.Errorf("err: %v", err)
	// 		}
	// 		s.AddrList = append(s.AddrList, &solana.Wallet{PrivateKey: PrivateKey})
	// 	}
	// 	return nil
	// })
	//
	// if err != nil {
	// 	logx.Errorf("err: %v", err)
	// }

}

var kacp = keepalive.ClientParameters{
	Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
	Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
	PermitWithoutStream: true,             // send pings even without active streams
}

func NewGrpcCli(address string, plaintext bool) pb.GeyserClient {
	var opts []grpc.DialOption
	if plaintext {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		pool, _ := x509.SystemCertPool()
		creds := credentials.NewClientTLSFromCert(pool, "")
		opts = append(opts, grpc.WithTransportCredentials(creds))
	}

	opts = append(opts, grpc.WithKeepaliveParams(kacp))

	log.Println("Starting grpc client, connecting to", address)
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	return pb.NewGeyserClient(conn)
}
