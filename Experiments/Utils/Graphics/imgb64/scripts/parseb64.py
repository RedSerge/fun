# This script uses imgb64 utility (made by me) to generate a bunch of
# base64-formatted strings containing images, as well as
# their corresponding names, for further embedding in a code.

# Loading necessary functions from modules
from os import system, walk
from os.path import join as fullpath

# Images are placed in subfolder with the matching name
fullroot=fullpath(".","Images")

# The 'first' flag controls whether the file with strings should be
# rewritten (when the first image file has been converted)
first=True

# Parsing directory, applying the utility to the content
for root, folders, files in walk(fullroot):
	for name in files:
		system(''.join(("imgb64.exe ",'"',fullpath(fullroot,name),'" >','' if first else '>',' base64.txt')))
		first=False

# Write names separately
	with open("names.txt","w") as f:
		f.write('\n'.join(files))
		f.write('\n')
		
# Subfolders are not included
	break