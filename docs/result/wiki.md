# Wiki 模块测试结果

## 1. 环境与鉴权

- CLI：`./xfchat_cli`
- 私有化 OpenAPI：`https://open.xfchat.iflytek.com`
- 应用：`cli_a9084c0706b8d379`
- User：`王琦钊 qzwang9`（`ou_6b90260ffbda660ec3e47f27c0871ec9`）
- Bot：`星火派-test`
- User 已授权 Wiki scope：`wiki:node:read`

## 2. 测试对象与素材

- 知识库链接：`https://yf2ljykclb.xfchat.iflytek.com/wiki/WG34wO7Nei7MIkkykdCrscUIz1t`
- Wiki 节点 token：`WG34wO7Nei7MIkkykdCrscUIz1t`
- 解析出的真实文档 token：`doxrze7hMi3Hrd6QMU6uRopj6ih`
- 解析出的对象类型：`docx`
- 标题：`AI应用知识库门户`
- space_id：`7434440606543250293`

## 3. Shortcut / 原子 API 覆盖

当前 CLI 的 Wiki 模块仅暴露原子能力：

- `wiki spaces get_node`

对应接口：

- `GET /open-apis/wiki/v2/spaces/get_node`

Schema 显示的 scope：

- `wiki:wiki`
- `wiki:wiki:readonly`
- `wiki:node:read`

本次实际打通所需最小 scope：

- `wiki:node:read`

## 4. 测试结果

### 4.1 User 读取真实 wiki 节点

- 权限：`wiki:node:read`
- 身份：`user`
- 输入命令：

```bash
./xfchat_cli wiki spaces get_node --as user --params '{"token":"WG34wO7Nei7MIkkykdCrscUIz1t"}'
```

- 关键输出：
  - `code: 0`
  - `node.node_token: WG34wO7Nei7MIkkykdCrscUIz1t`
  - `node.obj_type: docx`
  - `node.obj_token: doxrze7hMi3Hrd6QMU6uRopj6ih`
  - `node.title: AI应用知识库门户`
  - `node.space_id: 7434440606543250293`
- 结论：`通过`

### 4.2 User 显式传 `obj_type=wiki`

- 权限：`wiki:node:read`
- 身份：`user`
- 输入命令：

```bash
./xfchat_cli wiki spaces get_node --as user --params '{"token":"WG34wO7Nei7MIkkykdCrscUIz1t","obj_type":"wiki"}'
```

- 关键输出：
  - `code: 0`
  - 返回内容与默认查询一致
  - `obj_type` 仍解析为 `docx`
- 结论：`通过`

### 4.3 User 错误地把 wiki token 当作 docx token 查询

- 权限：`wiki:node:read`
- 身份：`user`
- 输入命令：

```bash
./xfchat_cli wiki spaces get_node --as user --params '{"token":"WG34wO7Nei7MIkkykdCrscUIz1t","obj_type":"docx"}'
```

- 关键输出：
  - `code: 131005`
  - `message: resource not found: document not found by token WG34wO7Nei7MIkkykdCrscUIz1t`
- 结论：`失败`

### 4.4 Bot 读取同一 wiki 节点

- 权限：Bot 控制台需具备 Wiki 相关读取能力
- 身份：`bot`
- 输入命令：

```bash
./xfchat_cli wiki spaces get_node --as bot --params '{"token":"WG34wO7Nei7MIkkykdCrscUIz1t"}'
```

- 关键输出：
  - `code: 131006`
  - `message: permission denied: node permission denied`
- 结论：`缺权限`

### 4.5 User / Bot 对 fake wiki token 的行为

- 输入命令：

```bash
./xfchat_cli wiki spaces get_node --as user --params '{"token":"wikcn_test"}'
./xfchat_cli wiki spaces get_node --as bot --params '{"token":"wikcn_test"}'
```

- 关键输出：
  - 两侧都返回 `131005 resource not found`
- 结论：`通过`
  - 说明私有化 Wiki 接口存在且正常工作，不是 404 或 registry 缺失

### 4.6 User 用真实 DocX token 反查 Wiki 节点

- 权限：`wiki:node:read`
- 身份：`user`
- 输入命令：

```bash
./xfchat_cli wiki spaces get_node --as user --params '{"token":"doxrze7hMi3Hrd6QMU6uRopj6ih","obj_type":"docx"}'
```

- 关键输出：
  - `code: 0`
  - `node.node_token: WG34wO7Nei7MIkkykdCrscUIz1t`
  - `node.obj_token: doxrze7hMi3Hrd6QMU6uRopj6ih`
  - `node.title: AI应用知识库门户`
  - `node.space_id: 7434440606543250293`
- 结论：`通过`
  - 说明同一个知识库对象既可由 `wiki token` 查询，也可由真实 `docx token + obj_type=docx` 反查

### 4.7 User 用真实 DocX token 但不传 `obj_type`

- 权限：`wiki:node:read`
- 身份：`user`
- 输入命令：

```bash
./xfchat_cli wiki spaces get_node --as user --params '{"token":"doxrze7hMi3Hrd6QMU6uRopj6ih"}'
```

- 关键输出：
  - `code: 131005`
  - `message: resource not found`
- 结论：`失败`
  - 说明当输入是真实文档 token 时，必须显式传 `obj_type`

### 4.8 Bot 用真实 DocX token 反查 Wiki 节点

- 权限：Bot 控制台需具备 Wiki 相关读取能力
- 身份：`bot`
- 输入命令：

```bash
./xfchat_cli wiki spaces get_node --as bot --params '{"token":"doxrze7hMi3Hrd6QMU6uRopj6ih","obj_type":"docx"}'
```

- 关键输出：
  - `code: 131006`
  - `message: permission denied: node permission denied`
- 结论：`缺权限`

## 5. 关键留痕与对象

- 知识库节点：`WG34wO7Nei7MIkkykdCrscUIz1t`
- 文档对象：`doxrze7hMi3Hrd6QMU6uRopj6ih`
- 群内留痕消息：`om_x100b53933a43bca0385eceeffe68def`
- 留痕群：`oc_010df6b42b975ce056cc7c2e717abde8`

## 6. 私有化差异、CLI 问题与结论

- 私有化 Wiki `get_node` 接口可用，不是 404
- User 侧仅用 `wiki:node:read` 即可打通本次主链路
- Bot 侧对真实知识库节点返回 `131006 node permission denied`
- Wiki token 与 DocX token 不能混用：
  - Wiki token 默认按 `wiki` 语义查询可成功
  - Wiki token 强行按 `docx` 查询会报 `131005`
  - 真实 DocX token 只有在显式传 `obj_type=docx` 时才能反查回知识库节点
  - 真实 DocX token 不带 `obj_type` 会报 `131005`

当前未发现 CLI 兼容问题；现阶段差异主要在权限与对象类型理解上。

## 7. 覆盖完成度

- 已覆盖：`wiki spaces get_node` 的 user / bot、真实 token / fake token、`obj_type` 差异
- 未覆盖：空间创建、节点创建、移动、复制、成员管理
- 原因：当前 CLI 尚未暴露这些 Wiki 命令

结论：当前 CLI 的 Wiki 模块在私有化环境下已完成现有命令面的有效测试，主链路可用边界清晰。
