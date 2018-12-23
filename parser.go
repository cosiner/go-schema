package schema

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type Type interface {
	DataType() interface{}
	Encode(v interface{}) (string, error)
	Decode(v string) (interface{}, error)
}

type fieldSource struct {
	Source string
	Name   string
}

func (f fieldSource) String() string {
	return fmt.Sprintf("(%s, %s)", f.Source, f.Name)
}

type fieldInfo struct {
	Sources  []fieldSource
	Field    reflect.StructField
	IsSlice  bool
	Encoding Type
}

type structureInfo struct {
	fields []fieldInfo
}

// format: sources[;flags], sources: source[,source]*, flags: [inline]
type FieldOptions struct {
	Sources []string
	Inline  bool
}

type Parser struct {
	optionsTag    string
	validSources  []string
	nameConverter func(string) string

	supportTypes map[reflect.Type]Type

	mu         sync.RWMutex
	structures map[reflect.Type]*structureInfo
}

// optionsTag is used for retrieve field options, each source can have it's own name, if not specified, use name of first
// source or converted field name by default.
func NewParser(optionsTag string, validSources []string, fieldNameConverter func(string) string) (*Parser, error) {
	p := Parser{
		optionsTag:    optionsTag,
		validSources:  validSources,
		nameConverter: fieldNameConverter,
		supportTypes:  make(map[reflect.Type]Type),

		structures: make(map[reflect.Type]*structureInfo),
	}
	if p.optionsTag == "" {
		return nil, fmt.Errorf("empty sources tag")
	}
	if len(p.validSources) == 0 {
		return nil, fmt.Errorf("empty valid sources")
	}
	if fieldNameConverter == nil {
		return nil, fmt.Errorf("nil field name converter")
	}
	return &p, nil
}

func (p *Parser) RegisterTypes(types ...Type) error {
	for _, t := range types {
		dt := reflect.TypeOf(t.DataType())
		_, has := p.supportTypes[dt]
		if has {
			return fmt.Errorf("duplicated types registered: %s", dt)
		}

		p.supportTypes[dt] = t
	}
	return nil
}

func (p *Parser) parseFieldOptions(val string) FieldOptions {
	if val == "" || val == "-" {
		return FieldOptions{}
	}
	var options FieldOptions
	secs := strings.SplitN(val, ";", 2)
	l := len(secs)
	if l > 0 {
		options.Sources = splitNonEmptyAndTrim(secs[0], ",")
	}
	if l > 1 {
		flags := splitNonEmptyAndTrim(secs[1], ";")
		for _, flag := range flags {
			switch flag {
			case "inline":
				options.Inline = true
			}
		}
	}
	return options
}
func (p *Parser) isFieldSourceExist(t *structureInfo, src fieldSource) bool {
	for _, f := range t.fields {
		for _, s := range f.Sources {
			if src == s {
				return true
			}
		}
	}
	return false
}

func (p *Parser) isFieldSourcesValid(sources []string) bool {
	for _, s := range sources {
		if !hasString(p.validSources, s) {
			return false
		}
	}
	return true
}

func (p *Parser) newContext(context, name string) string {
	if context == "" {
		return name
	}
	if name == "" {
		return context
	}
	return context + "." + name
}

func (p *Parser) isSupportedOrBySlice(t reflect.Type) (isSlice bool, enc Type, ok bool) {
	enc, has := p.supportTypes[t]
	if has {
		return false, enc, true
	}
	if t.Kind() == reflect.Slice {
		enc, has := p.supportTypes[t.Elem()]
		if has {
			return true, enc, true
		}
	}
	return false, nil, false
}
func (p *Parser) newIndex(parent, index []int) []int {
	if len(parent) > 0 {
		nindex := make([]int, 0, len(parent)+len(index))
		nindex = append(nindex, parent...)
		nindex = append(nindex, index...)
		index = nindex
	}
	return index
}
func (p *Parser) parse(typ reflect.Type) (*structureInfo, error) {
	type parseNode struct {
		Type    reflect.Type
		Index   []int
		Context string
	}
	var (
		typeInfo   structureInfo
		parseQueue = []parseNode{{Type: typ, Context: ""}}
	)

	for {
		l := len(parseQueue)
		if l <= 0 {
			break
		}
		node := parseQueue[0]
		copy(parseQueue, parseQueue[1:])
		parseQueue = parseQueue[:l-1]

		for i := 0; i < node.Type.NumField(); i++ {
			f := node.Type.Field(i)
			name := p.nameConverter(f.Name)
			if name == "" {
				continue
			}
			options := p.parseFieldOptions(f.Tag.Get(p.optionsTag))

			isSlice, enc, ok := p.isSupportedOrBySlice(f.Type)
			if !ok {
				if f.Type.Kind() == reflect.Struct {
					if f.Anonymous || options.Inline {
						parseQueue = append(parseQueue, parseNode{Type: f.Type, Context: node.Context, Index: p.newIndex(node.Index, f.Index)})
					} else {
						parseQueue = append(parseQueue, parseNode{Type: f.Type, Context: p.newContext(node.Context, name), Index: p.newIndex(node.Index, f.Index)})
					}
				} else if len(options.Sources) > 0 {
					return nil, fmt.Errorf("unsupported field type: %s, %s: %s", p.newContext(typ.String(), node.Context), name, f.Type.String())
				}
				continue
			}
			if len(options.Sources) == 0 {
				continue
			}
			if !p.isFieldSourcesValid(options.Sources) {
				return nil, fmt.Errorf("invalid source: field: %s, options.Sources: %v", f.Name, options.Sources)
			}

			fieldSources := make([]fieldSource, 0, len(options.Sources))
			for i, src := range options.Sources {
				val := f.Tag.Get(src)
				if val == "" {
					if i == 0 {
						val = p.newContext(node.Context, name)
					} else {
						val = fieldSources[0].Name
					}
				}
				source := fieldSource{Name: val, Source: src}
				if p.isFieldSourceExist(&typeInfo, source) {
					return nil, fmt.Errorf("duplicated field name: %+v", source)
				}
				fieldSources = append(fieldSources, source)
			}

			f.Index = p.newIndex(node.Index, f.Index)
			typeInfo.fields = append(typeInfo.fields, fieldInfo{
				Sources:  fieldSources,
				Field:    f,
				IsSlice:  isSlice,
				Encoding: enc,
			})
		}
	}
	return &typeInfo, nil
}

func (p *Parser) Parse(t reflect.Type) (*structureInfo, error) {
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("source type isn't structure: %s", t.String())
	}
	p.mu.RLock()
	info, has := p.structures[t]
	p.mu.RUnlock()
	if has {
		return info, nil
	}
	info, err := p.parse(t)
	if err != nil {
		return nil, fmt.Errorf("invalid type schema: %s, %s", t.String(), err.Error())
	}
	p.mu.Lock()
	p.structures[t] = info
	p.mu.Unlock()
	return info, nil
}
