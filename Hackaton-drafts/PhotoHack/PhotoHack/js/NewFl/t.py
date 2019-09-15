from flask import Flask
app = Flask(__name__)
@app.route('/')
def hello_world():
	with open("demo/dN.htm") as f:
		return f.read()