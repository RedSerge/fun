import os
from flask import Flask, request, redirect, send_from_directory #, url_for
#from werkzeug.utils import secure_filename

app = Flask(__name__)
token='/home/TZTeam/mysite/req'
key_activate='e5e4f43a7dda8e0968cf1e4cc0cbe828018993db36de10883a370592da98657ce03c150d7aec7e956b558fcea318a13440562217bf43793350f95fea7c2e56f8273c4dd3a4a3e4c40a31f8d29e7aee63860550a1e226446309a0d7cc1b1d55ccb0d260cef2eeea18860385a1a4119c2dd9506491489264304b27a665ac1589b505ca47543d98bb6cce526f5fdba415ca5b6eb2a631b3031c465e08f3776eb0c598e930070041e541835c31466ebaae04db69c963008b51da19a3f8cd5a4edc4fbbfe51aa6c6bdad260378c7386d17c848125cc11cbbe1424e7e157afda847bb7cf186e5190c92d65bc18677c4ec9670317aeaa00e3f760e984a3a8fa0b32ae'
key_stopwork='b2385a41e9d9e4a7194a8dae01b871fe287992947f127473677d96b72d1bf5c737fe80c2e7a9109e53e6eb3fd06a3e41e4179b5b4a0eb173c80f1f989ccd63fc256c6d765cbb11d4e2d9c50c10cb1341aa5cfce9d83af45f4b65b999c87b801963705566a86444a98c9d21d5bf6878fcc78f68be1d3b1d3f4d00dd04badb7983aa6eb64353457e5a3ef1ef08ad47db4bd558f9a5d0b6a1464e6f40712e448f70e3fef4ffb60a826e6ffaa136de15cdc58d64444d4e487391bcbbba0e77fb4098987266a44a0767bd1c886640f933727fada66db9751bd245b252412bc7e6b41e8f58e9f7c182590cab6c81f4376409923dd1ce10b74ab55642d97ddcef690f'
state=0

UPLOAD_FOLDER = '/home/TZTeam/mysite/work/'
ALLOWED_EXTENSIONS = set(['jpg', 'jpeg'])

app = Flask(__name__)
app.config['UPLOAD_FOLDER'] = UPLOAD_FOLDER

def allowed_file(filename):
    return '.' in filename and \
           filename.rsplit('.', 1)[1].lower() in ALLOWED_EXTENSIONS

def mode(n):
    global state
    # 0 - clear state
    # 1 - set state
    # 2 - check state
    if n!=2:
        state=n
    else:
        return state

@app.route('/uploads/<filename>')
def uploaded_file(filename):
    return send_from_directory(app.config['UPLOAD_FOLDER'], filename)

@app.route('/', methods=['GET', 'POST'])
def upload_file():
    if mode(2)==1:
        return 'Server is currently busy with the task (multiusers currently are not supported)'
    if request.method == 'POST':
        file = request.files.get('file',None)
        if not file or not file.filename:
            return redirect(request.url)
        if file and allowed_file(file.filename):
            #filename = secure_filename(file.filename)
            file.save(os.path.join(app.config['UPLOAD_FOLDER'], 'work.jpg'))
            mode(1)
            return 'Photo is saved and will be converted soon.'
    return '''
    <!doctype html>
    <title>Upload new File</title>
    <h1>Upload new File</h1>
    <form method=post enctype=multipart/form-data>
      <input type=file name=file><br>
      <input type=submit value=Upload>
    </form>
    '''

@app.route('/tokenize') #FOR DEBUG
def tokener():
    return """
    <form action="/token">
    <input name="keyid" type="hidden" value="#key#">
    <input type="submit" value="StoreRequest">
    </form>
    """.replace("#key#",key_activate)

@app.route('/token', methods=['GET'])
def newtoken():
    expected_key=str(request.args.get('keyid',''))
    if expected_key==key_activate:
        mode(1)
    elif expected_key==key_stopwork:
        mode(0)
    return str(mode(2))

@app.route('/untoken') #FOR DEBUG
def untoken():
    return """
    <form action="/token">
    <input name="keyid" type="hidden" value="#key#">
    <input type="submit" value="DerailRequest">
    </form>
    """.replace("#key#",key_stopwork)

@app.route('/requested')
def requested():
    return 'busy' if mode(2)==1 else 'free'
