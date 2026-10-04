package main

import (
      "bytes"
      "encoding/json"
      "io"
      "log"
      "net/http"
      "os"
)

var ollamaURL string

func main() {
      ollamaURL = os.Getenv("OLLAMA_URL")
      if ollamaURL == "" {
              ollamaURL = "http://localhost:11434"
      }

      http.HandleFunc("/health", healthHandler)
      http.HandleFunc("/infer", inferHandler)

      log.Println("edge-ai-demo starting on :8080")
      log.Fatal(http.ListenAndServe(":8080", nil))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
      w.Header().Set("Content-Type", "application/json")
      json.NewEncoder(w).Encode(map[string]string{
              "status":  "ok",
              "service": "edge-ai-demo",
      })
}

func inferHandler(w http.ResponseWriter, r *http.Request) {
      if r.Method != http.MethodPost {
              http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
              return
      }

      var req map[string]string
      if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
              http.Error(w, "invalid json", http.StatusBadRequest)
              return
      }

      input, ok := req["input"]
      if !ok {
              http.Error(w, "missing input field", http.StatusBadRequest)
              return
      }

      result, err := callOllama(input)
      if err != nil {
              http.Error(w, "ollama error: "+err.Error(), http.StatusInternalServerError)
              return
      }

      w.Header().Set("Content-Type", "application/json")
      json.NewEncoder(w).Encode(map[string]string{
              "input":  input,
              "result": result,
              "model":  "tinyllama",
      })
}

func callOllama(prompt string) (string, error) {
      body, _ := json.Marshal(map[string]interface{}{
              "model":  "tinyllama",
              "prompt": prompt,
              "stream": false,
      })

      resp, err := http.Post(ollamaURL+"/api/generate", "application/json", bytes.NewBuffer(body))
      if err != nil {
              return "", err
      }
      defer resp.Body.Close()

      data, _ := io.ReadAll(resp.Body)
      var result map[string]interface{}
      json.Unmarshal(data, &result)

      return result["response"].(string), nil
}