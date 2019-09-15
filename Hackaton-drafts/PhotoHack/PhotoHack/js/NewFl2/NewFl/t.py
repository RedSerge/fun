from flask import Flask, request

app = Flask(__name__)


#import threading
#import recod

@app.route('/')
def hello_world():
	#x = threading.Thread(target=recod.freeze_support, args=None)
	#x.start()
	with open("d.htm") as f:
		return f.read()


#@app.route('/gotcha',methods=['POST'])
#def gotcha():
	#return request.form.get('png')
#x = threading.Thread(target=recod.freeze_support, args=None)x.start()
