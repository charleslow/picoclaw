# Telegram

Telegram Channel 通过 Telegram 机器人 API 使用长轮询实现基于机器人的通信。它支持文本消息、媒体附件（照片、语音、音频、文档）、通过 Groq Whisper 进行语音转录以及内置命令处理器。

## 配置

```json
{
  "channels": {
    "telegram": {
      "enabled": true,
      "token": "123456789:ABCdefGHIjklMNOpqrsTUVwxyz",
      "allow_from": ["123456789"],
      "proxy": ""
    }
  }
}
```

| 字段       | 类型   | 必填 | 描述                                                      |
| ---------- | ------ | ---- | --------------------------------------------------------- |
| enabled    | bool   | 是   | 是否启用 Telegram 频道                                    |
| token      | string | 是   | Telegram 机器人 API Token                                 |
| allow_from | array  | 否   | 用户ID白名单，空表示允许所有用户                          |
| proxy      | string | 否   | 连接 Telegram API 的代理 URL (例如 http://127.0.0.1:7890) |

## 设置流程

1. 在 Telegram 中搜索 `@BotFather`
2. 发送 `/newbot` 命令并按照提示创建新机器人
3. 获取 HTTP API Token
4. 将 Token 填入配置文件中
5. (可选) 配置 `allow_from` 以限制允许互动的用户 ID (可通过 `@userinfobot` 获取 ID)

## 多机器人配置

如果需要在同一个 picoclaw 网关中运行多个 Telegram 机器人，可以使用 `telegram_bots` 列表代替单个 `telegram` 配置。每个机器人拥有独立的 Token、会话和配置。

**注意：** `telegram` 和 `telegram_bots` 不能同时使用。如需多个机器人，请将所有机器人移到 `telegram_bots` 列表中。

```json
{
  "channels": {
    "telegram_bots": [
      {
        "id": "alice",
        "enabled": true,
        "token": "TOKEN_A",
        "allow_from": ["111"],
        "placeholder": { "enabled": true, "text": "思考中... 💭" }
      },
      {
        "id": "bob",
        "enabled": true,
        "token": "TOKEN_B",
        "allow_from": ["222"]
      }
    ]
  },
  "bindings": [
    { "match": { "channel": "telegram:alice" }, "agent_id": "agent-alice" },
    { "match": { "channel": "telegram:bob" },   "agent_id": "agent-bob" }
  ]
}
```

| 字段 | 类型   | 必填 | 描述                                                                  |
| ---- | ------ | ---- | --------------------------------------------------------------------- |
| id   | string | 是   | 机器人唯一标识符，只能包含小写字母、数字、下划线和连字符 (`[a-z0-9_-]`) |

其他字段（`token`、`allow_from`、`proxy`、`base_url` 等）与单机器人配置相同。每个机器人可以有独立的 `placeholder`、`group_trigger`、`typing` 和 `reasoning_channel_id` 配置。

每个机器人注册的通道名称为 `telegram:<id>`（例如 `telegram:alice`），可以在 `bindings` 中使用此名称将不同机器人路由到不同的 agent。

### 通过环境变量配置多机器人

如果不想使用配置文件，也可以通过环境变量快速配置多个机器人。将多个 `id:token` 对用逗号分隔，设置到 `PICOCLAW_CHANNELS_TELEGRAM_TOKEN` 环境变量中：

```bash
export PICOCLAW_CHANNELS_TELEGRAM_TOKEN="alice:TOKEN_A,bob:TOKEN_B"
```

启动时，系统会自动将其展开为 `telegram_bots` 配置，等同于：

```json
{
  "channels": {
    "telegram_bots": [
      { "id": "alice", "enabled": true, "token": "TOKEN_A" },
      { "id": "bob",   "enabled": true, "token": "TOKEN_B" }
    ]
  }
}
```

每个机器人会继承环境变量中设置的 `base_url`、`proxy` 和 `allow_from` 等共享配置（即 `PICOCLAW_CHANNELS_TELEGRAM_BASE_URL`、`PICOCLAW_CHANNELS_TELEGRAM_PROXY` 等）。

**注意：**
- 格式为 `id:token`，其中 `id` 是机器人标识符，`token` 是 BotFather 提供的 API Token。
- 如果 Token 中不包含逗号（即只有一个 Token），则按传统单机器人模式处理。
- 如需为每个机器人设置不同的 `placeholder`、`group_trigger` 等配置，请使用配置文件中的 `telegram_bots` 列表。
