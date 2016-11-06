# Gossecer
Ossec syslog aggregator written in go

####Requires redis installed. Find more in the Wiki page.

####Sample `ini` file:

```
[ossec]
ConfFile = /home/girishg/ossec.conf
[redis]
Server = localhost
Port = 6379

# Any number of filters can be added under the below section.
# The key name doesn't matter. The value must be a valid regex, and
# the matching pattern will be removed.
[filters]
ip = (?:[0-9]{1,3}\.){3}[0-9]{1,3}
#datetime field that looks like Oct  1 03:29:36
datetime = [A-Z]{1}[a-z]{2}?\s+\d?\s+(\d{2}\:){2}\d{2}
port = port?\s+[0-9]+\s
tag = [a-zA-Z]+\[[0-9]+\]

#expire sets individual expire times for alerts from individual rule id. Defaults to 300 seconds.
[expire]
#RuleID = expire ttl(in seconds)

5501 = 600
5402 = 600
5502 = 600
5710 = 600

# Set individual threshold based on rule id. Defaults to 10 numbers.
# Eg: If you want an alert for rule 5501 when the threshold crosses 10 in the last 10 minutes,
# set [expire] to 600 for 5501, and [threshold] to 10 for 5501.

[threshold]

5402 = 10
5501 = 10
5710 = 5

[alert]
host = localhost
port = 8888


```

####Sample Output as read with netcat:
```
$> nc -ul -p 8888

{"EventOccurance":16,"EventThreshold":15,"Hostname":"myserver-\u003e127.0.0.1","Message":"Oct 13 12:59:15 myserver sshd[8918]: Failed password for invalid user a from 127.0.0.1 port 60293 ssh2","RuleID":5712,"Syslogcrit":10,"TimesAlerted":1,"TotalEventOccurance":16}


```

Consuming script can use a modulus on ```TimesAlerted``` value to decide on the frequency of triggering alerts.
