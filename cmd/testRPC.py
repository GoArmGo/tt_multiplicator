import requests
import json
import time

SERVICE_URL = "http://localhost:64333/get"

NUM_REQUESTS = 11000

print("Запуск теста RTP")
print(f"Отправится {NUM_REQUESTS} запросов на {SERVICE_URL}")

total_sum = 0
successful_requests = 0

start_time = time.time()

for i in range(NUM_REQUESTS):
    try:
        response = requests.get(SERVICE_URL)
        response.raise_for_status()
        
        data = json.loads(response.text)
        
        multiplier = data["result"]
        
         # Если multiplier > 1.0 (то есть, это 10000.0), считаем его за "успех"
        if multiplier > 1.0:
            total_sum += 1
        
        successful_requests += 1

    except requests.exceptions.RequestException as e:
        print(f"Ошибка при выполнении запроса: {e}")
        break

end_time = time.time()
duration = end_time - start_time

# Расчет RTP на основе результатов теста
if successful_requests > 0:
    calculated_rtp = total_sum / successful_requests
    print("\n--- Результаты теста ---")
    print(f"Количество успешных запросов: {successful_requests}")
    print(f"Количество 'успешных' мультипликаторов: {total_sum}")
    print(f"Рассчитанный RTP (отношение 'успехов' к общему числу запросов): {calculated_rtp:.4f}")
    print(f"Тест занял: {duration:.2f} секунды")
else:
    print("Не удалось выполнить ни одного запроса")
