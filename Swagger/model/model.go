package model

// ErrorResponse defines the structure for error responses
type ErrorResponse struct {
    Message string `json:"message"`
}

type SendCall struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	PathwayID   string `json:"pathway_id" binding:"required"`
}

type CallResponse struct {
	Status string `json:"status"`
	CallID string `json:"call_id"`
}

// Struct to represent the request body
type AnalyzeCallRequest struct {
	Goal      string     `json:"goal" binding:"required"`
	Questions [][]string `json:"questions" binding:"required"`
}

// Struct to represent the response from the external API
type AnalyzeCallResponse struct {
	Status      string   `json:"status"`
	Message     string   `json:"message"`
	Answers     []string `json:"answers"`
	CreditsUsed float64  `json:"credits_used"`
}


// CallDetail represents the structure of the call details response
type CallDetail struct {
	CallID               string             `json:"call_id" example:"f0300301-b066-47a0-83ce-895cb1b63a9a"`
	CallLength           float64            `json:"call_length" example:"120.5"`
	BatchID              *string            `json:"batch_id,omitempty" example:"batch123"`
	To                   string             `json:"to" example:"+14155552671"`
	From                 string             `json:"from" example:"+14155552672"`
	RequestData          RequestData        `json:"request_data"`
	Completed            bool               `json:"completed" example:"true"`
	CreatedAt            string             `json:"created_at" example:"2024-09-26T12:34:56Z"`
	Inbound              bool               `json:"inbound" example:"false"`
	QueueStatus          string             `json:"queue_status" example:"completed"`
	EndpointURL          string             `json:"endpoint_url" example:"https://api.bland.ai/callback"`
	MaxDuration          float64            `json:"max_duration" example:"3600"`
	ErrorMessage         *string            `json:"error_message,omitempty" example:"Call failed due to timeout"`
	Variables            map[string]string  `json:"variables" example:"{\"user_id\":\"12345\", \"campaign\":\"welcome\"}"`
	AnsweredBy           string             `json:"answered_by" example:"human"`
	Record               bool               `json:"record" example:"true"`
	RecordingURL         *string            `json:"recording_url,omitempty" example:"https://api.bland.ai/recordings/12345"`
	Metadata             map[string]string  `json:"metadata" example:"{\"source\":\"ads\", \"region\":\"US\"}"`
	Summary              string             `json:"summary" example:"Call summary notes..."`
	Price                float64            `json:"price" example:"0.05"`
	StartedAt            string             `json:"started_at" example:"2024-09-26T12:34:56Z"`
	LocalDialing         bool               `json:"local_dialing" example:"false"`
	CallEndedBy          string             `json:"call_ended_by" example:"customer"`
	PathwayLogs          *string            `json:"pathway_logs,omitempty" example:"Log details here..."`
	AnalysisSchema       *string            `json:"analysis_schema,omitempty" example:"analysis-schema-123"`
	Analysis             *string            `json:"analysis,omitempty" example:"Detailed analysis of the call..."`
	ConcatenatedTranscript string           `json:"concatenated_transcript" example:"Full transcript of the call..."`
	Transcripts          []Transcript       `json:"transcripts"`
	Status               string             `json:"status" example:"completed"`
	CorrectedDuration    string             `json:"corrected_duration" example:"2m 30s"`
	EndAt                string             `json:"end_at" example:"2024-09-26T12:36:56Z"`
}

// RequestData represents the structure for the request data field in the response
type RequestData struct {
	PhoneNumber string `json:"phone_number"`
	Wait        bool   `json:"wait"`
	Language    string `json:"language"`
}

// Transcript represents a single transcript entry from the call
type Transcript struct {
	ID        int     `json:"id"`
	CreatedAt string  `json:"created_at"`
	Text      string  `json:"text"`
	User      string  `json:"user"`
}

// CreateFolderRequest represents the structure of the request body for creating a folder
type CreateFolderRequest struct {
	Name           string `json:"name" binding:"required"`
	ParentFolderID string `json:"parent_folder_id,omitempty"`
}

// CreateFolderData represents the actual folder data returned inside the "data" field
type CreateFolderData struct {
	FolderID       string  `json:"folder_id"`
	Name           string  `json:"name"`
	ParentFolderID *string `json:"parent_folder_id,omitempty"`
}

// CreateFolderResponse represents the full response including "data" and "errors" fields
type CreateFolderResponse struct {
	Data   CreateFolderData `json:"data"`
	Errors *string          `json:"errors,omitempty"`
}

// CreatePathwayRequest represents the structure of the request body for creating a pathway
type CreatePathwayRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description,omitempty"`
}

// CreatePathwayResponse represents the structure of the response for pathway creation
type CreatePathwayResponse struct {
	Status    string `json:"status"`
	PathwayID string `json:"pathway_id"`
}

// MovePathwayRequest represents the structure of the request body for moving a pathway
type MovePathwayRequest struct {
	PathwayID string `json:"pathway_id" binding:"required"` // Pathway ID to move
	FolderID  string `json:"folder_id,omitempty"`           // Folder ID where the pathway will be moved
}

// MovePathwayData represents the nested data structure in the move pathway response
type MovePathwayData struct {
	PathwayID    string  `json:"pathway_id"`
	OldFolderID  *string `json:"old_folder_id,omitempty"`
	NewFolderID  *string `json:"new_folder_id,omitempty"` // This is the field you want, not FolderID
}

