#Simple script comparing two files to be sure that they're equal
import sys

NAME1=sys.argv[1]
NAME2=sys.argv[2]

with open(NAME1, "rb") as f1:
	with open(NAME2, "rb") as f2:
		while True:
			byte1=f1.read(1)
			byte2=f2.read(1)
			if byte1 != byte2:
				print("Not equal!")
				exit()
			if not byte1 or not byte2:
				break
print("Equal!")
