package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

func conexionBD() (conexion *sql.DB) {
	Driver := "mysql"
	DBUsuario := "root"
	DBSenha := ""
	DBBanco := "catalogo"

	conn, err := sql.Open(Driver, DBUsuario+":"+DBSenha+"@tcp(127.0.0.1)/"+DBBanco)
	if err != nil {
		panic(err)
	}
	return conn
}

var modelo = template.Must(template.ParseGlob("modelo/*"))

func main() {
	http.HandleFunc("/", Inicio)
	http.HandleFunc("/buscar", Buscar)
	http.HandleFunc("/relatorio", Relatorio)
	http.HandleFunc("/criar", Criar)
	http.HandleFunc("/inserir", Inserir)
	http.HandleFunc("/apagar", Apagar)
	http.HandleFunc("/editar", Editar)
	http.HandleFunc("/atualizar", Atualizar)
	log.Println("Servidor Rodando")
	http.ListenAndServe(":3000", nil)
}

type Certificado struct {
	Id             int
	Cliente        string
	Url            string
	Telefone       string
	DataEmissao    string
	DataVencimento string
}

func Apagar(w http.ResponseWriter, r *http.Request) {
	idRegistro := r.URL.Query().Get("id")
	connEstabelecida := conexionBD()
	apagarRegistros, err := connEstabelecida.Prepare("DELETE FROM certificado WHERE id=?")

	if err != nil {
		panic(err.Error())
	}
	apagarRegistros.Exec(idRegistro)
	http.Redirect(w, r, "/", 301)
}

func Inicio(w http.ResponseWriter, r *http.Request) {
	connEstabelecida := conexionBD()
	registros, err := connEstabelecida.Query("SELECT * FROM certificado")

	if err != nil {
		panic(err.Error())
	}

	certificado := Certificado{}
	arrayCertificado := []Certificado{}

	for registros.Next() {
		var id int
		var cliente, url, telefone, dataemissao, datavencimento string
		err = registros.Scan(&id, &cliente, &url, &telefone, &dataemissao, &datavencimento)

		if err != nil {
			panic(err.Error())
		}
		certificado.Id = id
		certificado.Cliente = cliente
		certificado.Url = url
		certificado.Telefone = telefone
		certificado.DataEmissao = dataemissao
		certificado.DataVencimento = datavencimento

		arrayCertificado = append(arrayCertificado, certificado)
	}

	modelo.ExecuteTemplate(w, "inicio", arrayCertificado)
}

func Editar(w http.ResponseWriter, r *http.Request) {
	idRegistro := r.URL.Query().Get("id")
	fmt.Println(idRegistro)

	connEstabelecida := conexionBD()
	registro, err := connEstabelecida.Query("SELECT * FROM certificado WHERE ID=?", idRegistro)
	if err != nil {
		panic(err.Error())
	}
	certificado := Certificado{}
	for registro.Next() {
		var id int
		var cliente, url, telefone, dataemissao, datavencimento string
		err = registro.Scan(&id, &cliente, &url, &telefone, &dataemissao, &datavencimento)

		if err != nil {
			panic(err.Error())
		}
		certificado.Id = id
		certificado.Cliente = cliente
		certificado.Url = url
		certificado.Telefone = telefone
		certificado.DataEmissao = dataemissao
		certificado.DataVencimento = datavencimento
	}

	modelo.ExecuteTemplate(w, "editar", certificado)
}

func Criar(w http.ResponseWriter, r *http.Request) {
	modelo.ExecuteTemplate(w, "criar", nil)
}

func Inserir(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("uploadArquivo")
	if err != nil {
		fmt.Println("Erro ao recuperar o arquivo")
		fmt.Println(err)
		return
	}

	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	dst, err := os.Create(handler.Filename)
	defer dst.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == "POST" {
		cliente := r.FormValue("cliente")
		telefone := r.FormValue("telefone")
		uploadArquivo := (handler.Filename)
		dataVencimento := r.FormValue("datavencimento")

		connEstabelecida := conexionBD()
		inserirRegistros, err := connEstabelecida.Prepare("INSERT INTO certificado(cliente, telefone, url, data_vencimento) VALUES(?, ?, ?, ?)")

		if err != nil {
			panic(err.Error())
		}
		inserirRegistros.Exec(cliente, telefone, uploadArquivo, dataVencimento)
		http.Redirect(w, r, "/", 301)
	}

}

func Atualizar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		cliente := r.FormValue("cliente")
		telefone := r.FormValue("telefone")

		connEstabelecida := conexionBD()
		atualizarRegistros, err := connEstabelecida.Prepare("UPDATE certificado SET cliente=?,telefone=? WHERE id=?")

		if err != nil {
			panic(err.Error())
		}

		atualizarRegistros.Exec(cliente, telefone, id)
		http.Redirect(w, r, "/", 301)
	}
}

func Buscar(w http.ResponseWriter, r *http.Request) {
	modelo.ExecuteTemplate(w, "buscar", nil)
}

func Relatorio(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		idBusca := r.FormValue("id")

		connEstabelecida := conexionBD()
		buscaRegistros, err := connEstabelecida.Query("SELECT * FROM certificado WHERE data_vencimento=?", idBusca)

		if err != nil {
			panic(err.Error())
		}

		certificado := Certificado{}
		arrayBusca := []Certificado{}
		for buscaRegistros.Next() {
			var id int
			var cliente, url, telefone, dataemissao, datavencimento string
			err = buscaRegistros.Scan(&id, &cliente, &url, &telefone, &dataemissao, &datavencimento)

			if err != nil {
				panic(err.Error())
			}
			certificado.Id = id
			certificado.Cliente = cliente
			certificado.Url = url
			certificado.Telefone = telefone
			certificado.DataEmissao = dataemissao
			certificado.DataVencimento = datavencimento

			arrayBusca = append(arrayBusca, certificado)
		}

		modelo.ExecuteTemplate(w, "relatorio", arrayBusca)
	}

}
