# icalm
Ip - CIDR Annotation Lookup Microservice

A microservice to lookup annotations for IP Addresses based on CIDR mappings


Annotations can be loaded from a cvs file. Such cvs file has 2 columns: *NETWORK*,*ANNOTATION*

### Example CSV

```
192.168.0.0/16,Annotation for RFC1918 Network (192.168.0.0/16)
10.0.0.0/8,Annotation for RFC1918 Network (10.0.0.0/8)
172.16.0.0/12,Annotation for RFC1918 Network (172.16.0.0/12)
192.0.2.0/24,(TEST-NET-1)
198.51.100.0/24,(TEST-NET-2)
203.0.113.0/24,(TEST-NET-3)
2001:DB8::0/32,v6 Example
ff00::/8,v6 Multicast
```

Lookups can be done via the line-protocol or http

## line-protocol
The line-protocol is the easiest way to query icalm. It used via a TCP or UNIX socket.

 - A client just sends an IP-Address terminated with a NEWLINE.
 - icalm replies with the annotation of the matching network.
 - If there is no hit in icalms lookup table, the response is an empty line.


### Example usage
#### Run server
```
# bin/icalm-server -networks test.csv -line-listen 127.0.0.1:4226
IP: CIDR annotation lookup microservice
2024/05/21 18:15:04 Loaded lookuptable with 6 entries
2024/05/21 18:15:04 Listening on 127.0.0.1:4226 for lineproto requests
2024/05/21 18:19:35 Shutting down
```

#### simulate client with netcat
```
# nc localhost 4226
192.168.1.1
Annotation for RFC1918 Network (192.168.0.0/16)
10.2.34.5
Annotation for RFC1918 Network (10.0.0.0/8)
```


## http protcol
Lookups via http api are currently *not* implemented. :(

