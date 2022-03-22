package qp

import (
	"bytes"
	"context"
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

const testQuery = `
CREATE TABLE users (
 id INTEGER PRIMARY KEY AUTOINCREMENT,
 username TEXT UNIQUE NOT NULL,
 password TEXT NOT NULL,
 email TEXT UNIQUE NOT NULL,
 created NUMERIC NOT NULL,
 updated NUMERIC
);
INSERT INTO users (username, password, email, created) VALUES ('alice', 'passw0rd', 'alice@example.com', datetime('2017-12-05'));`

func TestPrint(t *testing.T) {
	tests := []struct {
		stmt string
		want string
	}{
		{
			"SELECT * FROM users WHERE username = 'alice'",
			`+----+----------+----------+-------------------+---------------------+---------+
| id | username | password |       email       |       created       | updated |
+----+----------+----------+-------------------+---------------------+---------+
|  1 | alice    | passw0rd | alice@example.com | 2017-12-05 00:00:00 | <nil>   |
+----+----------+----------+-------------------+---------------------+---------+
(1 row)
`},
	}
	ctx := context.Background()
	for _, tt := range tests {
		db, err := sql.Open("sqlite3", ":memory:")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := db.ExecContext(ctx, testQuery); err != nil {
			t.Fatal(err)
		}

		{
			buf := new(bytes.Buffer)
			if _, err := New(db, Context(ctx), Out(buf)).Print(tt.stmt); err != nil {
				t.Fatal(err)
			}
			if got := buf.String(); got != tt.want {
				t.Errorf("got\n%v\nwant\n%v", got, tt.want)
			}
		}

		{
			tx, err := db.BeginTx(ctx, nil)
			if err != nil {
				t.Fatal(err)
			}
			buf := new(bytes.Buffer)
			if _, err := New(tx, Context(ctx), Out(buf)).Print(tt.stmt); err != nil {
				t.Fatal(err)
			}
			if got := buf.String(); got != tt.want {
				t.Errorf("got\n%v\nwant\n%v", got, tt.want)
			}
		}
	}
}
