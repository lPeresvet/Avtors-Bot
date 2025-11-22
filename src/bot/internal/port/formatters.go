package port

import (
	"avtor.ru/bot/client"
	"fmt"
)

func FormatZone(zone *client.ZoneDetails) string {
	return fmt.Sprintf("*Кадастровый номер:* %v\n*Формат собственности:* %v\n*Вид использования:* %v\n*Адрес:* %v\n*Площадь:* %v кв. м.", zone.Id, zone.PropertyType, zone.PermittedUsage, *zone.Address, *zone.Square)
}
