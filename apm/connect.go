package apm

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

)


type ConnectReply struct {
	RunID                 string        `json:"agent_run_id"`
	EntityGUID            string            `json:"entity_guid"`
}


func UnmarshalConnectReply(body []byte) (*ConnectReply, error) {
	var reply struct {
		Reply *ConnectReply `json:"return_value"`
	}
	err := json.Unmarshal(body, &reply)
	if nil != err {
		return nil, fmt.Errorf("unable to parse connect reply: %v", err)
	}
	return reply.Reply, nil
}


func Connect(cmd RpmCmd, cs RpmControls) (string, string) {

	pid := os.Getpid()
	headers := map[string]string{
		"User-Agent":       "NewRelic-Go-Agent/3.35.1",
		"Accept-Encoding":  "deflate",
		"Content-Encoding": "gzip", 
		"Content-Type":     "application/json",
	}
	AppName := os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
	data := []map[string]interface{}{
		{
			"pid":           pid,
			"language":      "go",
			"agent_version": "3.35.1",
			"host":          "AWS Lambda",
			"app_name":      []string{AppName},
			"identifier":     AppName,
		},
	}

	urlStr := RpmURL(cmd, cs)

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return "nil", "nil"
	}

	var compressedData bytes.Buffer
	writer := gzip.NewWriter(&compressedData)
	_, err = writer.Write(jsonData)
	if err != nil {
		fmt.Println("Error writing gzip data:", err)
		return "nil", "nil"
	}

	writer.Close()

	req, err := http.NewRequest("POST", urlStr, &compressedData)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "nil", "nil"
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error executing request:", err)
		return "nil", "nil"
	}
	defer resp.Body.Close()

	r := newRPMResponse(nil).AddStatusCode(resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if r.GetError() == nil {
		r.SetError(err)
	}
	r.AddBody(body)
	connectReponse, _:= UnmarshalConnectReply(body)
	if connectReponse != nil  {
		if connectReponse.EntityGUID == "" && connectReponse.RunID == "" {
			fmt.Println("Connect Response Unsuccessful")
			
		}
	}
	return connectReponse.RunID, connectReponse.EntityGUID

}
