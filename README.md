# solXen mint工具 [solXen Mint Tool]

## 特点 [Features]
```shell
支持多线程并发模式 [Supports multithreaded concurrent mode] done
支持获取优先级费用 [Supports automatic priority fee setting]
支持自动从子钱包转账 [Supports automatic transfer from sub-wallets]
```

## 0. 下载 [Download]

```
https://github.com/mmc-98/solxen-tx/releases
```

Windows (x64):
```shell
solxen-tx-windows-amd64.zip
```

Linux (x64):
```shell
solxen-tx-linux-amd64.tar.gz
```

Linux (arm64):
```shell
solxen-tx-linux-arm64.tar.gz
```

macOS (x64):
```shell
solxen-tx-darwin-amd64.tar.gz
```

macOS (arm64):
```shell
solxen-tx-darwin-arm64.tar.gz
```
 
## 1. 配置 [Configuration]


```shell
# sol配置 [Configuration]
Name: solxen-tx
Sol:
  Url: "https://api.devnet.solana.com"                          # rpc地址 [rpc address]
  Mnemonic: ""                                                  # 助记词 [mnemonic phrase]
  Num: 1                                                        # 并发数 [concurrency]
  Fee: 3000                                                     # 优先级费用 [priority fee]
  ToAddr: "0x4A7766a5BD50DFAB5665d27eDfA25471b194E204"          # eth空投地址 [eth address for receiving xn airdrop]
  ProgramID: "7LBe4g8Q6hq8Sk1nT8tQUiz2mCHjsoQJbmZ7zCQtutuT"     # solxen合约地址 [solxen contract address]
  Time: 1000                                                    # 间隔时间(单位毫秒) [interval time (milliseconds)]
  HdPAth: m/44'/501'/0'/0'                                      # 钱包地址路径 [wallet derivation path]
```
 

## 2. 运行 [Run]

```shell
./solxen-tx 
```

## 3. mint  
./mint