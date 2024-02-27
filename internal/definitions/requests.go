package definitions

type IndexRequest struct {
	Limit int `json:"limit"`
	Skip  int `json:"skip"`
}

type KeyRequest struct {
	Key string `json:"id"`
}
