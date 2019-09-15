//THE HTML GRAND ADVENTURE QUIZ! - the adventure game for students learning the basics of web-programming (HTML/JS/CSS)

package main

import (
	"time"
	"fmt"
	"runtime"
	"os"
	"os/exec"
	"path/filepath"
	"net"
	"net/http"
	"net/url"
	"strings"
	"strconv"
	"io/ioutil"	
	"math/rand"
	"encoding/base64"
	"bytes"
	"image"	
    _ "image/png"
)

//A little bit of JavaScript - short-coded to minimize the resulting executable.

//Initialization, draw the triangles, client-side
const V_html=`<html>
<head>
<script>
function fillpoly(ctx,a,c){
	ctx.strokeStyle="white"
	ctx.beginPath()
	ctx.moveTo(a[0],a[1])
	var x=a.length;
	for(i=2;i<x;i+=2){
		ctx.lineTo(a[i],a[i+1])
	}
	ctx.closePath()
	ctx.fillStyle=c
	ctx.fill()
}
function rnd(a,b){
  a=Math.ceil(a)
  b=Math.floor(b)
  return Math.floor(Math.random()*(b-a))+a
}
function shuffle(a){
    for(var i=a.length-1;i>0;i--){
        var j=Math.floor(Math.random()*(i+1));
        var temp=a[i];
        a[i]=a[j];
        a[j]=temp;
    }
}
function main(){
	var colors=["firebrick","springgreen","gray","royalblue","salmon","orange","aqua","olive","gold"]
	shuffle(colors)
	var x=document.getElementById("canvas")
	var ctx=x.getContext("2d")
	ctx.fillStyle="white"
	ctx.strokeStyle="white"
	ctx.fillRect(0,0,200,200)	
	var dx=3,dy=3
	var rx=200/dx,ry=200/dy
	var s=0,s2=rx*ry/4
	var coords=[],overall=[]
	for (var j=0;j<dy;j++){
		for (var i=0;i<dx;i++){
			var px0=rx*i+4,px1=rx*(i+1)-4,py0=ry*j+4,py1=ry*(j+1)-4
			do{
				coords=[]
				for(k=0;k<3;k++){
					var kpx=rnd(px0,px1),kpy=rnd(py0,py1)					
					coords.push(kpx)
					coords.push(kpy)
				}
				side1=Math.sqrt(Math.pow(coords[0]-coords[2],2)+Math.pow(coords[1]-coords[3],2))
				side2=Math.sqrt(Math.pow(coords[2]-coords[4],2)+Math.pow(coords[3]-coords[5],2))
				side3=Math.sqrt(Math.pow(coords[4]-coords[0],2)+Math.pow(coords[5]-coords[1],2))
				var p=(side1+side2+side3)/2
				s=Math.sqrt(p*(p-side1)*(p-side2)*(p-side3))			
			}while(isNaN(s)||s<s2)
			overall.push(coords)			
			fillpoly(ctx,coords,colors[j*dy+i])
		}
	}
	document.getElementById("id0").setAttribute("value",x.toDataURL())
	document.getElementById("id1").setAttribute("value",overall)
	document.getElementById("idf").submit()
}
</script>
</head>
<body onload="main()">
<form id="idf" method="POST"><input id="id1" name="coords" type="hidden"><input id="id0" name="field" type="hidden"></form>
<canvas id="canvas" style="display:none" width="200" height="200"></canvas>
</body>
</html>`

//Redraw the "shiny" triangles, client-side
const V_happyshine_part1=`<html>
<head>
<script>
function fillpoly(ctx,a,c){
	ctx.strokeStyle="white"	
	ctx.beginPath()	
	ctx.moveTo(a[0],a[1])
	var x=a.length;
	for(i=2;i<x;i+=2){
		ctx.lineTo(a[i],a[i+1])
	}
	ctx.closePath()
	ctx.fillStyle=c
	ctx.fill()
}
function connpoly(ctx,a){
	ctx.strokeStyle="#4d0000"
	ctx.beginPath()
	ctx.setLineDash([2,1])
	ctx.moveTo(a[0],a[1])
	var x=a.length;
	for(i=2;i<x;i+=2){
		ctx.lineTo(a[i],a[i+1])
	}
	ctx.closePath()	
	ctx.stroke()
	ctx.setLineDash([])
}
function m(a,w){
	return a.reduce(function(a,b){return w(a,b);})
}
function main(){
	var x=document.getElementById("canvas")
	var ctx=x.getContext("2d")
	ctx.fillStyle="white"
	ctx.strokeStyle="white"
	ctx.fillRect(0,0,200,200)
`

