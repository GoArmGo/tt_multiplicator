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

// Таблица выплат с заданными вероятностями.
// Сумма произведений (Multiplier * Probability) должна быть равна RTP.
// RTP = (0.0 * 0.9) + (2.0 * 0.08) + (5.0 * 0.015) + (100.0 * 0.0049) + (1000.0 * 0.0001) = 0.8
var payouts = []payout{
	{Multiplier: 0.0, Probability: 0.9},
	{Multiplier: 2.0, Probability: 0.08},
	{Multiplier: 5.0, Probability: 0.015},
	{Multiplier: 100.0, Probability: 0.0049},
	{Multiplier: 1000.0, Probability: 0.0001},
}

// MultiplierResponse - структура для вывода JSON
type MultiplierResponse struct {
	Result float64 `json:"result"`
}

// getHandler генерирует мультипликатор на основе глобальной переменной rtp
func getHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// Функция для генерации случайного мультипликатора на основе таблицы выплат.
	generateMultiplier := func() float64 {
		p := rand.Float64()
		cumulativeProbability := 0.0
		for _, payout := range payouts {
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

	// Создаем новый ServeMux
	mux := http.NewServeMux()
	mux.HandleFunc("/get", getHandler)

	// Определяем адрес сервера
	serverAddr := "localhost:64333"
	fmt.Printf("Запуск HTTP сервиса по адресу http://%s с RTP=%.2f\n", serverAddr, rtp)

	// Запускаем HTTP-сервер
	if err := http.ListenAndServe(serverAddr, mux); err != nil {
		log.Fatal("Ошибка ListenAndServe: ", err)
	}
}
