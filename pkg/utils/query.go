package utils

import (
	"errors"
	"fmt"
	"strings"
)

func ValidateAndReturnSortQuery(sortBy string) (string, error) {
	splits := strings.Split(sortBy, "|")
	if len(splits) != 2 {
		return "", errors.New("malformed sortBy query parameter, should be field.orderdirection i. e. id|asc or id|desc")
	}
	field, order := splits[0], splits[1]
	if strings.ToLower(order) != "desc" && strings.ToLower(order) != "asc" {
		return "", errors.New("malformed orderdirection in sortBy query parameter, should be asc or desc")
	}
	// if !stringInSlice(userFields, field) {
	// 	return "", errors.New("unknown field in sortBy query parameter")
	// }
	return fmt.Sprintf("%s %s", field, strings.ToLower(order)), nil
}

func stringInSlice(strSlice []string, s string) bool {
	for _, v := range strSlice {
		if v == s {
			return true
		}
	}
	return false
}

func ValidateAndReturnFilterMap(filter string) (map[string]string, error) {
	splits := strings.Split(filter, ":")
	if len(splits) != 2 {
		return nil, errors.New("malformed filter query parameter, should be field.value i. e. name.john or email:john@example.com")
	}
	field, value := splits[0], splits[1]

	// if !stringInSlice(userFields, field) {
	// 	return nil, errors.New("unknown field in filter query parameter")
	// }

	return map[string]string{field: value}, nil
}
