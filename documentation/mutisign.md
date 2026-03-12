# spiderd（Ignite + Cosmos SDK）多机多签（2/3）最终操作文档（denom=usc）

目标：实现 **2/3 多签账户**，签名人 **alice（机器A）**、**bob（机器A）**、**carol（机器B）**，测试交易由 **alice + carol** 联合签名完成。  
要求：**carol 私钥只存机器B**；币种单位使用 **`usc`**。

> 适用 Cosmos SDK v0.53.x（你报错版本 v0.53.5）。该版本在 `tx sign --multisig` 时需要本机 keyring 能找到对应 multisig 公钥记录，因此机器B也必须导入 multisig 公钥（不含私钥）。

---

## 1. 前置与约定

- 二进制：`spiderd`
- 链 ID：假设 `spider`（按实际替换）
- RPC：机器A跑节点，机器B可访问机器A的 `26657`
- keyring：示例使用 `test`（生产别用）

两台机器都设置：

```bash
export CHAIN_ID=spider
export KEYRING_BACKEND=test
```

机器A（本地节点）：

```bash
export NODE=http://127.0.0.1:26657
```

机器B（访问机器A节点 RPC）：

```bash
export NODE=http://<机器A可访问IP>:26657
```

---

## 2. 机器A：启动链（Ignite）

机器A在链工程目录：

```bash
ignite chain serve
```

确保机器B能访问机器A的 RPC `26657`（网络/防火墙/监听地址）。

---

## 3. 创建签名人账号（私钥分布在两台机器）

### 3.1 机器A：创建 alice、bob

```bash
spiderd keys add alice --keyring-backend $KEYRING_BACKEND
spiderd keys add bob   --keyring-backend $KEYRING_BACKEND
```

记录地址：

```bash
spiderd keys show alice -a --keyring-backend $KEYRING_BACKEND
spiderd keys show bob   -a --keyring-backend $KEYRING_BACKEND
```

### 3.2 机器B：创建 carol（私钥只在机器B）

```bash
spiderd keys add carol --keyring-backend $KEYRING_BACKEND
spiderd keys show carol -a --keyring-backend $KEYRING_BACKEND
```

导出 carol 公钥（给机器A用；不含私钥）：

```bash
spiderd keys show carol -p --keyring-backend $KEYRING_BACKEND > carol.pub.json
cat carol.pub.json
```

把 `carol.pub.json` 传给机器A（scp/复制粘贴皆可）。

---

## 4. 机器A：导入 carol 公钥并创建 multisig（2/3）

### 4.1 机器A：导入 carol 公钥为只读条目

```bash
spiderd keys add carol-pub \
  --pubkey "$(cat carol.pub.json)" \
  --keyring-backend $KEYRING_BACKEND
```

### 4.2 机器A：创建 2/3 多签 key（ms2of3）

```bash
spiderd keys add ms2of3 \
  --multisig alice,bob,carol-pub \
  --multisig-threshold 2 \
  --keyring-backend $KEYRING_BACKEND
```

获取多签地址：

```bash
export MS_ADDR=$(spiderd keys show ms2of3 -a --keyring-backend $KEYRING_BACKEND)
echo $MS_ADDR
```

---

## 5. 机器A：给 multisig 地址充值 usc（用于转账/手续费）

用一个有钱账户（示例用 alice；若你链上初始资金不在 alice，请换成实际有余额账户）：

```bash
spiderd tx bank send alice $MS_ADDR 1000000usc \
  --chain-id $CHAIN_ID \
  --node $NODE \
  --keyring-backend $KEYRING_BACKEND \
  --gas auto --gas-adjustment 1.3 \
  -y
```

查询余额：

```bash
spiderd query bank balances $MS_ADDR --node $NODE
```

---

## 6. 关键修复：机器B也必须拥有 ms2of3 的“公钥记录”

你遇到的错误本质是：机器B keyring 找不到 multisig key（按地址/引用解析失败），所以无法 `tx sign --multisig ...`。

### 6.1 机器A：导出 ms2of3 公钥（给机器B）

```bash
spiderd keys show ms2of3 -p --keyring-backend $KEYRING_BACKEND > ms2of3.pub.json
cat ms2of3.pub.json
```

把 `ms2of3.pub.json` 传给机器B。

### 6.2 机器B：导入 ms2of3 公钥（仅公钥，不含私钥）

