#	AliAsM: two python scripts adding ability to use the "alias" metaoperator,
#	in order to make Assembler code a little bit more fluent and understandable.
#	I made this feature, first and foremost, for myself, as I missed it a little bit.

# 	- - - Part 1: the parsing module - - -

#Import function used for crossplatform aliases
from platform import system as platform

#...In other words, "_0123456789aAbBcCdDeEfFgGhHiIjJkKlLmMnNoOpPqQrRsStTuUvVwWxXyYzZ" :
ALPHABET=''.join(("_", ''.join((str(i) for i in range(10))), ''.join((''.join((k,k.upper())) for k in (chr(i + ord('a')) for i in range(ord('z') - ord('a') + 1))))))

#Other constants
WHITESPACE=" "
DBL_WHITESPACE="  "
TO_STRIP=" \n"
TAB="\t"

#Global objects containing parameterized (functional) and ordinary aliases correspondingly
Fns={}
Aliases={}

#Strip alias string from unsupported characters
def clear_alias(string):
	return ''.join((c for c in string if c in ALPHABET))

#Detect if found word is not a substring of another one (check the "borders" of the word in a string)
def borders(string, start, finish):
	return ((start == 0) or (string[start - 1] not in ALPHABET)) and ((finish == len(string) - 1) or (string[finish + 1] not in ALPHABET))

#Find substring, return source string length instead of '-1' value if substring not found
def findsubstr(source, substr, pos):
	test = source.find(substr, pos+1)
	return len(source) if test < 0 else test

#Update given word (pattern) in source string with another replacement 
def replace_pattern(source, pattern, replaceby):
	pattern_len=len(pattern)
	source_len=len(source)
	x=-1
	pairs=[]
	while pattern:
		x = findsubstr(source, pattern, x)
		# no matched substring was found from given position, finish the loop
		if x >= source_len:
			break
		# final index of found substring
		y = x + pattern_len - 1
		# store the pair of indices
		pairs+=[(x, y)]
		# next time start from the next word
		x=y
	# replace found substrings moving backwards
	real_pairs=[p for p in pairs if borders(source, p[0], p[1])]
	real_len=len(real_pairs) - 1
	for i in range(real_len, -1, -1):
		source=replaceby.join((source[:real_pairs[i][0]], source[real_pairs[i][1] + 1:]))
	# return updated source string
	return source

#Convert aliases into actual Assembler code 
def ParseAliases(code):
	for a in Aliases:
		code=replace_pattern(code, a, Aliases[a])
	return code

#Convert functional aliases into actual Assembler code, taking passed argument values into consideration  
def ParseFn(name, argvals):
	# load given functional alias
	if not name in Fns:
		return
	args, code = Fns[name]
	# parse parameters as local subaliases, apply global ordinary alises later
	SubAliases=(dict(zip(args, argvals)), Aliases)
	for elem in SubAliases:
		for a in elem:
			code=replace_pattern(code, a, elem[a])
	# return updated code string
	return code

#Extract functional alias names for arguments and the alias itself
def FnStruct(statement, args_index):
	return (statement[:args_index].strip(' '),
	[elem.strip(' ') for elem in statement[args_index : statement.rfind(')') + 1].strip('() ').split(',')])

#Correct parsing result storage (inside a functional alias or as a pure code)
def StoreResult(Statement, Fn_buffer):
	if Fn_buffer:
		Fn_buffer[2]+=[Statement]
		return None
	else:
		return ParseAliases(Statement)

#This function converts the AliAsM code (provided in a file) into usual (Intel style, GCC-compatible) one. 
def ParseAsm(filename):
	# reset global objects
	global Fns
	global Aliases	
	Fns={}
	Aliases={}
	
	# check the count of curly braces to ignore unused platform-related code
	Ignore=0
	# temporary buffer for the functional aliases
	Fn_buffer=None
	# name of the substituted function in the original source code
	Substitution=None
	# actual Assembler code
	Results=[]
	
	# read file line by line, strip unnecessary whitespaces and comments at the end of the line ('#')
	with open (filename) as f:
		while True:
			# assuming the worst :)
			Result=None
			statement=f.readline()
			# EOF
			if not statement:
				break
			statement=statement.replace(TAB, WHITESPACE).strip(TO_STRIP)
			while DBL_WHITESPACE in statement:
				statement=statement.replace(DBL_WHITESPACE, WHITESPACE)
			comment_index=statement.rfind("#")		
			if comment_index != -1:
				statement=statement[:comment_index]
			# ignore empty lines
			if not statement:
				continue
			# count the curly brackets inside platform-related aliases:
			if Ignore:
				if statement == "}":
					Ignore-=1
				elif statement.endswith('{'):
					Ignore+=1
				continue
			# the end of a functional alias: store it into related global object and flush the buffer
			if statement == '}':
				if Fn_buffer:
					Fns[Fn_buffer[0]] = (Fn_buffer[1], '\n'.join(Fn_buffer[2]))
					Fn_buffer=None
			# the metaoperator 'alias' has been found:
			elif statement.startswith('alias'):
				statement=statement[5:]
				# functional or platform-related alias
				if statement.endswith('{'):				
					platform_index=statement.find('[')
					if platform_index != -1:
						# platform-related -> start ignoring if this block does not match the current platform
						platform_name=statement[platform_index : statement.rfind(']') + 1].strip('[] ').lower()
						if platform().lower() != platform_name:
							Ignore=1
						continue
					args_index=statement.find('(')
					if args_index != -1:
						# functional -> prepare the buffer to store the commands coming along the code flow
						fn_parameters=FnStruct(statement, args_index)
						Fn_buffer=[clear_alias(fn_parameters[0]), fn_parameters[1], []]
				# this is simple alias, left part (alias name) represents the right one (pure non-aliased Assembler code)
				elif ',' in statement and not '(' in statement:
					parts=[elem.strip(' ') for elem in statement.split(',')]
					parts[0]=clear_alias(parts[0])
					# statement "alias ,<...>" clears the whole alias dictionary
					if not parts[0]:
						Aliases={}
					# statement "alias Name," erases the alias with name "Name"
					elif not parts[1]:
						if parts[0] in Aliases:
							del Aliases[parts[0]]
					# statement "alias Name, Value" sets corresponding alias variable "Name" to the related value "Value"
					else:
						Aliases[parts[0]]=parts[1]
				# another alises: functional alias call and the substitution declaration
				else:
					args_index=statement.find('(')
					if args_index != -1:
						# functional alias -> 'call' & 'output'
						fn_data=FnStruct(statement, args_index)
						Result=StoreResult(ParseFn(fn_data[0], fn_data[1]), Fn_buffer)
					# substitution declaration -> save
					elif not Substitution:
						Substitution=statement.strip(' ')
			else:
				# ordinary instruction, 'output' the pure code
				Result=StoreResult(statement, Fn_buffer)
			# save the resulting code line
			if Result:
				Results+=[Result]
	# return the saved substitutuion value and the resulting pure code lines 
	return (Substitution, Results)

#This is a module, not a usual script
if __name__=='__main__':
	exit()