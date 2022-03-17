# GoSpray
 Tool to bruteforce different network protocols.
 
 Supports: ssh, ftp, http basic auth
 
>go run . -ul testUsernames.txt -pl testPasswords.txt -p ftp -t 192.168.56.102:21 -w 10
>
>---------------+
>Success: user:123
>-------------------

Protocols (-p):

  ssh
  ftp
  httpbasic

