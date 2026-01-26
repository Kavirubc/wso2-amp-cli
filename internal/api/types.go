package api

import "time"

type AgentResponse struct {
	UUID           string         `json:"uuid,omitempty"`
	Name           string         `json:"name"`
	DisplayName    string         `json:"displayName,omitempty"`
	Description    string         `json:"description,omitempty"`
	ProjectName    string         `json:"projectName"`
	Status         string         `json:"status,omitempty"`
	CreatedAt      time.Time      `json:"createdAt"`
	Provisioning   *Provisioning  `json:"provisioning,omitempty"`
	AgentType      *AgentTypeInfo `json:"agentType,omitempty"`
	RuntimeConfigs *RuntimeConfig `json:"runtimeConfigs,omitempty"`
	Language       string         `json:"language,omitempty"`
}

// Provisioning contains agent source configuration
type Provisioning struct {
	Type       string            `json:"type"`
	Repository *RepositoryConfig `json:"repository,omitempty"`
}

// RepositoryConfig holds git repository details for agent source
type RepositoryConfig struct {
	URL     string `json:"url"`
	Branch  string `json:"branch"`
	AppPath string `json:"appPath,omitempty"`
}

// AgentTypeInfo describes the agent's type classification
type AgentTypeInfo struct {
	Type    string `json:"type"`
	SubType string `json:"subType,omitempty"`
}

// RuntimeConfig holds runtime environment settings
type RuntimeConfig struct {
	Language        string                `json:"language,omitempty"`
	LanguageVersion string                `json:"languageVersion,omitempty"`
	RunCommand      string                `json:"runCommand,omitempty"`
	Env             []EnvironmentVariable `json:"env,omitempty"`
}

// OrganizationResponse represents an organization
type OrganizationResponse struct {
	Name        string    `json:"name"`
	DisplayName string    `json:"displayName,omitempty"`
	Description string    `json:"description,omitempty"`
	Namespace   string    `json:"namespace,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

// ProjectResponse represents a project
type ProjectResponse struct {
	UUID               string    `json:"uuid,omitempty"`
	Name               string    `json:"name"`
	OrgName            string    `json:"orgName"`
	DisplayName        string    `json:"displayName,omitempty"`
	Description        string    `json:"description,omitempty"`
	DeploymentPipeline string    `json:"deploymentPipeline,omitempty"`
	CreatedAt          time.Time `json:"createdAt"`
}

type ProjectListResponse struct {
	Projects []ProjectResponse `json:"projects"`
	Limit    int               `json:"limit"`
	Offset   int               `json:"offset"`
	Total    int               `json:"total"`
}

// BuildResponse represents a build
type BuildResponse struct {
	Name        string     `json:"name"`
	AgentName   string     `json:"agentName"`
	ProjectName string     `json:"projectName"`
	CommitID    string     `json:"commitId"`
	Status      string     `json:"status"`
	Branch      string     `json:"branch,omitempty"`
	StartedAt   time.Time  `json:"startedAt"`
	EndedAt     *time.Time `json:"endedAt,omitempty"`
}

// BuildDetailsResponse extends BuildResponse with step information and progress
type BuildDetailsResponse struct {
	BuildResponse              // Embed base build fields
	Percent         float64     `json:"percent,omitempty"`
	Steps           []BuildStep `json:"steps,omitempty"`
	DurationSeconds int         `json:"durationSeconds,omitempty"`
}

// BuildStep represents a single step in the build process
type BuildStep struct {
	Type       string `json:"type"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	StartedAt  string `json:"startedAt,omitempty"`
	FinishedAt string `json:"finishedAt,omitempty"`
}

// LogsResponse contains logs with metadata (used for both build and runtime logs)
type LogsResponse struct {
	Logs       []LogEntry `json:"logs"`
	TotalCount int        `json:"totalCount"`
	TookMs     float64    `json:"tookMs"`
}

// LogEntry represents a single log line
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Log       string `json:"log"`
	LogLevel  string `json:"logLevel"`
}

// DeploymentResponse represents a deployment (simple)
type DeploymentResponse struct {
	AgentName      string    `json:"agentName"`
	ProjectName    string    `json:"projectName"`
	Status         string    `json:"status"`
	Environment    string    `json:"environment"`
	LastDeployedAt time.Time `json:"lastDeployedAt"`
}

// DeploymentDetails contains detailed deployment info for an environment
type DeploymentDetails struct {
	ImageID                    string               `json:"imageId"`
	Status                     string               `json:"status"`
	LastDeployed               *time.Time           `json:"lastDeployed,omitempty"`
	Endpoints                  []DeploymentEndpoint `json:"endpoints,omitempty"`
	EnvironmentDisplayName     string               `json:"environmentDisplayName,omitempty"`
	PromotionTargetEnvironment *PromotionTarget     `json:"promotionTargetEnvironment,omitempty"`
}

