# Gossecer
Ossec syslog aggregator written in go

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
{"EventOccurance":363,"EventThreshold":15,"Hostname":"SERVER1-\u003e127.0.0.1","Message":"Oct  7 10:27:00 SERVER1 sshd[1960]: Accepted password for SOMEUSER from 10.10.10.10 port 21064 ssh2","RuleID":5715,"Syslogcrit":3,"TimesAlerted":349}
{"EventOccurance":2633,"EventThreshold":10,"Hostname":"SERVER-1-\u003e127.0.0.1","Message":"Oct  7 10:27:01 SERVER-1 CROND[24743]: pam_unix(cron:session): session opened for user root by (uid=0)","RuleID":5501,"Syslogcrit":3,"TimesAlerted":2624}
```

Consuming script can use a modulus on ```TimesAlerted``` value to decide on the frequency of triggering alerts.