package dto

type CreateAliasRequest struct {
	URL   string `json:"url"`
	Alias string `json:"alias,omitempty"`
}
