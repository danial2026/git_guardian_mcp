package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

// Tool represents an MCP tool
type Tool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Handler     func(params json.RawMessage) (interface{}, error)
}

// Server represents the MCP server
type Server struct {
	tools  map[string]*Tool
	logger *log.Logger
	reader *bufio.Reader
	writer io.Writer
}

// Request represents an MCP request
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// Response represents an MCP response
type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}

// Error represents an MCP error
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewServer creates a new MCP server
func NewServer(logger *log.Logger) *Server {
	return &Server{
		tools:  make(map[string]*Tool),
		logger: logger,
		reader: bufio.NewReader(os.Stdin),
		writer: os.Stdout,
	}
}

// RegisterTool registers a new tool
func (s *Server) RegisterTool(name, description string, handler func(params json.RawMessage) (interface{}, error)) {
	s.tools[name] = &Tool{
		Name:        name,
		Description: description,
		Handler:     handler,
	}
}

// Start starts the MCP server
func (s *Server) Start() error {
	s.logger.Println("MCP server starting...")

	for {
		line, err := s.reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("read error: %w", err)
		}

		var req Request
		if err := json.Unmarshal(line, &req); err != nil {
			s.sendError(nil, -32700, fmt.Sprintf("Parse error: %v", err))
			continue
		}

		s.handleRequest(&req)
	}
}

func (s *Server) handleRequest(req *Request) {
	// Only log non-routine methods to reduce noise
	if req.Method != "resources/list" && req.Method != "prompts/list" && req.Method != "ping" {
		s.logger.Printf("Handling method: %s", req.Method)
	}

	switch req.Method {
	case "initialize":
		s.handleInitialize(req)
	case "initialized", "notifications/initialized":
		// Notification - no response needed
		return
	case "tools/list":
		s.handleToolsList(req)
	case "tools/call":
		s.handleToolCall(req)
	case "resources/list":
		s.handleResourcesList(req)
	case "prompts/list":
		s.handlePromptsList(req)
	case "ping":
		s.sendResponse(req.ID, map[string]interface{}{})
	default:
		// Don't send error responses for notifications (methods without IDs)
		if req.ID == nil {
			s.logger.Printf("Ignoring notification: %s", req.Method)
			return
		}
		s.logger.Printf("Unknown method: %s", req.Method)
		s.sendError(req.ID, -32601, fmt.Sprintf("Method not found: %s", req.Method))
	}
}

func (s *Server) handleInitialize(req *Request) {
	s.logger.Println("MCP server initialized")
	result := map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"serverInfo": map[string]interface{}{
			"name":    "git-guardian-mcp",
			"version": "1.0.0",
		},
		"capabilities": map[string]interface{}{
			"tools":     map[string]interface{}{},
			"resources": map[string]interface{}{},
			"prompts":   map[string]interface{}{},
		},
	}
	s.sendResponse(req.ID, result)
}

func (s *Server) handleToolsList(req *Request) {
	tools := make([]map[string]interface{}, 0, len(s.tools))
	for _, tool := range s.tools {
		tools = append(tools, map[string]interface{}{
			"name":        tool.Name,
			"description": tool.Description,
			"inputSchema": map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		})
	}
	s.sendResponse(req.ID, map[string]interface{}{
		"tools": tools,
	})
}

func (s *Server) handleResourcesList(req *Request) {
	// No resources for now
	s.sendResponse(req.ID, map[string]interface{}{
		"resources": []interface{}{},
	})
}

func (s *Server) handlePromptsList(req *Request) {
	// No prompts for now
	s.sendResponse(req.ID, map[string]interface{}{
		"prompts": []interface{}{},
	})
}

func (s *Server) handleToolCall(req *Request) {
	var params struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		s.sendError(req.ID, -32602, fmt.Sprintf("Invalid params: %v", err))
		return
	}

	s.logger.Printf("Executing tool: %s", params.Name)

	tool, exists := s.tools[params.Name]
	if !exists {
		s.sendError(req.ID, -32602, fmt.Sprintf("Tool not found: %s", params.Name))
		return
	}

	result, err := tool.Handler(params.Arguments)
	if err != nil {
		s.logger.Printf("ERROR: Tool %s failed: %v", params.Name, err)
		s.sendError(req.ID, -32603, fmt.Sprintf("Tool execution error: %v", err))
		return
	}

	s.sendResponse(req.ID, map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": toJSON(result),
			},
		},
	})
}

func (s *Server) sendResponse(id interface{}, result interface{}) {
	resp := Response{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
	s.writeResponse(resp)
}

func (s *Server) sendError(id interface{}, code int, message string) {
	resp := Response{
		JSONRPC: "2.0",
		ID:      id,
		Error: &Error{
			Code:    code,
			Message: message,
		},
	}
	s.writeResponse(resp)
}

func (s *Server) writeResponse(resp Response) {
	data, err := json.Marshal(resp)
	if err != nil {
		s.logger.Printf("Failed to marshal response: %v", err)
		return
	}
	data = append(data, '\n')
	if _, err := s.writer.Write(data); err != nil {
		s.logger.Printf("Failed to write response: %v", err)
	}
}

func toJSON(v interface{}) string {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("%v", v)
	}
	return string(data)
}
