import os
import os.path

#We need to change the color format representation to be the same to simplify our model in regards of tech details :)
PARAMETER_TO_CHANGE = "screen_format"
CORRECT_UNIFIED_COLOR_SCHEME_VALUE = "CRCGCB\n"
PATH_TO_SCENARIOS = "scenarios"

#Transformer lambda function which does not touch unrelated lines of configuration file :)
transform = lambda y:' = '.join((PARAMETER_TO_CHANGE, CORRECT_UNIFIED_COLOR_SCHEME_VALUE)) if y.strip(" \t").startswith(PARAMETER_TO_CHANGE) else y

#Search for configs
for root, _, files in os.walk(PATH_TO_SCENARIOS):
	cfg_files = [os.path.join(root,f) for f in files if f.endswith(".cfg")]
	#Read, replace, store back
	for cfg in cfg_files:
		with open(cfg, encoding="utf-8") as f:
			content=[transform(line) for line in f.readlines()]
		with open(cfg, "w", encoding="utf-8") as f:
			f.writelines(content)
	#Folder parsing must be at root level only
	break
