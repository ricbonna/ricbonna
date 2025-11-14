package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Produto struct {
	CdProduto   int    `json:"cd_produto"`
	NomeProduto string `json:"nome_produto"`
}

var db *sql.DB

func main() {
	var err error

	// adjust to your DB: user:password@tcp(host:port)/dbname
	db, err = sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/test")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/produtos", produtosHandler)
	http.HandleFunc("/produto/", produtoHandler)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// ---- LIST + CREATE ----
func produtosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		rows, err := db.Query("SELECT Cd_Produto, Nome_Produto FROM Tb_produto")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()

		var produtos []Produto
		for rows.Next() {
			var p Produto
			rows.Scan(&p.CdProduto, &p.NomeProduto)
			produtos = append(produtos, p)
		}

		json.NewEncoder(w).Encode(produtos)

	case http.MethodPost:
		var p Produto
		json.NewDecoder(r.Body).Decode(&p)

		_, err := db.Exec("INSERT INTO Tb_produto (Nome_Produto) VALUES (?)", p.NomeProduto)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.WriteHeader(201)
		w.Write([]byte("Produto criado"))
	}
}

// ---- READ + UPDATE + DELETE ----
func produtoHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/produto/"):]
	id, _ := strconv.Atoi(idStr)

	switch r.Method {

	case http.MethodGet: // READ
		var p Produto
		err := db.QueryRow("SELECT Cd_Produto, Nome_Produto FROM Tb_produto WHERE Cd_Produto = ?", id).
			Scan(&p.CdProduto, &p.NomeProduto)
		if err != nil {
			http.Error(w, "Not found", 404)
			return
		}
		json.NewEncoder(w).Encode(p)

	case http.MethodPut: // UPDATE
		var p Produto
		json.NewDecoder(r.Body).Decode(&p)

		_, err := db.Exec("UPDATE Tb_produto SET Nome_Produto=? WHERE Cd_Produto=?", p.NomeProduto, id)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Write([]byte("Produto atualizado"))

	case http.MethodDelete: // DELETE
		_, err := db.Exec("DELETE FROM Tb_produto WHERE Cd_Produto=?", id)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Write([]byte("Produto deletado"))
	}
}
