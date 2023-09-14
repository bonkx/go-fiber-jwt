package helpers

import (
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
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

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("could not hash password %w", err)
	}
	return string(hashedPassword), nil
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
