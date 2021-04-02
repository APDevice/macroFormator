package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)
//TODO add CML interface that allows for input of file name, template, and separator
//TODO separate filename from path
//TODO custom template variable locator
func main () {
	fName := "test.txt"
	format := "\t{}: {},\n"
	sep := ","

	file, err := os.Open(fName)
	check(err)

	defer file.Close() //make sure file is closed afterwords

	pText, err := parse(file, format, sep)
	check(err)

	fileName := strings.Split(fName, ".")
	err = ioutil.WriteFile(strings.Join(fileName, "_parsed."), pText, 0666)
	check(err)

}

func check (err error) { //handles errors
	if err != nil {
		log.Fatal(err)
	}
}

//parse takes in file and macro template and outputs file text
func parse (f *os.File, form string, div string) ([]byte, error) {
	scanner := bufio.NewScanner(f)
	var output []byte

	// split the form so that the columns can be inserted in between
	formPieces := strings.Split(form, "{}")
	formBody, formEnd := formPieces[:len(formPieces) - 1], formPieces[len(formPieces) - 1]



	for idx := 1; scanner.Scan(); idx++ {
		vals := strings.Split(scanner.Text(), div)

		//TODO add CLI conditions to handle missing or extra values

		// throw error if the number of values does not match the template
		if len(vals) < len(formPieces) - 1 {
			err := fmt.Errorf("parse error: missing value(s) on row %v", idx)

			return nil, err
		}

		for idx, piece := range formBody {
			output = append(output, []byte(piece + vals[idx])...)
		}
		output = append(output, []byte(formEnd)...)
	}

	return output, nil
}

