package health

import (
	paseto "aidanwoods.dev/go-paseto"
)

type Service interface{}

type service struct {
	authKey paseto.V4SymmetricKey
}

func newService(authKey paseto.V4SymmetricKey) Service {
	return &service{authKey: authKey}
}
