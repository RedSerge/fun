// SeaShore - the version of the Battlefield game for HTML/JS/CSS education course,
// where the battle between the players is wadged by the JavaScript code (of their corresponding authorship). 

package main

import (
	"fmt"
	"strings"
	"os"
	"os/exec"
	"io/ioutil"
	"runtime"
	"net"
	"net/http"
	"strconv"
	"log"
	"bytes"
)

//Global state of the game; -20 at the start (10 ships to set by each of 2 players)
var GameState int = -20
//Two maps, one per player
var Overseas [] *Sea = [] *Sea {NewSea(), NewSea()}
//Who is currently playing? (1st or 2nd player)
var PlayerTurn int = 1
//Should the program summarize (finish) the game or keep playing?
var LastWords = false
//Base path to the scripts of both players
var ScriptName string = `fleet`
//Player names
var Name1, Name2 string

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

//Load file and represent it as a diary string (see below)
func File2Diary(filename string) string {	
	var diary string
	f, e:=os.Open(filename)
	if e==nil {
		content, err:=ioutil.ReadAll(f)
		if err==nil {
			diary=string(content)
			if len(diary)>2048 {diary=string([]rune(diary)[:2048])}
		}
	}
	f.Close()
	return diary
}

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

//Find Integer Maximum between a and b
func IMax (a, b int) int {
	if b > a {return b}
	return a
}

//Find [Integer] Minimum between a and b
func IMin (a, b int) int {
	if b < a {return b}
	return a
}

//Class describing sea map and corresponding states
type Sea struct {
	//Quantity of unused ships of different kind
	Ships []int
	//Description of the battlefield
	Field []rune
	//Battlefield represented as a string value
	FieldV string
	//Additional information describing attempts to hit the ships
	Trace []rune
	//Coordinates X, Y of the last attacked square on the map
	X int
	Y int
	//Ships left to destroy by enemy on this map
	Left int
	//Player score
	Score int
	//Diary - random string containing up to 2048 symbols
	Diary2048 string
} //NewSea() - constructor with default values for corresponding fields:
func NewSea() *Sea {return &Sea {[]int{4, 3, 2, 1}, []rune(strings.Repeat("0", 100)), "", []rune(strings.Repeat("0", 100)), -1, -1, 0, 0, ""}}

//Check if the piece of ship at (a, b) doesn't intersect with another one
func (this *Sea) CheckPiece (a, b int) bool {
	for i:=-1; i <= 1; i++ {
		for j:=-1; j <= 1; j++ {
			y, x := a + i, b + j
			if x < 0 || x > 9 || y < 0 || y > 9 {continue}			
			if this.Field[y * 10 + x] == '1' {return false}
		}
	}
	return true
}

//Put a piece of ship (1) on the map at (a, b), marking the space around the piece as 'surroundings' (2)
func (this *Sea) PutPiece (a, b int) {
	for i:=-1; i <= 1; i++ {
		for j:=-1; j <= 1; j++ {
			y, x := a + i, b + j
			if x < 0 || x > 9 || y < 0 || y > 9 {continue}
			pos := y * 10 + x
			if this.Field[pos] == '1' {continue}
			if i == 0 && j == 0 {this.Field[pos] = '1'} else {this.Field[pos] = '2'}
		}
	}
}

//Create full ship from the pieces placed at provided coordinates
func (this *Sea) BuildShip (y0, x0, y1, x1 int) int {
	//Reorder coordinates in the right direction, assigning them as y/x (s)tart and (f)inish
	x:=[]int{-IMin(y0, x0) - 1, IMax(y0, x0) - 1,-IMin(y1, x1) - 1, IMax(y1, x1) - 1}	
	ys, xs, yf, xf := IMin(x[0], x[2]), IMin(x[1], x[3]), IMax(x[0], x[2]), IMax(x[1], x[3])
	//Detect y-axis/x-axis (l)ength 
	yl, xl := yf - ys + 1, xf - xs + 1
	//At least one dimension should have length equal to 1, coordinates should not be out of the map bounds
	if yl != 1 && xl != 1 {return 5}
	for i:=range(x){if x[i] < 0 || x[i] > 9 {return 5}}
	//True, "z-axis" (l)ength of the ship
	zl:=xl + yl - 1
	//Ship can't have length more than 4 (pieces)
	if zl > 4 {return 5}
	//The ship of this type/length should be allowed to be set
	if this.Ships[zl - 1] <= 0 {return zl}
	//There should be enough room for ship
	for i:=ys; i <= yf; i++ {
		for j:=xs; j <= xf; j++ {
			if !this.CheckPiece(i, j){return zl}
		}
	}
	//If all criteria are satisfied, start building the ship
	for i:=ys; i <= yf; i++ {
		for j:=xs; j <= xf; j++ {
			this.PutPiece(i, j)
		}
	}
	//This kind of ship has been built, there's one more ship in the player fleet 
	this.Ships[zl - 1]-=1
	this.Left+=1
	return 0
}

