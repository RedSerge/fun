package main

import (
	"fmt"
	"./sqlite3"
)

func main () {
	//Connect to database
	sql3.Db("Check.db")
	//Run one SQL sentence
	sql3.Run("CREATE TABLE IF NOT EXISTS t (id INTEGER PRIMARY KEY, cnt BLOB);")
	//Prepare SQL statement
	sql3.Sql("INSERT INTO t(cnt) VALUES (?1)")
	//Bind the value
	s:=[]byte("Hello, 世界");
	sql3.Bind(1, s)
	//Apply the statement
	sql3.Run("")
	//Finalize the operation
	sql3.Sql("")
	//SELECT operator
	sql3.Sql("SELECT * FROM t")
	//Loop, row by row
	for {
		//One row from result
		r:=sql3.Run("")
		//Stop if it is last row
		if sql3.End(r) {
			break
		}
		//Get the id
		i:=sql3.IUnbind(0)
		//Get the length
		l:=sql3.DbLen(1)
		//Convert value to string
		x:=string(sql3.Unbind(1,l))
		//Print id, length, value
		fmt.Println(i,l,x)
	}
	//Finalize
	sql3.Sql("")
	//Close database
	sql3.Db("")
}