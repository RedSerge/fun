# Hardcoded resolution
RESOLUTION = (30, 45)

# Hardcoded memory size (keep it small)
REPLAY_SIZE = 10000
BATCH_SIZE = 32
# Discount factor
GAMMA = 0.99
# Minimum number of experiences
# before we start training
SAMPLES_TILL_TRAIN = 1000
# How often model is updated
# (in terms of agent steps)
UPDATE_RATE = 4
# How often target network is updated
# (in terms of agent steps)
TARGET_UPDATE_RATE = 1000
# How often we should save the model
# (in terms of agent steps)
SAVE_MODEL_EVERY_STEPS = 500
# Number of frames in frame stack
# (number of successive frames provided to agent)
FRAME_STACK_SIZE = 2
# Learning rate for Adam optimizer
# (What is Adam and why we use it? See this blog post and its figures:
#  http://ruder.io/optimizing-gradient-descent/ )
LEARNING_RATE = 0.00025
# Framerate for the action
FRAMERATE = 4
# "Greedy" Epsilon parameter, speed of the decay, 
# limit of the decay (Epsilon value can't be lesser than the limit).
GREEDY_EPS = 0.9
DECAY_SPEED = 0.05
DECAY_LIM = 0.02
