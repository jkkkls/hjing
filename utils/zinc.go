package utils

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
)

func WriteZinc(domain, username, password, index string, data []byte) error {
	//http://localhost:4080/api/games3/_doc
	//http://localhost:5080/api/default/quickstart1/_json
	req, err := http.NewRequest("POST", fmt.Sprintf("%v/api/default/%v/_json", domain, index), bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// log.Println(resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if !bytes.Contains(body, []byte("\"successful\":1")) {
		fmt.Println("WriteZinc error", string(body))
	}
	return nil
}
