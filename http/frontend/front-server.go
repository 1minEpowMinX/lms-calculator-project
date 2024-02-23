package main

import (
	"net/http"
)

func main() {
	// Установка маршрута для обслуживания статических файлов из каталога "static"
	fs := http.FileServer(http.Dir("http/frontend/static"))
	http.Handle("/", fs)

	// Запуск сервера на порту 8081
	http.ListenAndServe(":8081", nil)
}
