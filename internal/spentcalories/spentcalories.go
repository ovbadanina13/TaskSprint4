package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	slice := strings.Split(data, ",")
	if len(slice) != 3 {
		return 0, "", 0, fmt.Errorf("wait for 3 params, got %d", len(slice))
	}

	steps, err := strconv.Atoi(slice[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("error convert to int: %v", err)
	}
	if steps <= 0 {
		return 0, "", 0, errors.New("amount of steps mustn't be 0")
	}

	timeTrain, err := time.ParseDuration(slice[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("error to convert time: %v", err)
	}
	if timeTrain <= 0 {
		return 0, "", 0, errors.New("duration must be positive")
	}

	typeActivity := slice[1]

	return steps, typeActivity, timeTrain, nil
}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	distance := (stepLength * float64(steps)) / mInKm
	return distance
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	distance := distance(steps, height)
	meanSpeed := distance / duration.Hours()
	return meanSpeed

}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, typeActivity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	finalDistance := distance(steps, height)
	finalSpeed := meanSpeed(steps, height, duration)
	var calories float64
	switch typeActivity {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	default:
		return "", errors.New("неизвестный тип тренировки")
	}
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", typeActivity, duration.Hours(), finalDistance, finalSpeed, calories), nil

}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("amount of steps mustn't be 0")
	}
	if weight <= 30 || weight >= 300 {
		return 0, fmt.Errorf("weight must be between 30 and 300 kg, got %.2f", weight)
	}
	if height <= 0.5 || height >= 3.0 {
		return 0, fmt.Errorf("height must be between 0.5 and 3.0 m, got %.2f", height)
	}
	if duration <= 0 {
		return 0, errors.New("duration mustn't be 0")
	}

	mSpeed := meanSpeed(steps, height, duration)

	kkal := (weight * mSpeed * duration.Minutes()) / minInH
	return kkal, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("amount of steps mustn't be 0")
	}
	if weight <= 30 || weight >= 300 {
		return 0, fmt.Errorf("weight must be between 30 and 300 kg, got %.2f", weight)
	}
	if height <= 0.5 || height >= 3.0 {
		return 0, fmt.Errorf("height must be between 0.5 and 3.0 m, got %.2f", height)
	}
	if duration <= 0 {
		return 0, errors.New("duration mustn't be 0")
	}

	mSpeed := meanSpeed(steps, height, duration)
	kkal := (weight * mSpeed * duration.Minutes()) / minInH
	return kkal * walkingCaloriesCoefficient, nil
}
