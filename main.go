package main

import (
	"database/sql"
	"fmt"

	//"io/ioutil"
	"log"
	"net/http"

	//"os"
	"text/template"

	"io/ioutil"

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
	http.HandleFunc("/buscarCliente", BuscarCliente)
	http.HandleFunc("/relatorio", Relatorio)
	http.HandleFunc("/relatorioCliente", RelatorioCliente)
	http.HandleFunc("/criar", Criar)
	http.HandleFunc("/inserir", Inserir)
	http.HandleFunc("/apagar", Apagar)
	http.HandleFunc("/editar", Editar)
	http.HandleFunc("/atualizar", Atualizar)
	log.Println("Servidor Rodando")
	http.ListenAndServe(":3010", nil)
}

type Certificado struct {
	Id             int
	Cliente        string
	Doc            string
	Url            string
	CertPass       string
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
		var cliente, doc, url, certPass, telefone, dataemissao, datavencimento string
		err = registros.Scan(&id, &cliente, &doc, &url, &certPass, &telefone, &dataemissao, &datavencimento)

		if err != nil {
			panic(err.Error())
		}
		certificado.Id = id
		certificado.Cliente = cliente
		certificado.Doc = doc
		certificado.Url = url
		certificado.CertPass = certPass
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
		var cliente, doc, url, CertPass, telefone, dataemissao, datavencimento string
		err = registro.Scan(&id, &cliente, &doc, &url, &CertPass, &telefone, &dataemissao, &datavencimento)

		if err != nil {
			panic(err.Error())
		}
		certificado.Id = id
		certificado.Cliente = cliente
		certificado.Doc = doc
		certificado.Url = url
		certificado.CertPass = CertPass
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

	//aqui
	tempFile, err := ioutil.TempFile("certificado", "upload-*.pfx")
	if err != nil {
		fmt.Println(err)
	}
	nome := tempFile.Name()

	defer tempFile.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	tempFile.Write(fileBytes)
	//ate aqui
	if r.Method == "POST" {
		cliente := r.FormValue("cliente")
		doc := r.FormValue("doc")
		uploadArquivo := nome
		certPass := r.FormValue("certPass")
		telefone := r.FormValue("telefone")
		dataEmissao := r.FormValue("dataEmissao")
		dataVencimento := r.FormValue("dataVencimento")

		connEstabelecida := conexionBD()
		inserirRegistros, err := connEstabelecida.Prepare("INSERT INTO certificado(cliente, doc, url, certPass, telefone, data_emissao, data_vencimento) VALUES(?, ?, ?, ?, ?, ?, ?)")

		if err != nil {
			panic(err.Error())
		}
		inserirRegistros.Exec(cliente, doc, uploadArquivo, certPass, telefone, dataEmissao, dataVencimento)
		http.Redirect(w, r, "/", 301)
	}

}

func Atualizar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		cliente := r.FormValue("cliente")
		doc := r.FormValue("doc")
		//url := r.FormValue("url")
		certPass := r.FormValue("certPass")
		telefone := r.FormValue("telefone")
		data_emissao := r.FormValue("dataEmissao")
		data_vencimento := r.FormValue("dataVencimento")

		connEstabelecida := conexionBD()
		atualizarRegistros, err := connEstabelecida.Prepare("UPDATE certificado SET cliente=?, doc=?, certPass=?, telefone=?, data_emissao=?,  data_vencimento=? WHERE id=?")

		if err != nil {
			panic(err.Error())
		}

		atualizarRegistros.Exec(cliente, doc, certPass, telefone, data_emissao, data_vencimento, id)
		http.Redirect(w, r, "/", 301)
	}
}

func Buscar(w http.ResponseWriter, r *http.Request) {
	modelo.ExecuteTemplate(w, "buscar", nil)
}

func BuscarCliente(w http.ResponseWriter, r *http.Request) {
	modelo.ExecuteTemplate(w, "buscarCliente", nil)
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
			var cliente, doc, url, certPass, telefone, dataemissao, datavencimento string
			err = buscaRegistros.Scan(&id, &cliente, &doc, &url, &certPass, &telefone, &dataemissao, &datavencimento)

			if err != nil {
				panic(err.Error())
			}
			certificado.Id = id
			certificado.Cliente = cliente
			certificado.Doc = doc
			certificado.Url = url
			certificado.CertPass = certPass
			certificado.Telefone = telefone
			certificado.DataEmissao = dataemissao
			certificado.DataVencimento = datavencimento

			arrayBusca = append(arrayBusca, certificado)
		}

		modelo.ExecuteTemplate(w, "relatorio", arrayBusca)
	}

}

func RelatorioCliente(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		idBusca := r.FormValue("id")

		connEstabelecida := conexionBD()
		buscaRegistros, err := connEstabelecida.Query("SELECT * FROM certificado WHERE doc=?", idBusca)

		if err != nil {
			panic(err.Error())
		}

		certificado := Certificado{}
		arrayBusca := []Certificado{}
		for buscaRegistros.Next() {
			var id int
			var cliente, doc, url, certPass, telefone, dataemissao, datavencimento string
			err = buscaRegistros.Scan(&id, &cliente, &doc, &url, &certPass, &telefone, &dataemissao, &datavencimento)

			if err != nil {
				panic(err.Error())
			}
			certificado.Id = id
			certificado.Cliente = cliente
			certificado.Doc = doc
			certificado.Url = url
			certificado.CertPass = certPass
			certificado.Telefone = telefone
			certificado.DataEmissao = dataemissao
			certificado.DataVencimento = datavencimento

			arrayBusca = append(arrayBusca, certificado)
		}

		modelo.ExecuteTemplate(w, "relatorioCliente", arrayBusca)
	}

}
