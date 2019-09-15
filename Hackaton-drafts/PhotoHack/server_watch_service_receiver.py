import urllib.request
import time
import os.path

def getsignal():
	u2 = urllib.request.urlopen('http://tzteam.pythonanywhere.com/requested')
	return u2.read()

state=0
fname="work.jpg"

while True:
	time.sleep(2)
	if getsignal()!=b'free':
		u2 = urllib.request.urlopen('http://tzteam.pythonanywhere.com/uploads/work.jpg')
		if not os.path.isfile(fname):
			with open(fname,"wb") as f:
				f.write(u2.read())
		urllib.request.urlopen('http://tzteam.pythonanywhere.com/token?keyid=b2385a41e9d9e4a7194a8dae01b871fe287992947f127473677d96b72d1bf5c737fe80c2e7a9109e53e6eb3fd06a3e41e4179b5b4a0eb173c80f1f989ccd63fc256c6d765cbb11d4e2d9c50c10cb1341aa5cfce9d83af45f4b65b999c87b801963705566a86444a98c9d21d5bf6878fcc78f68be1d3b1d3f4d00dd04badb7983aa6eb64353457e5a3ef1ef08ad47db4bd558f9a5d0b6a1464e6f40712e448f70e3fef4ffb60a826e6ffaa136de15cdc58d64444d4e487391bcbbba0e77fb4098987266a44a0767bd1c886640f933727fada66db9751bd245b252412bc7e6b41e8f58e9f7c182590cab6c81f4376409923dd1ce10b74ab55642d97ddcef690f')
