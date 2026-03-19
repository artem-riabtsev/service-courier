# Профилирование микросервиса с pprof

##  Доступные эндпоинты профилирования

PProf сервер запускается на **localhost:6060** (только локальный доступ):

- `http://localhost:6060/debug/pprof/` -	Главная страница со списком всех профилей
- `http://localhost:6060/debug/pprof/profile` - CPU профиль (по умолчанию 30 секунд)
- `http://localhost:6060/debug/pprof/heap` - Профиль использования памяти (heap)
- `http://localhost:6060/debug/pprof/goroutine` - Дамп всех активных горутин
- `http://localhost:6060/debug/pprof/allocs` - Все аллокации памяти
- `http://localhost:6060/debug/pprof/threadcreate` - Создание потоков ОС
- `http://localhost:6060/debug/pprof/block` - Блокировки (contention)
- `http://localhost:6060/debug/pprof/mutex` - Мьютексы
- `http://localhost:6060/debug/pprof/trace` - Трассировка выполнения (trace)
- `http://localhost:6060/debug/pprof/symbol` - Симантизация адресов

##  Как использовать

- curl http://localhost:6060/debug/pprof/

- curl http://localhost:6060/debug/pprof/goroutine?debug=1

- go tool pprof http://localhost:6060/debug/pprof/profile?seconds=10

- go tool pprof -http=:8081 http://localhost:6060/debug/pprof/profile?seconds=30

- go tool pprof http://localhost:6060/debug/pprof/heap

- go tool pprof http://localhost:6060/debug/pprof/allocs
