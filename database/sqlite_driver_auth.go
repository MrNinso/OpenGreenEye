//+build auth sqlite

package database

import (
	"OpenGreenEye/enum"
	"OpenGreenEye/utils"
)

func (s SqliteDriver) CheckLoginUser(login, senha string) (token string) {
	r, err := s.Query("SELECT PASSWORD FROM USER WHERE LOGIN = ? LIMIT 1", login)

	defer func() {
		_ = r.Close()
	}()

	if err != nil {
		//TODO Log error
		return ""
	}

	if r.Next() {
		hash := ""

		if err = r.Scan(&hash); err != nil {
			//TODO Log error
			return ""
		}

		_ = r.Close()

		if valid, newHash := utils.Validate(hash, senha); valid {
			token = utils.NewToken()

			if newHash != "" {
				if _, err = s.Exec("UPDATE USER SET PASSWORD = ?, TOKEN = ? WHERE LOGIN = ? LIMIT 1", newHash, token, login); err != nil {
					//TODO Log Error
					return ""
				}
			} else {
				if _, err = s.Exec("UPDATE USER SET TOKEN = ? WHERE LOGIN = ?", token, login); err != nil {
					//TODO Log Error
					return ""
				}
			}

			return token
		}
	}

	return ""
}

func (s SqliteDriver) RegisterUser(login, senha string) (token string, error uint8) {
	r, err := s.Query("SELECT ID FROM USER WHERE LOGIN = ? LIMIT 1", login)

	defer func() {
		_ = r.Close()
	}()

	if err != nil {
		//TODO Log error
		return "", enum.ErrorDatabase
	}

	if r.Next() {
		return "", enum.ErrorUserAlreadyExists
	}

	_ = r.Close()

	hash, err := utils.Encrypt(senha)

	if err != nil {
		//TODO Log error
		return "", enum.ErrorInternal
	}

	token = utils.NewToken()

	if _, err = s.Exec("INSERT INTO USER (LOGIN, PASSWORD, TOKEN) VALUES (?, ?, ?)", login, hash, token); err != nil {
		//TODO Log error
		return "", enum.ErrorDatabase
	}

	return token, 0
}