//Try to build the ship at provided coordinates
func (this *Sea) SetShip (y0, x0, y1, x1 int) int {
	//Try to build the ship and send the message:
	fmt.Print("(", y0, ",", x0, ")-(", y1, ",", x1, ") : ")
	err := this.BuildShip(y0, x0, y1, x1)
	//-5 points for incorrect ship, (length-5) points for every failed ship of known type, 0 points for success
	var mess string
	var pts int
	if err == 5 {mess = "Ship building error (-5 points)"; pts = -5
	}else if err == 0 {mess = "Ship has been built successfully"; pts = 0
	}else {pts = err - 5; mess = fmt.Sprint("Ship building error; ship length = ", err, " pieces (", pts, " points)")}
	fmt.Println(mess)
	//Return resulting points
	return pts
}

//Print "visible" side of the map
func (this *Sea) Print () {
	for i:=0; i < 10; i++ {fmt.Println(string(this.Field[i * 10 : (i + 1) * 10]))}
}

//Print "hidden" side of the map
func (this *Sea) Expose () {
	for i:=0; i < 10; i++ {fmt.Println(string(this.Trace[i * 10 : (i + 1) * 10]))}
}

//Generate the string representing the map, substituting "surroundings" for an "empty space"
func (this *Sea) Redescribe () {
	this.FieldV = strings.Replace(string(this.Field), "2", "0", -1)
}

//Check if the fire at (a, b) was reasonable; penalize it otherwise
func (this *Sea) ReasonableHit (a, b int) int {
	//Storage for coordinates of attacked zones 
	hitzone:=[]int{}
	//Flag to determine whether the unwise resource wasting test is positive
	test:=false
	for y:=0; y < 10; y++ {
		for x:=0; x < 10; x++ {
			pos:=y * 10 + x
			//If this position has already been shot...
			if this.Trace[pos] == '1' {
				test = true
				//Check the surroundings, as player could unreasonably aim there
				for i:=-1; i <= 1; i += 2 {
					for j:=-1; j <= 1; j += 2 {
						x0, y0 := x + j, y + i
						if x0 < 0 || x0 > 9 || y0 < 0 || y0 > 9 {continue}
						if y0 == a && x0 == b {
							//Player must be penalized in that case
							fmt.Println("A priori missed shot (-7 points)")
							return -7
						}
						//Mark the hit area and all its surroundings for further testing
						hitzone = append(hitzone, y0 * 10 + x0)
					}
				}
				//Axis loop
				for x_or_y:=0; x_or_y < 2; x_or_y++ {
					//Sign loop
					for sign:=-1; sign <= 1; sign+=2 {
						//Maybe player fired reasonably, near the previously damaged (not defeated yet) ship, after all
						x_test, y_test := x, y
						if x_or_y == 0 {x_test+=sign} else {y_test+=sign}
						if x_test < 0 || x_test > 9 || y_test < 0 || y_test > 9 {continue}
						pos_test:=y_test * 10 + x_test
						wise_shot:=true
						for hitpoint:=range(hitzone) {
							if hitzone[hitpoint] == pos_test {
								wise_shot = false
								break
							}
						}
						if wise_shot {if y_test == a && x_test == b {return 0}}
					}
				}
			}
		}
	}
	//The decision was not wise enough; penalize the player
	if test {
		fmt.Println("Unwise resource wasting, fire far away from the recently hit ship (-6 points)")
		return -6
	}
	//Otherwise, no penalty (strictly 0 points)
	return 0
}

//Detect the borders of the ship which piece is located at (a, b) in the given direction along given axis
func (this *Sea) CheckBorders (a, b, offset int, axis bool) int {
	origin:=a
	var pos int
	for {
		destination:=origin + offset
		if destination < 0 || destination > 9 {return origin}
		if axis {pos = destination * 10 + b} else {pos = b * 10 + destination}
		//No ship pieces left
		if this.Field[pos] != '1' {return origin}
		origin = destination
	}
}

