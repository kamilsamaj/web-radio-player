package main

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
)

//go:embed templates/*
var templates embed.FS

var tpl *template.Template

var cmd *exec.Cmd

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := stopIfPlaying()
		if err != nil {
			log.Fatalln(err)
		}
		err = tpl.ExecuteTemplate(w, "home.gohtml", nil)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		log.Printf("unsupported method %s on %s", r.Method, r.URL)
	}
}

func rockRadioHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := playStream("https://ice.abradio.cz/rockradiomorava64.aac")
		if err != nil {
			log.Fatalln(err)
		}
		err = showRadioPage(w)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		log.Printf("unsupported method %s on %s", r.Method, r.URL)
	}
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func stopHandler(w http.ResponseWriter, r *http.Request) {
	var outb, errb bytes.Buffer
	log.Println("running pkill mpv")
	cmd = exec.Command("pkill", "mpv")
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	fmt.Fprintln(w, "killed all mpv processes")
	fmt.Fprintln(w, outb.String(), errb.String())
}

func showRadioPage(w http.ResponseWriter) error {
	// return HTML page
	err := tpl.ExecuteTemplate(w, "radio.gohtml", nil)
	if err != nil {
		return err
	}
	return nil
}

func playStream(streamUrl string) error {
	cmd = exec.Command("mpv", streamUrl, "--no-terminal")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Start()
	if err != nil {
		return err
	}
	log.Printf("running mpv in subprocess %d", cmd.Process.Pid)
	return nil
}

func stopIfPlaying() error {
	if cmd != nil {
		log.Printf("stopping process %d", cmd.Process.Pid)
		err := cmd.Process.Kill()
		if err != nil {
			return err
		}
		_ = cmd.Wait()
		log.Printf("process %d successfully stopped", cmd.Process.Pid)
		cmd = nil
	}
	return nil
}

func init() {
	var err error
	tpl, err = template.ParseFS(templates, "templates/*.gohtml")
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.HandleFunc("/rock-radio", rockRadioHandler)
	http.HandleFunc("/stop", stopHandler)
	fmt.Println("listening on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln(err)
	}
}
