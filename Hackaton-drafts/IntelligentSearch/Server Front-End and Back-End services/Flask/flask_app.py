template='''<html>
<head>
<style>
	form.req {
		display:inline-block;
	}
	div.hidden {
		display: none;
	}
</style>
<script>
	var block_ajax = false
	function F_AJAX (value) {
		if (block_ajax) return
		block_ajax = true
		document.getElementById("MainBlock").innerHTML="Загрузка..."
		var indeed=document.getElementById("Wolfram").checked?"merge":""
		var xhttp=new XMLHttpRequest()
		xhttp.open("GET", ["http://blagospherateam.pythonanywhere.com/",indeed,"request/!",value].join(''), true)
		xhttp.send()
		xhttp.onreadystatechange=function(){if(this.readyState==4 && this.status==200){F_ConfirmAJAX(this.responseText)}}
	}
	function F_ConfirmAJAX (value) {
		document.getElementById("MainBlock").innerHTML=value
		block_ajax = false
	}
	function formreq_sent (e) {
		F_AJAX (e.value)
	}
</script>
</head>
<body>
	Введите интересующий термин: <form class="req" onsubmit="formreq_sent(document.getElementById('request'));return false;"><input id="request" size=30 autofocus></form><input id="Wolfram" type="checkbox">Экспорт из Wolfram Alpha<br>
	<div id="MainBlock"><div class="hidden">
		<div id="WindowTrue"><span>Есть:</span><div id="ContentTrue"></div></div>
		<div id="WindowFalse"><span>Нет:</span><div id="ContentFalse"></div></div>
		<div id="percents"></div>
	</div></div>
</body>
</html>'''

from flask import Flask, send_from_directory
import sqlite3

app = Flask(__name__)

@app.route('/<path:path>')
def send_static(path):
    return send_from_directory('/home/blagospherateam/', path)

@app.route('/')
def mainpage():
    with open('/home/blagospherateam/index.html') as f:
        content=f.read()
    return content #os.getcwd()

@app.route('/engine')
def engine():
    return template

def anyrequest(value,dbname,stripped=True):
    value=value[1:]
    addr="https://www.wikidata.org/wiki/Q"
    db = sqlite3.connect(dbname)
    cursor = db.cursor()
    cursor.execute('''SELECT id FROM PTable WHERE value=?''',(value.lower(),))
    Pval=cursor.fetchone()
    Qs=[]
    if Pval:
    	cursor.execute('''SELECT qid FROM QtoP WHERE pid=?''',(Pval[0],))
    	Qs=cursor.fetchall()
    result=[]
    for q in Qs:
    	cursor.execute('''SELECT value FROM QTable WHERE id=?''',(q[0],))
    	names=cursor.fetchone()
    	result+=[(q[0],names[0])]
    cursor.execute('''SELECT * FROM QTable''')
    Qs=cursor.fetchall()
    db.close()
    QsL=list(set(Qs)-set(result))
    v=(len(Qs),len(QsL),len(result))
    vpro=(v[1]/v[0]*100,v[2]/v[0]*100)
    vres=''.join(("Нет: ",str(v[1])," (","%.3f" % float(vpro[0]),"%) ",
    "<br>Есть: ",str(v[2])," (","%.3f" % float(vpro[1]),"%) "))

    y=lambda x:'<br>'.join([''.join(("<a href='",addr,str(i[0]),"'>",i[1].title(),"</a>")) for i in x])


    if not stripped:
        a=['''<div>
        		<div id="WindowTrue"><span>Есть:</span><div id="ContentTrue">''',y(result),'''</div></div>
        		<div id="WindowFalse"><span>Нет:</span><div id="ContentFalse">''',y(QsL),'''</div></div>
        		<div id="percents">''',vres,'''</div>
        	</div>''']
    else:
        a=[y(result),'<hr>',y(QsL),'<hr>',vres]
    return '\n'.join(a)

@app.route('/request/<value>')
def request(value):
    return anyrequest(value,'/home/blagospherateam/mysite/wikidata')

@app.route('/mergerequest/<value>')
def mergerequest(value):
    return anyrequest(value,'/home/blagospherateam/mysite/wolfdata')