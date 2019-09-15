from os.path import isfile as model_exists
from DQN_DOOM import main

from default_reward import reward_fn
CONFIG_NAME="scenarios/defend_the_center.cfg"
MODEL_NAME="center_model"

#If model has not been generated yet...
if not model_exists(MODEL_NAME):
	#...TRAIN, otherwise...
	main(MODEL_NAME, CONFIG_NAME, reward_fn, 200000, False, False)
else:
	#...EVALUATE!
	main(MODEL_NAME, CONFIG_NAME, reward_fn, 1000, True, True)

