//////////////////////////////////////////////////////////////////////////////
// file: xmlParse.go
//         A Go SAX utility that parses and displays
//         well-formed xml files.
// author: John Schwartzman, Forte Systems, Inc.
// VERSION_NUMBER = "0.1.0"
// last revision:	03/13/2019
//////////////////////////////////////////////////////////////////////////////
package main

import (
	"container/list"
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

// define some constant color escape sequences
const red = "\033[1;31;40m"    //    red on black - for element name
const green = "\033[1;32;40m"  //  green on black - for element data
const blue = "\033[1;34;40m"   //   blue on black - for attributes
const cyan = "\033[1;36;40m"   //   cyan on black - for directives
const white = "\033[0;37;40m"  //  white on black - normal printing
const violet = "\033[1;35;40m" // violet on black - for future use
const black = "\033[0;30;40m"  //  black on black - for spaces
const yellow = "\033[1;33;40m" // yellow on black - for coments

// define some color macros
const spaces = black + "%s" + white
const startname = red + "%s " + white
const endname = red + "/%s\n" + white
const parenstmnt = blue + "(%s = %s) " + white
const chardata = green + "= %s " + white
const comment = yellow + "%s" + white
const elementdata = green + "%s " + white
const directivedata = cyan + "%s\n" + white

var nLastWritePos = 0 // must be defined before use
var bShowComments = true
var nFileArg = 1

func writeNewLine() { // advance the cursor row
	fmt.Printf("\n")
}

func writeSpaces(pos int, chars string) { // position the cursor column
	for i := 0; i < pos; i++ {
		fmt.Printf(spaces, chars)
	}
}

func writeComment(data string) { // write data at current x,y
	fmt.Printf(comment, data)
}

func writeStartName(name string) { // write startName at current x,y
	fmt.Printf(startname, name)
}

func writeEndName(name string) { // write endName at current x,y
	fmt.Printf(endname, name)
}

func writeCharacterData(data string) { // write data at current x,y
	fmt.Printf(elementdata, data)
}

func writeElementData(data string) { // write data at current x,y
	fmt.Printf(elementdata, data)
}

func writeDirective(data string) { // write directive at current x,y
	fmt.Printf(directivedata, data)
}

func writeAttribute(attrName string, attrValue string) { // write attribute
	fmt.Printf(parenstmnt, attrName, attrValue) // write name-value pair
}

// add an element to the end of the list (top of the stack)
func push(s *list.List, name string) int {
	pos := s.Len()   // use the index before push
	s.PushBack(name) // push it onto the stack
	return pos
}

// remove the element at the end of the list (top of the stack)
func pop(s *list.List, name string) int {
	e := s.Back() // get the last element in the list
	if e.Value == name {
		s.Remove(e) //pop it from the stack
	} else {
		fmt.Printf("%s\nError: %s was not at the top of the stack.\n\n",
			white, name)
		os.Exit(4)
	}
	return s.Len() // use the index after pop
}

func usage(exitCode int) {
	fmt.Printf("\nUSAGE: xmlParse [-h || --help]||[-i || --ignore_comments] xmlFile\n\n")
	os.Exit(exitCode)
}

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 { // there must be 2 or 3 arguments
		fmt.Printf("\nYou have provided an incorrect number of arguments.\n")
		usage(1)
	}
	if os.Args[1] == "-h" || os.Args[1] == "--help" { // want help?
		usage(0)
	}
	if len(os.Args) != 3 && (os.Args[1] == "-i" || os.Args[1] == "--ignore_comments") {
		fmt.Printf("\nYou didn't provide the name of the xml file you want to parse.\n")
		usage(2)
	}

	optStr := os.Args[1] // we're not using an option parser so check manually
	if strings.HasPrefix(optStr, "-") {
		if optStr != "-h" && optStr != "--help" && optStr != "-i" && optStr != "--ignore_comments" {
			fmt.Printf("\nYou have entered an unknown option.\n")
			usage(3)
		}
	}
	if os.Args[1] == "-i" || os.Args[1] == "--ignore_comments" { // want comments?
		bShowComments = false
		nFileArg = 2
	}
	xmlFile, e := os.Open(os.Args[nFileArg]) // os.Args[1 or 2] is xml file
	if e != nil {
		fmt.Printf("\nProblem reading %s: %s\n\n", os.Args[1], e)
		os.Exit(2)
	}
	decoder := xml.NewDecoder(xmlFile) // create and initializethe decoder
	elementStack := list.New()         // create the stack

	for { // while there are tokens, stay in for loop
		// get a new token
		t, err := decoder.Token()
		if err != nil && err.Error() != "EOF" {
			fmt.Printf("Error: %s\n", err)
			os.Exit(3)
		}
		if t == nil {
			// we've reached the end of the document
			break // exit the for loop
		}

		// Inspect the type of the token
		switch se := t.(type) {

		case xml.StartElement: // we've encountered a startElement
			pos := push(elementStack, se.Name.Local) // push it onto the stack
			writeNewLine()
			writeSpaces(pos, "   ") // write 3 spaces per index position
			writeStartName(se.Name.Local)
			for _, a := range se.Attr { // don't need index so use dummy var
				writeAttribute(a.Name.Local, a.Value)
			}
			nLastWritePos = pos

		case xml.EndElement: // we've encountered an endElement
			pos := pop(elementStack, se.Name.Local) // pop it from the stack
			if pos < nLastWritePos {                // write name at current x pos?
				writeSpaces(pos, "   ") // set x position to write end element
			}
			writeEndName(se.Name.Local)

		case xml.CharData: // we've encountered element data
			// remove any surronding whitespace
			data := strings.TrimSpace(string(t.(xml.CharData)))
			if data != "" {
				writeCharacterData(data) // write it at current x,y
			}

		case xml.Comment: // we've encountered a comment
			if bShowComments {
				data := string(t.(xml.Comment)) // write it at current x,y
				writeComment(data)
			}

		case xml.Directive: // we've encountered a directive
			data := string(t.(xml.Directive))
			writeDirective(data)

		} // end of switch statement
	} // end of for loop

	fmt.Printf(white) // restore normal screen formatting
	writeNewLine()
}
