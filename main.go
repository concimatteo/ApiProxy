package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	// Estrae il sito dal path (tutto dopo /get2post/)
	path := strings.TrimPrefix(r.URL.Path, "/get2post/")
	if path == "" || path == r.URL.Path {
		http.Error(w, "Sito non specificato. Usa /get2post/<sito>", http.StatusBadRequest)
		return
	}

	// Prepara l'URL di destinazione
	targetURL := path
	if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
		targetURL = "https://" + targetURL
	}

	// Converte i parametri GET in form data per POST
	formData := url.Values{}
	for key, values := range r.URL.Query() {
		for _, value := range values {
			formData.Add(key, value)
		}
	}

	// Log della richiesta
	log.Printf("Inoltro richiesta a: %s con parametri: %v", targetURL, formData)

	// Crea la richiesta POST
	resp, err := http.PostForm(targetURL, formData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Errore nella richiesta POST: %v", err), http.StatusInternalServerError)
		log.Printf("Errore: %v", err)
		return
	}
	defer resp.Body.Close()

	// Copia gli header della risposta
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Imposta lo status code
	w.WriteHeader(resp.StatusCode)

	// Copia il body della risposta
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("Errore nella copia della risposta: %v", err)
	}

	log.Printf("Risposta ricevuta con status: %d", resp.StatusCode)
}

func main() {
	http.HandleFunc("/get2post/", proxyHandler)

	// Usa PORT da variabile d'ambiente, default 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	log.Printf("Server in ascolto su porta %s", port)
	log.Printf("Esempio: http://localhost:8080/get2post/dati.meteotrentino.it/service.asmx/datiRealtimeUnaStazione?stazione=T0135&h=1", port)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
