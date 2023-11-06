package adminpage

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"text/template"
)

//페이지 구조체
type Page struct {
	Title string
}

//html 구획 나누기
var templ = func() *template.Template {

	//template 에 function 넣기
	funcadd := template.FuncMap{
		"add": func(a int, b int) int {
			return a + b
		},
	}

	// t := template.New("")
	t := template.New("").Funcs(funcadd)
	err := filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".html") {
			fmt.Println(path)
			_, err = t.ParseFiles(path)
			if err != nil {
				fmt.Println(err)
			}
		}
		return err
	})

	if err != nil {
		fmt.Print(err)
	}
	return t
}()

func Login_Handler(w http.ResponseWriter, r *http.Request) {
	templ.ExecuteTemplate(w, "base_header", &Page{Title: "LoginPage"})
}
