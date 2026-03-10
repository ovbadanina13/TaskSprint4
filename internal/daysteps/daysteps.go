package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	slice := strings.Split(data, ",")
	if len(slice) != 2 {
		return 0, 0, fmt.Errorf("wait for 2 params, got %d", len(slice))
	}
	steps, err := strconv.Atoi(slice[0])
	if err != nil {
		return 0, 0, fmt.Errorf("error convert to int: %v", err)
	}
	if steps <= 0 {
		return 0, 0, errors.New("amount of steps mustn't be 0")
	}

	timeTrain, err := time.ParseDuration(slice[1])
	if err != nil {
		return 0, 0, fmt.Errorf("error to convert time: %v", err)
	}
	if timeTrain <= 0 {
		return 0, 0, errors.New("duration must be positive")
	}

	return steps, timeTrain, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, timeTrain, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}
	if steps <= 0 {
		return ""
	}

	distance := stepLength * float64(steps)
	countKm := distance / mInKm

	kkal, err := spentcalories.WalkingSpentCalories(steps, weight, height, timeTrain)

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, countKm, kkal)

}
