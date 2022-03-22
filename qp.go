package qp

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type Querier interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

func Print(db Querier, stmt string) (int, error) {
	return New(db).Print(stmt)
}

type QueryPrinter struct {
	out io.Writer
	db  Querier
	ctx context.Context
}

func New(db Querier, opts ...Option) *QueryPrinter {
	qp := &QueryPrinter{
		out: os.Stderr,
		db:  db,
		ctx: context.Background(),
	}
	for _, opt := range opts {
		if err := opt(qp); err != nil {
			panic(err)
		}
	}
	return qp
}

func (qp *QueryPrinter) Print(stmt string) (int, error) {
	r, err := qp.query(stmt)
	if err != nil {
		return 0, err
	}
	return fmt.Fprint(qp.out, string(r))
}

func (qp *QueryPrinter) query(stmt string) ([]byte, error) {
	if !strings.HasPrefix(strings.ToUpper(stmt), "SELECT") {
		// exec
		_, err := qp.db.ExecContext(qp.ctx, stmt)
		if err != nil {
			return nil, err
		}
		return nil, nil
	} else {
		// query
		rows := []map[string]interface{}{}
		r, err := qp.db.QueryContext(qp.ctx, stmt)
		if err != nil {
			return nil, err
		}
		defer r.Close()
		columns, err := r.Columns()
		if err != nil {
			return nil, err
		}
		types, err := r.ColumnTypes()
		if err != nil {
			return nil, err
		}
		for r.Next() {
			row := map[string]interface{}{}
			vals := make([]interface{}, len(columns))
			valsp := make([]interface{}, len(columns))
			for i := range columns {
				valsp[i] = &vals[i]
			}
			if err := r.Scan(valsp...); err != nil {
				return nil, err
			}
			for i, c := range columns {
				switch v := vals[i].(type) {
				case []byte:
					s := string(v)
					t := strings.ToUpper(types[i].DatabaseTypeName())
					if strings.Contains(t, "CHAR") || t == "TEXT" {
						row[c] = s
					} else {
						num, err := strconv.Atoi(s)
						if err != nil {
							return nil, err
						}
						row[c] = num
					}
				default:
					row[c] = v
				}
			}
			rows = append(rows, row)
		}
		if err := r.Err(); err != nil {
			return nil, err
		}
		buf := new(bytes.Buffer)
		table := tablewriter.NewWriter(buf)
		table.SetHeader(columns)
		table.SetAutoFormatHeaders(false)
		table.SetAutoWrapText(false)
		for _, r := range rows {
			row := make([]string, 0, len(columns))
			for _, c := range columns {
				row = append(row, fmt.Sprintf("%v", r[c]))
			}
			table.Append(row)
		}
		table.Render()
		c := len(rows)
		if c == 1 {
			_, _ = fmt.Fprintf(buf, "(%d row)\n", len(rows))
		} else {
			_, _ = fmt.Fprintf(buf, "(%d rows)\n", len(rows))
		}
		return buf.Bytes(), nil
	}
}
