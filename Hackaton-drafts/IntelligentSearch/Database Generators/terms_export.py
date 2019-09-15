with open("fetchpropnames.txt") as f:
	d=f.read()
e=d.split("\n")
for i in range(0,len(e),3):
	print(e[i+2])