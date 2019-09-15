#	Google Drive + excessive "driving" == G[oogle] RoadHog

# This program transforms any input file into a bunch of CSV files
# which can be uploaded at Google Drive for free (0 bytes of storage!) as "Google Docs".

import sys

#Name of the source file (from the command prompt)
SOURCE_FILE=sys.argv[1]

#Constants:
MAX_STRING=50000
MAX_COLUMNS=3
MAX_ROWS=122

#Constant-based constant values:
MAX_STRING2=int(MAX_STRING / 2)
MAX_COLUMNS2=MAX_COLUMNS - 1
MAX_ROWS2=MAX_ROWS - 1

#Store encoded data in CSV:
def WriteToCSV(f):
	#Current CSV file index
	CSVNum=0
	while True:
		with open(''.join((str(CSVNum), ".csv")), "w") as file_csv:
			for row in range(MAX_ROWS):
				for col in range(MAX_COLUMNS):
					for _ in range(MAX_STRING2):
						#Read one byte
						onebyte=f.read(1)
						#EOF? Halt.
						if not onebyte: return
						#Get a hex code of the byte
						hex_code='%02x' % ord(onebyte)
						#Apply the Caesar cipher (avoid treating symbols like numbers on Google side) 
						code="".join([chr(int(i, 16) + ord('a')) for i in hex_code])
						#Store the result in the CSV file block:
						file_csv.write(code)
					#Add 'tab' symbol if neccessary (end of cell)
					if col != MAX_COLUMNS2:
						file_csv.write("\t")
				#Add 'new line' symbol if neccessary (end of row)
				if row != MAX_ROWS2:
					file_csv.write("\n")
		#This file block is full, proceed to another 
		CSVNum+=1

def EncodeFile(filename):
	with open(filename,"rb") as f: WriteToCSV(f)

EncodeFile(SOURCE_FILE)