// DeploymentEndpoint represents an endpoint for a deployed agent
type DeploymentEndpoint struct {
	Name       string `json:"name"`
	URL        string `json:"url"`
	Visibility string `json:"visibility"`
}

// PromotionTarget represents the target environment for promotion
type PromotionTarget struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName,omitempty"`
}

// EndpointResponse represents an endpoint configuration
type EndpointResponse struct {
	URL          string          `json:"url"`
	EndpointName string          `json:"endpointName"`
	Visibility   string          `json:"visibility"`
	Schema       *EndpointSchema `json:"schema,omitempty"`
}

// EndpointSchema holds the OpenAPI schema content
type EndpointSchema struct {
	Content string `json:"content"`
}

// AgentListResponse wraps the paginated agents response
type AgentListResponse struct {
	Agents []AgentResponse `json:"agents"`
	Limit  int             `json:"limit"`
	Offset int             `json:"offset"`
	Total  int             `json:"total"`
}

// OrganizationListResponse wraps the paginated orgs response
type OrganizationListResponse struct {
	Organizations []OrganizationResponse `json:"organizations"`
	Limit         int                    `json:"limit"`
	Offset        int                    `json:"offset"`
	Total         int                    `json:"total"`
}

// DeploymentPipelineResponse represents a deployment pipeline
type DeploymentPipelineResponse struct {
	Name           string          `json:"name"`
	DisplayName    string          `json:"displayName,omitempty"`
	Description    string          `json:"description,omitempty"`
	OrgName        string          `json:"orgName"`
	CreatedAt      time.Time       `json:"createdAt"`
	PromotionPaths []PromotionPath `json:"promotionPaths,omitempty"`
}

// PromotionPath represents a promotion path in a deployment pipeline
type PromotionPath struct {
	SourceEnvironmentRef  string      `json:"sourceEnvironmentRef"`
	TargetEnvironmentRefs []TargetRef `json:"targetEnvironmentRefs"`
}

// TargetRef represents a target environment reference
type TargetRef struct {
	Name string `json:"name"`
}

// DeploymentPipelineListResponse wraps paginated pipelines response
type DeploymentPipelineListResponse struct {
	DeploymentPipelines []DeploymentPipelineResponse `json:"deploymentPipelines"`
	Limit               int                          `json:"limit"`
	Offset              int                          `json:"offset"`
	Total               int                          `json:"total"`
}

// Environment represents an environment in the platform
type Environment struct {
	UUID         string    `json:"uuid,omitempty"`
	Name         string    `json:"name"`
	DisplayName  string    `json:"displayName,omitempty"`
	DataplaneRef string    `json:"dataplaneRef,omitempty"`
	IsProduction bool      `json:"isProduction"`
	DNSPrefix    string    `json:"dnsPrefix,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
}

// EnvironmentListResponse wraps paginated environments response
type EnvironmentListResponse struct {
	Environments []Environment `json:"environments"`
	Limit        int           `json:"limit"`
	Offset       int           `json:"offset"`
	Total        int           `json:"total"`
}

// DataPlane represents a data plane in the platform
type DataPlane struct {
	Name        string    `json:"name"`
	DisplayName string    `json:"displayName,omitempty"`
	Description string    `json:"description,omitempty"`
	OrgName     string    `json:"orgName,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

// DataPlaneListResponse wraps paginated data planes response
type DataPlaneListResponse struct {
	DataPlanes []DataPlane `json:"dataPlanes"`
	Limit      int         `json:"limit"`
	Offset     int         `json:"offset"`
	Total      int         `json:"total"`
}

// --- Request Types ---

// CreateOrganizationRequest for POST /orgs
type CreateOrganizationRequest struct {
	Name string `json:"name"`
}

// CreateProjectRequest for POST /orgs/{org}/projects
type CreateProjectRequest struct {
	Name               string  `json:"name"`
	DisplayName        string  `json:"displayName"`
	Description        *string `json:"description,omitempty"`
	DeploymentPipeline string  `json:"deploymentPipeline"`
}

// DeployAgentRequest for POST /orgs/{org}/projects/{proj}/agents/{agent}/deployments
type DeployAgentRequest struct {
	ImageId string                `json:"imageId"`
	Env     []EnvironmentVariable `json:"env,omitempty"`
}

// EnvironmentVariable represents a key-value pair for deployment config
type EnvironmentVariable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// CreateAgentRequest for POST /orgs/{org}/projects/{proj}/agents
type CreateAgentRequest struct {
	Name           string          `json:"name"`
	DisplayName    string          `json:"displayName"`
	Description    string          `json:"description,omitempty"`
	Provisioning   Provisioning    `json:"provisioning"`
	AgentType      AgentTypeInfo   `json:"agentType"`
	RuntimeConfigs *RuntimeConfig  `json:"runtimeConfigs,omitempty"`
	InputInterface *InputInterface `json:"inputInterface,omitempty"`
}

// InputInterface defines the agent's input interface configuration
type InputInterface struct {
	Type     string        `json:"type"`
	Port     int           `json:"port,omitempty"`
	BasePath string        `json:"basePath,omitempty"`
	Schema   *SchemaConfig `json:"schema,omitempty"`
}

// SchemaConfig holds OpenAPI schema configuration
type SchemaConfig struct {
	Path string `json:"path"`
}

// TokenRequest for POST /orgs/{org}/projects/{proj}/agents/{agent}/token
type TokenRequest struct {
	ExpiresIn string `json:"expiresIn,omitempty"`
}

// TokenResponse represents a generated agent token
type TokenResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expiresAt"`
	IssuedAt  int64  `json:"issuedAt"`
	TokenType string `json:"tokenType"`
}

