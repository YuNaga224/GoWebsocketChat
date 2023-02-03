package trace

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("Newからの戻り地がnilです")
	} else {
		tracer.Trace("こんにちは、トレースパッケージ")
		if buf.String() != "こんにちは、トレースパッケージ\n" {
			t.Errorf("'%s'という誤った文字列が検出されました.", buf.String())
		}
	}

}

func TestOff(t *testing.T) {
	var silentTracer = Off()
	silentTracer.Trace("データ")
}
