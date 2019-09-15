#	AliAsM: two python scripts adding ability to use the "alias" metaoperator,
#	in order to make Assembler code a little bit more fluent and understandable.
#	I made this feature, first and foremost, for myself, as I missed it a little bit.

# 	- - - Part 2: the calling code - - -

#Import the 1st part (the parsing module)
from asmparse import ParseAsm

#Import function to launch corresponding commands using OS-provided tools
from os import system as run

#Import functions related to path manipulation and command line arguments
import os.path
import sys

#Not enough arguments
if len(sys.argv) < 2:
	print("Please provide path to configuration file")
	exit()

#Configuration file name
CFG_FILE=sys.argv[1]
#Default parameters
GCC_BASE_PARAMS="-m64 -masm=intel -S"
#Alias extension
ALIAS_EXT=".asm"

#Use pattern template, like '%12', to bind values to the corresponding arguments %1, %2
def substitute(string, substitution, values):
	if len(substitution) > 1:
		val_len=len(values)
		for i, s in enumerate(substitution[1:]):
			if i>=val_len:
				break
			string=string.replace(''.join((substitution[0], s)), values[i])
	return string

#Try to detect configuration file and change the working directory
CFG_FILE=os.path.realpath(CFG_FILE)
if not os.path.isfile(CFG_FILE):
	exit()
os.chdir(os.path.dirname(CFG_FILE))

# Configuration file may contain following attributes:

	#-	default 		- 	default parameters to 'transpile' the given C code into assembler
	#-	main 			-	name of the related C/Assembler code without extension
	#-	gcc				-	compiler path/name ('gcc' by default)
	#-	substitution	-	substitution template
	#-	buildargs		-	parameters allowing to compile the resulting assembler code (with substitution support)
	#-	parameters		-	additional non-default parameters for compiler (with substitution support)
	#-	postrun			-	post-compilation command (with substitution support)
	#-	aliases			-	name of the folder containing alias files

#Default configuration dictionary object
cfg={'default':GCC_BASE_PARAMS, 'main':'', 'gcc':'gcc', 'substitution':'', 'buildargs':'', 'parameters':'', 'postrun':'', 'aliases':''}

#Load configuration file
with open(CFG_FILE) as f:
	for line in f:
		line=[l.strip(' \t\n') for l in line.split('=')]
		cfg[line[0].lower()]=line[1]

#'Transpile' C code into Assembler
source=''.join((cfg['main'], ".c"))
destination=''.join((cfg['main'], ".s"))
run(' '.join([cfg['gcc'], cfg['default'], substitute(cfg['parameters'], cfg['substitution'], (source, destination))]))

#Failure, exit
if not os.path.isfile(destination):
	exit()

#Load parsed (pure) aliased code and related substitution value describing which part of the original
#source code is supposed to be substituted with the corresponding alias 
Aliases={}

for root, _, files in os.walk(cfg['aliases']):
	for f in files:
		if os.path.splitext(f)[1].lower() == ALIAS_EXT:
			#Parse every .asm file in the corresponding folder		
			subalias=ParseAsm(os.path.join(root, f))
			if subalias[0]:
				Aliases[subalias[0]]=subalias[1]
	#Root level only (author's whim)
	break

#If we seem to have something to substitute, let's try to change the original source
if Aliases:
	Transform=[]
	with open(destination) as f:
		#Flag to drop the original funcion
		PurgeFn=False
		while True:
			statement=f.readline()
			if not statement:
				break
			statement=statement.strip('\n')
			statement_check=statement.strip(' \t')
			if PurgeFn:
				#End of the function, stop 'dropping', keep going
				if statement_check == "ret":
					PurgeFn=False
				continue
			elif statement_check.endswith(":"):
				for a in Aliases:
					#Replace the original function with the related aliased content
					if statement_check[:-1] == a:
						statement='\n'.join(Aliases[a])
						PurgeFn=True
						break
			#Store the result into memory, line by line
			Transform+=[statement]
	#Rewrite original file containing Assembler code with the modified content
	with open(destination, "w") as f:
		f.write('\n'.join(Transform))
		f.write('\n')

#Compile & run the resulting Assembler code (on-demand)
if cfg['buildargs']:
	run(' '.join([cfg['gcc'], substitute(cfg['buildargs'], cfg['substitution'], (destination, cfg['main']))]))
if cfg['postrun']:
	run(' '.join([substitute(cfg['postrun'], cfg['substitution'], (destination, cfg['main']))]))