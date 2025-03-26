package main

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "io"
    "net/http"
    "os"
    "os/exec"
)

var (
    secret    = os.Getenv("WEBHOOK_SECRET")
    port      = ":9000"
    playbook  = "test-build-run.yml"
    inventory = "hosts.ini"
)

func verifySignature(secret, body []byte, signature string) bool {
    mac := hmac.New(sha256.New, secret)
    mac.Write(body)
    expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
    return hmac.Equal([]byte(expected), []byte(signature))
}

func handler(w http.ResponseWriter, r *http.Request) {
    signature := r.Header.Get("X-Hub-Signature-256")
    if signature == "" {
        http.Error(w, "Missing signature", http.StatusBadRequest)
        return
    }

    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Error reading body", http.StatusInternalServerError)
        return
    }

    if !verifySignature([]byte(secret), body, signature) {
        http.Error(w, "Invalid signature", http.StatusForbidden)
        return
    }

    go func() {
        cmd := exec.Command("ansible-playbook", "-i", inventory, playbook)
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
        cmd.Run()
    }()

    w.WriteHeader(http.StatusOK)
    fmt.Fprintln(w, "Build triggered")
}

func main() {
    if secret == "" {
        fmt.Println("WEBHOOK_SECRET must be set")
        os.Exit(1)
    }

    http.HandleFunc("/", handler)
    fmt.Println("Listening on", port)
    if err := http.ListenAndServe(port, nil); err != nil {
        fmt.Println("Server error:", err)
    }
}
