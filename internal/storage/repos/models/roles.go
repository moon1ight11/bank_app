package models

type Role string

const (
	RoleBasic       Role = "Basic"
	RoleVerificator Role = "Verificator"
	RoleAdmin       Role = "Admin"
)
