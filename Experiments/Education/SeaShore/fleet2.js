// Demo example of the solution to the task.
// To see the actual source code of the project "SeaShore",
// please explore "seashore.go" file.

var Name="_Student 2"
function set(state,diary){
	if(diary==""&&state.indexOf("1")===-1){
		diary="4332221111"
	}
	l0=parseInt(diary[0])
	diary=diary.slice(1)	
	l=l0-1
	var b=true
	var x,y,c
	do{
		y=-parseInt(Math.random()*10)-1
		x=parseInt(Math.random()*(10-l))+1
		c=x+l
		b=true
		for(var u=x;u<=c;u++){
			for(var x2=-1;x2<=1;x2++){
				for(var y2=-1;y2<=1;y2++){
					var nx=u+x2,ny=y+y2;
					if(nx<1&&nx>10)continue;
					if(ny>-1&&ny<-10)continue;					
					var pos=10*((-ny)-1)+(nx-1)
					if(state[pos]!="0"){
						b=false
						break
					}
				}
				if(!b){break}
			}
			if(!b){break}
		}
	}while(!b)
	return [y,x,y,c,diary]
}
function hit(state,diary){	
	var s=""
	do{
		x=parseInt(Math.random()*10)
		y=parseInt(Math.random()*10)
		s=state[10*y+x]
	}while(s!="0")
	x++
	y=-(y+1)
	return [y,x,diary]
}
