---
name: xfchat-module-test
description: "用于 xfchat_cli 各模块的标准化联调与回归测试。覆盖重新构建、鉴权来源校验、scope 检查、私有化差异判断、结果落盘，以及在权限缺失时输出给人类申请 scope 的清单。"
---

# xfchat-module-test

## 适用场景

当需要对 `xfchat_cli` 的某个模块做系统化测试时使用本 skill。适用于：

- 从源码重新构建后验证最新 CLI
- 检查当前鉴权是否来自 `xfchat_cli`，避免混用官方 `lark-cli`
- 按模块逐项测试 shortcut 和原子 API
- 区分 Bot / User 身份差异
- 在 i 讯飞私有化飞书环境下识别“权限缺失 / 私有化缺能力 / CLI 兼容问题 / registry 未收录”的差异
- 将结果写入 `docs/result/<module>.md`
- 当权限缺失时，给出明确的 scope 申请清单，交由人类开通后再重试

## 强制规则

### 1. 只能使用 `xfchat_cli`

- 不要使用全局 `lark-cli`
- 不要假设 `~/.lark-cli` 的配置可用
- 先确认当前 CLI 二进制为仓库内构建产物

推荐检查：

```bash
pwd
ls -l ./xfchat_cli
./xfchat_cli auth status
```

### 2. 先确认环境，再做模块测试

测试前必须先确认：

- 已完成最新构建
- 当前 endpoint 指向私有化环境
- `auth status` 中 user token 有效
- Bot 信息可正常获取

最低检查集：

```bash
make build
./xfchat_cli doctor --offline
./xfchat_cli auth status
./xfchat_cli api GET /open-apis/bot/v3/info --as bot
```

### 3. 每条测试都必须记录

每一条测试结果都必须至少包含：

- 用到的权限
- 身份：`user` / `bot`
- 输入命令
- 关键输出
- 结论：`通过` / `失败` / `私有化不支持` / `CLI 问题` / `缺权限`

### 4. 失败时不能直接下结论

发生失败后，必须按下面顺序判断：

1. CLI 输入是否错误
2. `docs/scope.json` 中是否存在对应 scope
3. 当前授权是否已包含该 scope
4. registry 中是否有该方法
5. 私有化环境是否返回 `404`、`99991668`、静默失败等能力缺失迹象
6. 是否属于 CLI 兼容问题

## 工作流

### 第一步：重新构建并确认鉴权来源

```bash
make build
./xfchat_cli auth status
```

需要确认：

- 当前 appId 正确
- 当前身份与用户名正确
- 使用的是 `~/.xfchat_cli` 配置，而不是官方 CLI 旧配置

### 第二步：选择模块并建立结果文件

结果文件路径固定为：

```text
docs/result/<module>.md
```

例如：

- `docs/result/im.md`
- `docs/result/drive.md`
- `docs/result/doc.md`

结果文件不要按“第几轮测试”累计。最终文档应按以下结构整理：

1. 环境与鉴权
2. 测试对象与素材
3. Shortcut 覆盖
4. 原子 API 覆盖
5. 关键留痕与对象
6. 私有化差异、CLI 问题与结论
7. 覆盖完成度

### 第三步：按模块拆测试矩阵

每个模块至少拆成四类：

- 主链路成功用例
- User / Bot 身份差异用例
- 异常输入 / 边界用例
- 私有化差异确认用例

如果模块有 shortcut 和原子 API，两类都要测。

### 第四步：失败时检查 scope

#### A. 先查本地权限字典

使用：

```bash
rg -n '<scope-name>|<关键词>' docs/scope.json -S
```

或者按模块关键词检索：

```bash
rg -n 'im:|drive:|doc:|calendar:' docs/scope.json -S
```

#### B. 判断权限状态

如果失败与权限有关，必须区分：

- `scope 存在，但当前 token 未授权`
- `scope 在私有化环境中不存在`
- `scope 名称变动，需要找平替`
- `不是权限问题，是私有化能力缺失或 CLI 问题`

#### C. 寻找同义或可平替权限

如果原 scope 不支持或未收录，要尝试检索同义 scope，例如：

- `read` vs `readonly`
- `write` vs `write_only`
- `message` vs `message:readonly` / `message.reactions:*` / `message.pins:*`
- `chat` vs `chat:read` / `chat:update` / `chat.members:*`

这一步必须写进结果文档。

## 人类申请权限流程

当判断为“缺权限”或“可能需要平替 scope”时，必须输出一段给人类执行的申请清单。

格式固定：

### 待申请权限

- 原需求权限：`<scope-a>`
- 可能平替权限：`<scope-b>`
- 用途：`<一句话说明该权限用于什么测试>`
- 触发命令：`<导致报错的命令>`
- 失败现象：`<错误码 / 错误消息>`

### 给人类的操作提示

- 请在私有化飞书开放平台中为应用开通以上权限
- 开通后重新执行登录授权
- 然后重新运行以下命令复测：

```bash
./xfchat_cli auth login --web --scope '<scope 列表>'
<复测命令>
```

### 什么时候要明确要求人类申请 scope

满足以下任一条件时，要明确让人类申请：

- `docs/scope.json` 中存在该 scope，但当前 `auth status` 未授予
- 业务错误码或错误消息明显指向权限缺失
- 当前 scope 命名不匹配，但能找到高概率平替权限

### 什么时候不要误导人类申请 scope

以下情况不要直接要求人类申请权限：

- HTTP 404，且更像接口不存在
- `99991668 user access token not support`
- CLI multipart 参数拼装错误
- registry 根本没有该命令或方法
- 返回静默异常，且无法证明是权限导致

这些要标记为：

- `私有化不支持`
- `CLI 问题`
- `registry 未收录`

## 结果落盘规范

结果文档中每项必须尽量保留：

- 完整命令
- 关键输出中的 message_id / chat_id / reaction_id / pin / share_link
- 成功与失败边界
- 对同一能力的 user / bot 差异
- 对私有化环境的明确结论

不要只写总结，不要只写“报错了”。

## 输出判定标签

统一使用以下标签之一：

- `通过`
- `失败`
- `私有化不支持`
- `CLI 问题`
- `缺权限`
- `registry 未收录`

## IM 模块的特殊规则

测试 IM 时，额外遵守：

- 同时覆盖群聊、多成员群、私聊
- 同时覆盖留痕能力：文本、Markdown、文件、图片、回复、线程、表情、Pin、转发
- 尽量把成功路径留在群里或私聊里，便于人工复核
- 对 `+chat-messages-list --user-id ...` 这类依赖 P2P 解析的命令，要特别记录是否依赖未公开接口

## 推荐执行模板

当用户说“开始测试 <module> 模块”时，按下面顺序执行：

1. 重构建 CLI
2. 检查 `auth status`
3. 检查 Bot 基本信息
4. 列出该模块 shortcut 与原子 API
5. 建立测试对象与测试素材
6. 跑成功路径
7. 跑 user/bot 差异
8. 跑异常与边界
9. 查 `docs/scope.json`
10. 输出待申请权限清单（如需要）
11. 将最终结构化结果写入 `docs/result/<module>.md`

## 成功标准

一个模块可视为“测试完成”，必须同时满足：

- 当前 CLI 已暴露命令都至少有结果
- 主链路成功项已验证
- user / bot 差异已验证
- 失败项已归因到“缺权限 / 私有化不支持 / CLI 问题 / 未收录”之一
- 结果文档已落盘
