package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

// Глобальная переменная для хранения значения rtp
var rtp float64

// Структура для хранения мультипликатора и его вероятности
type payout struct {
	Multiplier  float64
	Probability float64
}

// MultiplierResponse - структура для вывода JSON
type MultiplierResponse struct {
	Result float64 `json:"result"`
}

// Handler - структура, которая хранит таблицу выплат и реализует интерфейс http.Handler
type Handler struct {
	payouts []payout
}

// ServeHTTP реализует интерфейс http.Handler
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Функция для генерации случайного мультипликатора на основе таблицы выплат.
	generateMultiplier := func() float64 {
		p := rand.Float64()
		cumulativeProbability := 0.0
		for _, payout := range h.payouts {
			cumulativeProbability += payout.Probability
			if p < cumulativeProbability {
				return payout.Multiplier
			}
		}
		// Запасной вариант, если что-то пойдет не так
		return 0.0
	}

	// Генерируем мультипликатор
	multiplier := generateMultiplier()

	// Создаем объект ответа
	response := MultiplierResponse{
		Result: multiplier,
	}

	// Кодируем объект ответа в JSON и записываем в writer ответа
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Ошибка при кодировании JSON ответа", http.StatusInternalServerError)
		return
	}

	// Логируем сгенерированный мультипликатор для отладки
	log.Printf("Сгенерированный мультипликатор: %.2f\n", multiplier)
}

// generatePayouts создает таблицу выплат на основе заданного RTP.
func generatePayouts(rtp float64) []payout {
	// Базовые выигрыши и их относительные веса.
	// Сумма этих весов не важна, важны только их пропорции.
	basePayouts := []struct {
		Multiplier float64
		Weight     float64
	}{
		{Multiplier: 2.0, Weight: 80.0},
		{Multiplier: 5.0, Weight: 15.0},
		{Multiplier: 100.0, Weight: 4.9},
		{Multiplier: 1000.0, Weight: 0.1},
	}

	// Расчет общей суммы "стоимости" всех выигрышей и их весов
	totalPayoutValue := 0.0
	totalWeight := 0.0
	for _, p := range basePayouts {
		totalPayoutValue += p.Multiplier * p.Weight
		totalWeight += p.Weight
	}

	// Рассчитываем общую вероятность выигрыша (totalWinProbability)
	// RTP = totalWinProbability * (totalPayoutValue / totalWeight)
	// => totalWinProbability = RTP * totalWeight / totalPayoutValue
	if totalPayoutValue == 0 {
		return []payout{{Multiplier: 0.0, Probability: 1.0}}
	}
	totalWinProbability := rtp * totalWeight / totalPayoutValue

	// Важная проверка: если общая вероятность выигрыша превышает 1.0, это означает,
	// что заданный RTP не может быть достигнут с текущими множителями,
	// так как RTP = 1.0 - это максимальное значение.
	if totalWinProbability >= 1.0 {
		log.Printf("Внимание: Заданный RTP (%.2f) равен или выше максимального. Вероятность проигрыша установлена в 0.0.", rtp)
		totalWinProbability = 1.0
	}

	// Создаем финальную таблицу выплат
	payouts := make([]payout, len(basePayouts)+1)

	// Добавляем проигрыш (0.0)
	payouts[0] = payout{Multiplier: 0.0, Probability: 1.0 - totalWinProbability}

	// Добавляем выигрыши с рассчитанными вероятностями
	for i, p := range basePayouts {
		winProbability := totalWinProbability * (p.Weight / totalWeight)
		payouts[i+1] = payout{Multiplier: p.Multiplier, Probability: winProbability}
	}

	// Логируем финальные вероятности
	log.Printf("Сгенерирована таблица выплат для RTP=%.2f:", rtp)
	for _, p := range payouts {
		log.Printf("  Множитель: %.2f, Вероятность: %.4f", p.Multiplier, p.Probability)
	}

	return payouts
}

func main() {
	// Устанавливаем сид для генератора случайных чисел, чтобы получать разные числа при каждом запуске
	rand.Seed(time.Now().UnixNano())

	// Определяем флаг для значения rtp со значением по умолчанию 0.8
	flag.Float64Var(&rtp, "rtp", 0.8, "Желаемое значение RTP (0 < rtp <= 1.0)")
	flag.Parse()

	// Проверяем значение rtp
	if rtp <= 0 || rtp > 1.0 {
		fmt.Println("Ошибка: Значение rtp должно быть в диапазоне (0, 1.0]")
		os.Exit(1)
	}

	// Генерируем таблицу выплат на основе заданного RTP
	payouts := generatePayouts(rtp)

	// Создаем новый ServeMux
	mux := http.NewServeMux()
	mux.Handle("/get", &Handler{payouts: payouts})

	// Определяем адрес сервера
	serverAddr := "localhost:64333"
	fmt.Printf("Запуск HTTP сервиса по адресу http://%s с RTP=%.2f\n", serverAddr, rtp)

	// Запускаем HTTP-сервер
	if err := http.ListenAndServe(serverAddr, mux); err != nil {
		log.Fatal("Ошибка ListenAndServe: ", err)
	}
}
