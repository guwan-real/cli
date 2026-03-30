# IM 模块测试记录

测试时间：2026-03-30  
测试环境：i 讯飞私有化飞书  
OpenAPI：`https://open.xfchat.iflytek.com`  
OAuth：`https://accounts.xfchat.iflytek.com`  
应用：`cli_a9084c0706b8d379` / `星火派-test`

## 1. 环境与鉴权

### 1.1 CLI 与配置

初始化命令：

```bash
./xfchat_cli config init --app-id cli_a9084c0706b8d379 --app-secret-stdin --brand feishu
```

输入：

```text
<AppSecret 已通过 stdin 输入，不落盘明文>
```

输出摘要：

```json
{
  "appId": "cli_a9084c0706b8d379",
  "appSecret": "****",
  "brand": "feishu"
}
```

补充配置：

```json
{
  "endpoints": {
    "open": "https://open.xfchat.iflytek.com",
    "accounts": "https://accounts.xfchat.iflytek.com"
  }
}
```

### 1.2 用户授权

授权命令：

```bash
./xfchat_cli auth login --web --scope 'im:chat:read im:chat:update im:chat.members:read im:chat.members:write_only offline_access'
```

实际授权后 `auth status` 摘要：

```json
{
  "appId": "cli_a9084c0706b8d379",
  "brand": "feishu",
  "defaultAs": "auto",
  "identity": "user",
  "scope": "auth:user.id:read contact:user.base:readonly im:chat.members:read im:chat.members:write_only im:chat:read im:chat:update im:message im:message:readonly search:message offline_access",
  "tokenStatus": "valid",
  "userName": "王琦钊 qzwang9",
  "userOpenId": "ou_6b90260ffbda660ec3e47f27c0871ec9"
}
```

结论：

- 登录成功
- 私有化环境最终授予的 scope 比输入 scope 多，额外出现了 `auth:user.id:read`、`contact:user.base:readonly`、`im:message`、`im:message:readonly`、`search:message`

### 1.3 环境检查

检查命令：

```bash
./xfchat_cli doctor --offline
./xfchat_cli api GET /open-apis/bot/v3/info --as bot
./xfchat_cli contact +get-user --format json
```

关键输出：

```json
{
  "app_resolved": "cli_a9084c0706b8d379 (feishu)",
  "bot.app_name": "星火派-test",
  "user.open_id": "ou_6b90260ffbda660ec3e47f27c0871ec9"
}
```

结论：环境可用，Bot 与 User 身份都已就绪。

## 2. 测试对象与素材

### 2.1 群聊对象

- `oc_c839c827acf6c34ca6287ea171a52258`：`xfchat-im-test-20260330-a`
- `oc_010df6b42b975ce056cc7c2e717abde8`：`xfchat-im-test-20260330-b-updated`
- `oc_c7d9abd30a87e6d16976966818e2c19a`：`xfchat-im-test-20260330-members`

### 2.2 私聊对象

- 用户：`ou_6b90260ffbda660ec3e47f27c0871ec9`
- P2P `chat_id`：`oc_934891fa0c30e80e4d7e7d5d7496d5e4`

### 2.3 多成员群成员

群 `oc_010df6b42b975ce056cc7c2e717abde8` 当前成员：

- `ou_6b90260ffbda660ec3e47f27c0871ec9`：王琦钊 qzwang9
- `ou_4dce96856e654228d968f1ae74e29ecd`：戴堃 kundai2
- `ou_981bb29d80e31c4aebe79225e866d8ff`：汪兵 bingwang32
- `ou_84c79a9284092a69283066a89e549251`：付涛 taofu5

群 `oc_c7d9abd30a87e6d16976966818e2c19a` 最终成员：

- `ou_6b90260ffbda660ec3e47f27c0871ec9`：王琦钊 qzwang9
- `ou_4dce96856e654228d968f1ae74e29ecd`：戴堃 kundai2
- `ou_981bb29d80e31c4aebe79225e866d8ff`：汪兵 bingwang32
- `ou_84c79a9284092a69283066a89e549251`：付涛 taofu5

### 2.4 测试素材

- 文件：`/Users/wangqizhao/Developer/iflytek/cli/tmp/im-test/sample.txt`
- 图片：`/Users/wangqizhao/Developer/iflytek/cli/tmp/im-test/sample.png`

下载校验：

```text
file_cmp:0
image_cmp:0
```

说明：下载回来的文件和图片与原始测试素材一致。

## 3. Shortcut 覆盖

