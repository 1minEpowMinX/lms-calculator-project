package main

import (
	"log"
	"os/exec"
)

func main() {
	// Пути к файлам бекенда и фронта
	backendPath := "http/backend/back-server.go"
	frontendPath := "http/frontend/front-server.go"

	// Запуск бекенда
	go runBackend(backendPath)

	// Запуск фронта
	go runFrontend(frontendPath)

	// Бесконечный цикл, чтобы программа не завершилась сразу после запуска
	select {}
}

func runBackend(backendPath string) {
	cmd := exec.Command("go", "run", backendPath)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Ошибка при запуске бекенда: %v", err)
	}
}

func runFrontend(frontendPath string) {
	cmd := exec.Command("go", "run", frontendPath)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Ошибка при запуске фронта: %v", err)
	}
}
