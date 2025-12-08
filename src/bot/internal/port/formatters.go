package port

import (
	"avtor.ru/bot/client"
	"fmt"
)

func FormatZone(zone *client.ZoneDetails) string {
	zoneVerdict := "Строительство жилья разрешено ✅"
	if !*zone.ZoneOK {
		zoneVerdict = "Строительство жилья запрещено ❌️"
	}

	return fmt.Sprintf("*Кадастровый номер:* %v\n*Формат собственности:* %v\n*Вид использования:* %v\n*Адрес:* %v\n*Площадь:* %v кв. м.\n*Функциональная зона:* %v\n*Территориальная зона:* %v\n\n*%s*", zone.Id, zone.PropertyType, zone.PermittedUsage, *zone.Address, *zone.Square, *zone.FunctionalZone, *zone.TerritorialZone, zoneVerdict)
}
