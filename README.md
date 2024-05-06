# SolXen mint工具

## 特点
```shell
支持多线程并发模式 done
支持获取优先级费用
支持自动从子钱包转账
```
## 0. 下载
windows:
```shell
https://github.com/mmc-98/solxen-tx/releases/download/v0.05/solxen-tx-v0.05-windows-amd64.zip
```
linux:
```shell
https://github.com/mmc-98/solxen-tx/releases/download/v0.05/solxen-tx-v0.05-linux-amd64.tar.gz
```
mac Apple芯片:
```shell
https://github.com/mmc-98/solxen-tx/releases/download/v0.05/solxen-tx-v0.05-darwin-arm64.tar.gz
```
 
## 1. 配置


```shell
# sol配置
Name: solxen-tx
Sol:
  Url: "https://api.devnet.solana.com"                          # rpc地址
  Key: ""                                                       # 私钥
  Num: 1                                                        # 并发数
  Fee: 3000                                                     # 优先级费用
  ToAddr: "0x4A7766a5BD50DFAB5665d27eDfA25471b194E204"          # eth空投地址
  ProgramID: "64SYet8RCT5ayZpMGbhcpk3vmt8UkwjZq8uy8Sd6V46A"     # solxen合约地址
  Time: 1000                                                    #  [间隔时间(单位毫秒)]
```
 

## 2. 运行

```shell
./solxen-tx 
```
 
