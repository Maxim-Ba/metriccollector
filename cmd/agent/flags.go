package main

import (
	"flag"
	"fmt"
	"os"
)

// неэкспортированная переменная flagRunAddr содержит адрес и порт для запуска сервера
var flagRunAddr string
var flagReportInterval int
var flagPollInterval int

// parseFlags обрабатывает аргументы командной строки 
// и сохраняет их значения в соответствующих переменных
func parseFlags() {
	// регистрируем переменную flagRunAddr 
	// как аргумент -a со значением :8080 по умолчанию
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&flagReportInterval, "r", 10, "seconds interval to send report")
	flag.IntVar(&flagPollInterval, "p", 2, "seconds interval to collect metrics")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()
	
	// Set a custom usage function to handle unknown flags
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	// Parse the flags
	flag.Parse()

	// Check for unknown flags
	for _, arg := range flag.Args() {
		if arg[0] == '-' {
			fmt.Fprintf(os.Stderr, "Unknown flag: %s\n", arg)
			flag.Usage()
			os.Exit(1)
		}
	}
} 
