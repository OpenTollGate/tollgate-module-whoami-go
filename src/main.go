package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

type WhoamiResult struct {
	Success      bool
	ErrorMessage string
	Ip           string
	Mac          string
}

func getMacByIp(ip string) (string, error) {
	cmd := exec.Command("ndsctl", "json", ip)

	commandOutput, err := cmd.CombinedOutput()
	var commandOutputString = string(commandOutput)
	if err != nil {
		fmt.Println(err, "Error when trying to find mac address for ip "+ip+"\n Command output: "+commandOutputString)

		if strings.Contains(commandOutputString, "ndsctl: opennds probably not yet started (Error: Connection refused)") {
			return "", errors.New("Internal Error: OpenNDS not running")
		}

		return "", err
	}

	var jsonMap map[string]interface{}
	err = json.Unmarshal([]byte(commandOutputString), &jsonMap)
	if err != nil {
		return "", err
	}

	var mac = jsonMap["mac"]

	if mac == nil {
		log.Println("ndsctl did not return a MAC address for ip " + ip)
		return "", errors.New("Internal Error: OpenNDS did not return a MAC address for ip " + ip)
	}

	return mac.(string), nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	var ip = getIP(r)
	var mac, err = getMacByIp(ip)

	log.Println(ip)

	var result WhoamiResult
	if err != nil {
		result = WhoamiResult{
			Success:      false,
			ErrorMessage: err.Error(), // bad practice long term, nice for debugging
		}
	} else {
		result = WhoamiResult{
			Success:      true,
			ErrorMessage: "",
			Ip:           ip,
			Mac:          mac,
		}
	}

	var responseJson, _ = json.Marshal(result)

	fmt.Fprintf(w, string(responseJson))
}

func main() {
	var port = ":2122"
	fmt.Println("Starting Tollgate - Whoami")
	fmt.Println("Listening on port", port)

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(port, nil))

	fmt.Println("Shutting down Tollgate - Whoami")
}

func getIP(r *http.Request) string {
	// Check if the IP is set in the X-Real-Ip header
	ip := r.Header.Get("X-Real-Ip")
	if ip != "" {
		return ip
	}

	// Check if the IP is set in the X-Forwarded-For header
	ips := r.Header.Get("X-Forwarded-For")
	if ips != "" {
		return strings.Split(ips, ",")[0]
	}

	// Fallback to the remote address, removing the port
	ip = r.RemoteAddr
	if colon := strings.LastIndex(ip, ":"); colon != -1 {
		ip = ip[:colon]
	}

	return ip
}
