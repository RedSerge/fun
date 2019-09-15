from matplotlib import pyplot
from collections import deque
import random
import numpy as np
import keras
from CONFIG import *
from cv2 import resize
import vizdoom as vzd
import datetime

def svg_name(name):
	return ''.join((name, "_", datetime.datetime.now().strftime("%y-%m-%d-%H-%M-%S_%f"), ".svg"))

class ReplayMemory:
	"""Simple implementation of replay memory for DQN

	Stores experiences (s, a, r, s') in circulating
	buffer
	"""
	def __init__(self, capacity, state_shape):
		self.capacity = capacity
		# Original state
		self.s1 = np.zeros((capacity, ) + state_shape, dtype=np.uint8)
		# Successor state
		self.s2 = np.zeros((capacity, ) + state_shape, dtype=np.uint8)
		# Action taken
		self.a = np.zeros(capacity, dtype=np.int)
		# Reward gained
		self.r = np.zeros(capacity, dtype=np.float)
		# If s2 was terminal or not
		self.t = np.zeros(capacity, dtype=np.uint8)
		# Current index in circulating buffer,
		# and total number of items in the memory
		self.index = 0
		self.num_total = 0

	def add_experience(self, s1, a, r, s2, t):
		# Turn states into uint8 to save some space
		self.s1[self.index] = (s1 * 255).astype(np.uint8)
		self.a[self.index] = a
		self.r[self.index] = r
		self.s2[self.index] = (s2 * 255).astype(np.uint8)
		self.t[self.index] = t

		self.index += 1
		self.num_total = max(self.index, self.num_total)
		# Return to beginning if we reach end of the buffer
		if self.index == self.capacity:
			self.index = 0

	def sample_batch(self, batch_size):
		"""Return batch of batch_size of random experiences

		Returns experiences in order s1, a, r, s2, t.
		States are already normalized
		"""
		# Here's a small chance same experience will occur twice
		indexes = np.random.randint(0, self.num_total, size=(batch_size,))
		# Normalize images to [0, 1] (networks really don't like big numbers).
		# They are stored in buffers as uint8 to save space.
		return [
			self.s1[indexes] / 255.0,
			self.a[indexes],
			self.r[indexes],
			self.s2[indexes] / 255.0,
			self.t[indexes],
		]

def update_target_model(model, target_model):
	"""Update target_model with the weights of model"""
	target_model.set_weights(model.get_weights())

def build_models(input_shape, num_actions):
	"""Build Keras models for predicting Q-values

	Returns two models: The main model and target model
	"""
	
	model = keras.models.Sequential([
		#  - Conv2D layer with 16 filters, kernel size 6, stride 3
		#  - Conv2D layer with 16 filters, kernel size 3, stride 1
		#  - Flatten
		#  - Dense layer with 32 units

		keras.layers.Conv2D(16, kernel_size=(6, 6), strides=(3, 3), activation='relu', input_shape=input_shape),
		keras.layers.Conv2D(16, kernel_size=(3, 3), strides=(1, 1), activation='relu'),
		keras.layers.Flatten(),
		keras.layers.Dense(32, activation='relu'),

		# Output layer, no activation here since Q-values can be
		# anything
		
		keras.layers.Dense(num_actions, activation=None, kernel_initializer=keras.initializers.Constant(value=0), bias_initializer=keras.initializers.Constant(value=1))
	])

	# Create target network and load 
	# current model's parameters to it
	target_model = keras.models.clone_model(model)
	update_target_model(model, target_model)

	# Compile models
	model.compile(optimizer=keras.optimizers.Adam(lr=LEARNING_RATE), loss="mse")
	target_model.compile(optimizer=keras.optimizers.Adam(lr=LEARNING_RATE), loss="mse")

	return model, target_model

