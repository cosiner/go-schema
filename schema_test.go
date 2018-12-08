package schema_test

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/cosiner/go-schema"
)

var (
	p *schema.Parser
	d *schema.Decoder
	e *schema.Encoder
)

func init() {
	p, err := schema.NewParser("schema", []string{"path", "query", "form", "header", "body"}, func(v string) string { return v })
	if err != nil {
		log.Fatal(err)
	}
	err = p.RegisterTypes(schema.BuiltinTypes()...)
	if err != nil {
		log.Fatal(err)
	}
	d, err = schema.NewDecoder(p)
	if err != nil {
		log.Fatal(err)
	}
	e, err = schema.NewEncoder(p)
	if err != nil {
		log.Fatal(err)
	}
}

type Sources map[string]url.Values

func (s Sources) Get(source, name string) []string {
	return s[source][name]
}

func (s Sources) Set(source, name string, vals []string) (bool, error) {
	svals, has := s[source]
	if !has {
		svals = make(url.Values)
		s[source] = svals
	}
	svals[name] = vals
	return true, nil
}

func TestSchema(t *testing.T) {
	type EmbedInline struct {
		Inline string `schema:"header"`
	}
	type EmbedAnonymous struct {
		Anonymous string `schema:"header"`
	}
	type Embed struct {
		Embed string `schema:"header"`
	}
	type TestDecoderStruct struct {
		String      string      `schema:"query"`
		Bool        bool        `schema:"query"`
		Int         int         `schema:"query"`
		Int8        int8        `schema:"query"`
		Int16       int16       `schema:"query"`
		Int32       int32       `schema:"query"`
		Int64       int64       `schema:"query"`
		Uint        uint        `schema:"query"`
		Uint8       uint8       `schema:"query"`
		Uint16      uint16      `schema:"query"`
		Uint32      uint32      `schema:"query"`
		Uint64      uint64      `schema:"query"`
		Float32     float32     `schema:"query"`
		Strings     []string    `schema:"body"`
		Bools       []bool      `schema:"body"`
		Ints        []int       `schema:"body"`
		Int8s       []int8      `schema:"body"`
		Int16s      []int16     `schema:"body"`
		Int32s      []int32     `schema:"body"`
		Int64s      []int64     `schema:"body"`
		Uints       []uint      `schema:"body"`
		Uint8s      []uint8     `schema:"body"`
		Uint16s     []uint16    `schema:"body"`
		Uint32s     []uint32    `schema:"body"`
		Uint64s     []uint64    `schema:"body"`
		Float32s    []float32   `schema:"body"`
		EmbedInline EmbedInline `schema:";inline"`
		EmbedAnonymous
		Embed Embed `schema:"embed"`
	}
	query := url.Values{
		"String":  []string{"1"},
		"Bool":    []string{"true"},
		"Int":     []string{"1"},
		"Int8":    []string{"1"},
		"Int16":   []string{"1"},
		"Int32":   []string{"1"},
		"Int64":   []string{"1"},
		"Uint":    []string{"1"},
		"Uint8":   []string{"1"},
		"Uint16":  []string{"1"},
		"Uint32":  []string{"1"},
		"Uint64":  []string{"1"},
		"Float32": []string{"1"},
	}
	body := url.Values{
		"Strings":  []string{"1", "2", "3"},
		"Bools":    []string{"true", "false", "true"},
		"Ints":     []string{"1", "2", "3"},
		"Int8s":    []string{"1", "2", "3"},
		"Int16s":   []string{"1", "2", "3"},
		"Int32s":   []string{"1", "2", "3"},
		"Int64s":   []string{"1", "2", "3"},
		"Uints":    []string{"1", "2", "3"},
		"Uint8s":   []string{"1", "2", "3"},
		"Uint16s":  []string{"1", "2", "3"},
		"Uint32s":  []string{"1", "2", "3"},
		"Uint64s":  []string{"1", "2", "3"},
		"Float32s": []string{"1.1", "2.2", "3"},
	}
	header := url.Values{
		"Inline":      []string{"Inline"},
		"Anonymous":   []string{"Anonymous"},
		"Embed.Embed": []string{"Embed.Embed"},
	}
	src := Sources{
		"query":  query,
		"body":   body,
		"header": header,
	}

	var data TestDecoderStruct
	err := d.Decode(src, &data)
	if err != nil {
		t.Fatal(err)
	}

	expectData := TestDecoderStruct{
		String:         "1",
		Bool:           true,
		Int:            1,
		Int8:           1,
		Int16:          1,
		Int32:          1,
		Int64:          1,
		Uint:           1,
		Uint8:          1,
		Uint16:         1,
		Uint32:         1,
		Uint64:         1,
		Float32:        1,
		Strings:        []string{"1", "2", "3"},
		Bools:          []bool{true, false, true},
		Ints:           []int{1, 2, 3},
		Int8s:          []int8{1, 2, 3},
		Int16s:         []int16{1, 2, 3},
		Int32s:         []int32{1, 2, 3},
		Int64s:         []int64{1, 2, 3},
		Uints:          []uint{1, 2, 3},
		Uint8s:         []uint8{1, 2, 3},
		Uint16s:        []uint16{1, 2, 3},
		Uint32s:        []uint32{1, 2, 3},
		Uint64s:        []uint64{1, 2, 3},
		Float32s:       []float32{1.1, 2.2, 3},
		EmbedInline:    EmbedInline{Inline: "Inline"},
		EmbedAnonymous: EmbedAnonymous{Anonymous: "Anonymous"},
		Embed:          Embed{Embed: "Embed.Embed"},
	}
	if !reflect.DeepEqual(data, expectData) {
		t.Fatalf("unexpected decode result: %+v", data)
	}

	var dst = make(Sources)
	err = e.Encode(expectData, dst)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(src, dst) {
		t.Fatalf("unexpected encode result: %+v\n", dst)
	}
}

