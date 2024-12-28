package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type FocusedWindow struct {
	AppName          string `json:"app_name"`
	WindowTitle      string `json:"window_title"`
	BundleIdentifier string `json:"bundle_identifier"`
}

func getFile(filePrefix, day string) string {
	return fmt.Sprintf("%s_%s.txt", filePrefix, day)
}

func storeEvent(window FocusedWindow) {
	currentTime := time.Now()
	formattedTime := strings.Split(currentTime.Format("2006-01-02 03:04:05"), " ")
	day, time := formattedTime[0], formattedTime[1]
	filePath := getFile("logs", day)
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	logLine := fmt.Sprintf("[%s] Currently focusing on application `%s (%s)` titled `%s`\n", time, window.AppName, window.BundleIdentifier, window.WindowTitle)
	_, err = writer.WriteString(logLine)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	writer.Flush()
}

func handleEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var focusedWindow FocusedWindow
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&focusedWindow); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	storeEvent(focusedWindow)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/event", handleEvent)
	port := "6969"
	fmt.Printf("Server listening on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
