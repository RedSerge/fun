//	Google Drive + excessive "driving" == G[oogle] RoadHog

// This program transforms any input file into a bunch of CSV files
// which can be uploaded at Google Drive for free (0 bytes of storage!) as "Google Docs".

//This version is written in C++17, as it works faster than Python-based prototype.

#include <iostream>
#include <iomanip>
#include <fstream>
#include <string>

//Constants

const std::string ALPHABET {"0123456789abcdefghijklmnop"};
const unsigned int MAX_STRING {50000 / 2};
const unsigned int MAX_COLUMNS {3};
const unsigned int MAX_ROWS {122};

//Encode one byte into buffer
void encode_byte (char* buf, char byte) {
	//Get hexademical representation of the byte 
	snprintf(buf, 3, "%02x", (unsigned char) byte);
	//Apply the Caesar cipher
	buf[0]=ALPHABET[ALPHABET.find(buf[0]) + 10];
	buf[1]=ALPHABET[ALPHABET.find(buf[1]) + 10];
}

int main (int argc, char *argv[]) {
	//Wrong quantity of parameters
	if (argc < 2) {
		std::cout << "Please provide the name of source file." << std::endl;
		return 0;
	}
	//Buffer for encoding, one more byte for '\0' symbol
	char buf[3];
	//Open input file
	std::string path_to_file = argv[1];
	std::ifstream source(path_to_file, std::ios_base::binary);
	//Current index of CSV file block
	unsigned long long current_file {0};
	//Current length of the cell
	unsigned int current_length {0};
	//Current index of the column
	unsigned int current_column {0};
	//Current index of the row
	unsigned int current_row {0};
	while (true) {
		//Generate correct filename for current CSV block
		std::string filename {std::to_string(current_file)};
		filename.append(".csv");
		//Print it out
		std::cout << filename << std::endl;
		{
			//Create new CSV file
			std::ofstream destination(filename);
			while (true) {
				//Read one byte
				source.get(buf[0]);
				//EOF? Enough.
				if (!source) {
					return 0;
				}
				//Otherwise, proceed with encoding
				encode_byte(buf, buf[0]);
				//Write the ciphered symbols
				destination << buf[0] << buf[1];
				//Check if current length, column, row or the whole file is filled
				++current_length;
				if (current_length >= MAX_STRING) {
					current_length = 0;
					++current_column;
					if (current_column >= MAX_COLUMNS) {
						current_column = 0;
						++current_row;
						if (current_row >= MAX_ROWS) {
							current_row = 0;
							++current_file;
							//Go to the next CSV file block
							break;
						//Add 'new line' symbol if neccessary
						} else destination << '\n';
					//Add 'tab' symbol if neccessary
					}  else destination << '\t';
				}
			}
		}
	}
}