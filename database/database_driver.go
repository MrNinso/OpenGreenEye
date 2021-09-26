//+build server

package database

type DriverGlobal interface {
	CheckUserKey(key string) bool
	CheckSensorKey(key, ip string) (id int)
}
