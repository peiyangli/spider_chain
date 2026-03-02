### spider_chain项目

spider项目的区块链

#### 项目生成与测试

使用工具[ignite](ignite.md)生成，生成步骤与测试可以参考[文档](stepbystep.md)

#### 模块

##### official

社区官方模块，提供部分模块权限分配

创世文件中配置参数【根管理员】，有权限创建其他模块的操作员。
1. Identity模块的绑定员地址，由中心分配绑定uid与用户身份；
2. snft命名空间定价员，通过命名空间前缀，管理nft-class-id创建价格；

##### identity
绑定用户uid与公钥，消息公钥由身份私钥持有者生成与修改

```proto
message Identity {
  string uid = 1; //用户uid
  string owner = 2; //用户地址
  bytes idkey = 3; //身份公钥
  bytes msgkey = 4; //消息公钥
  string creator = 5; //official模块设置的管理员，绑定者
}
```


##### snft

nft的铸造模块，管理nft命名空间，class-id所有者发布nft

nft由class-id+id一起确认唯一nft，参考[Ethereum ERC721 standard](https://ethereum.org/developers/docs/standards/tokens/erc-721)

nft类的所有者可以铸造相应的nft，铸造后由nft的所有者控制交易

###### 1.带前缀的class-id

class-id: 前缀.子id
定价：由模块操作员设置【前缀】的价格

    例子：
    {"avatar": 100usc}
    用户就可以注册avatar.alice为其所有的class-id，后续他可以发布nft



###### 2.不带前缀的class-id
由有权限的模块操作员创建


##### tokenfactory

token发行模块，用户支付一定的usc发行自有货币


##### loan

抵押借贷模块
