package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// main is the entry point of the program.
//
// It doesn't take any parameters and doesn't return anything.
func main() {
	// Получение текущей директории
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Ошибка при получении текущей директории: %v", err)
	}

	// Пути к файлам бекенда и фронта
	backendPath := filepath.Join(currentDir, "http", "backend", "back-server.go")
	frontendPath := filepath.Join(currentDir, "http", "frontend", "front-server.go")

	// Запуск бекенда
	go runServer(backendPath)

	// Запуск фронта
	go runServer(frontendPath)

	// Бесконечный цикл, чтобы программа не завершилась сразу после запуска
	select {}
}

// runServer runs the server using the specified file path.
//
// filePath string
func runServer(filePath string) {
	cmd := exec.Command("go", "run", filePath)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
