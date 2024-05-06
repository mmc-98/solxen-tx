package svc

import (
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/sha512"
	"math/big"

	"github.com/gagliardetto/solana-go"
	"github.com/mr-tron/base58"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/pbkdf2"
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

	logx.Infof("len: %v", s.Config.Sol.Num)
	for i := 0; i < s.Config.Sol.Num; i++ {
		// BIP-39
		mnemonic := s.Config.Sol.Key
		seed := pbkdf2.Key([]byte(mnemonic), []byte("mnemonic"), 2048, 64, sha512.New)

		// BIP-32
		h := hmac.New(sha512.New, []byte("ed25519 seed"))
		h.Write(seed)
		sum := h.Sum(nil)

		derivedSeed := sum[:32]
		chain := sum[32:]

		// BIP-44
		// m/44'/501'/index'/0'/0'
		// m/44'/501'/index'/1'/0'
		path := []uint32{Hardened + uint32(44), Hardened + uint32(501), Hardened + uint32(i), Hardened + uint32(0)}
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

		s.AddrList = append(s.AddrList, wallet)
	}

}
