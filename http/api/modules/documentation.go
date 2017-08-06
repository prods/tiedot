package modules

// TODO: Make documentation object swagger compatible

type APIRouteParameterDocumentation struct {
	Name string 		`json:"name"`
	Description string  `json:"description"`
	Type string 		`json:"type"`
	Required bool       `json:"required"`
}

type APIRouteResponseDocumentation struct {
	Name string 	`json:"name"`
	Description string `json:"description"`
	Code int 		   `json:"code"`
	CodeDescription string `json:"code_description"`
	ResponseBodySample map[string]interface{} `json:"response_sample"`
	IsError bool 	`json:"is_error"`
}

type APIRouteDocumentation struct {
	Name string 		`json:"name"`
	Description string  `json:"description"`
	Path string 		`json:"path"`
	Method string 		`json:"method"`
	RequestContentType []string `json:"request_content_type"`
	ResultContentType []string  `json:"response_content_type"`
	UrlParameters []APIRouteParameterDocumentation `json:"parameters"`
	RequestBodySample map[string]interface{}  `json:"request_body_sample"`
	Responses []APIRouteResponseDocumentation `json:"Responses"`
}

type APIModuleDocumentation struct {
	Root string  	`json:"root"`
	Name string 	`json:"name"`
	Routes []APIRouteDocumentation `json:"routes"`
}
