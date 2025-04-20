package models

type Gemini struct {
	Candidates []GeminiCandidateResp `json:"candidates"`
}

type GeminiCandidateReq struct {
	Content []GeminiContents `json:"contents"`
}

type GeminiCandidateResp struct {
	Content GeminiContents `json:"content"`
}

type GeminiContents struct {
	Parts []GeminiParts `json:"parts"`
	Role  string        `json:"role"`
}

type GeminiParts struct {
	Text string `json:"text"`
}