def update_model(model, target_model, replay_memory, batch_size=BATCH_SIZE):
	
	"""Run single update step on model and return loss"""

	# Get bunch of random variables
	s1, a, r, s2, t = replay_memory.sample_batch(batch_size)

	# Create target values (= the best action in succeeding state).
	# This here is the "bootstrapping" part you hear in RL literature
	
	s2_values = np.max(target_model.predict(s2), axis=1)
	
	# Do not include value of succeeding state if it is terminal
	# (= it is _the end_, hence it has no future and thus no value).
	target_values = r + GAMMA * s2_values * (1 - t)

	# Get the Q-values for first state. This will work as
	# the base for target output (the "y" in prediction task)
	s1_values = model.predict(s1)
	
	# Update Q-values of the actions we took with the new
	# target_values we just calculated. This is same as writing:
	# for i in range(batch_size):
	#	 s1_values[i, a[i]] = target_values[i]
	s1_values[np.arange(batch_size), a] = target_values

	# Finally, run the update through network
	loss = model.train_on_batch(s1, s1_values)
	return loss

def get_action(s1, model, num_actions):
	"""Return action to be taken in s1 according to Q-values from model"""

	q_values = model.predict(s1[None])[0]

	# Epsilon-greedy policy:
	# With a small probability (GREEDY_EPS), return random action (integer from interval [0, num_actions - 1]).
	# Otherwise return greedy action: take the one that has most promise (the one with highest value)
	action = random.randint(0,num_actions-1) if random.random()<=GREEDY_EPS else np.argmax(q_values)
	return action, q_values

def preprocess_state(state, stacker):
	"""Handle stacking frames, and return state with multiple consecutive frames"""
	stacker.append(state)
	# Create proper state to be used by the network	
	stacked_state = np.concatenate(stacker, axis=2)
	return stacked_state

def preprocess(frame):
	#Transpose frame to resize correctly	
	frame = frame.transpose([1, 2, 0])
	frame = resize(frame, RESOLUTION)
	#Conversion: RGB range -> [0;1] range (networks prefer small ranges)
	frame = frame.astype(np.float32) / 255.0
	#Take "gray" (RGB channel avg.) color
	frame = (frame[:,:,0:1]+frame[:,:,1:2]+frame[:,:,2:3])/3
	#Transpose again to support original VizDoom graphics representation
	frame = frame.transpose([1, 0, 2])
	return frame

