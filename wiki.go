package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"html/template"
	"regexp"
	"errors"
)
// variable global que parsea los templates que hasta el momento tenemos.
// si recibe valores nulos (aka: el template no existe) entonces le hace exit al programa.
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

// variable global para validar
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

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

func viewHandler(w http.ResponseWriter, r *http.Request, title string){
	if err != nil {
		return
	}
	p, err := loadPage(title) // como aca se retorna un &Page, esto le pasamos despues a fmt
	// si el error no es nulo, entonces que redireccione para poder crear el contenido
	if err != nil {
		http.Redirect(w, r, "/edit/" + title, http.StatusFound)
	}
	renderTemplate(w, "view.html", p)
}

// editHandler usa la plantilla html que hemos creado
func editHandler(w http.ResponseWriter, r *http.Request, title string){
	//title := r.URL.Path[len("/edit/"):] // obtenemos todo lo que esta despues de /edit/, ej: /edit/hola, return hola
	if err != nil{
		return
	}
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit.html", p)

}

func saveHandler(w http.ResponseWriter, r *http.Request, title string){
	//title := r.URL.Path[len("/save/"):] // lo mismo obtenemos lo que viene despues de save/...
	if err != nil {
		return
	}
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	error := p.save()
	if error != nil {
		http.Error(w, error.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

// generando funcion generica para la pega del renderizado y ejecuccion de las plantillas
func renderTemplate (w http.ResponseWriter, tmpl string, p *Page){
	// Este metodo ejecuta el template, escribiendo el html generado al http.ResponseWriter
	// aca en el error ejecutamos el template
	// obtenemos la variable global y ejecutamos segun corresponda en el template
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// usando function literals en go
func makeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// aca extraemos el titulo de la pagina desde el request
		// y llamamos al handler fn
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		//return m[2], nil // titulo ingresado en este caso es la segunda subexpresion
		fn(w, r, m[2])
		}
}

func main(){
	p1 := &Page{Title: "TestPage", Body: []byte("This is a sample page.")}
	p1.save()
	p2, _ := loadPage("TestPage")
	fmt.Println(string(p2.Body))

	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.ListenAndServe(":8080", nil)
}