// MovePathwayResponse represents the structure of the response for moving a pathway
type MovePathwayResponse struct {
	Data   MovePathwayData `json:"data"`
	Errors *string         `json:"errors,omitempty"`
}

// CombinedResponse represents the combined response of creating and moving the pathway
type CombinedResponse struct {
	CreatePathwayResponse
	MovePathwayData
}


// CreateChatRequest represents the structure of the request body for creating a chat
type CreateChatRequest struct {
	PathwayID   string `json:"pathway_id" binding:"required"`
	StartNodeID string `json:"start_node_id" binding:"required"`
}

// CreateChatResponseData represents the structure inside the "data" field of the response
type CreateChatResponseData struct {
	ChatID  string `json:"chat_id"`
	Message string `json:"message"`
}

// CreateChatResponse represents the entire response body from the API, including errors
type CreateChatResponse struct {
	Data   CreateChatResponseData `json:"data"`
	Errors *string                `json:"errors,omitempty"`
}
// GetPathwayResponse represents the structure of the response body for getting pathway information
type GetPathwayResponse struct {
	Name                   string        `json:"name"`
	Description            *string       `json:"description,omitempty"`
	Nodes                  []Node        `json:"nodes"`
	Edges                  []Edge        `json:"edges"`
	PublishedAt            *string       `json:"published_at,omitempty"`
	ProductionVersionNumber *string      `json:"production_version_number,omitempty"`
}

// Node represents a node in the pathway
type Node struct {
	ID   string    `json:"id"`
	Data NodeData  `json:"data"`
	Type string    `json:"type"`
}

// NodeData represents the data of a node
type NodeData struct {
	Name           string       `json:"name"`
	Active         bool         `json:"active"`
	Prompt         *string      `json:"prompt,omitempty"`
	GlobalPrompt   *string      `json:"globalPrompt,omitempty"`
	Condition      *string      `json:"condition,omitempty"`
	ModelOptions   ModelOptions `json:"modelOptions"`
	IsStart        bool         `json:"isStart"`
	IsGlobal       bool         `json:"isGlobal,omitempty"`
	GlobalLabel    *string      `json:"globalLabel,omitempty"`
	GlobalDescription *string   `json:"globalDescription,omitempty"`
}

// ModelOptions represents model options inside node data
type ModelOptions struct {
	ModelType        string  `json:"modelType"`
	Temperature      float64 `json:"temperature"`
	SkipUserResponse bool    `json:"skipUserResponse"`
	BlockInterruptions bool  `json:"block_interruptions"`
}

// Edge represents an edge in the pathway
type Edge struct {
	ID          string  `json:"id"`
	Label       *string `json:"label,omitempty"`
	Description *string `json:"description,omitempty"`
	Source      string  `json:"source"`
	Target      string  `json:"target"`
}


// UpdatePathwayRequest represents the request body for updating a pathway
type UpdatePathwayRequest struct {
    Name        string `json:"name" extensions:"x-order=1"`
    Description string `json:"description,omitempty" extensions:"x-order=2"`
    Nodes       []Node `json:"nodes" extensions:"x-order=3"`
    Edges       []Edge `json:"edges" extensions:"x-order=4"`
}

// UpdatePathwayResponse represents the response body after updating a pathway
type UpdatePathwayResponse struct {
    Status      string      `json:"status"`
    Message     string      `json:"message"`
    PathwayData PathwayData `json:"pathway_data"`
}

// PathwayData represents the detailed information of the updated pathway
type PathwayData struct {
    Name        string  `json:"name"`
    Description string  `json:"description,omitempty"`
    Nodes       []Node  `json:"nodes"`
    Edges       []Edge  `json:"edges"`
}


// DeletePathwayResponse represents the structure of the response body after deleting a pathway
type DeletePathwayResponse struct {
    Status    string `json:"status"`    // Status of the operation (success or error)
    Message   string `json:"message"`   // Message providing details (e.g., "Pathway deleted successfully")
	PathwayID  string `json:"pathway_id"`          // The ID of the pathway being used

}


// SendMessageRequest represents the request body for sending a message to the chat
type SendMessageRequest struct {
	Message string `json:"message" binding:"required"`
}

// ChatHistoryEntry represents a single entry in the chat history (user or assistant messages)
type ChatHistoryEntry struct {
	Role    string `json:"role"`    // "user" or "assistant"
	Content string `json:"content"` // The message content
}

// SendMessageResponseData represents the structure inside the "data" field of the response
type SendMessageResponseData struct {
	ChatID            string            `json:"chat_id"`             // ID of the chat session
	AssistantResponse string            `json:"assistant_response"`  // The assistant's latest response
	CurrentNodeID     string            `json:"current_node_id"`     // The current node ID in the conversation pathway
	CurrentNodeName   string            `json:"current_node_name"`   // The current node name
	ChatHistory       []ChatHistoryEntry `json:"chat_history"`       // Full history of the conversation
	PathwayID         string            `json:"pathway_id"`          // The ID of the pathway being used
	Variables         map[string]string `json:"variables"`           // Key-value pairs for any dynamic variables used in the conversation
}

// SendMessageResponse represents the entire response body from the API, including errors
type SendMessageResponse struct {
	Data   SendMessageResponseData `json:"data"`   // The main data of the response
	Errors *string                 `json:"errors"` // Any errors encountered during the request (optional)
}
