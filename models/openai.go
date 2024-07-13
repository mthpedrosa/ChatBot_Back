package models

type PostData struct {
	ThreadId          string `json:"threadId"`
	RunId             string `json:"runId"`
	FunctionResponses struct {
		ToolOutputs []struct {
			ToolCallId string `json:"tool_call_id"`
			Output     string `json:"output"`
		} `json:"tool_outputs"`
	} `json:"functionResponses"`
}

type MessageThread struct {
	ID          string      `json:"id"`
	Object      string      `json:"object"`
	CreatedAt   int64       `json:"created_at"`
	ThreadID    string      `json:"thread_id"`
	Role        string      `json:"role"`
	Content     []Content   `json:"content"`
	FileIDs     []string    `json:"file_ids"`
	AssistantID string      `json:"assistant_id"`
	RunID       string      `json:"run_id"`
	Metadata    interface{} `json:"metadata"`
}

type Content struct {
	Type string `json:"type"`
	Text Text   `json:"text"`
}

type Text struct {
	Value       string        `json:"value"`
	Annotations []interface{} `json:"annotations"`
}

type ThreadRun struct {
	ID             string `json:"id"`
	Object         string `json:"object"`
	CreatedAt      int64  `json:"created_at"`
	AssistantID    string `json:"assistant_id"`
	ThreadID       string `json:"thread_id"`
	Status         string `json:"status"`
	StartedAt      int64  `json:"started_at"`
	ExpiresAt      int64  `json:"expires_at"`
	CancelledAt    *int64 `json:"cancelled_at"`
	FailedAt       *int64 `json:"failed_at"`
	CompletedAt    *int64 `json:"completed_at"`
	RequiredAction struct {
		Type              string `json:"type"`
		SubmitToolOutputs struct {
			ToolCalls []struct {
				ID       string `json:"id"`
				Type     string `json:"type"`
				Function struct {
					Name      string `json:"name"`
					Arguments string `json:"arguments"`
				} `json:"function"`
			} `json:"tool_calls"`
		} `json:"submit_tool_outputs"`
	} `json:"required_action"`
	LastError    *interface{} `json:"last_error"`
	Model        string       `json:"model"`
	Instructions string       `json:"instructions"`
	Tools        []struct {
		Type     string `json:"type"`
		Function struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Parameters  struct {
				Type       string `json:"type"`
				Properties struct {
					Origem struct {
						Type        string `json:"type"`
						Description string `json:"description"`
					} `json:"origem"`
					Destino struct {
						Type        string `json:"type"`
						Description string `json:"description"`
					} `json:"destino"`
					Data struct {
						Type        string `json:"type"`
						Description string `json:"description"`
					} `json:"data"`
				} `json:"properties"`
				Required []string `json:"required"`
			} `json:"parameters"`
		} `json:"function"`
	} `json:"tools"`
	FileIDs  []interface{} `json:"file_ids"`
	Metadata struct{}      `json:"metadata"`
}

type ChatCompletion struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
	} `json:"choices"`
}

type ThreadResponse struct {
	ID        string                 `json:"id"`
	Object    string                 `json:"object"`
	CreatedAt int64                  `json:"created_at"`
	Metadata  map[string]interface{} `json:"metadata"`
}

type CallResponse struct {
	ToolCallID string `json:"tool_call_id"`
	OutPut     string `json:"output"`
}

type ThreadIds struct {
	ThreadId  string
	RunId     string
	MessageId string
}

// Tool represents the structure for the tools field in the JSON.
type Tool struct {
	Type string `json:"type"`
}

// ToolResources represents the structure for the tool_resources field in the JSON.
type ToolResources struct {
	VectorStoreIDs []string `json:"vector_store_ids"`
}

// Assistant represents the structure for the entire JSON.
type Assistant struct {
	ID             string                   `json:"id"`
	Object         string                   `json:"object"`
	CreatedAt      int64                    `json:"created_at"`
	Name           string                   `json:"name"`
	Description    *string                  `json:"description"`
	Model          string                   `json:"model"`
	Instructions   string                   `json:"instructions"`
	Tools          []Tool                   `json:"tools"`
	ToolResources  map[string]ToolResources `json:"tool_resources"`
	Metadata       map[string]interface{}   `json:"metadata"`
	TopP           float64                  `json:"top_p"`
	Temperature    float64                  `json:"temperature"`
	ResponseFormat string                   `json:"response_format"`
}

// AssistantDeleted represents the structure for the JSON.
type AssistantDeleted struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}
