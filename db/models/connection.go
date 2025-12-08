package models

import "time"

type ConnectionModel struct {
	Provider             string    `json:"provider"`
	ProviderUserId       string    `json:"providerUserId"`
	ProviderAccountEmail string    `json:"providerAccountEmail"`
	LinkedAt             time.Time `json:"linkedAt"`
}
