#include "sqlite3.h"

#define NULL 0
#define false 0 
#define true 1

sqlite3 *_db = NULL;
sqlite3_stmt* ppStmt = NULL;

extern _Bool linked () {
	return _db != NULL;
}

extern _Bool db (const char *path) {
	if (linked()) {
		sqlite3_close(_db);
		_db = NULL;
	}
	return (*path) ? !sqlite3_open(path, &_db) : true;
}

extern void sql (const char* sql) {
	sqlite3_finalize(ppStmt);
	ppStmt = NULL;
	if (*sql) sqlite3_prepare_v2(_db, sql, -1, &ppStmt, NULL);
}

extern int run (const char* sql) {
	return (*sql) ? sqlite3_exec(_db, sql, NULL, NULL, NULL) : sqlite3_step(ppStmt);
}

extern _Bool ok (int result) {
	return result == SQLITE_OK;
}

extern _Bool row (int result) {
	return result == SQLITE_ROW;
}

extern _Bool end (int result) {
	return result == SQLITE_DONE;
}

extern int bind (int n, const void* bytes, int length) {
	return sqlite3_bind_blob(ppStmt, n, bytes, length, SQLITE_STATIC);
}

extern int ibind (int n, int value) {
	return sqlite3_bind_int(ppStmt, n, value);
}

extern int fbind (int n, double value) {
	return sqlite3_bind_double(ppStmt, n, value);
}

extern const void* unbind (int n) {
	return sqlite3_column_blob(ppStmt, n);
}

extern int iunbind (int n) {
	return sqlite3_column_int(ppStmt, n);
}

extern double funbind (int n) {
	return sqlite3_column_double(ppStmt, n);
}

extern int len (int n) {
	return n<0 ? sqlite3_column_count(ppStmt) : sqlite3_column_bytes(ppStmt, n);
}

extern int err () {
	return sqlite3_errcode(_db);
}
