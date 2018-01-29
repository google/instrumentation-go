package zpages

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"go.opencensus.io/trace"
)

var (
	tid  = trace.TraceID{1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 4, 8, 16, 32, 64, 128}
	sid  = trace.SpanID{1, 2, 4, 8, 16, 32, 64, 128}
	sid2 = trace.SpanID{0, 3, 5, 9, 17, 33, 65, 129}
)

func TestTraceRows(t *testing.T) {
	now := time.Now()
	later := now.Add(2 * time.Second)
	data := traceDataFromSpans("foo", []*trace.SpanData{&trace.SpanData{
		SpanContext:  trace.SpanContext{TraceID: tid, SpanID: sid},
		ParentSpanID: sid2,
		Name:         "foo",
		StartTime:    now,
		EndTime:      later,
		Attributes:   map[string]interface{}{"stringkey": "stringvalue", "intkey": 42, "boolkey": true},
		Annotations: []trace.Annotation{
			{Time: now.Add(time.Millisecond), Message: "hello, world", Attributes: map[string]interface{}{"foo": "bar"}},
			{Time: now.Add(1500 * time.Millisecond), Message: "hello, world"},
		},
		MessageEvents: []trace.MessageEvent{
			trace.MessageEvent{Time: now, EventType: 2, MessageID: 0x3, UncompressedByteSize: 0x190, CompressedByteSize: 0x12c},
			trace.MessageEvent{Time: later, EventType: 1, MessageID: 0x1, UncompressedByteSize: 0xc8, CompressedByteSize: 0x64},
		},
		Status: trace.Status{
			Code:    1,
			Message: "d'oh!",
		},
	}})
	fakeTime := "2006/01/02-15:04:05.123456"
	for i := range data.Rows {
		data.Rows[i].Fields[0] = fakeTime
	}
	if want := (traceData{
		Name: "foo",
		Num:  1,
		Rows: []traceRow{
			traceRow{Fields: [3]string{fakeTime, "    2.000000", ""}, SpanContext: trace.SpanContext{TraceID: tid, SpanID: sid}, ParentSpanID: sid2},
			traceRow{Fields: [3]string{fakeTime, "", `Status{canonicalCode=CANCELLED, description="d'oh!"}`}},
			traceRow{Fields: [3]string{fakeTime, "", `Attributes:{boolkey=true, intkey=42, stringkey="stringvalue"}`}},
			traceRow{Fields: [3]string{fakeTime, "     .     0", "received message [400 bytes, 300 compressed bytes]"}},
			traceRow{Fields: [3]string{fakeTime, "     .  1000", `hello, world  Attributes:{foo="bar"}`}},
			traceRow{Fields: [3]string{fakeTime, "    1.499000", "hello, world"}},
			traceRow{Fields: [3]string{fakeTime, "     .500000", "sent message [200 bytes, 100 compressed bytes]"}}}}); !reflect.DeepEqual(data, want) {
		t.Errorf("traceRows: got %v want %v\n", data, want)
	}

	var buf bytes.Buffer
	writeTextTraces(&buf, data)
	if want := `When                       Elapsed(s)   Type
2006/01/02-15:04:05.123456     2.000000 trace_id: 01020304050607080102040810204080 span_id: 0102040810204080 parent_span_id: 0003050911214181
2006/01/02-15:04:05.123456              Status{canonicalCode=CANCELLED, description="d'oh!"}
2006/01/02-15:04:05.123456              Attributes:{boolkey=true, intkey=42, stringkey="stringvalue"}
2006/01/02-15:04:05.123456      .     0 received message [400 bytes, 300 compressed bytes]
2006/01/02-15:04:05.123456      .  1000 hello, world  Attributes:{foo="bar"}
2006/01/02-15:04:05.123456     1.499000 hello, world
2006/01/02-15:04:05.123456      .500000 sent message [200 bytes, 100 compressed bytes]
`; buf.String() != want {
		t.Errorf("writeTextTraces: got %q want %q\n", buf.String(), want)
	}
}