const V_happyshine_part2=`
var j=c1.length;
	connpoly(ctx,c2)
	for(var i=0;i<j;i+=6){
		var c3=[c1[i],c1[i+1],c1[i+2],c1[i+3],c1[i+4],c1[i+5]]
		var cx3=[c1[i],c1[i+2],c1[i+4]]
		var cy3=[c1[i+1],c1[i+3],c1[i+5]]		
		var grd=ctx.createLinearGradient(m(cx3,Math.min),m(cy3,Math.min),m(cx3,Math.max),m(cy3,Math.max));
		grd.addColorStop(0,"orange");
		grd.addColorStop(0.35,"red");
		grd.addColorStop(0.45,"white");
		grd.addColorStop(0.47,"honeydew");
		grd.addColorStop(0.49,"white");		
		grd.addColorStop(0.54,"yellow");
		grd.addColorStop(1,"gold");
		fillpoly(ctx,[c1[i],c1[i+1],c1[i+2],c1[i+3],c1[i+4],c1[i+5]],grd)
	}	
	document.getElementById("id0").setAttribute("value",x.toDataURL())
	document.getElementById("idf").submit()
}
</script>
</head>
<body onload="main()">
<form id="idf" method="POST"><input id="id0" name="field" type="hidden"></form>
<canvas id="canvas" style="" width="200" height="200"></canvas>
</body>
</html>`

//'Suspicious' obfuscated code to supress temptation of the students to partly deduct the answers easily just by inspecting the source code :)
const V_suspicious=`<script>var _0xdbb7=['stroke','beginPath','length','lineTo','closePath'];(function(_0x5d9c15,_0x19b974){var _0x4a6516=function(_0x1cf398){while(--_0x1cf398){_0x5d9c15['push'](_0x5d9c15['shift']());}};_0x4a6516(++_0x19b974);}(_0xdbb7,0x164));var _0x7dbb=function(_0x5b4826,_0x4a3682){_0x5b4826=_0x5b4826-0x0;var _0xd64a1a=_0xdbb7[_0x5b4826];return _0xd64a1a;};function suspicious(_0x320619,_0x4bb713){_0x320619[_0x7dbb('0x0')]();_0x320619['moveTo'](_0x4bb713[0x0],_0x4bb713[0x1]);var _0x6ff0dc=_0x4bb713[_0x7dbb('0x1')];for(i=0x2;i<_0x6ff0dc;i+=0x2){_0x320619[_0x7dbb('0x2')](_0x4bb713[i],_0x4bb713[i+0x1]);}_0x320619[_0x7dbb('0x3')]();_0x320619[_0x7dbb('0x4')]();}</script>`

//^ JavaScript 'const code' ^ is finished; back to normal message constants.

//When player level is too low:
const UnavailableMessage=`Error: you can't use this tool yet.`

//Messages #1-8 (from... Yes, I enjoyed playing DOOM for a while...)

const MessagesFromHell1=`Hello, world, the game is on! Out of nowhere, the evil rocket was launched to destroy you and the whole dimension you happened to stuck in! Let’s focus on different things, though. Like, where are you, exactly?
<br>Here comes the first task: you need to sense the environment, scrutinize it pixel by pixel. Literally. Find the boundaries of the place you will call your home, for a while. Or funeral, maybe. Nevertheless, this panel, at <a href="http://localhost:8080">http://localhost:8080</a>, always provides the fresh hints and the best tools to find a way out!
<br>The first such tool is a function “GetPixel” which allows to get the pixel color by coordinates all over the world. For example, to get the color (name, if available, or at least RGB code) of the pixel located at the point {x,y} (100, 100) you need to go to the address <a href="http://localhost:8080/GetPixel?x=100&y=100">http://localhost:8080/GetPixel?x=100&y=100</a>
<br>Don’t forget to come back and refresh this page!`

const MessagesFromHell2=`Great, one pixel is found among the whole lot of nothing! Now, reveal the whole map of the place! And I would really recommend to do it using AJAX.
<br>AJAX – “Asynchronous JavaScript And XML”. This technique allows to change a part of web content by request, without reloading the whole page. The way the modern mail services, messengers and other web applications usually work nowadays.
<br><br>How it works:<br><ul><li>the event requiring sending message to the server is happening;</li><li>the event is processed on the server side, result is returned to the client side with the help of the handler function catching the result;</li><li>this function modifies the part of a web content according to the received result.</li></ul>
<br>See “ajax.html” for example.
<br>Find a way to determine the dimension boundaries and go through every pixel on the map using AJAX. Come back when done.`

