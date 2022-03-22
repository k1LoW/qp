# qp

Printer that prints the results of a query to the database.

## Usage

``` go
package main

import (
	"database/sql"
	"log"

	"github.com/k1LoW/qp"
)

func main() {
	db, err := sql.Open("postgres", "user=root password=root host=localhost dbname=test sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, _, = qp.Print(db, "SELECT * FROM users WHERE name = 'alice';")
}

```