//Check if the ship is sunk after last attack at (a, b)
func (this *Sea) CheckDestruction (a, b int) bool {
	//Detect the ship borders
	y0:=this.CheckBorders(a, b, -1, true)
	y1:=this.CheckBorders(a, b, 1, true)
	x0:=this.CheckBorders(b, a, -1, false)
	x1:=this.CheckBorders(b, a, 1, false)
	//If at least one piece of ship remains undamaged, the whole ship is still alive
	for y:=y0; y <= y1; y++ {
		for x:=x0; x <= x1; x++ {
			if this.Trace[y * 10 + x] != '1' {
				return false
			}
		}
	}
	//Otherwise, sink the ship (set the state '2' on the hidden map)
	for y:=y0; y <= y1; y++ {
		for x:=x0; x <= x1; x++ {
			this.Trace[y * 10 + x] = '2'
		}
	}
	//One ship less
	this.Left-=1
	return true
}

//Try to hit the ship at provided coordinates
func (this *Sea) HitShip (y0, x0 int) int {
	//Set the game state to 'attacking'
	GameState = 0
	//Show message about attack
	fmt.Println(fmt.Sprint("Hit at (", y0, ",", x0, "):"))
	//Overall score after attack
	score:=0
	//Reorder coordinates
	x:=[]int{-IMin(y0, x0) - 1, IMax(y0, x0) - 1}
	//Repeated shot
	if this.Y == x[0] && this.X == x[1] {
		fmt.Println("Repeated shot (-20 points)")
		score-=20
	} else {
		//Update last shot coordinates
		this.Y, this.X = x[0], x[1]
	}
	//Fire out of the map bounds
	for i:=range(x) {
		if x[i] < 0 || x[i] > 9 {
			fmt.Println("Fire out of bounds (-8 points)")
			return score - 8
		}
	}
	//Check if the shot is wise and reasonable
	score+=this.ReasonableHit(x[0], x[1])
	pos:=x[0] * 10 + x[1]
	//Maybe the state of the position at given coordinates is already known? 
	recognize:=this.Trace[pos]
	//It was empty before
	if recognize == '3' {
		fmt.Println("Shot at knowingly empty field (-11 points)")
		score-=11
	//It's a shot at previously damaged or sunken ship
	} else if recognize == '2' {
		fmt.Println("Shot at dead ship (-10 points)")
		score-=10
	} else if recognize == '1' {
		fmt.Println("Shot at wrecks (-9 points)")
		score-=9
	//Unknown state, no penalty
	} else if recognize == '0' {
		s:=this.Field[pos]
		//We hurt the ship
		if s == '1' {
			//Set mode to 'keep turn', acknowledging the fact of a hit on the 'hidden' map
			GameState = 1
			this.Trace[pos] = '1'
			if this.CheckDestruction(x[0], x[1]) {
				//In case of victory
				if this.Left == 0 {
					score+=21
					fmt.Println("Victory! (+21 points)")
					//Set mode to 'victory'
					GameState = 2
				} else {
					fmt.Println("Ship is destroyed. This player has the next turn. Ships left: ", this.Left, ".")
				}
			} else {
				fmt.Println("Ship is damaged. This player has the next turn.")
			}
		//Missed shot
		} else {
			this.Trace[pos] = '3'
			fmt.Println("Miss.")
		}
	}
	return score
}

//The procedure creating the HTML tables for corresponding maps  
func BringTheTables(name1, name2 string) string {
	
	//Prepare the table units 1-4 & the captions with the names of the players
	u1:=[]string{`<table><caption class="names"><b>`, name1, `</b></caption>`}
	u2:=[]string{`<table><caption class="names"><b>`, name2, `</b></caption>`}
	u3:=[]string{`<br><table>`}
	u4:=[]string{`<table>`}
	
	for i:=0; i < 10; i++ {
		u1=append(u1, "<tr>")
		u2=append(u2, "<tr>")
		u3=append(u3, "<tr>")
		u4=append(u4, "<tr>")
		for j:=0; j < 10; j++ {
			pos:=i * 10 + j
			var player_pos1, player_pos2 string
			
			//Show the ships on the maps...
			if Overseas[0].Field[pos] == '1' {player_pos1 = "1"} else {player_pos1 = "0"}
			if Overseas[1].Field[pos] == '1' {player_pos2 = "1"} else {player_pos2 = "0"}
			u1=append(u1, fmt.Sprint(`<td class="c`, player_pos1, `"></td>`))
			u2=append(u2, fmt.Sprint(`<td class="c`, player_pos2, `"></td>`))
			
			//...as well as the hits performed by the players
			u3=append(u3, fmt.Sprint(`<td class="t`, string(Overseas[0].Trace[pos]), `"></td>`))
			u4=append(u4, fmt.Sprint(`<td class="t`, string(Overseas[1].Trace[pos]), `"></td>`))
			
		}
		u1=append(u1, "</tr>")
		u2=append(u2, "</tr>")
		u3=append(u3, "</tr>")
		u4=append(u4, "</tr>")
	}
	u1=append(u1, "</table>")
	u2=append(u2, "</table>")
	u3=append(u3, "</table>")
	u4=append(u4, "</table>")
	
	//Unite the tables
	return strings.Join([]string{strings.Join(u1, ""), strings.Join(u2, ""), strings.Join(u3, ""), strings.Join(u4, "")}, "")
}