const MessagesFromHell3=`Well, that map sure is pretty bizarre. A bunch of colored triangles. What are you supposed to do with them, exactly?.. Actually, these triangles are your only hope to survive the attack of the rocket slowly yet relentlessly flying somewhere above you. When merged, these triangles form an antiaircraft complex, capable to destroy the rocket before it reaches your head. So, your next task is to build the complex from its pieces, currently lying torn apart.
<br><br>First, it is necessary to determine the order of the pieces to connect them properly. The order of connection is actually the same as the order of creation of the pieces. You’re already familiar with “named” pixels on the map. The name of triangle color is the same as the ID (the value of “id” attribute) of that triangle in an SVG notation. The whole inner structure of this strange world is the elder SVG picture with rusty and abandoned elements. This SVG structure is lying somewhere under the surface, unreachable (directly). But the triangles are the children of the parent which is the said elder picture. And you should know the JavaScript DOM method allowing to get the access to the parent object, knowing its child. This time, no one gives you a tool. You should make it.
<br><br>Expand the world with your brand-new self-made JavaScript function! There’s folder “tools” in the root folder containing this program. You should put there every self-made tool you would have to make through your journey. This time, create file named “findparent.js” and write the function “DetectParent” there, accepting id of the child as an argument and returning id of its parent.
<br>In general, this function should look something like that:<br>
<pre>function DetectParent (child_id) {
	var x= /* detect the object by its id */
	// Look for the parent of object x
	return id_parent /* id_parent contains the parent id */
}</pre>
To check the function performance, navigate to the address http://localhost:8080/test?id=ColorName, where «ColorName» is an id (of any triangle). If your function is correct, the page will reveal the parent (SVG) ID! Otherwise you would see «undefined» as a sign of a failure. Don’t give up, check the function on your own examples! Come back as soon as you have that SVG ID.`

const MessagesFromHell4Intro=`Now you have the full access to the weird places underneath this world. Its skeleton mostly consists of the garbage and rusty wrecks. You need to find the parts related to the triangles and the order in which they appear within the structure. You may find AJAX useful to parse through the structure using DOM (mind the groups and hierarchy of the SVG elements). There’re other ways to solve this task, either. Anyway, good luck! The address of the structure:`

const MessagesFromHell4=`<br>There’s another problem to solve first, though. You don’t have the proper tool to connect triangles with closed polyline (polygon), it seems. You can’t even change SVG structure, adding the <polygon> tag, as it’s locked for you (read-only mode without visualization). Luckily, you still can draw the polygon on the canvas (using context methods: beginPath, closePath, moveTo, lineTo, stroke)...
<br>Create the tool file “polygon.js” and describe the function DrawPolygon (ctx, a) there. This function gets the canvas context (ctx) and an array (a) of the polygon vertex coordinates [x0,y0,x1,y1…xN,yN]. The function returns nothing, just creates unfilled closed polygon on the specified context with given vertexes (x0,y0)..(xN,yN); the fill and stroke styles of the context are remained unchanged.
<br>Check the function: <a href="http://localhost:8080/test">http://localhost:8080/test</a>. The message «Your function has been consumed…» means everything is ok, you may come back here.`

const MessagesFromHell5=`<br>Now, you have the function to build the polygon and unite the triangles in one single complex. To do so, you need to collect the array of ordered vertexes [x0,y0..x8,y8], where (x0,y0) – coordinates of a point with the color named after the most early created triangle, (x8,y8) – coordinates of the most lately created triangle (based on order within the SVG structure). The array should be sent using the POST method via the field named "colors" to the server address `

const MessagesFromHell7Intro=`The complex is activated! To lock the evenly moving rocket, you need to determine its trajectory first. Radars show that the rocket started its flight at the point `

const MessagesFromHell7=`; at the impact moment, 120 "epochs" (conventional time units) later, the rocket should arrive at the point (100,100). Detect rocket coordinates at each epoch [0..120] and create an array of truncated (integer) coordinates. Provide it to the Surface-to-Air Missile Forces using the GET method in the specified format: http://localhost:8080/SAM?rocketry=[x0,y0,..,x120,y120] , where xN & yN are the corresponding coordinates from the array at the N epoch from the rocket launch.
<br>Hint: if the initial (x0, y0) and final (xZ, yZ) points of the segment are known, and the point M divides the segment in a given ratio R = A / B, where A is the length of the segment from the beginning of the segment to M, and B is the length of the segment from M to the end of the segment, then the coordinates of the point M can be represented by the formula: Cm = Cs + (R * Cf) / (1 + R), where C stands for corresponding coordinate (X or Y) of a point, Cs – corresponding coordinate of the initial point, Cf – corresponding coordinate of the final point, Cm – corresponding coordinate of the point M.`

