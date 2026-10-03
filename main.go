package main

import (	
      "encoding/json"
      "log"
      "net/http"
)

func main() {
      http.HandleFunc("/health", healthHandler)
      http.HandleFunc("/infer", inferHandler)

      log.Println("CALISIYOR LAAAANNN :8080")
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

      w.Header().Set("Content-Type", "application/json")
      json.NewEncoder(w).Encode(map[string]string{
              "input":  input,
              "result": "inference placeholder: " + input,
              "model":  "edge-ai-v1",
      })
}