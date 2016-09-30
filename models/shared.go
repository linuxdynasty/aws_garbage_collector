package models

type DeleteStatus struct {
	Deleted bool   `json:"deleted"`
	Message string `json:"message"`
	ID      string `json:"id"`
	Name    string `json:"name"`
}

type DeleteResponseApi struct {
	Count   int            `json:"count"`
	Data    []DeleteStatus `json:"data"`
	Message string         `json:"message"`
}

type GetResponseApi struct {
	Count   int         `json:"count"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}
