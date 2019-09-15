#ifndef INC_SQLITE3_BIND
#define INC_SQLITE3_BIND
extern _Bool linked ();
extern _Bool db (const char *);
extern void sql (const char*);
extern int run (const char*);
extern _Bool ok (int);
extern _Bool row (int);
extern _Bool end (int);
extern int bind (int, const void*, int);
extern int ibind (int, int);
extern int fbind (int, double);
extern const void* unbind (int);
extern int iunbind (int);
extern double funbind (int);
extern int len (int);
extern int err ();
#endif