| 命令 | 场景 | 权限/身份 | 输入 | 输出摘要 | 结论 |
|---|---|---|---|---|---|
| `+chat-create` | 创建测试群 A | Bot, `im:chat:create` | `./xfchat_cli im +chat-create --as bot --name 'xfchat-im-test-20260330-a' --users 'ou_6b90260ffbda660ec3e47f27c0871ec9'` | `chat_id=oc_c839c827acf6c34ca6287ea171a52258` | 通过 |
| `+chat-create` | 创建测试群 B | Bot, `im:chat:create` | `./xfchat_cli im +chat-create --as bot --name 'xfchat-im-test-20260330-b' --users 'ou_6b90260ffbda660ec3e47f27c0871ec9'` | `chat_id=oc_010df6b42b975ce056cc7c2e717abde8` | 通过 |
| `+chat-search` | 按群名搜刚创建群 | User, `im:chat:read` | `./xfchat_cli im +chat-search --query 'xfchat-im-test-20260330' --as user --format json` | `total=0` | 失败，疑似索引延迟/可见性限制 |
| `+chat-search` | 按更新后群名搜索 | User, `im:chat:read` | `./xfchat_cli im +chat-search --query 'xfchat-im-test-20260330-b-updated' --as user --format json` | `total=0` | 失败 |
| `+chat-search` | 按成员组合搜索 | User, `im:chat:read` | `./xfchat_cli im +chat-search --member-ids 'ou_6b90260ffbda660ec3e47f27c0871ec9,ou_4dce96856e654228d968f1ae74e29ecd,ou_981bb29d80e31c4aebe79225e866d8ff' --as user --format json` | `total=0` | 失败 |
| `+messages-send` | 群文本消息 | Bot, `im:message:send_as_bot` | `... --chat-id oc_c839... --text 'xfchat text message test'` | `message_id=om_x100b5396df2cbca03851d2653299389` | 通过 |
| `+messages-send` | 群 Markdown 消息 | Bot, `im:message:send_as_bot` | `... --chat-id oc_c839... --markdown $'# xfchat markdown test\n- item 1\n- item 2'` | `message_id=om_x100b5396df2cd8a0385543fcbf33978` | 通过 |
| `+messages-send` | 群文件消息 | Bot, `im:message:send_as_bot`, `im:resource` | `... --chat-id oc_c839... --file ./sample.txt` | `message_id=om_x100b5396ddedc0a0386b90e18cce054` | 通过 |
| `+messages-send` | 群图片消息 | Bot, `im:message:send_as_bot`, `im:resource` | `... --chat-id oc_c839... --image ./sample.png` | `message_id=om_x100b5396dded3ca0386a2497bf613bf` | 通过 |
| `+messages-send` | 媒体绝对路径校验 | Bot | `--file /abs/path/sample.txt` / `--image /abs/path/sample.png` | 返回 validation，要求相对路径 | 符合预期 |
| `+chat-messages-list` | 群消息读取 | Bot, `im:message:readonly` | `./xfchat_cli im +chat-messages-list --chat-id oc_c839... --as bot --format json` | 返回文本/Markdown/文件/图片/系统消息 | 通过 |
| `+chat-messages-list` | 群消息读取 | User | `./xfchat_cli im +chat-messages-list --chat-id oc_c839... --as user --format json` | `99991668 user access token not support` | 私有化不支持 |
| `+chat-messages-list` | 多成员群分页 | Bot, `im:message:readonly` | `./xfchat_cli im +chat-messages-list --chat-id oc_010df6... --as bot --page-size 3 --sort desc --format json` | `has_more=true`，返回 mention 消息和系统消息 | 通过 |
| `+chat-messages-list` | 私聊读取 | Bot / User | `./xfchat_cli im +chat-messages-list --user-id ou_6b9026... --as bot --page-size 12 --sort desc`；`--as user` 同样 | `failed to parse chat_p2p response: invalid character 'p' after top-level value` | 失败，CLI 私聊解析问题 |
| `+chat-messages-list` | 指向别人私聊 | User | `./xfchat_cli im +chat-messages-list --user-id ou_4dce... --as user --page-size 12 --sort desc` | 同样报 `chat_p2p` 解析错误 | 当前不可用 |
| `+messages-mget` | 群消息批量读取 | Bot, `im:message:readonly` | `./xfchat_cli im +messages-mget --message-ids 'om_x100b5396df2cbca03851d2653299389,om_x100b5396df2cd8a0385543fcbf33978' --as bot` | 返回 2 条消息详情 | 通过 |
| `+messages-mget` | 群消息批量读取 | User | `./xfchat_cli im +messages-mget --message-ids 'om_x100b5396df2cbca03851d2653299389' --as user` | `230027 Permission denied` | 私有化受限 |
| `+messages-mget` | 私聊消息批量读取 | Bot, `im:message:readonly` | `./xfchat_cli im +messages-mget --message-ids 'om_x100b5396a72d64a0385803fd1b9dac0,om_x100b5396a72c90a0385bc4767a397be,om_x100b5396a72a94a03854c50c32a201f,om_x100b5396a72b4ca0385edc636d650d4' --as bot` | 返回文本、Markdown、文件、图片 4 条 | 通过 |
| `+messages-search` | 群消息搜索 | User, `search:message` | `./xfchat_cli im +messages-search --query 'xfchat' --chat-id oc_c839... --as user --format json` | HTTP 404 | 私有化接口不存在 |
| `+messages-reply` | 群主会话回复 | Bot, `im:message:send_as_bot` | `... --message-id om_x100b5396df2cbca03851d2653299389 --text 'xfchat reply test'` | `message_id=om_x100b5396db850ca438551a1eae1e111` | 通过 |
| `+messages-reply` | 群线程回复 | Bot, `im:message:send_as_bot` | `... --reply-in-thread` | `message_id=om_x100b5396db855ca03854c3712150623`，原消息出现 `thread_id=omt_1aafa1d7350f1429` | 通过 |
| `+messages-reply` | 私聊主会话回复 | Bot, `im:message:send_as_bot` | `./xfchat_cli im +messages-reply --message-id om_x100b5396a72d64a0385803fd1b9dac0 --text 'xfchat_cli 私聊留痕：主会话回复' --as bot` | `message_id=om_x100b5396a445b4a03860a9a6886e5c4` | 通过 |
| `+messages-reply` | 私聊线程回复 | Bot, `im:message:send_as_bot` | `./xfchat_cli im +messages-reply --message-id om_x100b5396a72d64a0385803fd1b9dac0 --text 'xfchat_cli 私聊留痕：线程回复' --reply-in-thread --as bot` | `message_id=om_x100b5396a44558a0385d100918bb009` | 通过 |
| `+threads-messages-list` | 群线程读取 | Bot, `im:message:readonly` | `./xfchat_cli im +threads-messages-list --thread omt_1aafa1d7350f1429 --as bot --format json` | 返回原消息与线程回复共 2 条 | 通过 |
| `+threads-messages-list` | 群线程读取 | User, `im:message:readonly` | `./xfchat_cli im +threads-messages-list --thread omt_1aafa1d7350f1429 --as user --format json` | `99991668 user access token not support` | 私有化不支持 |
| `+threads-messages-list` | 私聊线程读取 | Bot, `im:message:readonly` | `./xfchat_cli im +threads-messages-list --thread om_x100b5396a72d64a0385803fd1b9dac0 --as bot --format json` | 返回原消息和线程回复共 2 条，`thread_id=omt_1aafa657288f1429` | 通过 |
| `+messages-resources-download` | 群图片下载 | Bot, `im:resource` | `... --message-id om_x100b5396dded3ca0386a2497bf613bf --file-key img_v3_02109_72b74497-338b-4806-8b1b-5bfe4fb52fnh --type image --as bot --output tmp/im-downloads/downloaded.png` | 保存成功，`size_bytes=68` | 通过 |
| `+messages-resources-download` | 群文件下载 | Bot, `im:resource` | `... --message-id om_x100b5396ddedc0a0386b90e18cce054 --file-key file_v3_00109_8b77d5b6-860c-4339-a5f6-4849892bd8nh --type file --as bot --output tmp/im-downloads/downloaded.txt` | 保存成功，`size_bytes=20` | 通过 |
| `+messages-resources-download` | 群图片下载 | User, `im:resource` | `./xfchat_cli im +messages-resources-download --message-id om_x100b5396dded3ca0386a2497bf613bf --file-key img_v3_02109_72b74497-338b-4806-8b1b-5bfe4fb52fnh --type image --as user --output tmp/im-downloads/user-downloaded.png` | HTTP 400 / `99991668 user access token not support` | 私有化不支持 |
| `+messages-resources-download` | 私聊图片下载 | Bot, `im:resource` | `./xfchat_cli im +messages-resources-download --message-id om_x100b5396a72b4ca0385edc636d650d4 --file-key img_v3_02109_41ec59d3-1d36-4780-b7f3-c30adb6ba4nh --type image --as bot --output tmp/im-downloads/p2p-downloaded.png` | 保存成功，`size_bytes=68` | 通过 |
| `+messages-resources-download` | 私聊文件下载 | Bot, `im:resource` | `./xfchat_cli im +messages-resources-download --message-id om_x100b5396a72a94a03854c50c32a201f --file-key file_v3_00109_a7b37736-8c78-4b5e-a194-424f1baa60nh --type file --as bot --output tmp/im-downloads/p2p-downloaded.txt` | 保存成功，`size_bytes=20` | 通过 |
| `+chat-update` | 更新群描述 | User, `im:chat:update` | `./xfchat_cli im +chat-update --chat-id oc_010df6... --description 'updated by user shortcut' --as user` | 返回 `chat_id=oc_010df6...` | 通过 |

