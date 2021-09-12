//+build auth

package database

type DriverAuth interface {
	DriverGlobal
	CheckLoginUser(login, senha string) (token string)
	RegisterUser(login, senha string) (token string, error uint8)
}
