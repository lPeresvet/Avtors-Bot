package handlers

import (
	"avtor.ru/bot/analyse_service/internal/model"
	"avtor.ru/bot/server"
	"strings"
)

var (
	converter = map[string]server.PropertyType{
		"частная":       server.Private,
		"муниципальная": server.Municipal,
	}

	userConverter = map[model.Role]server.Role{
		model.Admin:    server.Admin,
		model.Analyser: server.Analyser,
	}
)

func ConvertOwnershipType(in string) server.PropertyType {
	res, ok := converter[strings.ToLower(in)]
	if !ok {
		return server.Undefined
	}

	return res
}

func ConvertRole(in model.Role) server.Role {
	lowerIn := strings.ToLower(string(in))

	res, ok := userConverter[model.Role(lowerIn)]
	if !ok {
		return server.UndefinedRole
	}

	return res
}
