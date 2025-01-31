package apm

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type PreconnectReply struct {
	Collector        string           `json:"redirect_host"`
}

func UnmarshalPreConnectReply(body []byte) (*PreconnectReply, error) {
	var preconnect struct {
		Preconnect PreconnectReply `json:"return_value"`
	}
	err := json.Unmarshal(body, &preconnect)
	if nil != err {
		return nil, fmt.Errorf("unable to parse pre-connect reply: %v", err)
	}
	return &preconnect.Preconnect, nil
}

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

type preconnectRequest struct {
	SecurityPoliciesToken string `json:"security_policies_token,omitempty"`
	HighSecurity          bool   `json:"high_security"`
}

func PreConnect(cmd RpmCmd, cs *RpmControls) string{
	preconnectData, _ := json.Marshal([]preconnectRequest{{
		SecurityPoliciesToken: "",
		HighSecurity:          false,
	}})
	cmd.Data = preconnectData
	cmd.Name = cmdPreconnect
	resp := CollectorRequest(cmd, cs)
	body, _ := io.ReadAll(resp.GetBody())
	preConnectReponse, _:= UnmarshalPreConnectReply(body)
	return preConnectReponse.Collector
}

func Connect(cmd RpmCmd, cs *RpmControls) (string, string) {

	pid := os.Getpid()
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
	cmd.Data, _ = json.Marshal(data)
	cmd.Name = cmdConnect
	resp := CollectorRequest(cmd, cs)

	body, _ := io.ReadAll(resp.GetBody())
	connectReponse, _:= UnmarshalConnectReply(body)
	if connectReponse != nil  {
		if connectReponse.EntityGUID == "" && connectReponse.RunID == "" {
			fmt.Println("Connect Response Unsuccessful")
			
		}
	}
	cs.SetRunId(connectReponse.RunID)
	cs.SetEntityGuid(connectReponse.EntityGUID)
	return connectReponse.RunID, connectReponse.EntityGUID

}