```bash
spiderd keys add ms2of3 \
  --pubkey "$(cat ms2of3.pub.json)" \
  --keyring-backend $KEYRING_BACKEND
```

校验机器B看到的地址与机器A一致：

```bash
spiderd keys show ms2of3 -a --keyring-backend $KEYRING_BACKEND
```

---

## 7. 机器A：生成待签名交易（unsigned tx.json，denom=usc）

示例：多签账户给 bob 转 `1usc`。

机器A：

```bash
export BOB_ADDR=$(spiderd keys show bob -a --keyring-backend $KEYRING_BACKEND)

spiderd tx bank send $MS_ADDR $BOB_ADDR 1usc \
  --from $MS_ADDR \
  --generate-only \
  --chain-id $CHAIN_ID \
  --node $NODE \
  > tx.json
```

把 `tx.json` 传给机器B（必须是同一份文件）：

```bash
scp tx.json user@machineB:/path/to/tx.json
```

---

## 8. 双机分别部分签名（alice + carol）

> 注意：这里 **`--multisig` 推荐传 key 名称 `ms2of3`**，不要传地址，减少 keyring 查找歧义。

### 8.1 机器A：alice 部分签名

```bash
spiderd tx sign tx.json \
  --from alice \
  --multisig ms2of3 \
  --chain-id $CHAIN_ID \
  --node $NODE \
  --keyring-backend $KEYRING_BACKEND \
  --output-document alice.sig.json
```

### 8.2 机器B：carol 部分签名（私钥不出机器B）

机器B在有 `tx.json` 的目录：

```bash
spiderd tx sign tx.json \
  --from carol \
  --multisig ms2of3 \
  --chain-id $CHAIN_ID \
  --node $NODE \
  --keyring-backend $KEYRING_BACKEND \
  --output-document carol.sig.json
```

把 `carol.sig.json` 传回机器A：

```bash
scp carol.sig.json user@machineA:/path/to/carol.sig.json
```

---

## 9. 机器A：聚合签名并广播

机器A（同目录有 `tx.json`、`alice.sig.json`、`carol.sig.json`）：

```bash
spiderd tx multisign tx.json ms2of3 alice.sig.json carol.sig.json \
  --chain-id $CHAIN_ID \
  --node $NODE \
  > tx.signed.json
```

广播：

```bash
spiderd tx broadcast tx.signed.json \
  --chain-id $CHAIN_ID \
  --node $NODE
```

验收（可选）：

```bash
spiderd query bank balances $MS_ADDR --node $NODE
spiderd query bank balances $BOB_ADDR --node $NODE
```

---

## 10. 常见问题与排查

### 10.1 报错：`key with address ... not found`（你遇到的）
原因：签名机器的 keyring 找不到 multisig 公钥记录。  
解决：按第6节在该机器导入 `ms2of3.pub.json`（只公钥）。

### 10.2 聚合/广播失败：sequence 不一致
原因：两人签的不是同一份交易，或签名期间 multisig 账户 sequence 发生变化。  
解决：
- 始终让所有人对同一个 `tx.json` 签名
- 如果多签账户期间又发过交易，重新生成 `tx.json` 再签

### 10.3 denom 不对
确认你的链 `usc` 是否为正确 denom（genesis/fee token）。若实际 denom 是别的（如 `stake`），把所有 `usc` 替换即可。

---

如果你把以下信息发我，我可以把文档里的占位符（IP、链ID、初始有钱账户、denom）改成你环境的一键可跑版本，并补上 scp 的完整命令与目录结构建议：
1) `spiderd status | jq .NodeInfo.network` 的 chain-id  
2) `spiderd query bank total` 或 genesis 中的 denom  
3) 你本地链默认的有钱账户名/地址（通常是 validator 或 faucet）




#### 自定义模块多签测试

```
spiderd tx official create-operator $(spiderd keys show ms2of3 -a) tokenfactory ms2of3 1 7 --from tom

spiderd query official list-operator
spiderd query tokenfactory list-denom

spiderd tx tokenfactory create-namespace foo 100usc  \
  --from $MS_ADDR \
  --generate-only \
  --chain-id $CHAIN_ID \
  --node $NODE \
  > tx.json

#其余过程同上各自签名与合并
spiderd query tokenfactory list-namespace
```