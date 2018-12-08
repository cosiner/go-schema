package schema

import (
	"fmt"
	"reflect"
)

type EncoderDestination interface {
	Set(source, field string, v []string) (bool, error)
}

type Encoder struct {
	parser *Parser
}

func NewEncoder(p *Parser) (*Encoder, error) {
	return &Encoder{parser: p}, nil
}

func (e *Encoder) encodeTypeToStrings(f *fieldInfo, v reflect.Value) ([]string, error) {
	if f.IsSlice {
		l := v.Len()
		if l == 0 {
			return nil, nil
		}

		vals := make([]string, 0, l)
		for i := 0; i < l; i++ {
			s, err := f.Encoding.Encode(v.Index(i).Interface())
			if err != nil {
				return nil, fmt.Errorf("encode field failed: %s", err.Error())
			}
			vals = append(vals, s)
		}
		return vals, nil
	}
	s, err := f.Encoding.Encode(v.Interface())
	if err != nil {
		return nil, fmt.Errorf("encode field failed: %s", err.Error())
	}
	if s == "" {
		return nil, nil
	}
	return []string{s}, nil
}

func (e *Encoder) encodeField(refv reflect.Value, field *fieldInfo) (v []string, err error) {
	fieldv := refv.FieldByIndex(field.Field.Index)
	return e.encodeTypeToStrings(field, fieldv)
}

func (e *Encoder) Encode(v interface{}, dst EncoderDestination) error {
	refv := reflect.ValueOf(v)
	if refv.Type().Kind() == reflect.Ptr {
		refv = refv.Elem()
	}
	reft := refv.Type()
	typInfo, err := e.parser.Parse(reft)
	if err != nil {
		return err
	}

	refv = reflect.Indirect(refv)
	for i := range typInfo.fields {
		field := &typInfo.fields[i]
		vals, err := e.encodeField(refv, field)
		if err != nil {
			return fmt.Errorf("encode field failed: %s, %s, %s", field.Field.Name, field.Sources[0], err.Error())
		}
		if len(vals) == 0 {
			continue
		}

		var ok bool
		for _, source := range field.Sources {
			ok, err = dst.Set(source.Source, source.Name, vals)
			if err != nil {
				return fmt.Errorf("set field failed: %s, %s, %v, %s", field.Field.Name, source, vals, err.Error())
			}
			if ok {
				break
			}
		}
		if !ok {
			return fmt.Errorf("cann't set to destination: %s, %s, %v", field.Field.Name, field.Sources[0], vals)
		}
	}
	return nil
}
