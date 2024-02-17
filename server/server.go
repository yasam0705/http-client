package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func Run(port string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/upload", handle)

	log.Println("server is listened", port)
	fmt.Println()

	if err := http.ListenAndServe(port, mux); err != nil {
		return err
	}
	return nil
}

func handle(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error receiving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileName := handler.Filename
	// fmt.Printf("File uploaded: %+v\n", handler.Filename)
	// fmt.Printf("File size: %+v\n", handler.Size)
	// fmt.Printf("MIME Header: %+v\n", handler.Header)

	dst, err := os.Create("./" + fileName)
	if err != nil {
		http.Error(w, "Error creating file on server", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Error copying file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully: %s", handler.Filename)

}
