# GoSpray
 Simple tool to bruteforce (spray actually) different network protocols.
 GoSpray also supports restoration of interrupted tasks ("-restore").
 
 GoSpray currently supports: **rdp, ssh, ftp, Windows LDAP, Windows Kerberos, http basic** and **digest authentication**

**Any requests for support of new protocols are welcomed!**

```
go run . -ul testUsernames.txt -pl testPasswords.txt -p ftp -tl targets.txt -w 10
---------------+
Success: 192.168.56.102:user:123
-------------------
```


```
go run . -ul testUsernames.txt -pl testPasswords.txt -p ftp -tl targets.txt -w 10
--------

CTRL+C

go run . -restore
-------+
Success: user:123
-------------------
```

-ul   Path to file with **usernames**

-ul   Path to file with **passwords**

-p   Protocol to brute ( winldap, rdp, ssh, ftp, httpbasic, httpdigest )

-tl   Path to file with **targets** (one target per line ex: "http://127.0.0.1:667/protected/folder/")

-w   Number of workers (threads)

-debug Run with debug information (for now only httpbasic) 

-restore use "progress.gob" to restore task
 
 
**Target formats:**

```
192.168.56.102 - for **ssh, rdp, ldap, ftp**

192.168.56.102:21 - for **ssh, rdp, ldap, ftp**

http://192.168.56.102:80/2 - for **basic** and **digest authentication**

test.local:88 - for **Windows Kerberos**

``` 
 
**Examples:**

```
spray.exe -ul testUsernames.txt -pl testPasswords.txt -p ssh -tl targets.txt -w 10

spray.exe -ul testUsernames.txt -pl testPasswords.txt -p httpbasic -tl targets.txt -w 10 -ru -rp

spray.exe -restore
```

