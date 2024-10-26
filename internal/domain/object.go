package domain

type Object struct {
	Transform Transform `json:"transform"`
	Name      string    `json:"name"`
}
