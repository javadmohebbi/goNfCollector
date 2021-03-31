# Go Netflow Collector
Is a go module that collect netflows from version 1 to 9 and also IPFIX.

It can export it to many other services like InfluxDB



# Test Netflow Dummy data
```docker run -it --rm networkstatic/nflow-generator -t 192.168.43.116 -p 6859```



# usage example
```
NFC_DEBUG="true" NFC_LISTEN_ADDRESS="0.0.0.0" NFC_LISTEN_PORT="6859" NFC_INFLUXDB_HOST="127.0.0.1" NFC_INFLUX_PORT="8086" NFC_INFLUXDB_TOKEN="JAD5kZ0n3GAQ3jdpe17NT5_NUg73GOvdjZjhxOMiJMx1cZyvLz-4DuR7K8xyRGlPcNQXLqrUTY20lWqbRiK--w==" NFC_INFLUXDB_BUCKET="nfCollector" NFC_INFLUXDB_ORG="MJMOHEBBI" NFC_IP_REPTATION_IPSUM="/opt/nfcollector/vendors/ipsum/ipsum.txt"  NFC_IP2L_ASN="/opt/nfcollector/vendors/ip2location/db/IP2LOCATION-LITE-ASN.IPV6.CSV/IP2LOCATION-LITE-ASN.IPV6.CSV" NFC_IP2L_IP="/opt/nfcollector/vendors/ip2location/db/IP2LOCATION-LITE-DB11.IPV6.BIN/IP2LOCATION-LITE-DB11.IPV6.BIN" NFC_IP2L_PROXY="/opt/nfcollector/vendors/ip2location/db/IP2PROXY-LITE-PX10.IPV6.CSV/IP2PROXY-LITE-PX10.IPV6.CSV" NFC_IP2L_LOCAL="/opt/nfcollector/vendors/ip2location/local-db/local.csv" go run cmd/collector/main.go
```