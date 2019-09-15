package sql3

/*
 #include <stdlib.h>
 #cgo CFLAGS: -Ioriginal
 #cgo LDFLAGS: -L. -l"sqlite3_bind"
 #include "sqlite3_bind.h"
 */
import "C"
import (
	"reflect"
	"unsafe"
)

func Linked () bool {
	return bool(C.linked())
}

func Db (path string) bool {
	convert:=C.CString(path)
	defer C.free(unsafe.Pointer(convert))
	return bool(C.db(convert))
}

func Sql (sql string) {
	convert:=C.CString(sql)
	defer C.free(unsafe.Pointer(convert))
	C.sql(convert)
}

func Run (sql string) int {
	convert:=C.CString(sql)
	defer C.free(unsafe.Pointer(convert))
	return int(C.run(convert))
}

func Ok (result int) bool {
	return bool(C.ok(C.int(result)))
}

func Row (result int) bool {
	return bool(C.row(C.int(result)))
}

func End (result int) bool {
	return bool(C.end(C.int(result)))
}

func Bind (number int, sequence []byte) {
	C.bind(C.int(number), unsafe.Pointer(&sequence[0]), C.int(len(sequence)))
}

func IBind (number int, value int) {
	C.ibind(C.int(number), C.int(value))
}

func FBind (number int, value float64) {
	C.fbind(C.int(number), C.double(value))
}

func Unbind (number int, length int) []byte {
	header:=&reflect.SliceHeader {
		Data: uintptr(C.unbind(C.int(number))),
		Len: length,
		Cap: length,
	}
	return *(*[]byte)(unsafe.Pointer(header))
}

func IUnbind (number int) int {
	return int(C.iunbind(C.int(number)))
}

func FUnbind (number int) float64 {
	return float64(C.funbind(C.int(number)))
}

func DbLen (number int) int {
	return int(C.len(C.int(number)))
}

func DbErr () int {
	return int(C.err())
}
