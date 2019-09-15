from ctypes import *
from platform import system as platform

name=''.join(("markov",".dll" if platform().lower()=="windows" else ".so" ))
dll=cdll.LoadLibrary(name)

class EulerStruct(Structure):
	_fields_ = [("first", POINTER(c_double)), ("second", POINTER(c_double)), ("third", c_int)]

eulerstruct = EulerStruct()
eulerstruct.first = (c_double * 30)(0.635578,0.000000,0.0,0.599658,0.726859,1.0,0.407514,0.316904,0.0,0.820063,0.860164,0.0,0.860408,0.690909,0.0,0.627155,0.270882,0.0,0.274728,0.080538,0.0,0.843379,0.617298,0.0,0.799341,0.168920,0.0,0.000000,0.531236,0.0);
eulerstruct.second = (c_double * 30)(0.1,0.1,0.05,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0)
eulerstruct.third = 10
dll.Euler(eulerstruct.first,eulerstruct.second,eulerstruct.third)
for i in range(30):
	print(eulerstruct.first[i],end="\n")