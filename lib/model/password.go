package model

type Password struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	Overwrite bool
}
