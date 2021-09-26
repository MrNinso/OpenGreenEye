//+build sensors sqlite

package database

import (
	"OpenGreenEye/enum"
	"OpenGreenEye/types"
	"OpenGreenEye/utils"
)

func (s SqliteDriver) RegisterSensor(sensor types.Sensor) (token string) {
	r, err := s.Query("SELECT ID FROM SENSORS WHERE NAME = ? LIMIT 1", sensor.Name)

	defer func() {
		_ = r.Close()
	}()

	if err != nil {
		//TODO Log error
		return ""
	}

	if r.Next() {
		//TODO Log error
		return ""
	}

	_ = r.Close()

	token = utils.NewToken()

	if _, err = s.Exec("INSERT INTO SENSORS (NAME, TYPE) VALUES (?, ?)", sensor.Name, sensor.SingType); err != nil {
		//TODO Log error
		return ""
	}

	r, err = s.Query("SELECT ID FROM SENSORS WHERE NAME = ? LIMIT 1", sensor.Name)

	if err != nil {
		//TODO Log error
		return ""
	}

	r.Next()
	id := -1

	if r.Scan(&id) != nil || id < 0 {
		return ""
	}

	_, err = s.Exec("DELETE FROM SENSORS_TAGS WHERE SENSOR_ID = ?", id)

	for _, tag := range sensor.Tags {
		if _, err = s.Exec("INSERT INTO SENSORS_TAGS (SENSOR_ID, TAG) VALUES (?, ?)", id, tag); err != nil {
			return ""
		}
	}

	if _, err = s.Exec("INSERT INTO SENSORS_IPS (SENSOR_ID, TOKEN, IP) VALUES (?, ?, ?)", id, token, sensor.IP); err != nil {
		return ""
	}

	return token
}

func (s SqliteDriver) RegisterSensorValue(id int, value float32) (error uint8) {
	if _, err := s.Exec("INSERT INTO SENSORS_VALUES (SENSOR_ID, VALUE) VALUES (?, ?, NOW())", id, value); err != nil {
		return enum.ErrorDatabase
	}

	return 0
}

func (s SqliteDriver) UpdateSensor(id int, sensor types.Sensor) (error uint8) {
	if _, err := s.Exec("UPDATE SENSORS SET NAME = ?, TYPE = ? WHERE ID = ?", sensor.Name, sensor.SingType, id); err != nil {
		return enum.ErrorDatabase
	}

	if _, err := s.Exec("UPDATE SENSORS SET NAME = ?, TYPE = ? WHERE ID = ?", sensor.Name, sensor.SingType, id); err != nil {
		return enum.ErrorDatabase
	}

	_, err := s.Exec("DELETE FROM SENSORS_TAGS WHERE SENSOR_ID = ?", id)

	for _, tag := range sensor.Tags {
		if _, err = s.Exec("INSERT INTO SENSORS_TAGS (SENSOR_ID, TAG) VALUES (?, ?)", id, tag); err != nil {
			return enum.ErrorDatabase
		}
	}

	if _, err = s.Exec("UPDATE SENSORS_IPS SET IP = ? WHERE SENSOR_ID = ?", sensor.IP, id); err != nil {
		return enum.ErrorDatabase
	}

	return 0
}
func (s SqliteDriver) ListSensors() []types.Sensor {
	sensorsRows, err := s.Query("SELECT s.ID, s.NAME, i.IP, s.TYPE  FROM SENSORS s LEFT JOIN SENSORS_IPS i ON s.ID = i.SENSOR_ID")

	defer func() {
		_ = sensorsRows.Close()
	}()

	sensors := make([]types.Sensor, 0)

	if err != nil {
		return sensors
	}

	for sensorsRows.Next() {
		id := -1

		sensor := types.Sensor{
			Tags: make([]string, 0),
		}

		if err = sensorsRows.Scan(&id, &sensor.Name, &sensor.IP, &sensor.SingType); err != nil {
			return make([]types.Sensor, 0)
		}

		tagsRows, err := s.Query("SELECT TAG FROM SENSORS_TAGS WHERE SENSOR_ID = ?", id)

		defer func() {
			_ = tagsRows.Close()
		}()

		if err != nil {
			return make([]types.Sensor, 0)
		}

		for tagsRows.Next() {
			tag := ""

			if err = tagsRows.Scan(&tag); err != nil {
				return make([]types.Sensor, 0)
			}

			sensor.Tags = append(sensor.Tags, tag)
		}

		sensors = append(sensors, sensor)
	}

	return sensors
}
