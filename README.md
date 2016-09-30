# Gossecer
Ossec syslog aggregator written in go

####Sample `ini` file:

`[ossec]
 ConfFile = /home/girishg/ossec.conf
 [redis]
 Server = localhost
 Port = 6379
 [filters]
 ip = ^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$
 
 #expire sets individual expire times for rule id
 [expire]`