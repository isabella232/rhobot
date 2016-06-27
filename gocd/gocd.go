package gocd

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
)

// PipelineConfig a GoCD structure that contains a pipeline and a group
type PipelineConfig struct {
	Group    string   `json:"group"`
	Pipeline Pipeline `json:"pipeline"`
}

// Pipeline a GoCD structure that represents a pipeline
type Pipeline struct {
	LabelTemplate         string        `json:"label_template"`
	EnablePipelineLocking bool          `json:"enable_pipeline_locking"`
	Name                  string        `json:"name"`
	Template              interface{}   `json:"template"`
	Parameters            []interface{} `json:"parameters"`
	EnvironmentVariables  []struct {
		Secure bool   `json:"secure"`
		Name   string `json:"name"`
		Value  string `json:"value"`
	} `json:"environment_variables"`
	Materials []struct {
		Type       string `json:"type"`
		Attributes struct {
			URL             string      `json:"url"`
			Destination     string      `json:"destination"`
			Filter          interface{} `json:"filter"`
			Name            interface{} `json:"name"`
			AutoUpdate      bool        `json:"auto_update"`
			Branch          string      `json:"branch"`
			SubmoduleFolder interface{} `json:"submodule_folder"`
		} `json:"attributes"`
	} `json:"materials"`
	Stages []struct {
		Name                  string `json:"name"`
		FetchMaterials        bool   `json:"fetch_materials"`
		CleanWorkingDirectory bool   `json:"clean_working_directory"`
		NeverCleanupArtifacts bool   `json:"never_cleanup_artifacts"`
		Approval              struct {
			Type          string `json:"type"`
			Authorization struct {
				Roles []interface{} `json:"roles"`
				Users []interface{} `json:"users"`
			} `json:"authorization"`
		} `json:"approval"`
		EnvironmentVariables []interface{} `json:"environment_variables"`
		Jobs                 []struct {
			Name                 string        `json:"name"`
			RunInstanceCount     interface{}   `json:"run_instance_count"`
			Timeout              interface{}   `json:"timeout"`
			EnvironmentVariables []interface{} `json:"environment_variables"`
			Resources            []interface{} `json:"resources"`
			Tasks                []struct {
				Type       string `json:"type"`
				Attributes struct {
					RunIf            []string    `json:"run_if"`
					OnCancel         interface{} `json:"on_cancel"`
					Command          string      `json:"command"`
					Arguments        []string    `json:"arguments"`
					WorkingDirectory string      `json:"working_directory"`
				} `json:"attributes"`
			} `json:"tasks"`
			Tabs       []interface{} `json:"tabs"`
			Artifacts  []interface{} `json:"artifacts"`
			Properties interface{}   `json:"properties"`
		} `json:"jobs"`
	} `json:"stages"`
	TrackingTool interface{} `json:"tracking_tool"`
	Timer        interface{} `json:"timer"`
}

// Server a representstion of a GoCD server
type Server struct {
	Host     string
	Port     string
	User     string
	Password string
}

// client returns a http client with longer timeout and skip verify
func client() *http.Client {
	timeout := time.Duration(120 * time.Second)
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{
		Timeout:   timeout,
		Transport: transCfg,
	}
}

// NewServerConfig Create a Server object from a config
func NewServerConfig(host string, port string, user string, password string) *Server {
	return &Server{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
	}
}

// URL returns the host of the GoCD server
func (server Server) URL() string {
	return fmt.Sprintf("%s:%s", server.Host, server.Port)
}

// readPipelineJSONFromFile reads a GoCD structure from a json file
func readPipelineJSONFromFile(path string) (pipeline Pipeline, err error) {
	data, err := ioutil.ReadFile(path)
	if err == nil {
		err = json.Unmarshal(data, &pipeline)
	}
	return
}

