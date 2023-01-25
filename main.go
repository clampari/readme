package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

func conexionBD() (conexion *sql.DB) {
	Driver := "mysql"
	Usuario := "root"
	Pass := ""
	Base := "trabajo"

	conexion, err := sql.Open(Driver, Usuario+":"+Pass+"@tcp(127.0.0.1)/"+Base)

	if err != nil {
		panic(err.Error())
	}
	return conexion
}

var plantillas = template.Must(template.ParseGlob("platillas/*"))

func main() {
	http.HandleFunc("/", inicio)
	http.HandleFunc("/crear", crear)
	http.HandleFunc("/insertar", insertar)

	http.HandleFunc("/editar", Editar)
	http.HandleFunc("/actualizar", Actualizar)

	http.HandleFunc("/borrar", Borrar)

	log.Println("Servidor corriendo")
	http.ListenAndServe(":8080", nil)
}

func Borrar(w http.ResponseWriter, r *http.Request) {
	idproduc := r.URL.Query().Get("id")
	//fmt.Println(idproduc)
	conexionEstablecida := conexionBD()
	borrarRegistro, err := conexionEstablecida.Prepare("DELETE FROM products WHERE sku=?")

	if err != nil {
		panic(err.Error())
	}
	borrarRegistro.Exec(idproduc)

	http.Redirect(w, r, "/", 301)
}

type Productos struct {
	Id      string
	Nombre  string
	Marca   string
	Tamanio string
	Precio  string
	Imagen  string
}

func inicio(w http.ResponseWriter, r *http.Request) {

	conexionEstablecida := conexionBD()
	registros, err := conexionEstablecida.Query("SELECT * FROM products")

	if err != nil {
		panic(err.Error())
	}

	produc := Productos{}
	arregloProdu := []Productos{}

	for registros.Next() {
		var id, tamanio, precio, nombre, marca, iman string
		err = registros.Scan(&id, &nombre, &marca, &tamanio, &precio, &iman)

		if err != nil {
			panic(err.Error())
		}

		produc.Id = id
		produc.Nombre = nombre
		produc.Marca = marca
		produc.Tamanio = tamanio
		produc.Precio = precio
		produc.Imagen = iman

		arregloProdu = append(arregloProdu, produc)

	}

	//fmt.Println(arregloProdu)

	plantillas.ExecuteTemplate(w, "inicio", arregloProdu)

}

func Editar(w http.ResponseWriter, r *http.Request) {
	idproduc := r.URL.Query().Get("id")
	fmt.Println(idproduc)

	conexionEstablecida := conexionBD()
	registro, err := conexionEstablecida.Query("SELECT * FROM products WHERE sku=?", idproduc)

	produc := Productos{}
	for registro.Next() {
		var id, tamanio, precio, nombre, marca, iman string
		err = registro.Scan(&id, &nombre, &marca, &tamanio, &precio, &iman)

		if err != nil {
			panic(err.Error())
		}

		produc.Id = id
		produc.Nombre = nombre
		produc.Marca = marca
		produc.Tamanio = tamanio
		produc.Precio = precio
		produc.Imagen = iman
	}

	fmt.Println(produc)
	plantillas.ExecuteTemplate(w, "editar", produc)

}

func crear(w http.ResponseWriter, r *http.Request) {
	plantillas.ExecuteTemplate(w, "crear", nil)
}

func insertar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("sku")
		nombre := r.FormValue("name")
		marca := r.FormValue("brand")
		tam := r.FormValue("size")
		pres := r.FormValue("price")
		imag := r.FormValue("img")

		conexionEstablecida := conexionBD()
		insertarRegistros, err := conexionEstablecida.Prepare("INSERT INTO products (sku, Name, Brand, Size, Price, Img) VALUES(?, ?, ?, ?, ?, ?)")

		if err != nil {
			panic(err.Error())
		}
		insertarRegistros.Exec(id, nombre, marca, tam, pres, imag)

		http.Redirect(w, r, "/", 301)
	}
}

func Actualizar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("sku")
		nombre := r.FormValue("name")
		marca := r.FormValue("brand")
		tam := r.FormValue("size")
		pres := r.FormValue("price")
		imag := r.FormValue("img")

		conexionEstablecida := conexionBD()
		ActualizarRegistros, err := conexionEstablecida.Prepare("UPDATE products SET Name=?, Brand=?, Size=?, Price=?, Img=? WHERE sku=?")

		if err != nil {
			panic(err.Error())
		}
		ActualizarRegistros.Exec(nombre, marca, tam, pres, imag, id)

		http.Redirect(w, r, "/", 301)
	}
}
