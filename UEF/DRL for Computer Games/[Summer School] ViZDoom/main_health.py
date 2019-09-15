from os.path import isfile as model_exists
from DQN_DOOM import main

from health_reward import reward_fn
#from default_reward import reward_fn #not good idea, actually
CONFIG_NAME="scenarios/health_gathering.cfg"
MODEL_NAME="health_model"

#If model has not been generated yet...
if not model_exists(MODEL_NAME):
	#...TRAIN, otherwise...
	main(MODEL_NAME, CONFIG_NAME, reward_fn, 100000, False, False)
else:
	#...EVALUATE!
	main(MODEL_NAME, CONFIG_NAME, reward_fn, 3000, True, True)

