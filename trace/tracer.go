package trace

import (
	"fmt"
	"io"
)

type tracer struct {
	out io.Writer
}

func (t *tracer) Trace(a ...interface{}) {
	t.out.Write([]byte(fmt.Sprint(a...)))
	t.out.Write([]byte("\n"))
}

//Tracerはコード内での出来事を記録できるオブジェクトを表すインタフェース

type Tracer interface {
	// ..interface{}は任意の型の引数を任意の個数受け取ることができるという意味
	Trace(...interface{})
}

func New(w io.Writer) Tracer {
	return &tracer{out: w}

}

type nilTracer struct{}

func (t *nilTracer) Trace(a ...interface{}) {}

// OffはTraceメソッドの呼び出しを無視するTracerを返す
func Off() Tracer {
	return &nilTracer{}
}
