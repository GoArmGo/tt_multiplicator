import requests
import sys
import time

def run_test(num_requests):
    """
    Отправляет заданное количество запросов на сервер и рассчитывает RTP.
    """
    url = "http://localhost:64333/get"
    total_winnings = 0.0
    win_count = 0
    
    print(f"Отправка {num_requests} запросов на {url}...")
    start_time = time.time()
    
    try:
        for _ in range(num_requests):
            response = requests.get(url)
            if response.status_code == 200:
                data = response.json()
                multiplier = data.get('result', 0.0)
                total_winnings += multiplier
                if multiplier > 0.0:
                    win_count += 1
            else:
                print(f"Ошибка при выполнении запроса: Статус-код {response.status_code}")
                return
    except requests.exceptions.RequestException as e:
        print(f"Ошибка при выполнении запроса: {e}")
        return

    end_time = time.time()
    
    print("\n--- Результаты теста ---")
    print(f"Количество запросов: {num_requests}")
    print(f"Количество выигрышных спинов: {win_count}")
    print(f"Общая сумма выигрышей: {total_winnings:.2f}")
    
    # Рассчитанный RTP
    if num_requests > 0:
        rtp = total_winnings / float(num_requests)
        print(f"Рассчитанный RTP: {rtp:.4f}")
    else:
        print("Рассчитанный RTP: N/A")
    
    print(f"Тест занял: {end_time - start_time:.4f} секунды")

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Использование: python inspector.py <количество запросов>")
    else:
        try:
            num_requests = int(sys.argv[1])
            run_test(num_requests)
        except ValueError:
            print("Ошибка: Количество запросов должно быть целым числом.")
