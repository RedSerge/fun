#	Google Drive + excessive "driving" == G[oogle] RoadHog

# This script recovers the downloaded GRoadHog blocks back to the original content.
# It looks for the folder containing all the neccessary blocks.
# It proceeds them all till the resulting file is recovered and ready to use.

import sys

#Name of the resulting file (from the command prompt)
DESTINATION_FILE=sys.argv[1]

from os import walk
from os.path import join as fullpath

# This part of code depends on "xlrd" library which can read data from Excel files -
# the file format for exported Google Docs
import xlrd
#This function reads our former CSV file block (now encoded as Excel Workbook) cell-by-cell,
#deciphers the block content and writes it back to the original file
def ReadWorkBook(f, path):
	#Open workbook and the first sheet
	wb=xlrd.open_workbook(path)
	sheet=wb.sheet_by_index(0)
	#Read every cell one-by-one
	for x in range(sheet.nrows):
		for y in range(sheet.ncols):
			cell_value=sheet.cell_value(rowx=x, colx=y)
			#Write decoded byte value back into original (f)ile
			[f.write(i) for i in decode(cell_value)]
#"xlrd"-dependent part is over

#Decode the cell
def decode(cell):
	#One byte encoded as a pair of hexademical digits
	pairs=[cell[i : i + 2] for i in range(0, len(cell), 2)]
	decoded=[]
	for p in pairs:
		#Those digits are encoded with the Caesar cipher 
		decoded += [bytes([int('%x%x' % 
		(ord(p[0]) - ord('a'), ord(p[1]) - ord('a')), 16)])]
	#The result is a set of decoded bytes from one cell
	return decoded

#Find a folder that contains no subfolders but files with '.xlsx' extension:
for root, folders, files in walk("."):
	if any([f.endswith('.xlsx') for f in files]) and not folders:
		#Sort the files by their number (0.xlsx, 1.xlsx, 2.xlsx, ..., 10.xlsx, ...)
		sorted_files=sorted(files, key=lambda item : int(item.split('.')[0]))
		break
		
#Start writing the destination file: read books, decode, fill the resulting file
with open(DESTINATION_FILE, "wb") as destination:
	for f in sorted_files:
		ReadWorkBook(destination, fullpath(root, f))