from flask import Flask, request
app = Flask(__name__)
@app.route('/')
def hello_world():
	with open("d.htm") as f:
		return f.read()

@app.route('/gotcha',methods=['POST'])
def gotcha():
	return request.form.get('png')