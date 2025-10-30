package handlers

import (
	"avtor.ru/bot/server"
	"strings"
)

var (
	converter = map[string]server.PropertyType{
		"частная":       server.Private,
		"муниципальная": server.Municipal,
	}
)

func ConvertOwnershipType(in string) server.PropertyType {
	res, ok := converter[strings.ToLower(in)]
	if !ok {
		return server.Undefined
	}

	return res
}
