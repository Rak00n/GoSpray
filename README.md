# GoSpray
 Simple tool to bruteforce (spray actually) different network protocols.
 
 GoSpray currently supports: **ssh, ftp, http basic** and **digest authentication**

```
go run . -ul testUsernames.txt -pl testPasswords.txt -p ftp -t 192.168.56.102:21 -w 10
---------------+
Success: user:123
-------------------
```

-ul   Path to file with **usernames**

-ul   Path to file with **passwords**

-p   Protocol to brute ( ssh, ftp, httpbasic, httpdigest )

-t   Target host. http://127.0.0.1:667/protected/folder/

-w   Number of workers (threads)

