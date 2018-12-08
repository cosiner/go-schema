# go-schema
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/cosiner/go-schema) 
[![Build Status](https://travis-ci.org/cosiner/go-schema.svg?branch=master&style=flat)](https://travis-ci.org/cosiner/go-schema)
[![Coverage Status](https://coveralls.io/repos/github/cosiner/go-schema/badge.svg?style=flat)](https://coveralls.io/github/cosiner/go-schema)
[![Go Report Card](https://goreportcard.com/badge/github.com/cosiner/go-schema?style=flat)](https://goreportcard.com/report/github.com/cosiner/go-schema)

go-schema is a simple library for go to bind source strings into structures, inspired by [gorilla/schema](https://github.com/gorilla/schema).

# install
```bash
go get github.com/cosiner/go-schema
``` 
# features
* implements builtin data types for almost all go primitive types: bool,string,int(8,16,32,64), uint, float...
* support slice
* support custom data type by implements specified interface
* support multiple data source such as url query params, path params, headers, and so on. user can add their own sources
  by implements specified interface.

# example
```Go

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

```

# license
MIT.