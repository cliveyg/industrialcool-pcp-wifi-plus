package main

import (
	"io"
	"log"
	"net/http"
	"strings"
)

// ----------------------------------------------------------------------------

func (a *App) getSystemStatus(w http.ResponseWriter, _ *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	args := []string{"wp_status", "200"}
	rc, err := a.ExecCmd("sudo cgi-bin/wifi-plus.sh", args)
	if err != nil {
		mess := `{"error": "` + err.Error() + `"}`
		if _, err := io.WriteString(w, mess); err != nil {
			log.Fatal(err)
		}
		return
	}

	mess := `{"message": "System running...", "return_code": "` + rc + `"}`
	if _, err := io.WriteString(w, mess); err != nil {
		log.Fatal(err)
	}

}

func (a *App) getWifiStatus(w http.ResponseWriter, _ *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var message string
	args := []string{"wlan0", "status"}
	retString, err := a.ExecCmd("/usr/local/etc/init.d/wifi", args)
	if err != nil {
		mess := `{"error": "` + err.Error() + `"}`
		if _, err := io.WriteString(w, mess); err != nil {
			log.Fatal(err)
		}
		return
	}
	if strings.Contains(retString, "wpa_supplicant running") {
		message = `{"command": "wifi status", "message": "wpa_supplicant running" }`
	} else {
		message = `{"command": "wifi status", "message": "wpa_supplicant not running"}`
	}
	if _, err := io.WriteString(w, message); err != nil {
		log.Fatal(err)
	}
}

func (a *App) getWifiSSID(w http.ResponseWriter, _ *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var message string
	args := []string{"-r"}
	SSID, err := a.ExecCmd("iwgetid", args)
	if err != nil {
		mess := `{"error": "` + err.Error() + `"}`
		if _, err := io.WriteString(w, mess); err != nil {
			log.Fatal(err)
		}
		return
	}
	if SSID == "" {
		message = `{, "message": "No SSID found" }`
	} else {
		message = `{ "SSID": "` + SSID + `" }`
	}
	if _, err := io.WriteString(w, message); err != nil {
		log.Fatal(err)
	}
}
