# MCP Weather Service

[中文](README.md) | English

## Introduction

MCP Weather Service is a weather forecast service based on the Gaode (AutoNavi) Map API, providing real-time weather conditions and weather forecast data. This service is part of the MCP (Model Control Protocol) ecosystem and can be easily integrated into applications that support MCP.

## Features

- **Real-time Weather Data**: Get current weather conditions for specified locations
- **Weather Forecasts**: Provide future weather forecast data
- **China City Coverage**: Support for weather queries for all cities and districts in China
- **Simple Integration**: Easy to integrate with other applications as an MCP service

## What is MCP?

MCP (Model Control Protocol) is a framework for building tools that can be used by AI models. MCP services can provide various functionalities, such as weather forecasts, search, calculations, etc., for AI models to call.

## Installation and Configuration

### Prerequisites

- Go 1.18 or higher
- Gaode Map API key

### Installation Steps

1. Clone the repository:
```bash
git clone https://github.com/tung/mcp.git
cd mcp
```

2. Create a `.env` file with your Gaode Map API key:
```
AMAP_API_KEY=your_api_key_here
```

> Note: You can obtain an API key by registering at the [Gaode Open Platform](https://lbs.amap.com/).

## Running the Service

### Direct Run

```bash
go run cmd/server/main.go
```

The service will run on `http://localhost:8080`.

### Run via MCP Framework

1. Create an MCP configuration file (see [Configuration Documentation](docs/mcp-config.md))
2. Start the service using the MCP framework:
```bash
mcp --config mcp-config.json
```

## API Usage

### Get Weather Forecast

#### Request Format

```
POST /weather
Content-Type: application/json

{
    "location": "Beijing"
}
```

#### Response Format

```json
{
    "location": "Beijing",
    "location_key": "110000",
    "country": "China",
    "current_conditions": {
        "temperature": {
            "value": 25.6,
            "unit": "C"
        },
        "weather_text": "Sunny",
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
            "weather_text": "Sunny",
            "precipitation_probability": 10,
            "precipitation_type": "None",
            "precipitation_intensity": "None"
        },
        // More hourly forecasts...
    ]
}
```

### Response Field Descriptions

#### Current Conditions (`current_conditions`)
- `temperature`: Current temperature (in Celsius)
- `weather_text`: Weather condition description
- `relative_humidity`: Relative humidity percentage
- `precipitation`: Whether precipitation is occurring
- `observation_time`: Observation time

#### Hourly Forecast (`hourly_forecast`)
- `relative_time`: Time relative to current time
- `temperature`: Forecasted temperature
- `weather_text`: Weather condition description
- `precipitation_probability`: Probability of precipitation
- `precipitation_type`: Type of precipitation
- `precipitation_intensity`: Intensity of precipitation

## Documentation

- [API Documentation](docs/weather.md)
- [MCP Configuration Documentation](docs/mcp-config.md)
- [Claude MCP Protocol Documentation](docs/claude-mcp.md)

## Development

### Project Structure

```
mcp/
├── cmd/
│   └── server/
│       └── main.go         # Main program entry
├── internal/
│   ├── bean/
│   │   ├── weather.go      # Data models
│   │   ├── mcp.go          # MCP protocol data models
│   │   └── amap_weather.go # Gaode Map API data models
│   ├── service/
│   │   ├── weather_service.go     # Original service layer (interacting with AccuWeather API)
│   │   └── amap_weather_service.go # Gaode Map service layer
│   ├── logic/
│   │   └── weather.go     # Business logic layer
│   └── handler/
│       ├── weather_handler.go # HTTP handler layer
│       └── mcp_handler.go     # MCP protocol handler layer
├── docs/
│   ├── weather.md         # API documentation
│   ├── mcp-config.md      # MCP configuration documentation
│   └── claude-mcp.md      # Claude MCP protocol documentation
├── .env.example           # Environment variable example
└── README.md             # Project documentation
```

## Using in MCP Configuration

In MCP configuration, you can use this service as follows:

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

## Gaode Map API

This service uses the [Gaode Map Weather Query API](https://lbs.amap.com/api/webservice/guide/api/weatherinfo) to obtain weather data. The Gaode Map API provides real-time weather and weather forecast functionality, supporting weather queries for all cities and districts in China.

Since the Gaode Map API does not provide hourly weather forecasts, this service simulates hourly forecast data based on daily forecast data to maintain compatibility with the original interface.

## Limitations

- API calls are subject to Gaode Map free account limitations
- Weather data update frequency depends on the Gaode Map API update frequency

## Contributing

Contributions via Pull Requests or Issues are welcome to improve this project.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgements

- [Gaode Map](https://lbs.amap.com/) for providing weather data
- [MCP Protocol](https://github.com/anthropics/anthropic-cookbook/tree/main/mcp) for the model control protocol framework 