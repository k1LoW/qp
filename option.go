package qp

import (
	"context"
	"io"
)

type Option func(*QueryPrinter) error

func Out(out io.Writer) Option {
	return func(qp *QueryPrinter) error {
		qp.out = out
		return nil
	}
}

func Context(ctx context.Context) Option {
	return func(qp *QueryPrinter) error {
		qp.ctx = ctx
		return nil
	}
}
