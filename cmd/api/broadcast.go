package main

import "organum/internal/domain"

type WsModelBroadcast struct {
	Method string            `json:"method"`
	Model  *domain.ModelJSON `json:"model"`
}

func NewWsModelBroadcast(model *domain.ModelJSON) *WsModelBroadcast {
	return &WsModelBroadcast{
		Method: "model",
		Model:  model,
	}
}
