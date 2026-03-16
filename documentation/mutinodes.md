### 多节点流程



#### 一、部署前准备

##### 机器与网络

- 固定公网 IP 或内网互通 IP
- 开放端口（至少）：
    - 26656/TCP：P2P（必须）
    - 26657/TCP：RPC（可选，建议内网或不对公网）
    - 1317/TCP：REST API（可选）
    - 9090/TCP：gRPC（可选）
- 时间同步：chrony/ntpd（很重要）

#####  版本统一

-  同一份**spiderd**二进制
-  同一份 genesis.json
-  同一 seeds/persistent_peers 配置

##### 节点角色
- 至少 1 台作为 协调机：负责收集所有验证人的 gentx，最终产出 genesis
- 每台验证人节点：生成自己的 key、节点 p2p node_key、validator gentx


#### 二、安装与初始化

##### 安装二进制与依赖
把 spiderd 放到每台机器（建议 /usr/local/bin/spiderd）
```
scp spiderd root@host:/usr/local/bin/
chmod +x /usr/local/bin/spiderd
spiderd version
```

##### 初始化 home（每台节点独立）


spiderd init -h
spiderd init [moniker] --chain-id spider

```
spiderd init spider001 --chain-id spider
spiderd init spider002 --chain-id spider

或在各个机器上定义变量
export MONIKER=sa99
export MONIKER=sa93
export MONIKER=sa92

spiderd init $MONIKER --chain-id spider
```

#### 三、生成验证人账户与 gentx（每个验证人节点做）

##### 创建验证人账户 key

spiderd keys add -h
spiderd keys add <name> [flags]

在各个机器上生成账户
```
# 可以把助记词记下来
spiderd keys add $MONIKER --keyring-backend file --keyring-dir $MONIKER

# 导出地址为变量，后面有用
MONIKER_ADDR=$(spiderd keys show $MONIKER -a --keyring-backend file --keyring-dir $MONIKER)
echo $MONIKER_ADDR

# 显示完整地址和公钥，公钥后续拿来多签
spiderd keys show $MONIKER --keyring-backend file --keyring-dir $MONIKER
```

##### 把地址和公钥发给协调机（用于加余额和生成多签）
```
- address: spider1cfte5kj58nsvdaqzqhutqw3de4ycjk5j7patrl
  name: sa99
  pubkey: '{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"A7gbymjqPHs8R7sg/TCw/KW+WDZ/dFYI6pZTGvxrk2Kp"}'
  type: local
- address: spider1r7l4uhhdfcevk5uafcnhz0k7l8424gmpdm24yt
  name: sa93
  pubkey: '{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"A8GhssjO0Z7NkGFzXibawapIAln6EB5Pmq45YN6IvX3K"}'
  type: local
- address: spider1jsul760jfqth6qk8kck2d2n45sjzydr2yaprc5
  name: sa92
  pubkey: '{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"Ajm6bx9yVDc7s5HRR3iZ1CuwSYKzblGvPQM5I0vPUTZC"}'
  type: local
```

- 生成多签

协调机器
```
#导入公钥
spiderd keys add sa99-pub --pubkey '{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"A7gbymjqPHs8R7sg/TCw/KW+WDZ/dFYI6pZTGvxrk2Kp"}' --keyring-backend file --keyring-dir sateam
spiderd keys add sa93-pub --pubkey '{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"A8GhssjO0Z7NkGFzXibawapIAln6EB5Pmq45YN6IvX3K"}' --keyring-backend file --keyring-dir sateam
spiderd keys add sa92-pub --pubkey '{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"Ajm6bx9yVDc7s5HRR3iZ1CuwSYKzblGvPQM5I0vPUTZC"}' --keyring-backend file --keyring-dir sateam
#生成多签公钥
spiderd keys add sateam --multisig sa99-pub,sa93-pub,sa92-pub --multisig-threshold 2 --keyring-backend file --keyring-dir sateam

#展示keyring中的所有key
spiderd keys list --keyring-backend file --keyring-dir sateam

#导出多签地址为变量
export MS_ADDR=$(spiderd keys show sateam -a --keyring-backend file --keyring-dir sateam)
echo $MS_ADDR
```

