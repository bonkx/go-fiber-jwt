package helpers

import (
	"encoding/base64"
	"strings"
)

func Encode(s string) string {
	data := base64.StdEncoding.EncodeToString([]byte(s))
	return string(data)
}

func Decode(s string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func FormatPhoneNumber(phone_number string) string {
	phone := phone_number
	// format value
	// replace 0 with +62
	// phone = strings.Replace(phone, "0", "+62", -1)
	// - remove white space
	phone = strings.ReplaceAll(phone, " ", "")
	// - remove (
	phone = strings.ReplaceAll(phone, "(", "")
	// - remove )
	phone = strings.ReplaceAll(phone, ")", "")
	// - remove -
	phone = strings.ReplaceAll(phone, "-", "")
	return phone
}
