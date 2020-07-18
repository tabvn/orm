package orm

import (
	"log"
	"reflect"
	"regexp"
	"strings"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

type RelationType string

const (
	OneToOne   RelationType = "OneToOne"
	ManyToOne               = "ManyToOne"
	ManyToMany              = "ManyToMany"
)

type Column struct {
	Primary  bool
	Name     string
	Type     string
	Default  string
	Nullable bool
	Index    bool
}
type Field struct {
	Name   string
	Type   reflect.Type
	Column *Column
}

type Relation struct {
	Model *Model
	Type  string
}
type Model struct {
	Name      string
	Table     string
	Fields    []*Field
	Relations map[string]*Relation
}

func SnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func getType(obj interface{}) reflect.Type {
	return reflect.ValueOf(obj).Type()
}

func getColumn(t reflect.StructField) *Column {
	log.Println("type", t.Name, t.Type.Kind())
	if t.Type.Kind() == reflect.String {
		log.Println("col", t.Name, t.Type)
	}
	if t.Type.Kind() == reflect.Ptr {
		log.Println("field nae", t.Name, t.Type)
	}
	col := &Column{
		Name:     SnakeCase(t.Name),
		Type:     "",
		Default:  "",
		Nullable: false,
		Index:    false,
	}
	ormTag, ok := t.Tag.Lookup("orm")
	if strings.ToLower(t.Name) == "id" {
		col.Primary = true
	}
	if ok {
		splitTags := strings.Split(ormTag, ";")
		for _, tag := range splitTags {
			splitLevel2 := strings.Split(strings.TrimSpace(tag), ":")
			for _, tag1 := range splitLevel2 {
				lowerTag := strings.ToLower(tag1)
				if lowerTag == "pk" {
					col.Primary = true
				}
				if lowerTag == "index" {
					col.Index = true
				}
			}
		}
	}

	return col
}

func NewModel(obj interface{}) (*Model, error) {
	t := getType(obj)
	kind := t.Kind()
	if kind == reflect.Ptr {
		t = t.Elem()
	}
	model := &Model{
		Name:      t.Name(),
		Table:     SnakeCase(t.Name()),
		Fields:    []*Field{},
		Relations: nil,
	}
	for i := 0; i < t.NumField(); i++ {
		field := &Field{
			Name: t.Field(i).Name,
			Type: t.Field(i).Type,
		}
		field.Column = getColumn(t.Field(i))
		model.Fields = append(model.Fields, field)
	}
	return model, nil
}
