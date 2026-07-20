package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/oneclickvirt/UnlockTests/transnation"
)

func main() {
	client := &http.Client{Timeout: 10 * time.Second}
	results := []struct {
		Name   string `json:"name"`
		Status string `json:"status"`
	}{
		{Name: "dola", Status: transnation.Dola(client).Status},
		{Name: "x", Status: transnation.X(client).Status},
	}
	encoded, _ := json.Marshal(results)
	fmt.Println(string(encoded))
}
