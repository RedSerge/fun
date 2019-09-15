/*	This program converts a picture to base64 web-compatible format
 *	and exports the result as a line to the standard output stream.
 */
 
package main

import (
	"fmt"
	"os"
	"bufio"
	"path/filepath"
	"bytes"
	"strings"	
	"encoding/base64"
)

func Base64 (filename string) string {
	//Trying to open file with given name 
	f,e:=os.Open(filename)
	//Failure results in an empty string
	if e != nil {
		return ""
	}
	//Close file after function execution automatically
	defer f.Close()
	//Retrieving the information about file to allocate the necessary amount of memory
	fileinfo,_:=f.Stat()
	buf:=make([]byte, fileinfo.Size())
	//Reading the file content
	content:=bufio.NewReader(f)
	content.Read(buf)
	//Preparing the buffer for the generated base64 string
	var result bytes.Buffer
	result.WriteString("data:image")
	//Extracting the extension, trying to append to the resulting string
	extension:=filepath.Ext(filename)
	if len(extension) > 0 {
		result.WriteString("/")
		result.WriteString(strings.ToLower(extension[1:]))
	}
	result.WriteString(";base64,")
	result.WriteString(base64.StdEncoding.EncodeToString(buf))
	return result.String()
}

func main () {
	// If filename is missing, warn and stop.
	if len(os.Args)<2 {
		fmt.Println("The filename has not been specified.")
		return
	}
	// Output the result
	fmt.Println(Base64(os.Args[1]))
}