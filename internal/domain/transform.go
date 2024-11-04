package domain

type Transform struct {
	Position Vector3    `json:"position"`
	Rotation Quaternion `json:"rotation"`
	Scale    Vector3    `json:"scale"`
}
