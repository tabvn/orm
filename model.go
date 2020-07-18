package orm

import (
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

type Field struct {
	Name    string
	Type    reflect.Type
	Primary bool
	Index   bool
	Default string
	NotNull bool
}
type Relation struct {
	Model *Model
	Type  string
}
type Model struct {
	Name      string
	Table     string
	Fields    []*Field
	Relations []*Relation
}

func SnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func GetType(obj interface{}) reflect.Type {
	return reflect.ValueOf(obj).Type()
}

func NewModel(obj interface{}) (*Model, error) {
	t := GetType(obj)
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
			Name:    t.Field(i).Name,
			Type:    t.Field(i).Type,
			Index:   false,
			Default: "",
			NotNull: false,
		}
		ormTag, ok := t.Field(i).Tag.Lookup("orm")
		if strings.ToLower(field.Name) == "id" {
			field.Primary = true
		}
		if ok {
			splitTags := strings.Split(ormTag, ";")
			for _, tag := range splitTags {
				splitLevel2 := strings.Split(strings.TrimSpace(tag), ":")
				for _, tag1 := range splitLevel2 {
					lowerTag := strings.ToLower(tag1)
					if lowerTag == "pk" {
						field.Primary = true
					}
					if lowerTag == "index" {
						field.Index = true
					}
				}
			}
		}
		model.Fields = append(model.Fields, field)
	}
	return model, nil
}
