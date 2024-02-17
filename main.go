package main

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"test-tasks/http-multipart/server"
)

const (
	port     = ":3030"
	hostname = "http://localhost"
	endpoint = "/upload"

	fileName = "files/microservices.pdf"
)

func main() {
	go func() {
		if err := server.Run(port); err != nil {
			log.Fatal(err)
		}
	}()

	err := upload(fileName)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("File uploaded successfully.")
}

func upload(fileName string) error {
	file, _ := os.Open(fileName)
	defer file.Close()

	body, pipeWriter := io.Pipe()
	writer := multipart.NewWriter(pipeWriter)

	go func() {
		defer pipeWriter.Close()
		defer writer.Close()

		part, err := writer.CreateFormFile("file", file.Name())
		if err != nil {
			fmt.Println("Error creating form file:", err)
			return
		}

		if _, err := io.Copy(part, file); err != nil {
			fmt.Println("Error copying file contents to part:", err)
			return
		}
	}()

	u, err := url.Parse(hostname)
	if err != nil {
		fmt.Println("Error parsing the host:", err)
		return err
	}
	u.Host += port
	u.Path = endpoint

	r, err := http.NewRequest(http.MethodPost, u.String(), body)
	if err != nil {
		return err
	}
	r.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}

	resp, err := client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: %s", resp.Status)
	}

	bb, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("=====", string(bb))

	return nil
}
