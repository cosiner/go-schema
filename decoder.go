package schema

import (
	"fmt"
	"reflect"
)

type DecoderSource interface {
	Get(source, field string) []string
}

type Decoder struct {
	parser *Parser
}

func NewDecoder(p *Parser) (*Decoder, error) {
	return &Decoder{parser: p}, nil
}

func (d *Decoder) decodeStringsToType(f *fieldInfo, v []string) (reflect.Value, bool, error) {
	l := len(v)
	if f.IsSlice {
		refv := reflect.MakeSlice(f.Field.Type, 0, l)
		for _, v := range v {
			val, err := f.Encoding.Decode(v)
			if err != nil {
				return reflect.Value{}, false, err
			}
			refv = reflect.Append(refv, reflect.ValueOf(val))
		}
		return refv, true, nil
	}

	if l != 1 {
		return reflect.Value{}, false, fmt.Errorf("multiple values of non-slice field is not allowed: %v", v)
	}
	if v[0] == "" {
		return reflect.Value{}, false, nil
	}
	val, err := f.Encoding.Decode(v[0])
	if err != nil {
		return reflect.Value{}, false, err
	}
	return reflect.ValueOf(val), true, nil
}

func (d *Decoder) decodeField(refv reflect.Value, field *fieldInfo, v []string) (bool, error) {
	val, ok, err := d.decodeStringsToType(field, v)
	if err != nil || !ok {
		return false, err
	}
	if val.Type() != field.Field.Type {
		return false, fmt.Errorf("different decoded value type: expect %s, but got %s", field.Field.Type, val.Type())
	}

	fieldv := refv.FieldByIndex(field.Field.Index)
	fieldv.Set(val)
	return true, nil
}

func (d *Decoder) Decode(s DecoderSource, v interface{}) error {
	refv := reflect.ValueOf(v)
	if refv.Type().Kind() != reflect.Ptr {
		return fmt.Errorf("decode destination type isn't pointer: %s", refv.Type().String())
	}
	refv = refv.Elem()
	typInfo, err := d.parser.Parse(refv.Type())
	if err != nil {
		return err
	}

	for i := range typInfo.fields {
		field := &typInfo.fields[i]

		var updatedFrom fieldSource
		for _, source := range field.Sources {
			v := s.Get(source.Source, source.Name)
			if len(v) == 0 {
				continue
			}
			if updatedFrom.Source != "" {
				return fmt.Errorf("duplicated field values from different sources: %s, %s", updatedFrom, source)
			}

			ok, err := d.decodeField(refv, field, v)
			if err != nil {
				return fmt.Errorf("invalid field values: %s, %s", source, err.Error())
			}
			if ok {
				updatedFrom = source
			}
		}
	}
	return nil
}
