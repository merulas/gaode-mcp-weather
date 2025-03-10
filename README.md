# MCP 天气服务

[English](README_EN.md) | 中文

## 项目简介

MCP 天气服务是一个基于高德地图 API 的天气预报服务，提供实时天气状况和天气预报数据。该服务作为 MCP (Model Control Protocol) 生态系统的一部分，可以轻松集成到支持 MCP 的应用程序中。

## 功能特点

- **实时天气数据**：获取指定位置的当前天气状况
- **天气预报**：提供未来天气预报数据
- **中国城市覆盖**：支持中国所有城市和区县的天气查询
- **简单集成**：作为 MCP 服务，易于与其他应用集成

## 什么是 MCP？

MCP（Model Control Protocol）是一个框架，用于构建可以被 AI 模型使用的工具。MCP 服务可以提供各种功能，如天气预报、搜索、计算等，供 AI 模型调用。

## 安装与配置

### 前提条件

- Go 1.18 或更高版本
- 高德地图 API 密钥

### 安装步骤

1. 克隆仓库：
```bash
git clone https://github.com/tung/mcp.git
cd mcp
```

2. 创建 `.env` 文件并添加你的高德地图 API 密钥：
```
AMAP_API_KEY=your_api_key_here
```

> 注：您可以通过在[高德开放平台](https://lbs.amap.com/)注册获取 API 密钥。

## 运行服务

### 直接运行

```bash
go run cmd/server/main.go
```

服务将在 `http://localhost:8080` 上运行。

### 通过 MCP 框架运行

1. 创建 MCP 配置文件（参见 [配置文档](docs/mcp-config.md)）
2. 使用 MCP 框架启动服务：
```bash
mcp --config mcp-config.json
```

## API 使用说明

### 获取天气预报

#### 请求格式

```
POST /weather
Content-Type: application/json

{
    "location": "北京"
}
```

#### 响应格式

```json
{
    "location": "北京市",
    "location_key": "110000",
    "country": "中国",
    "current_conditions": {
        "temperature": {
            "value": 25.6,
            "unit": "C"
        },
        "weather_text": "晴",
        "relative_humidity": 65,
        "precipitation": false,
        "observation_time": "2024-03-10T15:00:00+08:00"
    },
    "hourly_forecast": [
        {
            "relative_time": "+1 hour",
            "temperature": {
                "value": 26.1,
                "unit": "C"
            },
            "weather_text": "晴",
            "precipitation_probability": 10,
            "precipitation_type": "None",
            "precipitation_intensity": "None"
        },
        // 更多小时预报...
    ]
}
```

### 响应字段说明

#### 当前天气条件 (`current_conditions`)
- `temperature`: 当前温度（摄氏度）
- `weather_text`: 天气状况描述
- `relative_humidity`: 相对湿度百分比
- `precipitation`: 是否有降水
- `observation_time`: 观测时间

#### 逐小时预报 (`hourly_forecast`)
- `relative_time`: 相对于当前时间的时间
- `temperature`: 预计温度
- `weather_text`: 天气状况描述
- `precipitation_probability`: 降水概率
- `precipitation_type`: 降水类型
- `precipitation_intensity`: 降水强度

## 文档

- [API 文档](docs/weather.md)
- [MCP 配置文档](docs/mcp-config.md)
- [Claude MCP 协议文档](docs/claude-mcp.md)

## 开发

### 项目结构

```
mcp/
├── cmd/
│   └── server/
│       └── main.go         # 主程序入口
├── internal/
│   ├── bean/
│   │   ├── weather.go      # 数据模型
│   │   ├── mcp.go          # MCP协议数据模型
│   │   └── amap_weather.go # 高德地图API数据模型
│   ├── service/
│   │   ├── weather_service.go     # 原服务层（与AccuWeather API交互）
│   │   └── amap_weather_service.go # 高德地图服务层
│   ├── logic/
│   │   └── weather.go     # 业务逻辑层
│   └── handler/
│       ├── weather_handler.go # HTTP处理层
│       └── mcp_handler.go     # MCP协议处理层
├── docs/
│   ├── weather.md         # API文档
│   ├── mcp-config.md      # MCP配置文档
│   └── claude-mcp.md      # Claude MCP协议文档
├── .env.example           # 环境变量示例
└── README.md             # 项目说明文档
```

## 在 MCP 配置中使用

在 MCP 配置中，你可以这样使用这个服务：

```json
{
    "mcpServers": {
        "weather": {
            "command": "go",
            "args": ["run", "cmd/server/main.go"],
            "env": {
                "AMAP_API_KEY": "your_api_key_here"
            }
        }
    }
}
```

## 高德地图 API

本服务使用[高德地图天气查询API](https://lbs.amap.com/api/webservice/guide/api/weatherinfo)获取天气数据。高德地图API提供了实时天气和天气预报功能，支持全国城市和区县的天气查询。

由于高德地图API不提供每小时的天气预报，本服务根据每日预报数据模拟生成了每小时的预报数据，以保持与原接口的兼容性。

## 限制说明

- API 调用受高德地图免费账户限制
- 天气数据更新频率取决于高德地图API的更新频率

## 贡献指南

欢迎提交 Pull Request 或创建 Issue 来改进此项目。

## 许可证

本项目采用 MIT 许可证 - 详情请参阅 LICENSE 文件。

## 致谢

- [高德地图](https://lbs.amap.com/) 提供天气数据
- [MCP 协议](https://github.com/anthropics/anthropic-cookbook/tree/main/mcp) 提供模型控制协议框架 