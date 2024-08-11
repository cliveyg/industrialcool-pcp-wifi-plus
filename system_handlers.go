package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func (a *App) systemAction(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	sa := vars["action"]

	//TODO: Check input string more thoroughly

	var err error
	pr := WifiPlusResponse{
		Function: "sysAction",
		Action:   sa,
	}

	switch sa {
	case "config":
		if r.Method == http.MethodGet || r.Method == http.MethodPut {
			a.sysPCPConfig(&pr, r.Method, &err, nil)
		} else {
			pr.StatusCode = 400
			pr.Message = "Incorrect HTTP method for action"
		}
	case "picore":
		if r.Method == http.MethodGet {
			a.sysPiCoreDetails(&pr, &err)
		} else {
			pr.StatusCode = 400
			pr.Message = "Incorrect HTTP method for action"
		}
	case "reboot":
		a.sysReboot(w, &pr)
		return
	case "shutdown":
		a.sysShutdown(w, &pr)
		return
	case "status":
		if r.Method == http.MethodGet {
			a.sysStatus(&pr, &err)
		} else {
			pr.StatusCode = 400
			pr.Message = "Incorrect HTTP method for action"
		}
	default:
		// do nowt
		pr.StatusCode = 400
		pr.Message = "Action does not exist"
	}
	pr.ReturnResponse(w, err)
}

func (a *App) sysShutdown(w http.ResponseWriter, pr *WifiPlusResponse) {
	pr.StatusCode = 202
	pr.Message = "System shutting down"
	pr.Cmd = "sudo pcp sd"
	pr.ReturnResponse(w, nil)
	time.Sleep(2 * time.Second)
	rc, err := exec.Command("sh", "-c", "sudo pcp sd").Output()
	log.Debug(rc)
	if err != nil {
		pr.ReturnResponse(w, err)
	}
}

func (a *App) sysReboot(w http.ResponseWriter, pr *WifiPlusResponse) {
	pr.StatusCode = 202
	pr.Message = "System rebooting"
	pr.Cmd = "sudo pcp rb"
	pr.ReturnResponse(w, nil)
	time.Sleep(2 * time.Second)
	rc, err := exec.Command("sh", "-c", "sudo pcp rb").Output()
	log.Debug(rc)
	if err != nil {
		pr.ReturnResponse(w, err)
	}
}

func (a *App) sysStatus(pr *WifiPlusResponse, err *error) {

	var rcInt int
	var rc []byte
	pr.Cmd = "wifi-plus.sh wp_status 200"
	rc, *err = exec.Command("sh", "-c", "cd cgi-bin && ./wifi-plus.sh wp_status 200").Output()
	if *err != nil {
		return
	}
	rcInt, *err = strconv.Atoi(strings.TrimSpace(string(rc)))
	if *err != nil {
		return
	}
	pr.StatusCode = rcInt
	pr.Message = "System running"

}

func (a *App) sysPCPConfig(pr *WifiPlusResponse, hm string, err *error, sr *string) {

	var r []byte
	if hm == http.MethodGet {
		pr.Message = "Fetch pcp config settings"
		pr.Cmd = "./wifi-plus.sh wp_pcp_config"
		r, *err = exec.Command("sh", "-c", "cd cgi-bin; ./wifi-plus.sh wp_pcp_config read").Output()
		if *err != nil {
			return
		}
		pr.StatusCode = 200
		*sr = string(r)
		pr.Data = textToMap(string(r))
	} else if hm == http.MethodPut {
		log.Debug("Editing not implemented yet")
		pr.StatusCode = 501
		pr.Message = "Not implemented yet"
	}

}

func (a *App) sysPiCoreDetails(pr *WifiPlusResponse, err *error) {

	var rc []byte
	pr.Cmd = "wifi-plus.sh wp_picore_details"
	rc, *err = exec.Command("sh", "-c", "cd cgi-bin && sudo ./wifi-plus.sh wp_picore_details").Output()
	if *err != nil {
		return
	}

	pr.StatusCode = 200
	pr.Message = "piCore system details"
	picoreData := PiCoreSystemData{}

	*err = json.Unmarshal(rc, &picoreData)
	if *err != nil {
		return
	}
	pr.Data = picoreData

}
