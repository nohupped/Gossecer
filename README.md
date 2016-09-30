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
# The key name doesn't matter. The value must be a valid regex.
[filters]
ip = (?:[0-9]{1,3}\.){3}[0-9]{1,3}
#datetime field that looks like Oct  1 03:29:36
datetime = [A-Z]{1}[a-z]{2}?\s+\d?\s+(\d{2}\:){2}\d{2}

#expire sets individual expire times for rule id
[expire]

```