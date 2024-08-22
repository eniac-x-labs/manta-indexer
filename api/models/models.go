package models

type QueryDTParams struct {
	Page     int
	PageSize int
	Order    string
}

type DepositTokensResponse struct {
	Current int   `json:"Current"`
	Size    int   `json:"Size"`
	Total   int64 `json:"Total"`
}
