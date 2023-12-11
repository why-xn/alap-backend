package types

type RunMode string

type GetListInput struct {
	AccountId string `json:"accountId"`
	Search    string `json:"search"`
	Start     int    `json:"start"`
	Limit     int    `json:"limit"`
}

type SortParam struct {
	SortBy string // bson field name
	Type   int    // -1 or 1
}
