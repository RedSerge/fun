#query=input().strip(" \t")


#mass, length, wingspan

value="mass"

import sqlite3

addr="https://www.wikidata.org/wiki/Q"
db = sqlite3.connect('wikidata')
cursor = db.cursor()
cursor.execute('''SELECT id FROM PTable WHERE value=?''',(value.lower(),))
Pval=cursor.fetchone()
print(Pval);exit()
Qs=[]
if Pval:
	cursor.execute('''SELECT qid FROM QtoP WHERE pid=?''',(Pval[0],))
	Qs=cursor.fetchall()
result=[]
result2=[]
for q in Qs:
	cursor.execute('''SELECT value FROM QTable WHERE id=?''',(q[0],))
	names=cursor.fetchone()
	result+=[(q[0],names[0])]
cursor.execute('''SELECT * FROM QTable''')
Qs=cursor.fetchall()
db.close()
QsL=list(set(Qs)-set(result))
#print(QsL)
#print(result)
v=(len(Qs),len(QsL),len(result))
#print(v)
vpro=(v[1]/v[0]*100,v[2]/v[0]*100)
vres=''.join(("Нет: ",str(v[1])," (","%.3f" % float(vpro[0]),"%) ",
"Есть: ",str(v[2])," (","%.3f" % float(vpro[1]),"%) "))

y=lambda x:'<br>'.join([''.join(("<a href='",addr,str(i[0]),"'>",i[1].title(),"</a>")) for i in x])

a=['''<div class="hidden">
		<div id="WindowTrue"><span>Есть:</span><div id="ContentTrue">''',y(result),'''</div></div>
		<div id="WindowFalse"><span>Нет:</span><div id="ContentFalse">''',y(QsL),'''</div></div>
		<div id="percents">''',vres,'''</div>
	</div>''']
print('\n'.join(a))