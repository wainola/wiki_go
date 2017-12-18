package main

import (
	"log"
	"fmt"
	"io/ioutil"
	"net/http"
	"html/template"
)

type Page struct {
	Title string
	Body []byte
}

// metodo que guarda y opera sobre Page.
func (p *Page) save() error {
	// Generamos un archivo con el nombre del campo inicializado en el struct.
	// luego escribimos el archivo en el sistema. Nombre archivo => filename, contenido, Body, permisos para r y w
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)// el valor octal del ultimo parametro indica permisos r y w
}

// Funcion para cargar paginas
/*
Esta funcion esta haciendo la pega de los datos a obtener, puesto que tenemos que generar los archivos para poder ver la informacion cargada, posterior al retorno del struct
*/
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil // retornamos un puntero a una pagina en literal y nil porque en este linea no deberia haber error
}

func viewHandler(w http.ResponseWriter, r *http.Request){
	title := r.URL.Path[len("/view/"):]
	log.Println(title)
	p, err := loadPage(title) // como aca se retorna un &Page, esto le pasamos despues a fmt
	// si el error no es nulo, entonces que redireccione para poder crear el contenido
	if err != nil {
		http.Redirect(w, r, "/edit/" + title, http.StatusFound)
	}
	renderTemplate(w, "view.html", p)
}

// editHandler usa la plantilla html que hemos creado
func editHandler(w http.ResponseWriter, r *http.Request){
	title := r.URL.Path[len("/edit/"):] // obtenemos todo lo que esta despues de /edit/, ej: /edit/hola, return hola
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit.html", p)

}

func saveHandler(w http.ResponseWriter, r *http.Request){
	title := r.URL.Path[len("/save/"):] // lo mismo obtenemos lo que viene despues de save/...
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

// generando funcion generica para la pega del renderizado y ejecuccion de las plantillas
func renderTemplate (w http.ResponseWriter, tmpl string, p *Page){
	// Lee los contenidos de edit.html y retorna un *template.Template. 
	t, _ := template.ParseFiles(tmpl)
	// Este metodo ejecuta el template, escribiendo el html generado al http.ResponseWriter
	t.Execute(w,p)
}


func main(){
	p1 := &Page{Title: "TestPage", Body: []byte("This is a sample page.")}
	p1.save()
	p2, _ := loadPage("TestPage")
	fmt.Println(string(p2.Body))

	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.ListenAndServe(":8080", nil)
}