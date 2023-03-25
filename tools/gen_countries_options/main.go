package main

import (
	"fmt"

	"github.com/biter777/countries"
)

func main() {
	fmt.Println(`<select name="country">`)
	fmt.Println("\t<option value=\"unknown\" selected>Выберите страну...</option>")
	for _, c := range countries.AllInfo() {
		if len(c.Code.StringRus()) <= 25 {
			fmt.Printf("\t<option value=%q>%s</option>\n", c.Alpha2, c.Code.StringRus())
		}
	}
	fmt.Println("</select>")
}
