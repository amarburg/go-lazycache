package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
)

var PipelineFlag = flag.String("pipeline", "", "Id of pipeline to run")
var MessageFlag = flag.String("message", "Started by Go", "Message to attach to run")

func main() {

	flag.Parse()

	if len(*PipelineFlag) < 1 {
		fmt.Println("Need to set --pipeline option")
		return
	}

	url := "https://app.wercker.com/api/v3/runs/"
	client := &http.Client{}

	form := struct {
		PipelineId string `json:"pipelineId"`
		Message    string `json:"message"`
	}{
		PipelineId: *PipelineFlag,
		Message:    *MessageFlag,
	}

	buf, _ := json.Marshal(form)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(buf))
	req.Header.Add("Content-Type", " application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", WerckerBearerToken))
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	log.Println(resp)

	respBuffer := new(bytes.Buffer)
	respBuffer.ReadFrom(resp.Body)
	fmt.Println(respBuffer.String())

	defer resp.Body.Close()
}
