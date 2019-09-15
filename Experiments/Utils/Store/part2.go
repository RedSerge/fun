// Part 2 of the "Store" utility. For more information please read "build.py".
package main

import (
	"os"
	"path/filepath"
	"io/ioutil"
	"encoding/binary"
)

// The length of this part of executable after compilation (modified by build.py)
const ExecutableOffset=2023936

func main() {
	// Read the part of this executable file that is related to the data, not
	// the compiled alghoritm
	SourceName,_:=os.Executable()
	data,_:=ioutil.ReadFile(SourceName)
	data=data[ExecutableOffset:]
	// Parse blocks of data and store them in the corresponding outer files
	// placed in the relative paths
	var count int64 = int64(len(data))
	var index int64
	var nameindex int64
	// Stop if EOF has been reached
	for index < count {
		// We've reached the end of the line containing the name of the file?..
		if data[index] == byte('\n') {
			// Extract the corresponding name
			filename:=string(data[nameindex:index])
			// Make sure the whole path to the file exists
			fullpath,_:=filepath.Split(filename)
			os.MkdirAll(fullpath, os.ModePerm)
			// Extract the length of the block	
			index++
			firstbyte:=index + binary.MaxVarintLen64
			length,_:=binary.Varint(data[index:firstbyte])
			// Extract and save the block of data in the file
			ioutil.WriteFile(filename, data[firstbyte:firstbyte + length], 0644)
			// The next file name is supposed to start right after the block of extracted data 
			nameindex=firstbyte + length
			index=nameindex
			continue;
		}
		// ...Otherwise, keep going through the data
		index++
	}
}