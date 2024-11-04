package domain

import (
	"errors"
)

var (
	ErrObjectNotFound = errors.New("object not found")
)

type Scene struct {
	Transform Transform `json:"transform"`
	Checksum  string    `json:"checksum"`
	Objects   []*Object `json:"objects"`
	IsJoined  bool      `json:"isJoined"`
}

func (s *Scene) SetTransform(t Transform) {
	s.Transform = t
}

func (s *Scene) UpdateObject(name string, fn func(o *Object) (*Object, error)) error {
	for i, o := range s.Objects {
		if o.Name == name {
			updated, err := fn(o)
			if err != nil {
				return err
			}
			s.Objects[i] = updated
			return nil
		}
	}
	return ErrObjectNotFound
}
