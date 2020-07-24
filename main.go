package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
	"moul.io/http2curl"
)

type HttpError struct {
	StatusCode int    `json:"statusCode"`
	Msg        string `json:"msg"`
	Key        string `json:"key"`
}

func (e *HttpError) Error() string {
	return e.Msg
}

func CreateErrorWithMsg(status int, key string, msg string) error {
	return &HttpError{StatusCode: status, Msg: msg, Key: key}
}
func CreateError(status int, key string) error {
	return &HttpError{StatusCode: status, Msg: key, Key: key}
}

// DoHTTPRequest Sends generic http request
func DoHTTPRequest(method string, url string, headers map[string]string, body []byte) (responseBody string, err error) {
	httpClient := &http.Client{}
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Close = true
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	command, _ := http2curl.GetCurlCommand(req)

	response, err := httpClient.Do(req)
	if err != nil {
		log.Printf("ERROR error requesting with http: %s, error: %v\n", command, err)
		err = CreateError(500, "failed_do_get")
		return
	}
	bodyBytes, err := ioutil.ReadAll(response.Body)
	response.Body.Close()

	if err != nil {
		log.Printf("ERROR error requesting with http: %s, error: %v\n", command, err)
		err = CreateError(500, "failed_read_body")
		return
	}

	responseBody = string(bodyBytes)

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		log.Printf("ERROR error requesting with http: %s, status code: %v, response: %s\n", command, response.StatusCode, responseBody)
		err = CreateError(500, "invalid_status")
		return
	}

	return
}

type SSHOptions struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	User string `json:"user"`
}

const (
	// KataGoBin the bin file path
	KataGoBin string = "/content/katago"
	// KataGoWeightFile the default weight file
	KataGoWeightFile string = "/content/weight.bin.gz"
	// KataGoConfigFile the default config file
	KataGoConfigFile string = "/content/katago-colab/config/gtp_colab.cfg"
	// KataGoChangeConfigScript changes the config
	KataGoChangeConfigScript string = "/content/katago-colab/scripts/change_config.sh"
)

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		log.Printf("ERROR usage: colab-katago USER_NAME USER_PASSWORD")
		return
	}
	username := args[0]
	userpassword := args[1]
	var newConfig *string = nil
	if len(args) >= 3 {
		newConfig = &args[2]
	}
	log.Printf("INFO using user name: %s password: %s\n", username, userpassword)
	sshJSONURL := "https://kata-config.oss-cn-beijing.aliyuncs.com/" + username + ".ssh.json"
	response, err := DoHTTPRequest("GET", sshJSONURL, nil, nil)
	if err != nil {
		log.Printf("ERROR error requestting url: %s, err: %+v\n", sshJSONURL, err)
		return
	}
	log.Printf("ssh options\n%s", response)
	sshoptions := SSHOptions{}
	// parse json
	err = json.Unmarshal([]byte(response), &sshoptions)
	if err != nil {
		log.Printf("ERROR failed parsing json: %s\n", response)
		return
	}

	config := &ssh.ClientConfig{
		Timeout:         30 * time.Second,
		User:            sshoptions.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	config.Auth = []ssh.AuthMethod{ssh.Password(userpassword)}

	addr := fmt.Sprintf("%s:%d", sshoptions.Host, sshoptions.Port)
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Fatal("failed to create ssh client", err)
		return
	}
	defer sshClient.Close()

	configFile := KataGoConfigFile
	if newConfig != nil {
		// start the sesssion to do it
		session, err := sshClient.NewSession()
		if err != nil {
			log.Fatal("failed to create ssh session", err)
			return
		}
		defer session.Close()

		cmd := fmt.Sprintf("%s %s", KataGoChangeConfigScript, *newConfig)
		log.Printf("DEBUG running commad:%s\n", cmd)
		configFile = fmt.Sprintf("/content/gtp_colab_%s.cfg", *newConfig)
		session.Run(cmd)

	}

	session, err := sshClient.NewSession()
	if err != nil {
		log.Fatal("failed to create ssh session", err)
		return
	}

	defer session.Close()
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	cmd := fmt.Sprintf("%s gtp -model %s -config %s", KataGoBin, KataGoWeightFile, configFile)
	log.Printf("DEBUG running commad:%s\n", cmd)
	session.Run(cmd)
}