- 添加余额

协调机器
```
# 在这里添加所有初始账户余额
spiderd genesis add-genesis-account spider1jsul760jfqth6qk8kck2d2n45sjzydr2yaprc5 200000000usc
spiderd genesis add-genesis-account spider1r7l4uhhdfcevk5uafcnhz0k7l8424gmpdm24yt 200000000usc
spiderd genesis add-genesis-account spider1cfte5kj58nsvdaqzqhutqw3de4ycjk5j7patrl 300000000usc
spiderd genesis add-genesis-account spider1hkdc8ut955nwz6vxg7zs870qp8pfzc26ggleu8 500000000usc

# 顺便修改创世文件

# 修改app_name
- "app_name": "\u003cappd\u003e"
+ "app_name": "spiderd"

#把所有货币单位stake修改为usc，并搜索denom进行确认


# 添加参数：official模块管理员operator_map
        {
          "address": "spider1hkdc8ut955nwz6vxg7zs870qp8pfzc26ggleu8",
          "module": "official",
          "name": "sateam",
          "role": "512",
          "permissions": "255",
          "creator": ""
        }
# 其他模块参数
```

- 修改好后拷贝到各个验证节点


各个验证节点质押，生成gentx
```
spiderd genesis gentx $MONIKER 5000000usc --chain-id spider --keyring-backend file --keyring-dir $MONIKER
#生成的文件在目录里面： ~/spider/config/gentx
```

##### 协调机：制作最终 genesis（只在协调机做）

把所有机器上的gentx文件拷贝到协调机器上

```
spiderd genesis collect-gentxs
spiderd genesis validate-genesis
```

在各个机器上对照genesis.json里面的validator_address进行验证
```
spiderd keys show $MONIKER --bech val --keyring-backend file --keyring-dir $MONIKER
```


##### 分发最终 genesis 给所有节点
把最终 genesis.json 分发到每台机器的：
~/.spider/config/genesis.json
注意：所有节点必须一模一样（可以对比 sha256）：
```
sha256sum ~/.spider/config/genesis.json
```


#### P2P 互联配置（每台机器都做）

##### 获取每台节点的 node_id
```
# 同验证gentx那个文件名
spiderd tendermint show-node-id
```

##### 配 persistent_peers（推荐至少互相连到 1-2 个稳定节点）

编辑每台机器的：
~/.spider/config/config.toml
设置：

persistent_peers = "id1@ip1:26656,id2@ip2:26656"
或者设置 seeds（更适合有种子节点）：

seeds = "seedid@seedip:26656"

小规模（3~10 个验证人）直接互相 persistent_peers 最简单。

##### 确认端口与外网地址
config.toml 里常见项：

laddr = "tcp://0.0.0.0:26656"
（如在 NAT 后）需要正确的 external_address = "ip:26656"（没有 NAT 可不配）

##### 设置气费
~/.spider/config/app.toml

minimum-gas-prices = "0.00001usc"


#### 启动与检查（每台机器）

##### 启动

```
spiderd start
```

##### 检查是否出块

```
spiderd status | jq .sync_info.latest_block_height
```


```
#查看keys
spiderd keys list --keyring-backend file --keyring-dir sateam

#查看账户
spiderd query bank balances $addr

#转账测试(sa92 to sa93)
spiderd tx bank send sa92 spider1r7l4uhhdfcevk5uafcnhz0k7l8424gmpdm24yt 1000usc --chain-id spider --fees 2usc --keyring-backend file --keyring-dir sa92

#official
spiderd query official list-operator
```

在window的wsl中防火墙开启端口访问
```ps
New-NetFirewallRule -DisplayName "WSL 26656" -Direction Inbound -Action Allow -Protocol TCP -LocalPort 26656
```


#### 动态添加新的验证者
