FAILURE_TIME=300
POS_REWARD=1
NEG_REWARD=-0.01

def reward_fn(reward, state, terminal):
	reward_test = int(state.tic < FAILURE_TIME-1)
	success = terminal and reward_test
	reward = POS_REWARD if success else NEG_REWARD
	return reward, success
