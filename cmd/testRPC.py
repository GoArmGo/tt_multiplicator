import requests
import sys

# Проверка наличия аргументов командной строки
if len(sys.argv) < 2:
    print("Использование: python inspector.py <количество запросов>")
    sys.exit(1)

try:
    num_requests = int(sys.argv[1])
except ValueError:
    print("Количество запросов должно быть целым числом.")
    sys.exit(1)

server_url = "http://localhost:64333/get"
total_winnings = 0.0
successful_requests = 0

print(f"Отправка {num_requests} запросов на {server_url}...")

# Отправка запросов и подсчет результатов
for i in range(num_requests):
    try:
        response = requests.get(server_url)
        if response.status_code == 200:
            multiplier_data = response.json()
            multiplier = float(multiplier_data.get('result', 0.0))
            total_winnings += multiplier
            if multiplier > 0.0:  # Считаем успешными все выигрыши, кроме 0
                successful_requests += 1

    except requests.exceptions.RequestException as e:
        print(f"Ошибка при выполнении запроса: {e}")
        sys.exit(1)

# Расчет итогового RTP
# RTP = (общая_сумма_выигрышей / общее_количество_запросов)
final_rtp = total_winnings / num_requests

print("\n--- Результаты теста ---")
print(f"Количество запросов: {num_requests}")
print(f"Количество выигрышных спинов: {successful_requests}")
print(f"Общая сумма выигрышей: {total_winnings:.2f}")
print(f"Рассчитанный RTP: {final_rtp:.4f}")
print(f"Тест занял: {response.elapsed.total_seconds():.2f} секунды")
