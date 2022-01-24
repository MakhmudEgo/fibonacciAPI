# fibonacciAPI

#### Сервис предоставляет http и grpc api для генерации последовательности чисел Фибоначчи
## build:
#### 1. Чтобы запустить сервис убедитесь, что у вас установлен docker ([см](https://www.docker.com)) и запущен.
#### 2. в терминале введите следующую команду:
		make или make run
### 3.1 http: 
	curl 'http://localhost:8888/fibonacci?from=1&to=10000'
### 3.2 grpc: 
		./grpc_client 1 2000
