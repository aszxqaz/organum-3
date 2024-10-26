package domain

type Model struct {
	Checksum   string
	Bytes      []byte
	RoomID     string
	UploaderID string
}

type ModelJSON struct {
	Checksum   string  `json:"checksum"`
	RoomID     string  `json:"roomId"`
	SizeMb     float64 `json:"sizeMb"`
	UploaderID string  `json:"uploaderId"`
}

func NewModelJSON(m *Model) *ModelJSON {
	return &ModelJSON{
		Checksum:   m.Checksum,
		RoomID:     m.RoomID,
		SizeMb:     float64(len(m.Bytes)) / 1024.0 / 1024.0,
		UploaderID: m.UploaderID,
	}
}
