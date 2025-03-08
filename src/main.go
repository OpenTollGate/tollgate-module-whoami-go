package main

import (
	"encoding/json"
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

var tollgateMerchantPubkey string = "c1f4c025e746fd307203ac3d1a1886e343bea76ceec5e286c96fb353be6cadea"

func getMacAddress(ipAddress string) (string, error) {
	cmdIn := `cat /tmp/dhcp.leases | cut -f 2,3,4 -s -d" " | grep -i ` + ipAddress + ` | cut -f 1 -s -d" "`
	commandOutput, err := exec.Command("sh", "-c", cmdIn).Output()

	var commandOutputString = string(commandOutput)
	if err != nil {
		fmt.Println(err, "Error when getting client's mac address. Command output: "+commandOutputString)
		return "nil", err
	}

	return strings.Trim(commandOutputString, "\n"), nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	var ip = getIP(r)
	var mac, err = getMacAddress(ip)

	log.Println("ip", ip, "mac", mac)

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

func handlePubkey(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, tollgateMerchantPubkey)
}

func main() {
	var port = ":2122"
	fmt.Println("Starting Tollgate - Whoami")
	fmt.Println("Listening on port", port)

	http.HandleFunc("/", handler)
	http.HandleFunc("/pubkey", handlePubkey)
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
