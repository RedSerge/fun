// Part 1 of the "Store" utility. For more information please read "build.py".

package main

import (
	"fmt"
	"path/filepath"
	"os"
	"io/ioutil"
	"encoding/binary"
)

// The length of this part of executable after compilation (modified by build.py)
const ExecutableOffset=2280448

// Name of final executable that is supposed to be produced as a result
// of this utility usage; provided by user as 1st parameter
var TargetName string

// Full path to the launched executable file of this utility
var SourceName string

// Buffer of bytes to store in the produced file
var ByteBuffer []byte

// This function scrutinizes every item related to the provided path
func InsideDirectory(path string, info os.FileInfo, err error) error {
	// Ignore an item if it's a directory, this executable or the file
	// we're trying to create
	if info.IsDir() {
		return nil
	}
	fullpath,_:=filepath.Abs(path)
	if fullpath == TargetName || fullpath == SourceName {
		return nil
	}
	//Store related path to the item, as well as the length of the content
	ByteBuffer=append(ByteBuffer,[]byte(path)...)
	ByteBuffer=append(ByteBuffer,byte('\n'))
	length:=make([]byte, binary.MaxVarintLen64)
	binary.PutVarint(length, info.Size())
	ByteBuffer=append(ByteBuffer,length...)
	//Read the content and store it, too
	data,_:=ioutil.ReadFile(fullpath)
	ByteBuffer=append(ByteBuffer, data...)
	return nil
}

//Entry point
func main() {
	// Check if all necessary arguments are presented
	if len(os.Args)<3 {
		fmt.Println("Usage: 1st parameter = relative_path, 2nd = file_to_write.")
		return
	}
	// Read the path to parse, resulting file name and name of this executable
	PathToWalk:=os.Args[1]
	TargetName,_=filepath.Abs(os.Args[2])
	SourceName,_=os.Executable()
	SourceName,_=filepath.Abs(SourceName)
	// Read the appended part of the exectuable (the second part);
	// it serves as a preambule to extract the block of data correctly
	data,_:=ioutil.ReadFile(SourceName)
	ByteBuffer=append(make([]byte, 0, ExecutableOffset), data[ExecutableOffset:]...)
	// Collect information and store it into resulting executable file
	filepath.Walk(PathToWalk, InsideDirectory)
	ioutil.WriteFile(TargetName, ByteBuffer, 0644)
}