#Primary routine
def main(model_path, cfg_path, reward_fn, steps, evaluate=False, visible=True):
	
	#"Greedy" Epsilon
	#Once upon a time it was an ordinary constant, and then I've decided to punish it for its greediness and force it to DECAY over time... So now it's a weird CAPITAL-named variable... :)	
	global GREEDY_EPS
	
	#load game, prepare environment
	game = vzd.DoomGame()
	game.load_config(cfg_path)
	game.set_window_visible(visible)
	if not visible:
		game.set_mode(vzd.Mode.PLAYER)
	else:
		game.set_mode(vzd.Mode.ASYNC_PLAYER)
	game.init()
	
	#create deques for keeping states and rewards
	state_stacker = deque(maxlen=FRAME_STACK_SIZE)
	reward_stack = deque(maxlen=100)
	
	#prepare world-related parameters and models
	state_shape = RESOLUTION + (FRAME_STACK_SIZE,)
	model = None
	target_model = None
	action_space = game.get_available_buttons_size()
	
	#check the evaluation mode
	if not evaluate:
		# Construct new models
		model, target_model = build_models(
			state_shape,
			action_space
		)
		replay_memory = ReplayMemory(REPLAY_SIZE, state_shape)
	else:
		# Load existing model
		GREEDY_EPS=DECAY_LIM
		model = keras.models.load_model(model_path)
		replay_memory = None
		
	#The learning loop
	
	step_ctr = 0
	
	true_rewards=[]
	mod_rewards=[]

	while step_ctr < steps:
		
		terminal = False
		episode_reward = 0
		episode_q_sum = 0
		episode_q_num = 0
				
		# Keep track of losses
		losses = []

		# Reset frame stacker to empty frames
		state_stacker.clear()
		for i in range(FRAME_STACK_SIZE):
			state_stacker.append(np.zeros(RESOLUTION + (1,)))
		
		# Start new episode
		game.new_episode()
		
		# Preprocess 1st state
		s1 = preprocess(game.get_state().screen_buffer)
		s1 = preprocess_state(s1, state_stacker)
		
		# The episodic loop
		true_episode_reward = 0
		state = None
		while not terminal:
			#Prepare action and receive related Q-values
			action, q_values = get_action(s1, model, action_space)
			passed_action=[0]*action_space
			passed_action[action]=1
			
			#Store the parameters related to Q-values (for latter logging)
			episode_q_sum+=sum(q_values)
			episode_q_num+=len(q_values)
			
			#Activate prepared action and make an observation of the state
			current_episode_reward = game.make_action(passed_action, FRAMERATE)
			true_episode_reward += current_episode_reward
			statelast=state
			state=game.get_state()
			if not state:
				state=statelast
			terminal = game.is_episode_finished()
			
			# Preprocess 2nd state
			s2 = preprocess(state.screen_buffer)
			s2 = preprocess_state(s2, state_stacker)
			
			# Apply the reward rules :)
			reward, success = reward_fn(current_episode_reward, state, terminal)
			
			# Greedy epsilon decays over every victory to the certain limit,
			# expecting our model to be more relatable over time
			if success:
				GREEDY_EPS=max(GREEDY_EPS - DECAY_SPEED, DECAY_LIM)
			
			# Count the passed step
			step_ctr += 1
			# Count episodic reward
			episode_reward += reward
			
			# Skip training/replay memory stuff if we are evaluating
			if not evaluate:
				# Store the experience to replay memory
				replay_memory.add_experience(s1, action, reward, s2, terminal)
				# Check if we should do updates or saving model
				if (step_ctr % UPDATE_RATE) == 0:
					if replay_memory.num_total > SAMPLES_TILL_TRAIN:
						losses.append(update_model(model, target_model, replay_memory))
				if (step_ctr % TARGET_UPDATE_RATE) == 0:
					update_target_model(model, target_model)
				if (step_ctr % SAVE_MODEL_EVERY_STEPS) == 0:
					model.save(model_path)

			# s2 becomes s1 for the next iteration
			s1 = s2
		
		#Fill reward stack, count the vanilla and real (function-based) rewards, avg. reward
		true_rewards += [true_episode_reward]
		mod_rewards += [episode_reward]
		reward_stack.append(episode_reward)
		reward_avg=sum(reward_stack)/len(reward_stack)

		# To avoid div-by-zero
		if len(losses) == 0: losses.append(0.0)
		
		# Count avg. training loss and Q-value
		avgloss=sum(losses)/len(losses)
		qavg=episode_q_sum/episode_q_num

		#Log the info, print the message
		s = "Episode reward: {:.3f}\tSteps: {}\tAvg. loss: {}\tAvg. reward: {}\tAvg. Q: {}\tGreedinessless: {}\tVanilla reward: {:.3f}".format(
			episode_reward, step_ctr, avgloss, reward_avg, qavg, GREEDY_EPS, true_episode_reward
		)
		print(s)
	
	#Leave the environment, plot the rewards per game (vanilla & function-based)
	game.close()
	
	print("\nTrue pts:",true_rewards,"\nMod pts:",mod_rewards)
	
	pyplot.plot(true_rewards)
	pyplot.xlabel("Games")
	pyplot.ylabel("Episodic reward (true)")
	pyplot.savefig(svg_name("model_true"), format="svg")
	pyplot.close()
	
	pyplot.plot(mod_rewards)
	pyplot.xlabel("Games")
	pyplot.ylabel("Episodic reward (mod)")
	pyplot.savefig(svg_name("model_mod"), format="svg")
	pyplot.close()

#I use these parameters to launch the script properly and less verbose (just a reminder to myself =) ):
# conda activate uefdrl19;export TF_CPP_MIN_LOG_LEVEL=3