// Partially generated by curl-to-Go: https://mholt.github.io/curl-to-go
func (server Server) pipelineConfigPUT(pipeline Pipeline, etag string) (pipelineResult Pipeline, err error) {

	pipelineName := pipeline.Name

	payloadBytes, err := json.Marshal(pipeline)
	if err != nil {
		return
	}

	payloadBody := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("PUT", server.URL()+"/go/api/admin/pipelines/"+pipelineName, payloadBody)
	if err != nil {
		return
	}

	req.SetBasicAuth(server.User, server.Password)
	req.Header.Set("Accept", "application/vnd.go.cd.v1+json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("If-Match", etag)

	log.Debugf("Sending request: %v", req)
	resp, err := client().Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("Bad response code: %d, response: %s", resp.StatusCode, body)
		return
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "\t")
	if err != nil {
		log.Warn("Failed to prettify JSON: ", err)
	}

	log.Debug("pipelineConfig JSON:", string(prettyJSON.Bytes()))
	err = json.Unmarshal(body, &pipelineResult)
	return
}

// Generated by curl-to-Go: https://mholt.github.io/curl-to-go
func (server Server) pipelineConfigPOST(pipelineConfig PipelineConfig) (pipeline Pipeline, err error) {
	payloadBytes, err := json.Marshal(pipelineConfig)
	if err != nil {
		return
	}

	payloadBody := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", server.URL()+"/go/api/admin/pipelines", payloadBody)
	if err != nil {
		return
	}

	req.SetBasicAuth(server.User, server.Password)
	req.Header.Set("Accept", "application/vnd.go.cd.v1+json")
	req.Header.Set("Content-Type", "application/json")

	log.Debugf("Sending request: %v", req)
	resp, err := client().Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("Bad response code: %d with response: %s", resp.StatusCode, body)
		return
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "\t")
	if err != nil {
		log.Warn("Failed to prettify JSON: ", err)
	}

	log.Debug("pipelineConfig JSON: ", string(prettyJSON.Bytes()))
	err = json.Unmarshal(body, &pipeline)
	return
}

// Partially generated by curl-to-Go: https://mholt.github.io/curl-to-go
func (server Server) pipelineGET(pipelineName string) (pipeline Pipeline, etag string, err error) {
	req, err := http.NewRequest("GET", server.URL()+"/go/api/admin/pipelines/"+pipelineName, nil)
	if err != nil {
		return
	}

	req.SetBasicAuth(server.User, server.Password)
	req.Header.Set("Accept", "application/vnd.go.cd.v1+json")
	req.Header.Set("Content-Type", "application/json")

	log.Debugf("Sending request: %v", req)
	resp, err := client().Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("Bad response code: %d with response: %s", resp.StatusCode, body)
		return
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "\t")
	if err != nil {
		log.Warn("Failed to prettify JSON: ", err)
	}

	log.Debug("pipelineConfig JSON:", string(prettyJSON.Bytes()))

	etag = resp.Header.Get("ETag")
	err = json.Unmarshal(body, &pipeline)
	return
}

// Push takes a pipeline from a file and sends it to GoCD
func Push(server *Server, path string, group string) (err error) {
	pipeline, err := readPipelineJSONFromFile(path)
	if err != nil {
		return
	}

	etag, err := Exist(server, pipeline.Name)
	if err != nil {
		log.Info(err)
	}

	if etag == "" {
		pipelineConfig := PipelineConfig{group, pipeline}
		_, err = server.pipelineConfigPOST(pipelineConfig)
	} else {
		_, err = server.pipelineConfigPUT(pipeline, etag)
	}
	return
}

// Pull reads pipeline from a file, finds it on GoCD, and updates the file
func Pull(server *Server, path string) (err error) {
	pipeline, err := readPipelineJSONFromFile(path)
	if err != nil {
		return
	}

	name := pipeline.Name
	err = Clone(server, path, name)
	return
}

// Exist checks if a pipeline of a given name exist, returns it's etag or an empty string
func Exist(server *Server, name string) (etag string, err error) {
	_, etag, err = server.pipelineGET(name)
	return
}

// Clone finds a pipeline by name on GoCD and saves it to a file
func Clone(server *Server, path string, name string) (err error) {
	pipelineFetched, _, err := server.pipelineGET(name)
	if err != nil {
		return
	}

	pipelineJSON, _ := json.MarshalIndent(pipelineFetched, "", "    ")
	err = ioutil.WriteFile(path, pipelineJSON, 0666)
	return
}