const MessagesFromHell8Intro=`The evil rocket is locked! Its id: `

const MessagesFromHell8=`<br>One last step to wipe the rocket off the world! Write the CSS rule for the rocket ID in a tool file “sam.css”. This rule is supposed to make the rocket object visually disappear with all the space it takes (In other words, the object should not affect the layout of the page). Visit the page <a href="http://localhost:8080/SAM">http://localhost:8080/SAM</a> when you’re ready. As soon as you got the message that the rocket disappeared, come back to celebrate your well-earned victory!`

//Victory!
const MessageFromHeaven=`Congratulations! You won the game & successfully passed the whole adventure! (...and the exam, by the way! ;)) )`

//Main part of the project starts here

//Random Number Generator
var V_rnd *rand.Rand
//Coordinates of the triangles
var V_poly []int
//Coordinates of the link between triangles
var V_poly2 []int
//State of the game, player level
var V_state int
//Image of the map
var V_map image.Image
//(base64) representation of the map
var V_mapshadow string
//Is the pixel with given coordinates revealed to the player?
var V_pixelsmet [][]byte
//Total count of revealed pixels
var V_pixelscount int
//Buffer for the self-call via the fake-POST method
var V_fakeQuery string
//The ID of the giant structure underneath the world 
var V_ParentId string
//The content of the structure
var V_structure string
//The evil rocket launch coordinates (X & Y)
var V_posX int
var V_posY int
//The ID of the evil rocket
var V_evil string

//Available colors of the triangles, indices of the colors determine the sequence of the corresponding triangle creation
var AvailableColors = []string {"gold", "gray", "aqua", "olive", "royalblue", "springgreen", "orange", "firebrick", "salmon"}
var AvailableColors_index []int

//Launch corresponding browser
func Launch (addr string) {
	switch runtime.GOOS {
		case "windows":
			exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", addr).Start()
		case "linux":
			exec.Command("xdg-open", addr).Start()
		case "darwin":
			exec.Command("open", addr).Start()
	}
}

//Random integer number in range [x;y)
func rnd (x, y int) int {
	return (int)(V_rnd.Float64() * (float64)(y - x)) + x
}

//Generate long random id
func GenId() string {
	return fmt.Sprintf("id%x%x", V_rnd.Int63(), V_rnd.Int63())
}

//Generator of the structure "underneath the world"
func StructureGen () {
	//Player level must be below 4
	if V_state >= 4 {
		return
	}
	//Find total count of elements and 9 positions of the primary components
	count:=rnd(1000, 5001)
	a:=rand.Perm(count)[:9]
	//count of group tags to close
	group_tags:=0
	//Set id of the structure and start its generation
	V_ParentId=`svg` + GenId()
	s:=[]string{"<svg id=" + V_ParentId + ">"}
	//Set random order for the 9 primary components 
	AvailableColors_index=rand.Perm(9)
	//Count of the primary components passed
	components_passed:=0
	//Go through the hierarchy elements, one by one
	for i:=0; i < count; i++ {
		//Check if any primary component is supposed to be at this position
		is_primary_component:=false
		for j:=range(a) {
			if i==a[j] {
				s=append(s, "<polygon id='", AvailableColors[AvailableColors_index[components_passed]], "'></polygon>")
				components_passed+=1
				is_primary_component=true
				break
			}
		}
		//The rest of the code supports ordinary elements
		if is_primary_component {
			continue
		}
		//1 chance out of 10 that the next element is a group tag,
		//otherwise it's a wreckage element, drown underneath the world
		chances:=rnd(0, 10)
		if (chances < 9) {
			s=append(s,fmt.Sprint("<polygon id='wreckage",i,GenId(),"'></polygon>"))
		} else {
			//Direction of group tags: 0 to open, 1 to close
			var group_direction int
			//If up to 5 tags are opened, force closing mode
			if group_tags > 5 {
				group_direction=1
			} else if group_tags==0 {
			//If no tags are opened, force opening mode
				group_direction=0
			//Otherwise, set it randomly to opening or closing mode
			} else {
				group_direction=rnd(0, 2)
			}
			//Add appropriate group tag
			if group_direction==0 {
				group_tags+=1;
				s=append(s, "<g>")
			} else {
				group_tags-=1;
				s=append(s, "</g>")
			}
		}
	}
	//Close all group tags left unclosed
	for tags:=0; tags < group_tags; tags++ {
		s=append(s, "</g>")
	}
	//Save the generated structure and set player level to 4
	V_structure=strings.Join(s, "")
	V_state=4
}

