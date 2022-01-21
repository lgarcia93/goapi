package utils

import "fmt"

func WeekDayFromNumber(weekDay int8) string {

	switch weekDay {
	case 0:
		return "SUNDAY"

	case 1:
		return "MONDAY"

	case 2:
		return "TUESDAY"

	case 3:
		return "WEDNESDAY"

	case 4:
		return "THURSDAY"

	case 5:
		return "FRIDAY"

	case 6:
		return "SATURDAY"
	default:
		return ""
	}
}

func NumberFromWeekDay(weekDay string) (int8, error) {
	switch weekDay {
	case "SUNDAY":
		return 0, nil

	case "MONDAY":
		return 1, nil

	case "TUESDAY":
		return 2, nil

	case "WEDNESDAY":
		return 3, nil

	case "THURSDAY":
		return 4, nil

	case "FRIDAY":
		return 5, nil

	case "SATURDAY":
		return 6, nil
	default:
		return -1, fmt.Errorf("Parâmetro inválido para dia da semana: %g", weekDay)
	}
}
