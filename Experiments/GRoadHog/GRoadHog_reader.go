//	Google Drive + excessive "driving" == G[oogle] RoadHog

// This program reads exported "Google Docs" (".csv.xlsx") file blocks to recover the original file.
// Basically, the code does the same job as "GRoadHog_recover.py", but is written in Go and, therefore, runs faster. 
// Also, this method uses rather different approach, parsing XLSX file directly (as XML, not as Excel Workbook).
// For perfomance reasons, almost every part not related to the extraction process itself is eliminated.
// For example, this version of code doesn't search for the right folder with Excel files, supposing the program is already therein.

package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"strconv"
	"archive/zip"
)

//Alphabet for the Caesar cipher
const alphabet=`0123456789abcdefghijklmnop`

func main () {
	//Check if the name of the output file is provided
	if len(os.Args) < 2 {
		fmt.Println("Please provide the name of output file.")
		return
	}
	//Create the output file
	output_file, err := os.Create(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer output_file.Close()
	//Prepare the buffer. The sliding window of 32768 bytes is a format restriction.
	bytes:=make([]byte, 32768)
	//Start from part #0
	number:=0
	for {
		//Open the current part
		xl, err := zip.OpenReader(strings.Join([]string{strconv.Itoa(number), ".csv.xlsx"}, ""))
		//Stop if it does not exist
		if err != nil {
			return
		}
		//Prepare the buffers
		wholefile:=make([]byte, 0, 18306019)
		wholelist:=make([]byte, 0, 16549)
		//We need to find just 2 files in the archive, 0 found yet
		wholecount:=0
		for _, f := range xl.File {
			if f.Name=="xl/sharedStrings.xml" || f.Name=="xl/worksheets/sheet1.xml" {
				wholecount++
				//Open subfile either with values of shared strings or links towards them (marks from the 1st sheet)
				shared, err := f.Open()
				if err != nil {
					continue
				}
				//Save the whole file content to the corresponding buffer 
				for {
					count, err:=shared.Read(bytes)
					if f.Name=="xl/sharedStrings.xml" {
						wholefile=append(wholefile, bytes[:count]...)
					} else {
						wholelist=append(wholelist, bytes[:count]...)
					}
					if err != nil /* == io.EOF or in any other case */ {
						break
					}
				}
				shared.Close()
			}
			if wholecount >= 2 {
				//Search within the archive should be successfully finished at this point
				break
			}
		}
		xl.Close()
		//Parse the files and recover necessary information from tags;
		//here, we don't actually check the format properly, saving our time & believing Google :)
		values:=strings.Split(strings.ReplaceAll(strings.ReplaceAll(string(wholefile), "</t></si><si>", ""), "</t></si></sst>", ""), "<t>")[1:]
		marks:=strings.Split(string(wholelist), "<v>")[1:]
		//Buffer premade for one file; capacity is set based on constant values defined in the primary part of GRoadHog
		total_bytes:=make([]byte, 0, 9150000)
		//Decode information and store it into the buffer defined above
		for _,m:=range marks {
			//Index of the string
			index,_:=strconv.Atoi(strings.Split(strings.Trim(m, " "), "</v>")[0])
			//Corresponding string
			encoded_bytes:=values[index]
			total_count:=len(encoded_bytes)
			//Decode each 2 symbols using the Caesar cipher, store appropriate decoded value as a byte
			for i:=0; i < total_count; i+=2 {
				decoded_byte:=string([]byte{alphabet[strings.Index(alphabet, string(encoded_bytes[i])) - 10], alphabet[strings.Index(alphabet, string(encoded_bytes[i + 1])) - 10]})
				decoded_value,_:=strconv.ParseUint(decoded_byte, 16, 8)
				total_bytes=append(total_bytes, byte(decoded_value))
			}
		}
		//Save the bytes received from the block to the output file
		_,output_err:=output_file.Write(total_bytes)
		if output_err != nil {
			log.Fatal(err)
		}
		//Next block number
		number++
	}
}