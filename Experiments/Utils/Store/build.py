#	This utility ("Store") is allowing to produce self-extractable executable (not archive)
#	with relative path extraction. It may be useful to contain static parts of projects
#	in one portable file to deploy an environment of the project.

#	This script builds unified executable based on two other parts written in Go.

#	Usage example: store.exe "..\related_path" stored_executable.exe
#	Then, move resulting file "stored_executable.exe" to another folder & run.
#	The files would be extracted to the folder "related_path" created in the parent folder
#	relative to the folder containing stored_executable.exe

from os import system, name as os_name
from os.path import getsize

# Marker of changeable string in source code
MARKER="const ExecutableOffset="
MARKER_LENGTH=len(MARKER)

# Detect if exectuable needs an ".exe" extension based on OS
executable_extension="" if os_name=="posix" else ".exe"

# Read both parts of source code and create exectuables
# (It is supposed that Golang is presented and reachable)
with open("part1.go") as f:
	part1=f.readlines()
with open("part2.go") as f:
	part2=f.readlines()
system("go build part1.go")
system("go build part2.go")

# Prepare the names of corresponding executables for 1st part, 2nd part & merged (final) part
part1_filename=''.join(("part1", executable_extension))
part2_filename=''.join(("part2", executable_extension))
final_filename=''.join(("store", executable_extension))

# Detect the lengths of both executables and replace corresponding consts
# in the source code
lens=(str(getsize(part1_filename)), str(getsize(part2_filename)))
parts=(part1, part2)
for index, part in enumerate(parts):
	line_index=[line[0] for line in enumerate(part) if line[1].startswith(MARKER)][0]
	part[line_index]=''.join((part[line_index][:MARKER_LENGTH], lens[index],'\n'))
with open("part1.go","w") as f:
	f.writelines(part1)
with open("part2.go","w") as f:
	f.writelines(part2)

# Recreate executables according to the final changes
system("go build part1.go")
system("go build part2.go")

# Unite both parts into one executable
with open(final_filename, "wb") as final_build:
	with open(part1_filename, "rb") as f:
		final_build.write(f.read())
	with open(part2_filename, "rb") as f:
		final_build.write(f.read())