## 4. 原子 API 覆盖

### 4.1 Chats

| 方法 | 场景 | 权限/身份 | 输入 | 输出摘要 | 结论 |
|---|---|---|---|---|---|
| `chats.get` | 群信息 | Bot, `im:chat:read` | `./xfchat_cli im chats get --params '{"chat_id":"oc_c839...2258","user_id_type":"open_id"}' --as bot` | 返回群名称、人数、配置项 | 通过 |
| `chats.get` | 群信息 | User, `im:chat:read` | 同上 `--as user` | 成功返回群详情 | 通过 |
| `chats.list` | 群列表 | Bot, `im:chat:read` | `./xfchat_cli im chats list --params '{"page_size":20,"user_id_type":"open_id"}' --as bot` | 返回机器人所在群列表，包含测试群 | 通过 |
| `chats.list` | 群列表 | User, `im:chat:read` | 同上 `--as user` | 成功返回用户群列表，`has_more=true` | 通过 |
| `chats.link` | 周有效期 | Bot, `im:chat:read` | `./xfchat_cli im chats link --params '{"chat_id":"oc_c839...2258"}' --data '{"validity_period":"week"}' --as bot` | 返回分享链接和过期时间 | 通过 |
| `chats.link` | 周有效期 | User, `im:chat:read` | `./xfchat_cli im chats link --params '{"chat_id":"oc_010df6b42b975ce056cc7c2e717abde8"}' --data '{"validity_period":"week"}' --as user` | 返回分享链接和过期时间 | 通过 |
| `chats.link` | 年有效期 | Bot, `im:chat:read` | `./xfchat_cli im chats link --params '{"chat_id":"oc_010df6b42b975ce056cc7c2e717abde8"}' --data '{"validity_period":"year"}' --as bot` | 返回 `expire_time`，`is_permanent=false` | 通过 |
| `chats.link` | 永久有效 | Bot, `im:chat:read` | `./xfchat_cli im chats link --params '{"chat_id":"oc_010df6b42b975ce056cc7c2e717abde8"}' --data '{"validity_period":"permanently"}' --as bot` | 返回 `is_permanent=true` | 通过 |
| `chats.link` | 永久有效 | User, `im:chat:read` | `./xfchat_cli im chats link --params '{"chat_id":"oc_010df6b42b975ce056cc7c2e717abde8"}' --data '{"validity_period":"permanently"}' --as user` | 返回 `is_permanent=true` | 通过 |
| `chats.update` | 更新群名与描述 | Bot, `im:chat:update` | `./xfchat_cli im chats update --params '{"chat_id":"oc_010d...bde8"}' --data '{"name":"xfchat-im-test-20260330-b-updated","description":"updated by atomic api"}' --as bot` | `msg=success` | 通过 |

