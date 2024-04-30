# SolXen Blockchain(Testnet) Load Testing Tool 压力测试程序


 
## 1. Configuration [配置]


Configuration File Path [配置文件路径]: **build/etc/solxen-tx.yaml**

```shell
The explanation of the fields in this configuration file is as follows: [配置文件字段解释如下:]
 
# sol配置
Sol:
  Url: "https://api.devnet.solana.com"   # RPC Address [rpc地址]
  Key: ""                             # key [私钥]
  Num: 1                              # Concurrency [并发数]
  Fee: 3000
  ToAddr: "0x4A7766a5BD50DFAB5665d27eDfA25471b194E204"   #  填上你的eth地址
  ProgramID: "64SYet8RCT5ayZpMGbhcpk3vmt8UkwjZq8uy8Sd6V46A"
  Time: 1000             # Interval time (in milliseconds) [间隔时间(单位毫秒)]
```
 

## 2. Start Running [运行]

```shell
make start
```
 
