package domain

type Transform struct {
	Position Vector3 `json:"position"`
	Rotation Vector3 `json:"rotation"`
	Scale    Vector3 `json:"scale"`
}
