@echo off
gcc -c original\sqlite3.c -o sqlite3.obj
gcc -Wall -shared -Ioriginal -o sqlite3_bind.dll sqlite3_bind.c sqlite3.obj
copy sqlite3_bind.dll ".."