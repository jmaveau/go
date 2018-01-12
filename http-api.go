package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var (
	addr = flag.String("addr", ":8080", "http service address")
	data map[string]string
)

func main() {
	connectmariadb()
	flag.Parse()
	data = map[string]string{}
	r := httprouter.New()
	r.GET("/entry/:key", show)
	r.GET("/list", show)
	r.PUT("/entry/:key/:value", update)
	err := http.ListenAndServe(*addr, r)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func connectmariadb() {
	db,_ := sql.Open("mysql", "API_user:@/API_DB")
	defer db.Close()

	// Connect and check the server version
	var version string
	quer := db.QueryRow("SELECT VERSION()").Scan(&version)
	//("SELECT * FROM Animal")
	if quer != nil{
		fmt.Println("Connected to:", version)
		fmt.Println(quer)
	} else {
		fmt.Println("Perdu")
	}
}

func show(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	k := p.ByName("key")
	if k == "" {
		fmt.Fprintf(w, "Read list: %v\n", data)
		return 
	}
	fmt.Fprintf(w, "Read entry: data[%s] = %s\n", k, data[k])
}

func update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	k := p.ByName("key")
	v := p.ByName("value")

	data[k] = v
	fmt.Fprintf(w, "Updated: data[%s] = %s\n", k, data[k])
}