//Return ID of the structure underneath the world, if the player function works
func ParentId (id string) string {
	for i:=range (AvailableColors) {
		if id == AvailableColors[i] {
			StructureGen()
			return V_ParentId
		}
	}
	return "undefined"
}

//Get pixel at coordinates (x, y)
func PixelInterpreter (x, y int) string {
	channel_r, channel_g, channel_b, _:=V_map.At(x, y).RGBA()
	color_string:=fmt.Sprint("rgb(", uint8(channel_r), ",", uint8(channel_g), ",", uint8(channel_b), ")")
	if color_string=="rgb(0,0,0)"		{color_string="black"}
	if color_string=="rgb(255,255,255)"	{color_string="white"}
	if color_string=="rgb(255,215,0)"	{color_string="gold"}
	if color_string=="rgb(128,128,128)"	{color_string="gray"}
	if color_string=="rgb(0,255,255)"	{color_string="aqua"}
	if color_string=="rgb(128,128,0)"	{color_string="olive"}
	if color_string=="rgb(65,105,225)"	{color_string="royalblue"}
	if color_string=="rgb(0,255,127)"	{color_string="springgreen"}
	if color_string=="rgb(255,165,0)"	{color_string="orange"}
	if color_string=="rgb(178,34,34)"	{color_string="firebrick"}
	if color_string=="rgb(250,128,114)"	{color_string="salmon"}
	return color_string
}

//Load file and represent it as a tag
func File2Tag (filename, tagname string) string {
	var buffer bytes.Buffer
	buffer.WriteString("<")
	buffer.WriteString(tagname)
	buffer.WriteString(">")
	f, e:=os.Open(filename)
	if e == nil {
		content, err:=ioutil.ReadAll(f)
		if err == nil {
			buffer.WriteString(string(content))
		}
	}
	f.Close()
	buffer.WriteString("</")
	buffer.WriteString(tagname)
	buffer.WriteString(">")
	return buffer.String()
}

//Detect the evil rocket trajectory by coordinates to make the rocket finally disappear (later)
func RocketCheck (raw_coordinates string) string {
	//Clean the string with coordinates
	coordinates_string:=strings.Split(strings.Replace(strings.Replace(raw_coordinates, "[", "", -1), "]", "", -1), ",")
	coordinates:=[]int{}
	for i:=range(coordinates_string) {
		coordinate,_:=strconv.Atoi(coordinates_string[i])
		coordinates=append(coordinates, coordinate)
	}
	//If not enough coordinates are presented, stop it
	if len(coordinates)!=242 {return "undefined"}
	//Calculate the trajectory
	trajectory:=[]int{V_posX, V_posY}
	for i:=1; i < 120; i++ {
		fi:=float64(i)
		lambda:=fi / (120 - fi)
		x0:=(float64(V_posX) + lambda * 100) / (1 + lambda)
		y0:=(float64(V_posY) + lambda * 100) / (1 + lambda)
		trajectory=append(trajectory, int(x0), int(y0))
	}
	trajectory=append(trajectory, 100, 100)
	//Check if both provided and calculated trajectories are equal
	for trajectory_index:=range(trajectory) {
		if trajectory[trajectory_index] != coordinates[trajectory_index] {
			return "false"
		}
	}
	//Set the next player level and send corresponding message
	V_state=8
	return "The rocket has been identified."
}

//Load map from base64-formatted string
func GatherMap (data string) {
	//Ignore preamble, decode the rest into image
	i:=strings.Index(data, ",")
	decoded:=base64.NewDecoder(base64.StdEncoding, strings.NewReader(data[i+1:]))
	img,_,_:=image.Decode(decoded)
	//Store the info as the map and its data
	V_map=img
	V_mapshadow=data
}

//---*--- Handlers ---*---

