package port

import (
	"avtor.ru/bot/client"
	"fmt"
)

func FormatZone(zone *client.ZoneDetails) string {
	return fmt.Sprintf("Кадастровый номер: %v\nФормат собственности: %v\nВид использования: %v", zone.Id, zone.PropertyType, zone.PermittedUsage)
}