func hit_the_enemy (NameN, NameOther string, IndexN, IndexOther, y0, x0 int) {
	//Fire at the target
	fmt.Print("Turn of the player ", NameN, ". ")
	Overseas[IndexN].Score+=Overseas[IndexOther].HitShip(y0, x0)
	//Disqualification rule
	if Overseas[IndexN].Score < -1000 {
		fmt.Println("The player ", NameOther, " wins as the other player ", NameN, " is disqualified.")
		GameState = 2
	}
}

//[i]x/yN=0..1 - coordinates to place/hit piece of ship
//[i]dy - diary storage
//[i]ne - name of the player

//Handler interactively processing the game by the means of web browser 
func handler (w http.ResponseWriter, r *http.Request) {
	//Set headers
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	//Run certain actions for the POST method
	if(r.Method == "POST") {
		//Not all ships are placed on the sea surface yet 
		if GameState < 0 {
			//Read the POST method values (building the ship)
			y0,_:=strconv.Atoi(r.PostFormValue("y0"))
			x0,_:=strconv.Atoi(r.PostFormValue("x0"))
			y1,_:=strconv.Atoi(r.PostFormValue("y1"))
			x1,_:=strconv.Atoi(r.PostFormValue("x1"))
			diary:=r.PostFormValue("dy")
			if len(diary) > 2048 {diary = string([]rune(diary)[:2048])}
			//Decorate the names and show the messages about the fleet construction
			if(PlayerTurn == 1) {
				if Name1 == "" {
					Name1 = "1{" + r.PostFormValue("ne") + "}"
					fmt.Println("Player ", Name1, " starts building the fleet.")
				}
			} else {
				if Name2 == "" {
					Name2 = "2{" + r.PostFormValue("ne") + "}"
					fmt.Println("Player ", Name2, " starts building the fleet.")
				}
			}
			//Update information and try to apply the received info to build the ship
			Overseas[PlayerTurn - 1].Diary2048 = diary
			Overseas[PlayerTurn - 1].Score+=Overseas[PlayerTurn - 1].SetShip(y0, x0, y1, x1)
			Overseas[PlayerTurn - 1].Redescribe()
			//Update the global game state
			GameState+=1
		} else {
			//No victory yet
			if GameState != 2 {
				//Read the POST method values (firing at the target)
				y0,_:=strconv.Atoi(r.PostFormValue("y0"))
				x0,_:=strconv.Atoi(r.PostFormValue("x0"))
				diary:=r.PostFormValue("dy")
				if len(diary) > 2048 {diary = string([]rune(diary)[:2048])}
				Overseas[PlayerTurn - 1].Diary2048 = diary
				//Force the neutral battle mode and start brawling
				GameState = 0
				if PlayerTurn == 2 {
					hit_the_enemy(Name2, Name1, 1, 0, y0, x0)
				} else {
					hit_the_enemy(Name1, Name2, 0, 1, y0, x0)
				}
				//No hits? Pass the turn to the opponent
				if GameState == 0 {if PlayerTurn == 1 {PlayerTurn = 2} else {PlayerTurn = 1}}
			}
		}
	}
	//The Universal part of the handler function
	var buffer bytes.Buffer
	//No victory yet
	if GameState != 2 {
		//Give the next player chances to put the ships on the sea surface
		if GameState == -10 { PlayerTurn = 2}
		//Load corresponding scripts related to the logic of players
		buffer.WriteString(File2Tag(fmt.Sprint(ScriptName, PlayerTurn, ".js"), "script"))
		//If the battle has not yet begun, create the web form to build the ships
		if GameState < 0 {
			buffer.WriteString(form_part1)
			buffer.WriteString(Overseas[PlayerTurn - 1].FieldV)
			buffer.WriteString(`","`)
			buffer.WriteString(strings.Replace(Overseas[PlayerTurn - 1].Diary2048, `"`, `'`, -1))
			buffer.WriteString(form_part2)
		//Otherwise, create the web form to fire at the ships
		} else {
			var callback_parameters string
			if PlayerTurn == 1 {
				callback_parameters = string(Overseas[1].Trace)
			} else {
				callback_parameters = string(Overseas[0].Trace)
			}
			buffer.WriteString(form_part3)
			buffer.WriteString(callback_parameters)
			buffer.WriteString(`","`)
			buffer.WriteString(strings.Replace(Overseas[PlayerTurn - 1].Diary2048, `"`, `'`, -1))
			buffer.WriteString(form_part4)
		}
	} else {
		if !LastWords {
			//Victory - the last message summarizing the game session results
			LastWords = true
			fmt.Println(fmt.Sprint(Name1, " {", Overseas[0].Score, " points}"), fmt.Sprint(Name2, " {", Overseas[1].Score, " points}"))			
			fmt.Println("The game session is finished.")
			fmt.Println("Diary #1: ", Overseas[0].Diary2048)
			fmt.Println("Diary #2: ", Overseas[1].Diary2048)
			//If the player names do not contain magic '_' symbols at the start of their names, proceed with normal logging procedure
			if []rune(Name1)[2] != '_' && []rune(Name2)[2] != '_' {log.Fatal("Game over.")}
		}
	}
	//Output the whole HTML content
	w.Write([]byte(strings.Join([]string{upperTemplate, BringTheTables(fmt.Sprint(Name1, " {", Overseas[0].Score, " points}"), fmt.Sprint(Name2, " {", Overseas[1].Score, " points}")), "</div>", buffer.String(), "</body></html>"}, "")))
}