//Handler for the '/GetPixel' branch: receive the pixel value based on the coordinates on the surface of the world
func pixelhandler (w http.ResponseWriter, r *http.Request) {
	//Headers: 1st one describes the type and character set of the web page,
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	//the 2nd one is basically to allow embedding of the page in the IFRAME elements with full access to the content
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//Check if the tool is available
	if V_state < 1 {
		w.Write([]byte(UnavailableMessage))
	} else {
		//Read coordinates
		x, y := 0, 0
		GetQuery,_:=url.Parse("?" + r.URL.RawQuery)
		for key, values:=range(GetQuery.Query()) {
			if key == "x" {
				xi,_:=strconv.Atoi(values[0])
				x=xi
			} else if key == "y" {
				yi,_:=strconv.Atoi(values[0])
				y=yi
			}
		}
		//Out of range
		if x < 0 || x >= 200 || y < 0 || y >= 200 {
			w.Write([]byte("undefined"))
		} else {
			if V_state < 3 {
				//Pixel has not been revealed before?
				if V_pixelsmet[y][x] == 0 {
					//Mark the pixel as parsed
					V_pixelsmet[y][x]=1
					V_pixelscount+=1
					//Upgrade player level for parsing through all the pixels
					if V_pixelscount>=40000 {V_state=3}
				}
			}
			w.Write([]byte(PixelInterpreter(x, y)))
		}
		//Upgrade player level for using the tool for the first time
		if V_state==1 {V_state=2}
	}
}

//Handler for the '/test' branch: uploading and testing tools used by player
//to draw polygons and detect the parents of the elements
func polyhandler (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	//The player level is too low to use this tool
	if V_state < 3 {
		w.Write([]byte(UnavailableMessage))
	} else {
		if V_state==3 {			
			//Inner self-call
			if r.Method == "POST" {
				fakePOST:=V_fakeQuery
				V_fakeQuery=""
				if fakePOST != "" {
					if fakePOST == r.PostFormValue("id") {
						GetQuery,_:=url.Parse("?" + r.URL.RawQuery)
						for key, values:=range(GetQuery.Query()) {
							//The player function works fine; check the provided id and detect the real parent
							if key == "id" {
								w.Write([]byte(ParentId(values[0])))
							}
						}
					}
				}
			} else {
				//Check the player function by generating random IDs of a parent and a child,
				//submitting the fake form as a self-call via the POST method.
				V_fakeQuery=GenId()
				FakeId:=strings.Join([]string{`unique`, GenId()}, "")
				FindParentFunc:=File2Tag(filepath.FromSlash(`tools/findparent.js`), "script")
				var buffer bytes.Buffer
				buffer.WriteString(`<form method="POST" id="`)
				buffer.WriteString(V_fakeQuery)
				buffer.WriteString(`"><input id="`)
				buffer.WriteString(FakeId)
				buffer.WriteString(`" type="hidden" name="id"></form><script>document.getElementById("`)
				buffer.WriteString(FakeId)
				buffer.WriteString(`").setAttribute("value",DetectParent("`)
				buffer.WriteString(FakeId)
				buffer.WriteString(`"));document.getElementById("`)
				buffer.WriteString(V_fakeQuery)
				buffer.WriteString(`").submit()</script>`)
				w.Write([]byte(strings.Join([]string{"<html>", FindParentFunc, buffer.String(), "</html>"}, "")))
			}
		} else if V_state == 4 {
			//Inner self-call
			if r.Method == "POST" {
				fakePOST:=V_fakeQuery
				V_fakeQuery=""
				if fakePOST != "" {
					if fakePOST == r.PostFormValue("id") {
						//Everything is ok, improve player level
						w.Write([]byte("Your function has been consumed by the world successfully!"))
						V_state=5
					}
				}
			} else {
				//Again, check the correctness of the function, this time - the one that draws polygons by given points.
				V_fakeQuery=GenId()
				PolygonFunc:=File2Tag(filepath.FromSlash(`tools/polygon.js`), "script")
				var buffer bytes.Buffer
				//To check the function, the 'suspicious' code is drawing the same random polygon on one canvas,
				//as the player function draws it on another. If resulting base64 representations of both canvases are equal,
				//the player function is accepted.
				buffer.WriteString(V_suspicious)
				buffer.WriteString(`<canvas id="canvas_id" style="display:none" width=50 height=50></canvas><canvas id="canvas_1d" style="display:none" width=50 height=50></canvas><form id="idf" method="POST"><input type='hidden' id='idi' name="id"></form><script>function rnd(){var a=[];for(var i=0;i<18;i++){a.push(parseInt(Math.random()*50))};return a}</script><script>function mainwrap(){var x=document.getElementById("canvas_id"),y=document.getElementById("canvas_1d"),z=rnd();DrawPolygon(y.getContext("2d"),z);suspicious(x.getContext("2d"),z);if(x.toDataURL()==y.toDataURL())document.getElementById("idi").setAttribute("value","`)
				buffer.WriteString(V_fakeQuery)
				buffer.WriteString(`");document.getElementById("idf").submit();}mainwrap();</script>`)
				w.Write([]byte(strings.Join([]string{"<html>", PolygonFunc, buffer.String(), "</html>"}, "")))
			}
		}
	}
}

