with open("fetchpropnames.txt",encoding="utf-8") as f:
	names=f.read()
d=names.split("\n")
PNames={}
for i in range(0,len(d),3):
	PNames[d[i]]=d[i+2]
with open("fetchstore.txt",encoding="utf-8") as f:
	names=f.read()
d=names.split("\n")
QNames={}
for i in d:
	j=i.split('\t')
	k=j[0].split('/')[-1]
	QNames[k]=j[1]
with open("fetchprops.txt",encoding="utf-8") as f:
	names=f.read()
d=names.split("\n")
QandP={}
for i in d:
	j=i.split("\t\t")
	k=j[0]
	v=list(set(j[1].split("\t")))
	QandP[k]=v


QNames=[(i[1:],QNames[i].lower()) for i in QNames]
PNames=[(i[1:],PNames[i].lower()) for i in PNames]
QP=[]
for qp in QandP:
	qval=qp[1:]
	for p in QandP[qp]:
		pval=p[1:]
		QP+=[(qval,pval)]

import sqlite3

db = sqlite3.connect('wikidata')
cursor = db.cursor()
cursor.execute('''CREATE TABLE IF NOT EXISTS PTable(id INTEGER PRIMARY KEY, value TEXT)''')
db.commit()
cursor.execute('''CREATE TABLE IF NOT EXISTS QTable(id INTEGER PRIMARY KEY, value TEXT)''')
db.commit()
cursor.execute('''CREATE TABLE IF NOT EXISTS QtoP(id INTEGER PRIMARY KEY, qid INTEGER, pid INTEGER,
FOREIGN KEY(qid) REFERENCES QTable(id), FOREIGN KEY(pid) REFERENCES PTable(id))''')
db.commit()

db.executemany("INSERT INTO QTable(id,value) VALUES (?,?)",QNames)
db.commit()

db.executemany("INSERT INTO PTable(id,value) VALUES (?,?)",PNames)
db.commit()

db.executemany("INSERT INTO QtoP(qid,pid) VALUES (?,?)",QP)
db.commit()

db.close()