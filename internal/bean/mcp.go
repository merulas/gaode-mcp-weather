package bean

// MCPRequest Claude MCP请求结构
type MCPRequest struct {
	Name       string                 `json:"name"`
	Parameters map[string]interface{} `json:"parameters"`
}

// MCPResponse Claude MCP响应结构
type MCPResponse struct {
	Content interface{} `json:"content"`
	Type    string      `json:"type"`
}

// MCPErrorResponse Claude MCP错误响应结构
type MCPErrorResponse struct {
	Error string `json:"error"`
	Type  string `json:"type"`
}

// WeatherMCPRequest 天气MCP请求参数
type WeatherMCPRequest struct {
	Location string `json:"location"`
}

// NewMCPResponse 创建新的MCP响应
func NewMCPResponse(content interface{}) MCPResponse {
	return MCPResponse{
		Content: content,
		Type:    "application/json",
	}
}

// NewMCPErrorResponse 创建新的MCP错误响应
func NewMCPErrorResponse(err string) MCPErrorResponse {
	return MCPErrorResponse{
		Error: err,
		Type:  "error",
	}
}