type httpRequestSource struct {
	req *http.Request
}

func (h *httpRequestSource) Get(source, name string) []string {
	switch source {
	case "body":
		_ = h.req.ParseForm()
		return h.req.PostForm[name]
	case "query":
		return h.req.URL.Query()[name]
	case "header":
		return h.req.Header[name]
	default:
		return nil
	}
}

type DateType struct{}

func (DateType) DataType() interface{} { return time.Time{} }

func (DateType) Decode(s string) (val interface{}, err error) {
	t, err := time.Parse("2006/01/02", s)
	if err != nil {
		return t, err
	}
	return t, nil
}

func (DateType) Encode(val interface{}) (s string, err error) {
	v, ok := val.(time.Time)
	if !ok {
		return "", fmt.Errorf("invalid data type, expect time.Time, but got %s", reflect.TypeOf(val))
	}
	return v.Format("2006/01/02"), nil
}

func newDecoder() (*schema.Decoder, error) {
	p, err := schema.NewParser("schema", []string{"body", "query", "header"}, func(name string) string {
		if name == "" {
			return ""
		}
		return strings.ToLower(name[:1]) + name[1:]
	})
	if err == nil {
		err = p.RegisterTypes(schema.BuiltinTypes()...)
	}
	if err == nil {
		err = p.RegisterTypes(DateType{})
	}
	if err != nil {
		return nil, err
	}
	return schema.NewDecoder(p)
}

func ExampleDecoder_Decode() {
	type QueryRequest struct {
		Name        string    `schema:"body"`
		Date        time.Time `schema:"body"`
		AccessToken string    `schema:"header" header:"Authorization"`
		Page        uint32    `schema:"query" query:"p"`
	}
	u, err := url.Parse("http://localhost?p=3")
	if err != nil {
		log.Fatal(err)
	}
	httpreq := http.Request{
		URL: u,
		PostForm: url.Values{
			"name": []string{"Someone"},
			"date": []string{"2018/12/25"},
		},
		Header: http.Header{
			"Authorization": []string{"Token"},
		},
	}
	d, err := newDecoder()
	if err != nil {
		log.Fatal(err)
	}
	var reqData QueryRequest
	err = d.Decode(&httpRequestSource{&httpreq}, &reqData)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s %s %s %d\n", reqData.Name, reqData.Date.Format("2006-01-02"), reqData.AccessToken, reqData.Page)
	// Output: Someone 2018-12-25 Token 3
}