### 4.2 Chat Members

| 方法 | 场景 | 权限/身份 | 输入 | 输出摘要 | 结论 |
|---|---|---|---|---|---|
| `chat.members.get` | 单成员群成员 | Bot, `im:chat.members:read` | `./xfchat_cli im chat.members get --params '{"chat_id":"oc_c839...2258","member_id_type":"open_id","user_id_type":"open_id","page_size":50}' --as bot` | 返回用户 `ou_6b90...1ec9` | 通过 |
| `chat.members.get` | 单成员群成员 | User, `im:chat.members:read` | 同上 `--as user` | 成功返回成员列表 | 通过 |
| `chat.members.get` | 多成员群成员 | Bot, `im:chat.members:read` | `./xfchat_cli im chat.members get --params '{"chat_id":"oc_010df6b42b975ce056cc7c2e717abde8","member_id_type":"open_id","user_id_type":"open_id","page_size":100}' --as bot` | `member_total=4` | 通过 |
| `chat.members.get` | 多成员群成员 | User, `im:chat.members:read` | 同上 `--as user` | `member_total=4` | 通过 |
| `chat.members.create` | 重复拉单人 | Bot, `im:chat.members:write_only` | `./xfchat_cli im chat.members create --params '{"chat_id":"oc_010d...bde8","member_id_type":"open_id","succeed_type":1}' --data '{"id_list":["ou_6b90...1ec9"]}' --as bot` | `invalid_id_list=[]` | 通过 |
| `chat.members.create` | 重复拉单人 | User, `im:chat.members:write_only` | 同上 `--as user` | `invalid_id_list=[]` | 通过 |
| `chat.members.create` | 拉人到成员测试群 | User, `im:chat.members:write_only` | `./xfchat_cli im chat.members create --params '{"chat_id":"oc_c7d9abd30a87e6d16976966818e2c19a","member_id_type":"open_id","succeed_type":1}' --data '{"id_list":["ou_981bb29d80e31c4aebe79225e866d8ff"]}' --as user` | 返回 success | 通过，存在最终一致性延迟 |
| `chat.members.create` | 批量拉人 | Bot, `im:chat.members:write_only` | `./xfchat_cli im chat.members create --params '{"chat_id":"oc_c7d9abd30a87e6d16976966818e2c19a","member_id_type":"open_id","succeed_type":1}' --data '{"id_list":["ou_84c79a9284092a69283066a89e549251","ou_981bb29d80e31c4aebe79225e866d8ff"]}' --as bot` | 返回 success，无 invalid/not_existed | 通过 |
| `chat.members.create` | 重复拉人幂等 | Bot, `im:chat.members:write_only` | `./xfchat_cli im chat.members create --params '{"chat_id":"oc_c7d9abd30a87e6d16976966818e2c19a","member_id_type":"open_id","succeed_type":1}' --data '{"id_list":["ou_981bb29d80e31c4aebe79225e866d8ff"]}' --as bot` | 仍返回 success | 通过，幂等 |
| `chat.members.create` | 非法成员 ID | Bot, `im:chat.members:write_only` | `./xfchat_cli im chat.members create --params '{"chat_id":"oc_c7d9abd30a87e6d16976966818e2c19a","member_id_type":"open_id","succeed_type":1}' --data '{"id_list":["ou_invalid_test_id"]}' --as bot` | `not_existed_id_list=["ou_invalid_test_id"]` | 通过 |

### 4.3 Messages

