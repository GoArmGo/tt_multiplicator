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

// MultiplierResponse - структура для вывода JSON
type MultiplierResponse struct {
	Result float64 `json:"result"`
}

// getHandler генерирует мультипликатор на основе глобальной переменной rtp
func getHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// Генерируем случайное число от 0 до 1
	randomNumber := rand.Float64()

	var multiplier float64

	// Если заданный RTP равен 1.0 всегда возвращаем "успех"
	if rtp == 1.0 {
		multiplier = 10000.0
	} else if randomNumber < rtp {
		// Если случайное число меньше RTP возвращаем "успех"
		multiplier = 10000.0 // Наибольшее число в допустимом диапазоне
	} else {
		// В противном случае возвращаем "неудачу"
		multiplier = 1.0 // Наименьшее число в допустимом диапазоне
	}

	// Создаем объект ответа
	response := MultiplierResponse{
		Result: multiplier,
	}

	// Кодируем объект ответа в JSON и записываем в writer ответа
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Ошибка при кодировании JSON ответа", http.StatusInternalServerError)
		return
	}

	// Логируем сгенерированный мультипликатор и случайное число для отладки
	log.Printf("Сгенерированный мультипликатор: %.2f (случайное число: %.2f, rtp: %.2f)\n", multiplier, randomNumber, rtp)
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
