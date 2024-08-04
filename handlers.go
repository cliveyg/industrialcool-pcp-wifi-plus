package main

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"os/exec"
	"strings"
)

// ----------------------------------------------------------------------------

func (a *App) testTings(w http.ResponseWriter, _ *http.Request) {

	log.Debug("In testTings")
	rc, err := exec.Command("sh", "-c", "cd cgi-bin && sudo ./wifi-plus.sh wp_test").Output()
	pr := WifiPlusResponse{
		Cmd:        "testTings",
		StatusCode: 200,
		Message:    strings.TrimSpace(string(rc))}
	pr.FormatResponse(w, err)

}

func (a *App) getPiCoreDetails(w http.ResponseWriter, _ *http.Request) {

	log.Debug("In getPiCoreDetails")
	retData, err := exec.Command("sh", "-c", "cd cgi-bin && sudo ./wifi-plus.sh wp_picore_details").Output()
	pr := WifiPlusResponse{
		Cmd:        "getPiCoreDetails",
		StatusCode: 200,
		Message:    "piCore details",
		Data:       string(retData)}
	pr.FormatResponse(w, err)

}

func (a *App) getSystemStatus(w http.ResponseWriter, _ *http.Request) {

	log.Debug("In getSystemStatus")
	rc, err := exec.Command("sh", "-c", "cd cgi-bin && sudo ./wifi-plus.sh wp_status 200").Output()

	pr := WifiPlusResponse{
		Cmd:        "getSystemStatus",
		StatusCode: 200,
		Message:    strings.TrimSpace(string(rc))}
	pr.FormatResponse(w, err)

}

func (a *App) getWifiStatus(w http.ResponseWriter, _ *http.Request) {

	log.Debug("In getWifiStatus")
	args := []string{"wlan0", "status"}
	rc, err := a.ExecCmd("/usr/local/etc/init.d/wifi", args)
	pr := WifiPlusResponse{Cmd: "getWifiStatus"}

	if strings.Contains(rc, "wpa_supplicant running") {
		pr.Message = "wpa_supplicant running"
		pr.StatusCode = 200
	} else {
		pr.Message = "wpa_supplicant not running"
		pr.StatusCode = 404
	}
	pr.FormatResponse(w, err)
}

func (a *App) getWifiSSID(w http.ResponseWriter, _ *http.Request) {

	log.Debug("In getWifiSSID")
	args := []string{"-r"}
	SSID, err := a.ExecCmd("iwgetid", args)
	pr := WifiPlusResponse{Cmd: "getWifiSSID"}

	if SSID == "" {
		pr.StatusCode = 404
		pr.Message = "No SSID found"
	} else {
		pr.StatusCode = 200
		pr.Message = "SSID found"
		pr.Data = `"SSID": "` + SSID + `"`
	}
	pr.FormatResponse(w, err)
}