//Handler for the '/ParseSvg' branch: check the elder structures beneath the world
//and send the coordinates describing the order of connection between the parts of the complex
func svghandler (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//This tool is not available till the player reaches level 4
	if V_state < 4 {
		w.Write([]byte(UnavailableMessage))
	} else {
		GetQuery,_:=url.Parse("?" + r.URL.RawQuery)
		for key, values:=range(GetQuery.Query()){
			//Search for the provided id
			if key=="id" {
				if V_state==5 {
					//Receive the POST method values to launch the complex
					if r.Method=="POST" {
						//Convert the coordinates and check their correctness										
						coordinates:=strings.Split(r.PostFormValue("colors"), ",")
						int_coordinates:=[]int{}
						passed:=true
						for index:=range(coordinates) {
							one_coordinate,_:=strconv.Atoi(coordinates[index])
							//Out of range - not valid coordinate
							if one_coordinate<0 || one_coordinate>=200 {
								passed=false
								break
							}
							int_coordinates=append(int_coordinates, one_coordinate)
						}
						//Not exactly 18 coordinates - not valid set
						if len(int_coordinates) != 18 {
							passed=false
						}
						//Test is failed
						if !passed {
							w.Write([]byte("undefined"))
						} else {
							//Next test: check if the coordinates are provided in the right order
							for i:=range(AvailableColors_index) {
								//If the color at the given coordinates does not match the expected color, the test is failed
								if PixelInterpreter(int_coordinates[2*i], int_coordinates[2*i+1]) != AvailableColors[AvailableColors_index[i]] {
									passed=false
									break
								}
							}
							//Test is failed
							if(!passed){
								w.Write([]byte("false"))
							} else {
								//Test is passed, increase the player level
								V_poly2=int_coordinates
								w.Write([]byte("Coordinates are accepted. The complex is ready to launch the air defense."))
								V_state=6
							}
						}
						//That's all at this stage
						return
					}
				}
				//If the id is correct, provide the whole structure without any visible parts yet with warning
				if values[0]==V_ParentId {
					w.Write([]byte("Access is limited. Visualization is turned off.<br>"))
					w.Write([]byte(V_structure))
				}
			}
		}
	}
}

//Handler for the '/SAM' branch: determine the evil rocket trajectory and make it disappear
func samhandler (w http.ResponseWriter, r *http.Request) {
	if V_state != 7 {
		if V_state > 7 {
			//Inner self-call, again
			if r.Method == "POST" {
				fakePOST:=V_fakeQuery
				V_fakeQuery=""
				if fakePOST != "" {
					if fakePOST == r.PostFormValue("id") {
						V_state=9
						w.Write([]byte("The evil rocket disappeared from existence."))
					}
				}
			} else {
				//Can the player CSS style definition finally stop the evil rocket and save the world?! Let's find out!
				V_fakeQuery=GenId()				
				CSSCanSaveTheWorld:=File2Tag(filepath.FromSlash(`tools/sam.css`), "style")
				//The trick is, offsetParent is strictly equal to null when the object is invisible.
				//That way, we know for sure that the rocket disappeared.
				w.Write([]byte(strings.Join([]string{"<html>", CSSCanSaveTheWorld, fmt.Sprint(`<object`, GenId(), ` id="`, V_evil, `"></object2><form method="POST" id="idf"><input id="idi" type="hidden" name="id"></form><script>if(document.getElementById("`, V_evil, `").offsetParent===null)document.getElementById("idi").setAttribute("value","`, V_fakeQuery, `");document.getElementById("idf").submit()</script>`), "</html>"}, "")))
			}
		} else {
			//The player level is way too low to proceed
			w.Write([]byte(UnavailableMessage))
		}
	} else {
		//We're just trying to determine the rocket trajectory yet
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		v,_:=url.Parse("?" + r.URL.RawQuery)
		for key, values:=range(v.Query()) {
			if key == "rocketry" {
				w.Write([]byte(RocketCheck(values[0])))
			}
		}
	}
}