func main () {
	//Load the diaries and initialize the structure of the maps
	Overseas[0].Diary2048 = File2Diary(`diary1.txt`)
	Overseas[1].Diary2048 = File2Diary(`diary2.txt`)
	Overseas[0].Redescribe()
	Overseas[1].Redescribe()
	//Launch the server
	addr := "localhost:8080"
	fulladdr:=strings.Join([]string{"http://", addr}, "")
	network_connection, err := net.Listen("tcp", addr)
	if err != nil {
		return
	}
	http.HandleFunc("/", handler)
	fmt.Printf("Starting server at %s ...\n", fulladdr)
	Launch(fulladdr)
	http.Serve(network_connection, nil)
}

//The style definitions and other prelude info
const upperTemplate = `<!DOCTYPE html>
<html>
<head>
<style>
table {
	border-collapse:collapse;
	display:inline-block;
	margin:10px;
}
tr {
	height:30px;
}
td {
	width:30px;
	border:1px solid;
}
.middle {
	position:absolute;
	left:50%;
	top:50%;
	transform:translate(-50%,-50%);
}
.between {
	width:1px;
	display:inline-block;
}
.names {
	font-size:20px;
	margin:7px;
}
.c0 {
	background-color:white;
}
.c1 {
	background-color:royalblue;
}
.t0 {
	background-color:white;
}
.t1 {
	background-color:blue;
}
.t2 {
	background-color:red;
}
.t3 {
	background-color:silver;
}
</style>
</head>
<body>
<div class="middle">`

//The essential parts of the forms
const form_part1=`<form style="display:none" id="f" method="POST">
		<input id="iy0" type="hidden" name="y0">
		<input id="ix0" type="hidden" name="x0">
		<input id="iy1" type="hidden" name="y1">
		<input id="ix1" type="hidden" name="x1">
		<input id="idy" type="hidden" name="dy">
		<input id="ine" type="hidden" name="ne">
		</form>
		<script>var ret=set("`
const form_part2=`");	
		document.getElementById("iy0").value=ret[0];
		document.getElementById("ix0").value=ret[1];
		document.getElementById("iy1").value=ret[2];
		document.getElementById("ix1").value=ret[3];
		document.getElementById("idy").value=ret[4];
		document.getElementById("ine").value=Name;
		document.getElementById("f").submit();
		</script>`
const form_part3=`<form style="display:none" id="f" method="POST">
		<input id="iy0" type="hidden" name="y0">
		<input id="ix0" type="hidden" name="x0">
		<input id="idy" type="hidden" name="dy">
		</form>
		<script>var ret=hit("`
const form_part4=`");
		document.getElementById("iy0").value=ret[0];
		document.getElementById("ix0").value=ret[1];	
		document.getElementById("idy").value=ret[2];	
		document.getElementById("f").submit();
		</script>`