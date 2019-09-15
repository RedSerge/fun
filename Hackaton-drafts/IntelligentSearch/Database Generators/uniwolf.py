import shutil
shutil.copyfile("wikidata","wolfdata")
import sqlite3
db = sqlite3.connect('wolfdata')
with open("wolfram_result.txt") as f:
	d=f.read()
with open("fetchstore.txt",encoding="utf-8") as f:
	e=f.read()

mass="2067"
e=e.split("\n")
d=d.split("\n")
v=[]
for i in e:
	i=i.split("\t")
	i[0]=i[0].split('/Q')[-1]
	for j in d:
		j=j.split(" -> ")[0]
		if j==i[1]:
			#cursor = db.cursor()
			#cursor.execute('''SELECT qid,pid FROM QtoP WHERE qid=? AND pid=?''',(i[0],mass))
			v+=[(i[0],mass)]
			break

db.executemany("INSERT INTO QtoP(qid,pid) VALUES (?,?)",v)
db.commit()
db.close()