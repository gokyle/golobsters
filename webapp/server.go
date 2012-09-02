package webapp

import (
        "log"
        "net/http"
        "os"
)

func setupHttpHandlers() {
        http.HandleFunc("/", rootPage)
}

// Start the HTTP server
func HttpServer() {
    log.Println("[+] webapp starting server")
    setupHttpHandlers()
    port := os.Getenv("PORT")
    if port == "" {
            port = "8080"
    }

    log.Println("[+] webapp will listen on port ", port)
    err := http.ListenAndServe(":"+port, nil)
    if err != nil {
        log.Fatal("[!] webapp listener: ", err)
    } 

    log.Println("[+] webapp exiting")
}
