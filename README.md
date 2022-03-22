# qp

Printer that prints the results of a query to the database.

## Usage

``` go
package main

import (
	"database/sql"
	"log"

	"github.com/k1LoW/qp"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "path/to/db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, _, = qp.Print(db, "SELECT * FROM users WHERE username = 'alice'")
}
```

``` console
$ go run main.go
+----+----------+----------+-------------------+---------------------+---------+
| id | username | password |       email       |       created       | updated |
+----+----------+----------+-------------------+---------------------+---------+
|  1 | alice    | passw0rd | alice@example.com | 2017-12-05 00:00:00 | <nil>   |
+----+----------+----------+-------------------+---------------------+---------+
(1 row)
```

