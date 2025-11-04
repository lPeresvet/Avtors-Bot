package zones

import "strings"

func ValidateZone(zoneID string) bool {
	zoneParts := strings.Split(zoneID, ":")
	if len(zoneParts) != 4 {
		return false
	}

	if len(zoneParts[0]) != 2 && len(zoneParts[1]) != 2 && len(zoneParts[2]) != 6 {
		return false
	}

	return true
}
