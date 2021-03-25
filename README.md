# Go Netflow Collector
Is a go module that collect netflows from version 1 to 9 and also IPFIX.

It can export it to many other services like InfluxDB



# Test Netflow Dummy data
```docker run -it --rm networkstatic/nflow-generator -t 192.168.43.116 -p 6859```