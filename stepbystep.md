### 蜘蛛链SpiderChain设计

```
#生成项目
ignite scaffold chain spider_chain --address-prefix spider --no-module
#并配置config.yml

ignite chain serve
```

## id绑定

由官方分配uid，并绑定身份公钥与消息根公钥

```
#生成官方模块
ignite scaffold module official

```

修改config.yml
```yml
version: 1
validation: sovereign
default_denom: usc
accounts: 
- name: alice
  coins:
    - 100000000000usc
- name: bob
  coins:
    - 50000000000usc
client:
  openapi:
    path: docs/static/openapi.json
faucet:
  name: bob
  coins:
    - 5usc
validators:
  - name: alice
    bonded: 10000000000usc
genesis:
  app_state:
    staking:
      params:
        bond_denom: usc
    crisis:
      constant_fee:
        denom: usc
        amount: "1000"
    gov:
      params:
        min_deposit:
          - denom: usc
            amount: "10000000"
    mint:
      params:
        mint_denom: usc
```

##### 账号及权限管理 official 模块

<!-- 1. Manager
负责录入操作员Operator，可添加修改Operator及权限；
创世可配置，默认模块地址；
所属模块： official
```
ignite scaffold params manager:string
``` -->

2. Operator
操作员，负责写入uid->pubkey绑定

role = 0-0xf: manager
role = 0x10: uid->pubkey绑定、取消绑定
```
ignite scaffold map operator name:string role:uint32 module:string permissions:uint64 --index address --module official

#修改代码
#启动服务
ignite chain serve --reset-once

#测试接口
spiderd query official list-operator
#创建一个操作员
spiderd tx official create-operator $(spiderd keys show bob -a) official bob 0 1 --from alice
spiderd query official get-operator $(spiderd keys show bob -a)
#再创建一个操作员
spiderd tx official create-operator $(spiderd keys show bob -a) identity bob 1 7 --from alice
spiderd query official get-operator $(spiderd keys show bob -a)
#删除操作员
spiderd tx official delete-operator $(spiderd keys show bob -a) --from alice
#订单查询
spiderd q tx


#备注：配置地址，先生成地址获取助记词，地址填入配置。后续服务启动后通过如下恢复即可
spiderd keys add tom --recover
#  address spider1gnx609p4e8p3snwm54f204pnncmy99ltm2w3wr and mnemonic:
#  pitch please side evolve depart stay bid distance enter tired notable oxygen lunar choice replace wrap vacant street decorate craft focus salute inspire joke

#记得转账给tom
spiderd tx bank send alice $(spiderd keys show tom -a) 1000usc
spiderd query bank balances $(spiderd keys show tom -a)
```

##### 用户身份 identity 模块

huid=hash(uid)
idpub: 身份公钥
msgpub：消息公钥

```
#添加模块
ignite scaffold module identity
#添加绑定huid->pubkey
ignite scaffold map identity owner idkey:bytes msgkey:bytes --index uid --module identity
#添加管理员
spiderd tx official create-operator $(spiderd keys show bob -a) identity bob 1 7 --from tom --generate-only
#创建id->pub
spiderd tx identity create-identity uid001 $(spiderd keys show alice -a) $(printf 'hello' | base64 -w0) $(printf 'world' | base64 -w0) --from bob
#查询结果
spiderd query identity list-identity
#修改,本人owner只能修改msgkey,其他字段忽略
spiderd tx identity update-identity uid001 ccc $(printf 'good' | base64 -w0) $(printf 'byte' | base64 -w0) --from bob

ignite chain serve
```

## 资产交易

##### 铸币模块 tokenfactory

```
#模块创建
ignite scaffold module tokenfactory --dep auth,bank,official
#数据创建
ignite scaffold map namespace creation_fee:coin --signer creator --index namespace --module tokenfactory
ignite scaffold map denom description:string ticker:string precision:int url:string maxSupply:int supply:int canChangeMaxSupply:bool creationFee:coin --signer owner --index denom --module tokenfactory
#移除delete相关
make proto-gen

#发币，限owner
ignite scaffold message MintAndSend denom:string amount:int recipient:string --module tokenfactory --signer owner
#转移owner
ignite scaffold message UpdateOwner denom:string newOwner:string --module tokenfactory --signer owner

#参数不可以，之间改params.proto文件后make proto-gen
ignite scaffold params creation_fee:coin --module tokenfactory


ignite chain serve

# 测试
spiderd tx official create-operator $(spiderd keys show bob -a) tokenfactory bob 1 7 --from tom
# bob 创建成功， spiderd query official list-operator
spiderd tx tokenfactory create-denom xc "xid coin" IGNITE 6 "xid.spider.com" 1000000000 true 100usc --from bob
# alice创建失败
spiderd tx tokenfactory create-denom usc "My denom spider coin" IGNITE 6 "spider.com" 1000000000 true 100usc --from alice
# alice创建成功 tf/
spiderd tx tokenfactory create-denom tf/abc "My denom spider coin" IGNITE 6 "spider.alice.com" 1000000000 true 100usc --from alice

spiderd query tokenfactory list-denom

# mintandsend
spiderd tx tokenfactory mint-and-send xid 1200 $(spiderd keys show tom -a) --from bob
spiderd tx tokenfactory mint-and-send tf/abc 1200 $(spiderd keys show tom -a) --from alice
```

