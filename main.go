package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	listen = flag.String("listen", ":8080", "host:port to listen on")
)

func main() {
	flag.Parse()

	http.HandleFunc("/overlay", handleOverlay)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	fmt.Printf("Listening on %s\n", *listen)
	panic(http.ListenAndServe(*listen, nil))
}

func handleOverlay(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	imgURL := params.Get("url")
	dptNum := params.Get("dpt")

	if imgURL == "" {
		bye(w, "Missing url")
		return
	}

	if dptNum == "" {
		bye(w, "Missing dpt")
	}

	resp, err := http.Get(imgURL)
	if err != nil {
		bye(w, fmt.Sprintf("Could not get image at url %s", imgURL))
		log.Printf("Error GETing image: %s", imgURL)
		log.Printf("Error: %s", err)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "image/png")
	err = overlayPipe(resp.Body, dptNum, w)

	if err != nil {
		bye(w, "Error overlaying")
		log.Printf("Error overlaying: %s\nError: %s\n", imgURL, err)
		return
	}
}

func bye(w http.ResponseWriter, msg string) {
	w.WriteHeader(400)
	w.Write([]byte(msg))
}
