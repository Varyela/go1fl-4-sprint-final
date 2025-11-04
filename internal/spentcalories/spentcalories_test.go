package spentcalories

import (
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
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("неверный формат данных: ожидается 3 части, получено %d", len(parts))
	}

	stepsStr := strings.TrimSpace(parts[0])
	activityType := strings.TrimSpace(parts[1])
	durationStr := strings.TrimSpace(parts[2])

	// Проверяем количество шагов
	if stepsStr == "" {
		return 0, "", 0, fmt.Errorf("количество шагов не может быть пустым")
	}

	// Удаляем возможные пробелы и проверяем на наличие нечисловых символов
	stepsStr = strings.TrimSpace(stepsStr)
	if strings.ContainsAny(stepsStr, " \t\n") {
		return 0, "", 0, fmt.Errorf("количество шагов содержит пробелы")
	}

	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка преобразования шагов: %v", err)
	}

	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("количество шагов должно быть положительным")
	}

	// Проверяем тип активности
	if activityType == "" {
		return 0, "", 0, fmt.Errorf("тип активности не может быть пустым")
	}

	// Проверяем продолжительность
	if durationStr == "" {
		return 0, "", 0, fmt.Errorf("продолжительность не может быть пустой")
	}

	// Удаляем возможные пробелы
	durationStr = strings.TrimSpace(durationStr)
	if strings.ContainsAny(durationStr, " \t\n") {
		return 0, "", 0, fmt.Errorf("продолжительность содержит пробелы")
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка преобразования продолжительности: %v", err)
	}

	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("продолжительность должна быть положительной")
	}

	return steps, activityType, duration, nil
}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	distanceMeters := float64(steps) * stepLength
	return distanceMeters / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	
	dist := distance(steps, height)
	durationHours := duration.Hours()
	
	if durationHours == 0 {
		return 0
	}
	
	return dist / durationHours
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверка входных параметров
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть положительным")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть положительным")
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть положительным")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность должна быть положительной")
	}

	speed := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()
	
	calories := (weight * speed * durationMinutes) / minInH
	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверка входных параметров
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть положительным")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть положительным")
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть положительным")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность должна быть положительной")
	}

	speed := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()
	
	calories := (weight * speed * durationMinutes) / minInH
	calories *= walkingCaloriesCoefficient
	
	return calories, nil
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var calories float64
	var caloriesErr error

	switch activityType {
	case "Ходьба":
		calories, caloriesErr = WalkingSpentCalories(steps, weight, height, duration)
	case "Бег":
		calories, caloriesErr = RunningSpentCalories(steps, weight, height, duration)
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}

	if caloriesErr != nil {
		return "", caloriesErr
	}

	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)
	durationHours := duration.Hours()

	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", 
		activityType, durationHours, dist, speed, calories), nil
}