```
import "cosmos/base/v1beta1/coin.proto";
message Params {
cosmos.base.v1beta1.Coin creation_fee = 3 [(gogoproto.nullable) = false, (amino.dont_omitempty) = true];
}
```

##### nft支持

nft支持
```
ignite scaffold module snft --dep auth,bank,nft,official

<!-- import "cosmos/base/v1beta1/coin.proto";
message Params {
cosmos.base.v1beta1.Coin creation_fee = 3 [(gogoproto.nullable) = false, (amino.dont_omitempty) = true];
} -->

# name:string symbol:string description:string uri:string uri_hash:string
ignite scaffold map class_owner pending_owner:string creation_fee:coin --signer owner --index class_id --module snft

# 
ignite scaffold map class_namespace creation_fee:coin --signer creator --index namespace --module snft

// decision: 1 accept, 2 reject
ignite scaffold message RespondClassOwnerTransfer class_id:string decision:uint64 --module snft --signer owner
ignite scaffold message MintAndSend class_id:string nft_id:string uri:string uri_hash:string recipient:string --module snft --signer owner


//测试准备
ignite chain serve

#备注：配置地址，先生成地址获取助记词，地址填入配置。后续服务启动后通过如下恢复即可
spiderd keys add tom --recover
#  address spider1gnx609p4e8p3snwm54f204pnncmy99ltm2w3wr and mnemonic:
#  pitch please side evolve depart stay bid distance enter tired notable oxygen lunar choice replace wrap vacant street decorate craft focus salute inspire joke

#记得转账给tom
spiderd tx bank send alice $(spiderd keys show tom -a) 1000usc
spiderd query bank balances $(spiderd keys show tom -a)

spiderd tx official create-operator $(spiderd keys show bob -a) snft bob 1 7 --from tom
spiderd query official list-operator

//创建classnamespace
//namespace creation-fee
spiderd tx snft create-class-namespace foo 100usc --from tom

//classid pendingowner coin
spiderd tx snft create-class-owner xid "" 100usc --from tom

spiderd query snft list-class-namespace

//owner alice
spiderd tx snft mint-and-send foo:abc 123456 nft.alice.com hashalicecom $(spiderd keys show tom -a) --from alice
spiderd query nft nfts foo:abc
spiderd query nft nfts "" --owner=$(spiderd keys show tom -a)

spiderd tx nft send foo:abc 123456 $(spiderd keys show bob -a) --from tom

```

##### 抵押借贷模块 loan

抵押coin或者nft

```
# 模块
ignite scaffold module loan --dep bank,nft

# 基本结构 deadline为区块高度
ignite scaffold map loan status:uint64 seq:uint64 deadline:uint64 public_liquidation_delay:uint64 public_liquidation_reward amount fee lender:address collateral_type collateral_coin collateral_nft_class collateral_nft_id --index borrower --module loan --no-message

# 请求
ignite scaffold message request-loan term:uint64 deadline:uint64 amount fee collateral_type collateral_coin collateral_nft_class collateral_nft_id --module loan

# 同意借款, 三方清算奖励以及块高度
ignite scaffold message approve-loan seq:uint64 borrower:address public_liquidation_delay:uint64 public_liquidation_reward --module loan

# 取消
ignite scaffold message cancel-loan seq:uint64 --module loan

# 偿还, 代偿borrower
ignite scaffold message repay-loan seq:uint64 borrower:address --module loan

# 清算，三方清算奖励，目标清算后一定时期才行
ignite scaffold message liquidate-loan seq:uint64 borrower:address --module loan

# 测试
ignite chain serve
# 1. tokenfactory生成代币模块， 质押获取usc
spiderd query loan list-loan
spiderd tx loan request-loan 5000 500usc 50usc coin 1000tf/ac "" "" --from bob
spiderd tx loan cancel-loan 0 --from bob
spiderd tx loan approve-loan 1 $(spiderd keys show bob -a) 100 100tf/ac --from alice

spiderd tx loan liquidate-loan 1 $(spiderd keys show bob -a) --from alice

spiderd tx loan repay-loan 1 $(spiderd keys show bob -a) --from bob
```



##### 资产模块 xid