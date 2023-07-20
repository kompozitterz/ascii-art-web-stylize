package main

import (
	"bufio"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

type Er struct {
	Status int
	text   string
}

func ErrorHandler(w http.ResponseWriter, i int) {
	w.WriteHeader(i)
	tmpl, err := template.ParseFiles("/templates/*.html")
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	er := Er{i, http.StatusText(i)}
	err = tmpl.Execute(w, er)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
}

func posthandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if r.URL.Path != "/" {
			// w.WriteHeader(http.StatusNotFound)
			http.Error(w, "ERROR-404\nPage not found", http.StatusNotFound)
			return
		}
		templates.ExecuteTemplate(w, "index.html", nil)
	}

	if r.Method == "POST" {
		text1 := r.FormValue("text")
		font := r.FormValue("fonts")
		text := ""
		// vérifiez si appuyez sur Entrée pour une nouvelle ligne
		if strings.Contains(text1, "\r\n") {
			text = strings.ReplaceAll(text1, "\r\n", "\\n")
		} else {
			text = text1
		}
		// Vérifier si l'utilisateur tape dans l'art ascii approprié
		for _, v := range text {
			if !(v >= 32 && v <= 126) {
				// Si l'utilisateur ne rentre pas un caractère ASCII
				http.Error(w, "ERROR-400\nBad Request!! \nAssurez-vous d'entrer un texte ou de sélectionner une police.", http.StatusBadRequest)
				return
			}
		}

		// Gère erreur sever
		file, err := os.Open(font + ".txt")
		if err != nil {
			// Erreur Server obsolète parce que notre Server est hébergé en local
			http.Error(w, "ERROR-500\nInternal Server Error!", http.StatusInternalServerError)
			fmt.Fprintf(w, "Status 500: Internal Server Error")
			return
		}

		defer file.Close()
		// Lis le fichier
		Scanner := bufio.NewScanner(file)

		// Identifie les lettres avec le ASCII code
		var lines []string
		for Scanner.Scan() {
			lines = append(lines, Scanner.Text())
		}
		asciiChrs := make(map[int][]string)
		dec := 31

		for _, line := range lines {
			if line == "" {
				dec++
			} else {
				asciiChrs[dec] = append(asciiChrs[dec], line)
			}
		}
		c := ""
		for i := 0; i < len(text); i++ {
			if text[i] == 92 && text[i+1] == 110 {
				c = PrintArt(text[:i], asciiChrs) + PrintArt(text[i+2:], asciiChrs)
			}
		}
		if !strings.Contains(text, "\\n") {
			c = PrintArt(text, asciiChrs)
		}

		// pin := os.WriteFile("download.doc", []byte(c), 0o644)
		// if pin != nil {
		// 	panic(pin)
		// }
		// pin1 := os.WriteFile("download.txt", []byte(c), 0o644)
		// if pin1 != nil {
		// 	panic(pin1)
		// }

		templates.ExecuteTemplate(w, "index.html", c)
	}
}

// fonction qui va générer notre ASCII art
func PrintArt(n string, y map[int][]string) string {
	// Imprime horizontalement
	a := []string{}
	// Imprime horizontalement
	for j := 0; j < len(y[32]); j++ {
		for _, letter := range n {
			a = append(a, y[int(letter)][j])
		}
		a = append(a, "\n")
	}
	b := strings.Join(a, "")
	// fmt.Println("Votre résultat viens d'être téléchargé dans vos fichiers .doc et .txt respectifs")
	return b
}

// func download(w http.ResponseWriter, r *http.Request) {
// 	formatType := r.FormValue("fileformat")

// 	f, _ := os.Open("download." + formatType)
// 	if f == nil {
// 		http.Error(w, "ERROR-500\nInternal Server Error!", http.StatusInternalServerError)
// 		fmt.Fprintf(w, "Status 500: Internal Server Error")
// 		return
// 	}
// 	defer f.Close()

// 	file, _ := f.Stat()
// 	fsize := file.Size()

// 	sfSize := strconv.Itoa(int(fsize))

// 	w.Header().Set("Content-Disposition", "attachment; filename=asciiresults."+formatType)
// 	w.Header().Set("Content-Type", "text/html")
// 	w.Header().Set("Content-Length", sfSize)

// 	io.Copy(w, f)
// }
