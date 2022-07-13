package main

import (
	"bufio"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
)

type Sub struct {
	Username string
	Ascii    string
	Output   []byte
}

var t *template.Template

func main() {
	t, _ = template.ParseGlob("templates/*.html")

	http.HandleFunc("/", renderTemplate)
	http.HandleFunc("/ascii-art", processposthandler)
	http.ListenAndServe(":8080", nil)
}

func renderTemplate(w http.ResponseWriter, r *http.Request) {
	var s Sub
	t, err := template.ParseFiles("web.html")
	if err != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return

	}
	if r.URL.Path != "/" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	err = t.Execute(w, s)
	if err != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return

	}
}

func processposthandler(w http.ResponseWriter, r *http.Request) {
	var s Sub

	s.Username = r.FormValue("name")

	s.Ascii = r.FormValue("Ascii")

	if r.FormValue("name") == "" {
		http.Error(w, "400 bad request\nrequested character(s) are not in the valid range", http.StatusBadRequest)
		return
	}

	switch s.Ascii {
	case "standard":
	case "shadow":
	case "thinkertoy":
	default:
		s.Ascii = "standard"
	}

	r.ParseForm()
	var style string = s.Ascii

	if style != "standard" && style != "shadow" && style != "thinkertoy" {

		http.Error(w, "404 the banner file has not been found or you did not send any data\ntry to start from localhost:8080", http.StatusNotFound)
		return
	}

	data := Banner(style)
	output := s.Username
	bannerMap := Array(data)

	ascii_Output := Print(output, bannerMap)
	defer data.Close()

	s.Output = []byte(ascii_Output)

	t, err := template.ParseFiles("web.html")
	if err != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return

	}

	err = t.Execute(w, s)
	if err != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}
}

func Banner(style string) *os.File {
	b, err := os.Open(style + ".txt")
	if err != nil {
		panic("404 the banner file has not been found or you did not send any data\ntry to start from localhost:8080")
	}

	return b
}

func Array(file *os.File) map[rune][]string {
	bannerMap := make(map[rune][]string)
	var lines []string
	scanner := bufio.NewScanner(io.Reader(file))
	firstLine := false
	char := ' '
	for scanner.Scan() {
		if scanner.Text() == "" && firstLine {
			bannerMap[char] = lines
			lines = nil
			char++
		} else if firstLine {
			lines = append(lines, scanner.Text())
		}
		firstLine = true
	}
	bannerMap[char] = lines
	return bannerMap
}

func Print(str string, banner map[rune][]string) string {
	list := Split(str)
	ascii_output := ""
	Dec := 10

	for _, word := range list {
		if word == "" {
			ascii_output = ascii_output + string(rune(Dec))
		} else {
			for i := 0; i < 8; i++ {
				line := ""
				for _, r := range word {
					line = line + banner[r][i]
				}
				ascii_output = ascii_output + line + string(rune(Dec))
			}
		}
	}

	return ascii_output
}

func Split(str string) []string {
	answer := ""

	b := strings.Split(str, "\r\n")

	for _, j := range b {
		if j == "" {
			answer += "\n"
		}
	}
	return b
}