| 方法 | 场景 | 权限/身份 | 输入 | 输出摘要 | 结论 |
|---|---|---|---|---|---|
| `messages.read_users` | 群消息已读 | Bot, `im:message:readonly` | `./xfchat_cli im messages read_users --params '{"message_id":"om_x100b5396df2cbca03851d2653299389","user_id_type":"open_id"}' --as bot` | 返回已读用户 `ou_6b90...1ec9` | 通过 |
| `messages.read_users` | 私聊消息已读 | Bot, `im:message:readonly` | `./xfchat_cli im messages read_users --params '{"message_id":"om_x100b5396a72d64a0385803fd1b9dac0","user_id_type":"open_id"}' --as bot` | 返回已读用户 `ou_6b90260ffbda660ec3e47f27c0871ec9` | 通过 |
| `messages.forward` | 群 A -> 群 B | Bot, `im:message` | `./xfchat_cli im messages forward --params '{"message_id":"om_x100b5396df2cbca03851d2653299389","receive_id_type":"chat_id"}' --data '{"receive_id":"oc_010d...bde8"}' --as bot` | `message_id=om_x100b5396d43374a0385886280912a79` | 通过 |
| `messages.merge_forward` | 群 A -> 群 B | Bot, `im:message` | `./xfchat_cli im messages merge_forward --params '{"receive_id_type":"chat_id"}' --data '{"receive_id":"oc_010d...bde8","message_id_list":["om_x100b5396df2cbca03851d2653299389","om_x100b5396df2cd8a0385543fcbf33978"]}' --as bot` | `message_id=om_x100b5396d433cca8386580caaf49420` | 通过 |
| `messages.delete` | 撤回转发消息 | Bot, `im:message:recall` | `./xfchat_cli im messages delete --params '{"message_id":"om_x100b5396d43374a0385886280912a79"}' --as bot` | `msg=success` | 通过 |
| `messages.delete` | 私聊撤回 Bot 发给用户的消息 | User, `im:message:recall` | `./xfchat_cli im messages delete --params '{"message_id":"om_x100b5396a72d64a0385803fd1b9dac0"}' --as user` | `230026 No permission to recall this message` | 失败，用户不能撤回对方发出的 P2P 消息 |
| `messages.forward` | 群 B -> 群 A | Bot, `im:message` | `./xfchat_cli im messages forward --params '{"message_id":"om_x100b53968070c8a03859494c87de645","receive_id_type":"chat_id"}' --data '{"receive_id":"oc_c839c827acf6c34ca6287ea171a52258"}' --as bot` | `message_id=om_x100b53969c2300a03863c8f09d35f5a` | 通过 |
| `messages.forward` | 群 B Markdown -> 群 A | Bot, `im:message` | `./xfchat_cli im messages forward --params '{"message_id":"om_x100b5396807300a4386e3cdc0c477ae","receive_id_type":"chat_id"}' --data '{"receive_id":"oc_c839c827acf6c34ca6287ea171a52258"}' --as bot` | `message_id=om_x100b53969a973ca83864c9622dbf7b6` | 通过 |
| `messages.merge_forward` | 群 B 多消息 -> 群 A | Bot, `im:message` | `./xfchat_cli im messages merge_forward --params '{"receive_id_type":"chat_id"}' --data '{"receive_id":"oc_c839c827acf6c34ca6287ea171a52258","message_id_list":["om_x100b53968070c8a03859494c87de645","om_x100b5396807300a4386e3cdc0c477ae","om_x100b5396807300a038634cc1ada61f9"]}' --as bot` | `message_id=om_x100b53969a94f4a0386aa2d482df98f` | 通过 |
| `messages.forward` | 群 A -> 群 B 回转发 | Bot, `im:message` | `./xfchat_cli im messages forward --params '{"message_id":"om_x100b53969c2300a03863c8f09d35f5a","receive_id_type":"chat_id"}' --data '{"receive_id":"oc_010df6b42b975ce056cc7c2e717abde8"}' --as bot` | `message_id=om_x100b53969a9604a03857b5e694a28ab` | 通过 |
| `messages.forward` | 私聊自转发 | Bot, `im:message` | `./xfchat_cli im messages forward --params '{"message_id":"om_x100b5396a72d64a0385803fd1b9dac0","receive_id_type":"chat_id"}' --data '{"receive_id":"oc_934891fa0c30e80e4d7e7d5d7496d5e4"}' --as bot` | `message_id=om_x100b5396a2fb20a43868325a20b08ff` | 通过 |
| `messages.merge_forward` | 私聊合并转发 | Bot, `im:message` | `./xfchat_cli im messages.merge_forward --params '{"receive_id_type":"chat_id"}' --data '{"receive_id":"oc_934891fa0c30e80e4d7e7d5d7496d5e4","message_id_list":["om_x100b5396a72d64a0385803fd1b9dac0","om_x100b5396a72c90a0385bc4767a397be"]}' --as bot` | `message_id=om_x100b5396a2f874a0385a95fa6224813` | 通过 |
| 原生 `POST /open-apis/im/v1/messages` | 用户态发群消息 | User, `im:message` | `./xfchat_cli api POST /open-apis/im/v1/messages --as user ...` | 进程返回 0，但无响应体，回读群消息未见消息落地 | 未打通 |

### 4.4 Reactions

