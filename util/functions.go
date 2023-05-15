package util

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"
)

type T1 struct {
	A int
	B int
	C string
	D D
}
type D struct {
	Name string
	E    E
}
type E struct {
	Role string
}

type T2 struct {
	A int    //`print:"B"`
	B int    //`print:"A"`
	C string `print:"D.E.Role"`
}

func CopyCommonFields(srcP, destP interface{}) {
	destV := reflect.ValueOf(destP).Elem()
	srcV := reflect.ValueOf(srcP).Elem()

	destT := destV.Type()
	for i := 0; i < destT.NumField(); i++ {
		field := destT.Field(i)

		var fieldName string
		var value reflect.Value
		if len(field.Tag.Get("print")) > 0 {
			fieldName = field.Tag.Get("print")

			if strings.Contains(fieldName, ".") {
				value = srcV
				words := strings.Split(fieldName, ".")
				for idx, word := range words {
					fmt.Printf("Word %d is: %s\n", idx, word)
					value = value.FieldByName(word)
					fmt.Println(value)
				}
			} else {
				value = srcV.FieldByName(fieldName)
			}
		} else {
			fieldName = field.Name
			log.Println(fieldName)
			//value = srcV.FieldByName(fieldName)
		}

		if !value.IsValid() || !value.Type().AssignableTo(field.Type) {
			continue
		}
		destV.Field(i).Set(value)
	}
}

//func ConvertModels(s interface{}) T {
//	res := T{}
//	v := reflect.ValueOf(s)
//	t := v.Type()
//	for i := 0; i < t.NumField(); i++ {
//		f := t.Field(i)
//		fv := v.Field(i)
//		if f.Tag.Get("print") == "yes" {
//			//res.f.Name = fv.Interface(
//			//res. = 1
//			fmt.Println(f.Name)
//			fmt.Println(fv.Interface())
//		}
//	}
//	return res
//}

func DiffTimeHours(startDate string, endDate string) int {
	startTime, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		fmt.Println("Could not parse time:", err)
	}

	endTime, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		fmt.Println("Could not parse time:", err)
	}

	diff := endTime.Sub(startTime)

	return int(diff.Hours())
}