//The primary handler, with instructions and overall guidance
func handler (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	//Initialization of the map, the first message to player
	if V_state == 0 {
		if r.Method == "POST" {
			data:=r.PostFormValue("field")
			GatherMap(data)
			coordinates:=strings.Split(r.PostFormValue("coords"), ",")
			int_coordinates:=[]int{}
			for i:=range(coordinates) {
				one_coordinate,_:=strconv.Atoi(coordinates[i])
				int_coordinates=append(int_coordinates, one_coordinate)
			}
			//Store polygon coordinates data
			V_poly=int_coordinates
			w.Write([]byte(MessagesFromHell1))
			//Advance the state (player level)
			V_state=1
		} else {
			//Force fake POST self-submitting, generating the environment for the first time
			w.Write([]byte(V_html))
		}
	} else if V_state==1 {
		//Message #1
		w.Write([]byte(MessagesFromHell1))
	} else if V_state==2 {
		//Message #2
		w.Write([]byte(MessagesFromHell2))
	} else if V_state==3 {
		//Revealed map & Message #3
		w.Write([]byte(strings.Join([]string{"<img src='", V_mapshadow, "'/><br>", MessagesFromHell3}, "")))		
	} else if V_state==4 {
		//Map, message #4, link
		prelink:="http://localhost:8080/ParseSvg?id="
		w.Write([]byte(strings.Join([]string{"<img src='", V_mapshadow, "'/><br>", MessagesFromHell4Intro, "<a href='", prelink, V_ParentId, "'>", prelink, V_ParentId, "</a>", MessagesFromHell4}, "")))
	} else if V_state==5 {
		//Map, message #5, link
		prelink:="http://localhost:8080/ParseSvg?id="
		w.Write([]byte(strings.Join([]string{"<img src='", V_mapshadow, "'/><br>", MessagesFromHell4Intro, "<a href='", prelink, V_ParentId, "'>", prelink, V_ParentId, "</a>", MessagesFromHell5, prelink, V_ParentId}, "")))
	} else if V_state==6 {
		if r.Method == "POST" {
			//Complex has been activated
			data:=r.PostFormValue("field")
			GatherMap(data)
			w.Write([]byte(strings.Join([]string{"<img src='", V_mapshadow, "'/><br>", MessagesFromHell7Intro, "(", fmt.Sprint(V_posX, ",", V_posY), ")", MessagesFromHell7}, "")))
			//Next level
			V_state=7
		} else {
			//"Activate" the complex, redrawing the components (triangles) & send fake POST self-call upon completion
			w.Write([]byte(strings.Join([]string{V_happyshine_part1, "var c1=", strings.Replace(fmt.Sprint(V_poly), " ", ",", -1), ";var c2=", strings.Replace(fmt.Sprint(V_poly2), " ", ",", -1), V_happyshine_part2}, "")))
		}
	} else if V_state==7 {
		//Map, message, coordinates
		w.Write([]byte(strings.Join([]string{"<img src='", V_mapshadow, "'/><br>", MessagesFromHell7Intro, "(", fmt.Sprint(V_posX, ",", V_posY), ")", MessagesFromHell7}, "")))
	} else if V_state==8 {
		//Map, message, the rocket ID
		w.Write([]byte(strings.Join([]string{"<img src='", V_mapshadow, "'/><br>", MessagesFromHell8Intro, V_evil, MessagesFromHell8}, "")))
	} else if V_state>8 {
		//Victory!!!
		w.Write([]byte(strings.Join([]string{"<img src='", V_mapshadow, "'/><br><h1>", MessageFromHeaven, "</h1>"}, "")))
	}
}

//^ ---*--- END OF Handlers ---*--- ^

//Inititalize the startup parameters
func init() {
	//Launch random generator
	rand.Seed(time.Now().UnixNano())
	V_rnd=rand.New(rand.NewSource(time.Now().UnixNano()))
	//Prepare the storage to check if certain pixel has been investigated
	V_pixelsmet=make([][]byte, 200)
	for i:=0; i < 200; i++ {
		V_pixelsmet[i]=make([]byte, 200)
	}
	//Generate the evil rocket ID and its start coordinates
	V_evil=`rock`+GenId()
	for V_posX >= 0 && V_posX < 200 {
		V_posX=rnd(-10000, 10000)
	}
	for V_posY >= 0 && V_posY < 200 {
		V_posY=rnd(-10000, 10000)
	}
}

//Entry point
func main() {
	//Launch server
	addr:="localhost:8080"
	fulladdr:=strings.Join([]string{"http://", addr}, "")
	network_connection, err:=net.Listen("tcp", addr)
	if err != nil {
		return
	}
	//Connect the handlers
	http.HandleFunc("/", handler)
	http.HandleFunc("/test", polyhandler)
	http.HandleFunc("/GetPixel", pixelhandler)
	http.HandleFunc("/ParseSvg", svghandler)		
	http.HandleFunc("/SAM", samhandler)
	//Connect the server
	fmt.Printf("Starting server at %s ...\n", fulladdr)
	Launch(fulladdr)
	http.Serve(network_connection, nil)
}