package transcript

import (
	"audigo/models"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func Extract(m *models.Audio) error {
	apiKey := os.Getenv("ASSEMBLYAI_API_KEY")
	if apiKey == "" {
		fmt.Println("missing ASSEMBLYAI_API_KEY. skipping transcript extraction")
		return nil
	}
	const UPLOAD_URL = "https://api.assemblyai.com/v2/upload"

	data, err := os.ReadFile(m.Path)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, _ := http.NewRequest("Post", UPLOAD_URL, bytes.NewBuffer(data))
	req.Header.Set("authorization", apiKey)
	res, err := client.Do(req)
	defer res.Body.Close()

	if err != nil {
		log.Fatalln(err)
	}

	var result map[string]interface{}
	json.NewDecoder(res.Body).Decode(&result)

	fmt.Println(result["upload_url"])

	AUDIO_URL := fmt.Sprintf("%s", result["upload_url"])
	fmt.Println("AUDIO_URL: ", AUDIO_URL)
	const TRANSCRIPT_URL = "https://api.assemblyai.com/v2/transcript"

	values := map[string]string{"audio_url": AUDIO_URL}
	jsonData, err := json.Marshal(values)
	if err != nil {
		log.Fatalln(err)
	}

	client = &http.Client{}
	req, _ = http.NewRequest("POST", TRANSCRIPT_URL, bytes.NewBuffer(jsonData))
	req.Header.Set("content-type", "application/json")
	req.Header.Set("authorization", apiKey)
	res, err = client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(&result)

	fmt.Println(result["id"])
	resultId := fmt.Sprintf("%s", result["id"])
	POLLING_URL := TRANSCRIPT_URL + "/" + resultId

	for {
		client = &http.Client{}
		req, _ = http.NewRequest("GET", POLLING_URL, nil)
		req.Header.Set("content-type", "application/json")
		req.Header.Set("authorization", apiKey)
		res, err = client.Do(req)
		if err != nil {
			log.Fatalln(err)
		}

		defer res.Body.Close()

		json.NewDecoder(res.Body).Decode(&result)

		if result["status"] == "completed" {
			fmt.Println("Status is completes...")
			fmt.Println(result["text"])
			m.Metadata.Transcript = fmt.Sprintf("%s", result["text"])
			fmt.Println("m.Metadata.Transcript: ", m.Metadata.Transcript)

			break
		} else {
		}
	}
	return nil
}
