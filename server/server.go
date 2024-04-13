package server

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"lab2/client"
	"lab2/cryp"
	"net/http"
)

func RunServer() {
	http.HandleFunc("/verify", verifySignature)
	http.HandleFunc("/getsignmessage", getSignMessage)
	http.ListenAndServe(":8080", nil)
}

func verifySignature(w http.ResponseWriter, r *http.Request) {
	var data cryp.Data
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Не получилось распарсить данные из json", http.StatusBadRequest)
		return
	}
	fmt.Println("Сервер: полученное сообщение -", data.Message)
	err = cryp.VerifySignature(data)
	if err != nil {
		http.Error(w, "Подпись неверна", http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getSignMessage(w http.ResponseWriter, r *http.Request) {
	message, err := generateRandomMessage(10)
	if err != nil {
		fmt.Println("Сервер: Не удалось сгенерировать рандомное сообщение:", err)
	}
	fmt.Println("Сервер: сгенерированное сообщение -", message)
	privateKey, err := cryp.GenKeys("Сервер")
	if err != nil {
		fmt.Println("Сервер: Не удалось сгенерировать ключи:", err)
	}

	signatureHex, err := cryp.SignMessage(message, privateKey)
	if err != nil {
		fmt.Println("Сервер: Не получилось подписать сообщение:", err)
	}

	pubKeyPem := cryp.GetPubKeyPem(privateKey)

	json, err := client.CreateJson(message, signatureHex, pubKeyPem)
	if err != nil {
		fmt.Println("Сервер: Не получилось преобразовать данные в жысон:", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func generateRandomMessage(length int) (string, error) {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}
	for i := 0; i < length; i++ {
		buffer[i] = chars[int(buffer[i])%len(chars)]
	}
	return string(buffer), nil
}
