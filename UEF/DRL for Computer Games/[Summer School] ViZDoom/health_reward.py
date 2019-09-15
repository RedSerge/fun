POS_REWARD=1
NEG_REWARD=-0.01

Health1=None
Health2=None

def reward_fn(reward, state, terminal):
	global Health1
	global Health2
	reward=0
	if state:
		Health2=state.game_variables
		if Health1 is not None:
			reward=POS_REWARD if Health2>Health1 else NEG_REWARD
		Health1=Health2
	if terminal:
		Health1=None
	return reward, terminal
