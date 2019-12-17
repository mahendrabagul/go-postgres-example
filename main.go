package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

//HealthCheck healthcheck endpoint
func HealthCheck(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprint(w, "Ok")
}

func main() {
	router := httprouter.New()
	router.GET("/", HealthCheck)
	go CreateData()
	log.Fatal(http.ListenAndServe(":5000", router))
}

//CreateData create data in the DB.
func CreateData() {
	connstring := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_DATABASE"))
	db, err := sql.Open("postgres", connstring)
	if err != nil {
		log.Println(`Could not connect to db`)
		panic(err)
	}
	defer db.Close()

	t := time.Now()
	createstmt := fmt.Sprintf("CREATE TABLE TEST_TABLE_%s(counter varchar);", t.Format("20060102150405"))
	log.Printf("Creating Table TEST_TABLE_%s :", t.Format("20060102150405"))
	_, err = db.Exec(createstmt)
	if err != nil {
		panic(err)
	}

	for index := 0; index < 50; index++ {
		log.Println("Inserting records in the database")
		t2 := time.Now()

		insertstmt := fmt.Sprintf("INSERT INTO TEST_TABLE_%s(counter) values (%s);", t.Format("20060102150405"), t2.Format("20060102150405"))
		_, err = db.Exec(insertstmt)
		if err != nil {
			panic(err)
		}

		selectstmt := fmt.Sprintf("SELECT * FROM TEST_TABLE_%s;", t.Format("20060102150405"))

		rows, err := db.Query(selectstmt)
		if err != nil {
			panic(err)
		}

		var col1 string
		for rows.Next() {
			rows.Scan(&col1)
			log.Printf("Fetch records got. :%s\n", col1)
		}
		time.Sleep(3 * time.Second)
	}

}