| 方法 | 场景 | 权限/身份 | 输入 | 输出摘要 | 结论 |
|---|---|---|---|---|---|
| `reactions.create` | 群消息加 `SMILE` | Bot, `im:message.reactions:write_only` | `./xfchat_cli im reactions create --params '{"message_id":"om_x100b5396df2cbca03851d2653299389"}' --data '{"reaction_type":{"emoji_type":"SMILE"}}' --as bot` | `reaction_id=u3pf...TQT` | 通过 |
| `reactions.list` | 群消息读表情 | Bot, `im:message.reactions:read` | `./xfchat_cli im reactions list --params '{"message_id":"om_x100b5396df2cbca03851d2653299389","user_id_type":"open_id"}' --as bot` | 能查到 `SMILE` | 通过 |
| `reactions.delete` | 群消息删表情 | Bot, `im:message.reactions:write_only` | `./xfchat_cli im reactions delete --params '{"message_id":"om_x100b5396df2cbca03851d2653299389","reaction_id":"u3pf...TQT"}' --as bot` | 删除成功；后续 list 为空 | 通过 |
| `reactions.batch_query` | 群消息批量查表情 | Bot, `im:message.reactions:read` | `./xfchat_cli im reactions batch_query --params '{"user_id_type":"open_id"}' --data '{"queries":[{"message_id":"om_x100b5396df2cbca03851d2653299389"}],"page_size_per_message":10}' --as bot` | HTTP 404 | 私有化接口不存在 |
| `reactions.create` | 群消息加 `THUMBSUP` | User, `im:message.reactions:write_only` | `./xfchat_cli im reactions create --params '{"message_id":"om_x100b5396df2cbca03851d2653299389"}' --data '{"reaction_type":{"emoji_type":"THUMBSUP"}}' --as user` | `reaction_id=u3pf...Hg==` | 通过 |
| `reactions.list` | 群消息查 User 表情 | User, `im:message.reactions:read` | `./xfchat_cli im reactions list --params '{"message_id":"om_x100b5396df2cbca03851d2653299389","user_id_type":"open_id"}' --as user` | 能查到 `THUMBSUP` | 通过 |
| `reactions.delete` | 群消息删 User 表情 | User, `im:message.reactions:write_only` | `./xfchat_cli im reactions delete --params '{"message_id":"om_x100b5396df2cbca03851d2653299389","reaction_id":"u3pf...Hg=="}' --as user` | 删除成功 | 通过 |
| `reactions.create` | 私聊文本加 `THUMBSUP` | User, `im:message.reactions:write_only` | `./xfchat_cli im reactions create --params '{"message_id":"om_x100b5396a72d64a0385803fd1b9dac0"}' --data '{"reaction_type":{"emoji_type":"THUMBSUP"}}' --as user` | `reaction_id=u-4K...A==` | 通过 |
| `reactions.list` | 私聊文本查 User 表情 | Bot, `im:message.reactions:read` | `./xfchat_cli im reactions list --params '{"message_id":"om_x100b5396a72d64a0385803fd1b9dac0","user_id_type":"open_id"}' --as bot` | 能查到 `THUMBSUP` | 通过 |
| `reactions.delete` | 私聊文本删 User 表情 | User, `im:message.reactions:write_only` | `./xfchat_cli im reactions delete --params '{"message_id":"om_x100b5396a72d64a0385803fd1b9dac0","reaction_id":"u-4K3...A=="}' --as user` | 删除成功 | 通过 |
| `reactions.create` | 私聊 Markdown 加 `SMILE` | Bot, `im:message.reactions:write_only` | `./xfchat_cli im reactions create --params '{"message_id":"om_x100b5396a72c90a0385bc4767a397be"}' --data '{"reaction_type":{"emoji_type":"SMILE"}}' --as bot` | `reaction_id=hoc_hrww...Jq` | 通过 |
| `reactions.list` | 私聊 Markdown 查 Bot 表情 | Bot, `im:message.reactions:read` | `./xfchat_cli im reactions list --params '{"message_id":"om_x100b5396a72c90a0385bc4767a397be","user_id_type":"open_id"}' --as bot` | 能查到 `SMILE` | 通过 |
| `reactions.delete` | 私聊 Markdown 删 Bot 表情 | Bot, `im:message.reactions:write_only` | `./xfchat_cli im reactions delete --params '{"message_id":"om_x100b5396a72c90a0385bc4767a397be","reaction_id":"hoc_hrww...Jq"}' --as bot` | 删除成功 | 通过 |

### 4.5 Pins

