package validate

import (
	"OpenGreenEye/utils"
	"github.com/go-playground/validator/v10"
	"regexp"
)

var ipRegex *regexp.Regexp

func init() {
	var err error
	ipRegex, err = regexp.Compile("[0-9]{3}\\.[0-9]{3}\\.[0-9]{3}\\.[0-9]{3}")

	utils.Pie(err)
}

func ConfigValidate() *validator.Validate {
	v := validator.New()

	utils.Pie(v.RegisterValidation("ip", func(fl validator.FieldLevel) bool {
		ip := fl.Field().Bytes()

		if len(ip) != 15 {
			return false
		}

		return ipRegex.Match(ip)
	}))

	utils.Pie(v.RegisterValidation("signType", func(fl validator.FieldLevel) bool {
		signType := fl.Field().String()

		if len(signType) != 1 {
			return false
		}

		return signType == "A" || signType == "D"
	}))

	return v
}
