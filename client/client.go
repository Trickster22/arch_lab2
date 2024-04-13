package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"lab2/cryp"
	"net/http"
)

func Scenario1() {
	client := &http.Client{}
	fmt.Println("Клиент: Придумайте сообщение")
	var message string
	fmt.Scan(&message)
	fmt.Println("Клиент: Хешируем ваше сообщение")

	privateKey, err := cryp.GenKeys("Клиент")
	if err != nil {
		fmt.Println("Клиент: Не удалось сгенерировать ключи:", err)
	}

	signatureHex, err := cryp.SignMessage(message, privateKey)
	if err != nil {
		fmt.Println("Клиент: Не получилось подписать сообщение:", err)
	}

	pubKeyPem := cryp.GetPubKeyPem(privateKey)

	json, err := CreateJson(message, signatureHex, pubKeyPem)
	if err != nil {
		fmt.Println("Клиент: Не получилось преобразовать данные в жысон:", err)
		return
	}

	resp, err := client.Post("http://localhost:8080/verify", "application/json", bytes.NewBuffer(json))
	if err != nil {
		fmt.Println("Клиент: Не получилось выполнить POST запрос:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Клиент: Сервер сказал, что подпись верна :)")
	} else {
		fmt.Println("Клиент: Сервер сказал, что подпись неверна, вы жулик >:(")
	}
}

func Scenario2() {
	client := &http.Client{}
	resp, err := client.Post("http://localhost:8080/getsignmessage", "application/json", nil)
	if err != nil {
		fmt.Println("Клиент: Ошибка при отправке запроса:", err)
		return
	}
	defer resp.Body.Close()
	var data cryp.Data
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Println("Клиент: Не получилось распарсить данные из json:", err)
		return
	}
	fmt.Println("Клиент: полученное сообщение -", data.Message)
	err = cryp.VerifySignature(data)
	if err != nil {
		fmt.Println("Клиент: Подпись неверна >:(")
	} else {
		fmt.Println("Клиент: Подпись верна :)")
	}

}

func CreateJson(message, signatureHex, pubKeyPem string) ([]byte, error) {
	data := map[string]interface{}{
		"message":   message,
		"signature": signatureHex,
		"publickey": pubKeyPem,
	}

	json, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return json, nil
}