| 方法 | 场景 | 权限/身份 | 输入 | 输出摘要 | 结论 |
|---|---|---|---|---|---|
| `pins.create` | 群文本 Pin | Bot, `im:message.pins:write_only` | `./xfchat_cli im pins create --data '{"message_id":"om_x100b5396df2cbca03851d2653299389"}' --as bot` | 创建成功 | 通过 |
| `pins.list` | 群文本 Pin 列表 | Bot, `im:message.pins:read` | `./xfchat_cli im pins list --params '{"chat_id":"oc_c839...2258","page_size":20}' --as bot` | 能查到该 pin | 通过 |
| `pins.delete` | 群文本删 Pin | Bot, `im:message.pins:write_only` | `./xfchat_cli im pins delete --params '{"message_id":"om_x100b5396df2cbca03851d2653299389"}' --as bot` | 删除成功；后续 list 为空 | 通过 |
| `pins.create` | 群 Markdown Pin | User, `im:message.pins:write_only` | `./xfchat_cli im pins create --data '{"message_id":"om_x100b5396df2cd8a0385543fcbf33978"}' --as user` | 创建成功，操作者为 `ou_6b90...1ec9` | 通过 |
| `pins.list` | 群 Markdown Pin 列表 | User, `im:message.pins:read` | `./xfchat_cli im pins list --params '{"chat_id":"oc_c839c827acf6c34ca6287ea171a52258","page_size":20}' --as user` | `items=[]` | 异常，对自己不可见 |
| `pins.list` | 群 Markdown Pin 交叉验证 | Bot, `im:message.pins:read` | 同上 `--as bot` | 能查到用户刚创建的 Pin | 通过，说明数据已写入 |
| `pins.delete` | 群 Markdown 删 Pin | User, `im:message.pins:write_only` | `./xfchat_cli im pins delete --params '{"message_id":"om_x100b5396df2cd8a0385543fcbf33978"}' --as user` | 删除成功 | 通过 |
| `pins.create` | 私聊文本 Pin | Bot, `im:message.pins:write_only` | `./xfchat_cli im pins create --data '{"message_id":"om_x100b5396a72d64a0385803fd1b9dac0"}' --as bot` | 创建成功 | 通过 |
| `pins.list` | 私聊文本 Pin 列表 | Bot, `im:message.pins:read` | `./xfchat_cli im pins list --params '{"chat_id":"oc_934891fa0c30e80e4d7e7d5d7496d5e4","page_size":20}' --as bot` | 能查到该 pin | 通过 |
| `pins.create` | 私聊 Markdown Pin | User, `im:message.pins:write_only` | `./xfchat_cli im pins create --data '{"message_id":"om_x100b5396a72c90a0385bc4767a397be"}' --as user` | 创建成功 | 通过 |
| `pins.list` | 私聊 Pin 列表 | User, `im:message.pins:read` | `./xfchat_cli im pins list --params '{"chat_id":"oc_934891fa0c30e80e4d7e7d5d7496d5e4","page_size":20}' --as user` | 能看到 User Pin 和 Bot Pin 共 2 条 | 通过 |
| `pins.delete` | 私聊 Markdown 删 Pin | User, `im:message.pins:write_only` | `./xfchat_cli im pins delete --params '{"message_id":"om_x100b5396a72c90a0385bc4767a397be"}' --as user` | 删除成功 | 通过 |

### 4.6 Images

| 方法 | 场景 | 权限/身份 | 输入 | 输出摘要 | 结论 |
|---|---|---|---|---|---|
| `images.create` | 原子图片上传 | Bot, `im:resource` | `/Users/wangqizhao/Developer/iflytek/cli/xfchat_cli im images create --data '{"image_type":"message","image":"@./sample.png"}' --as bot` | `234001 Invalid request param` | 失败，CLI multipart 参数格式待修复 |

## 5. 关键留痕与消息对象

### 5.1 群聊可见留痕

群 `oc_010df6b42b975ce056cc7c2e717abde8` 留痕：

- 文本：`om_x100b53968070c8a03859494c87de645`
- Markdown：`om_x100b5396807300a4386e3cdc0c477ae`
- @提及：`om_x100b5396807300a038634cc1ada61f9`
- 主会话回复：`om_x100b5396818854a83865c9d088d7534`
- 线程回复：`om_x100b5396818b38a038697474f4c6c1a`
- 文件：`om_x100b53969ebd80a03852c73976c009a`
- 图片：`om_x100b53969ebcd8a0386cbb5346f1977`
- 转发：`om_x100b53968107c8a0385e81e0fcd8b3f`
- 合并转发：`om_x100b53968107bca038551a10ebce558`
- 用户表情 `THUMBSUP`：挂在 `om_x100b53968070c8a03859494c87de645`
- Bot 表情 `SMILE`：挂在 `om_x100b5396807300a4386e3cdc0c477ae`
- Bot Pin：挂在 `om_x100b53968070c8a03859494c87de645`
- User Pin：挂在 `om_x100b5396807300a4386e3cdc0c477ae`

群 `oc_c839c827acf6c34ca6287ea171a52258` 留痕：

- `B -> A` 文本转发：`om_x100b53969c2300a03863c8f09d35f5a`
- `B -> A` Markdown 转发：`om_x100b53969a973ca83864c9622dbf7b6`
- `B -> A` 合并转发：`om_x100b53969a94f4a0386aa2d482df98f`

成员测试群 `oc_c7d9abd30a87e6d16976966818e2c19a` 留痕：

- Bot 文本：`om_x100b539762fca8a03863168c3eb02fa`
- 系统消息：用户入群与批量拉人消息已保留

### 5.2 私聊可见留痕

私聊 `oc_934891fa0c30e80e4d7e7d5d7496d5e4` 留痕：