// RuntimeLogRequest for POST /orgs/{org}/projects/{proj}/agents/{agent}/runtime-logs
type RuntimeLogRequest struct {
	EnvironmentName string   `json:"environmentName"`
	StartTime       string   `json:"startTime,omitempty"`
	EndTime         string   `json:"endTime,omitempty"`
	Limit           int      `json:"limit,omitempty"`
	SortOrder       string   `json:"sortOrder,omitempty"`
	LogLevels       []string `json:"logLevels,omitempty"`
	SearchPhrase    string   `json:"searchPhrase,omitempty"`
}

// --- Trace Types ---

// TraceListOptions for GET traces with query parameters (not serialized to JSON)
type TraceListOptions struct {
	Environment string
	StartTime   string
	EndTime     string
	Limit       int
	Offset      int
	SortOrder   string
}

// Trace represents a trace summary in list responses (TraceOverview in OpenAPI)
type Trace struct {
	TraceID         string       `json:"traceId"`
	RootSpanID      string       `json:"rootSpanId"`
	RootSpanName    string       `json:"rootSpanName"`
	RootSpanKind    string       `json:"rootSpanKind,omitempty"` // llm, tool, agent, etc.
	StartTime       time.Time    `json:"startTime"`
	EndTime         time.Time    `json:"endTime"`
	DurationInNanos int64        `json:"durationInNanos,omitempty"`
	SpanCount       int          `json:"spanCount,omitempty"`
	TokenUsage      *TokenUsage  `json:"tokenUsage,omitempty"`
	Status          *TraceStatus `json:"status,omitempty"`
	Input           interface{}  `json:"input,omitempty"`
	Output          interface{}  `json:"output,omitempty"`
}

// TraceListResponse for list traces endpoint
type TraceListResponse struct {
	Traces     []Trace `json:"traces"`
	TotalCount int     `json:"totalCount"`
}

// TokenUsage for LLM token tracking
type TokenUsage struct {
	InputTokens  int `json:"inputTokens"`
	OutputTokens int `json:"outputTokens"`
	TotalTokens  int `json:"totalTokens"`
}

// TraceStatus indicates error state
type TraceStatus struct {
	ErrorCount int `json:"errorCount"`
}

// Span represents a single span in a trace
type Span struct {
	TraceID         string                 `json:"traceId"`
	SpanID          string                 `json:"spanId"`
	ParentSpanID    string                 `json:"parentSpanId,omitempty"`
	Name            string                 `json:"name"`
	Service         string                 `json:"service"`
	StartTime       time.Time              `json:"startTime"`
	EndTime         time.Time              `json:"endTime"`
	DurationInNanos int64                  `json:"durationInNanos"`
	Kind            string                 `json:"kind"`   // CLIENT, SERVER, PRODUCER, CONSUMER, INTERNAL
	Status          string                 `json:"status"` // OK, ERROR, UNSET
	Attributes      map[string]interface{} `json:"attributes,omitempty"`
}

// TraceDetailsResponse for single trace GET
type TraceDetailsResponse struct {
	Spans      []Span `json:"spans"`
	TotalCount int    `json:"totalCount"`
}

// FullTrace for export (includes all spans inline)
type FullTrace struct {
	TraceID         string       `json:"traceId"`
	RootSpanID      string       `json:"rootSpanId"`
	RootSpanName    string       `json:"rootSpanName"`
	RootSpanKind    string       `json:"rootSpanKind,omitempty"`
	StartTime       time.Time    `json:"startTime"`
	EndTime         time.Time    `json:"endTime"`
	DurationInNanos int64        `json:"durationInNanos,omitempty"`
	SpanCount       int          `json:"spanCount,omitempty"`
	TokenUsage      *TokenUsage  `json:"tokenUsage,omitempty"`
	Status          *TraceStatus `json:"status,omitempty"`
	Input           interface{}  `json:"input,omitempty"`
	Output          interface{}  `json:"output,omitempty"`
	Spans           []Span       `json:"spans"`
}

// TraceExportResponse for export endpoint
type TraceExportResponse struct {
	Traces     []FullTrace `json:"traces"`
	TotalCount int         `json:"totalCount"`
}
