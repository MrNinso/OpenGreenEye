//+build server sqlite
package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

func NewClient() DriverGlobal {
	path := os.Getenv("SQLITE_PATH")

	if path == "" {
		path = "./database.db"
	}

	db, err := sql.Open("sqlite3", path)

	if err != nil {
		panic(err)
	}

	if _, err = db.Exec("CREATE TABLE IF NOT EXISTS DB(VERSION MEDIUMINT)"); err != nil {
		panic(err)
	}

	s := SqliteDriver{db}

	s.updateVersion()

	return s
}

type SqliteDriver struct {
	*sql.DB
}

func (s SqliteDriver) updateVersion() {
	version := -1

	r, err := s.Query("SELECT VERSION FROM DB LIMIT 1")

	if err != nil {
		panic(err)
	}

	if r.Next() {
		if err = r.Scan(&version); err != nil {
			panic(err)
		}

		_ = r.Close()
	}

	if version < 0 {
		if _, err = s.Exec("CREATE TABLE IF NOT EXISTS USER(ID INTEGER PRIMARY KEY, LOGIN TEXT, PASSWORD TEXT, TOKEN CHAR(36))"); err != nil {
			panic(err)
		}

		if _, err = s.Exec("CREATE TABLE IF NOT EXISTS SENSORS(ID INTEGER PRIMARY KEY, NAME TEXT UNIQUE, TYPE CHAR(1))"); err != nil {
			panic(err)
		}

		if _, err = s.Exec("CREATE TABLE IF NOT EXISTS SENSORS_TAGS(SENSOR_ID INTEGER, TAG TEXT, UNIQUE (SENSOR_ID, TAG))"); err != nil {
			panic(err)
		}

		if _, err = s.Exec("CREATE TABLE IF NOT EXISTS SENSORS_IPS(SENSOR_ID INTEGER UNIQUE, TOKEN CHAR(36), IP CHAR(15))"); err != nil {
			panic(err)
		}

		if _, err = s.Exec("CREATE TABLE IF NOT EXISTS SENSORS_VALUES(SENSOR_ID INTEGER, VALUE REAL, TIME DATE)"); err != nil {
			panic(err)
		}

		if _, err = s.Exec("CREATE TABLE IF NOT EXISTS ACTUATORS(ID INTEGER PRIMARY KEY, NAME TEXT UNIQUE, TYPE CHAR(1))"); err != nil {
			panic(err)
		}

		if _, err = s.Exec("CREATE TABLE IF NOT EXISTS ACTUATOR_LOG(ACTUATOR_ID INTEGER, TAG TEXT, LIMIT_TYPE CHAR(1), TIME DATE)"); err != nil {
			panic(err)
		}

		if _, err = s.Exec("CREATE TABLE IF NOT EXISTS ACTUATOR_TAG_CONFIG(ACTUATOR_ID INTEGER, TAG TEXT, LOW_LIMIT REAL, HIGH_LIMIT REAL, LOW_AJUST REAL, HIGH_AJUST REAL, UNIQUE(ACTUATOR_ID, TAG) )"); err != nil {
			panic(err)
		}

		if _, err = s.Exec("CREATE TABLE IF NOT EXISTS ACTUATOR_IPS(ACTUATOR_ID INTEGER, TOKEN CHAR(36), IP CHAR(15))"); err != nil {
			panic(err)
		}

		if _, err = s.Exec("INSERT INTO DB VALUES (?)", 0); err != nil {
			panic(err)
		}

		version = 0
	}
}

func (s SqliteDriver) CheckSensorKey(key, ip string) (id int) {
	r, err := s.Query("SELECT SENSOR_ID FROM SENSORS_IPS WHERE TOKEN = ? AND IP = ? LIMIT 1", key, ip)

	defer func() {
		_ = r.Close()
	}()

	if err != nil {
		//TODO Log error
		return -1
	}

	if r.Next() {
		id = -1
		if err = r.Scan(&id); err != nil {
			//TODO Log error
			return -1
		}

		return id
	}

	return -1
}

func (s SqliteDriver) CheckUserKey(key string) bool {
	r, err := s.Query("SELECT ID FROM USER WHERE TOKEN = ? LIMIT 1", key)

	defer func() {
		_ = r.Close()
	}()

	if err != nil {
		//TODO Log error
		return false
	}

	if r.Next() {
		id := -1
		if err = r.Scan(&id); err != nil {
			//TODO Log error
			return false
		}

		return id >= 0
	}

	return false
}