- 文本：`om_x100b5396a72d64a0385803fd1b9dac0`
- Markdown：`om_x100b5396a72c90a0385bc4767a397be`
- 文件：`om_x100b5396a72a94a03854c50c32a201f`
- 图片：`om_x100b5396a72b4ca0385edc636d650d4`
- 主会话回复：`om_x100b5396a445b4a03860a9a6886e5c4`
- 线程回复：`om_x100b5396a44558a0385d100918bb009`
- 转发：`om_x100b5396a2fb20a43868325a20b08ff`
- 合并转发：`om_x100b5396a2f874a0385a95fa6224813`
- User 表情 `THUMBSUP`：挂在 `om_x100b5396a72d64a0385803fd1b9dac0`
- Bot 表情 `SMILE`：挂在 `om_x100b5396a72c90a0385bc4767a397be`
- Bot Pin：挂在 `om_x100b5396a72d64a0385803fd1b9dac0`
- User Pin：挂在 `om_x100b5396a72c90a0385bc4767a397be`

## 6. 私有化差异、CLI 问题与结论

### 6.1 私有化环境已确认差异

- User 身份读取群消息：`+chat-messages-list --as user` 返回 `99991668 user access token not support`
- User 身份读取群线程：`+threads-messages-list --as user` 返回 `99991668 user access token not support`
- User 身份下载群消息资源：`+messages-resources-download --as user` 返回 `99991668 user access token not support`
- User 身份搜索消息：`+messages-search --as user` 返回 HTTP 404
- `reactions.batch_query`：HTTP 404
- `+chat-search` 对测试群和成员组合均未命中，存在索引/可见性问题
- 用户态发消息：原生 `POST /open-apis/im/v1/messages --as user` 返回异常静默，未验证消息落地

### 6.2 CLI 兼容问题

- `images.create` 原子接口：CLI multipart 适配未打通，返回 `234001 Invalid request param`
- 群聊 `pins.list --as user`：用户自己创建的 Pin 在自己视角下返回空，但 Bot 视角可见，存在视图不一致
- 私聊按 `user_id` 解析 `chat_id`：依赖了未公开且在私有化环境中不存在的接口 `POST /open-apis/im/v1/chat_p2p/batch_query`

### 6.3 `chat_p2p/batch_query` 根因

当前 `+chat-messages-list --user-id ...` 的失败，不是消息读取权限本身被拒绝，而是前置私聊 `chat_id` 解析失败。

CLI 当前使用的接口：

```text
POST /open-apis/im/v1/chat_p2p/batch_query
```

实测：

```bash
./xfchat_cli api POST /open-apis/im/v1/chat_p2p/batch_query --as bot --params '{"chatter_id_type":"open_id"}' --data '{"chatter_ids":["ou_6b90260ffbda660ec3e47f27c0871ec9"]}'
./xfchat_cli api POST /open-apis/im/v1/chat_p2p/batch_query --as user --params '{"chatter_id_type":"open_id"}' --data '{"chatter_ids":["ou_6b90260ffbda660ec3e47f27c0871ec9"]}'
```

返回：

```text
HTTP 404: 404 page not found
```

代码位置：

- 私聊 `chat_id` 解析逻辑在 [shortcuts/im/helpers.go](/Users/wangqizhao/Developer/iflytek/cli/shortcuts/im/helpers.go)
- 当前实现直接假定该接口返回 JSON，并尝试解析 `p2p_chats`
- 私有化环境实际返回纯文本 `404 page not found`

因此出现报错：

```text
failed to parse chat_p2p response: invalid character 'p' after top-level value
```

这条报错的真实含义是“接口不存在但被误当成 JSON 解析”，不是标准业务权限报错。

补充确认：

- 官方公开文档中未检索到 `chat_p2p/batch_query`
- 当前本地 IM registry 中也未收录该接口
- 该路径目前只出现在本项目 IM shortcut 的内部解析逻辑和测试桩里

## 7. 覆盖完成度

### 7.1 已覆盖情况

- 按当前 CLI 已注册的 IM 命令面，除明确缺失的踢人接口外，现有命令均已覆盖
- 已覆盖 user/bot、群聊/私聊、成功/失败/异常输入三类边界
- 已保留关键会话留痕，便于人工复核

### 7.2 当前无法打通或无法继续覆盖的项

- 踢人
  当前 CLI registry 只注册了 `chat.members.create/get`，未收录 `delete/remove`
- 私聊历史按 `user_id` 直接读取
  当前被 `chat_p2p/batch_query` 卡住，未修复前无法继续
- User 群消息读取与搜索
  私有化环境已明确返回 `99991668` 或 404
- 原生 `images.create`
  当前是 CLI multipart 适配问题，待修复后可回归
- 原生 `POST /open-apis/im/v1/messages --as user`
  当前返回异常静默，未能确认消息落地

### 7.3 总结

- Bot 主链路已打通：建群、发消息、回消息、线程、批量查消息、下载资源、转发、合并转发、撤回、表情、Pin、群管理、私聊主链路
- User 可用能力已确认：`chats.get/list/link`、`chat.members.get/create`、`+chat-update`、群与私聊表情增删查、部分 Pin 能力
- User 不可用或受限能力已确认：群消息读取、群线程读取、群资源下载、消息搜索、`+messages-mget`、稳定的用户态发消息
- IM 模块当前可以按“Bot 负责消息面，User 负责部分群管理面与交互面”的模式在 i 讯飞私有化环境下落地
