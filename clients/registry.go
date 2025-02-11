package clients

import (
	"room-service/clients/config"
	clients "room-service/clients/user"
	config2 "room-service/config"
)

type ClientRegistry struct{}

type IClientRegistry interface {
	GetUser() clients.IUserClient
}

func NewClientRegistry() IClientRegistry {
	return &ClientRegistry{}
}

func (r *ClientRegistry) GetUser() clients.IUserClient {
	return clients.NewUserClient(
		config.NewClientConfig(
			config.WithBaseURL(config2.Config.InternalService.User.Host),
			config.WithSignatureKey(config2.Config.SignatureKey),
		),
	)
}
