package vvdatabase

import (
    "fmt"
    "database/sql"
    _"github.com/go-sql-driver/mysql"
)

var (
    // DBCon is the connection handle
    // for the database
    DBCon *sql.DB
)

func ConnectDatabase() {
    var err error
	fmt.Println("Database connecting.")
	DBCon,err = sql.Open("mysql","root:@tcp(127.0.0.1:3306)/retaildost")
	if err != nil {
        panic(err.Error())
    }
	// defer DBCon.Close()

	// results, err := db.Query("SELECT * FROM test")
    // if err != nil {
    //     panic(err.Error()) // proper error handling instead of panic in your app
    // }
    //
    // for results.Next() {
    //     var tag Tag
    //     // for each row, scan the result into our tag composite object
    //     err = results.Scan(&tag.id, &tag.name, &tag.address, &tag.contact)
    //     if err != nil {
    //         panic(err.Error()) // proper error handling instead of panic in your app
    //     }
    //             // and then print out the tag's Name attribute
    //     log.Printf(tag.name)
    // }
}
