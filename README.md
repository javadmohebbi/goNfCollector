<!-- # Go Netflow Collector
Is a go module that collect netflows from version 1 to 9 and also IPFIX.

It can export it to many other services like InfluxDB



# Test Netflow Dummy data
```docker run -it --rm networkstatic/nflow-generator -t 192.168.43.116 -p 6859```



# usage example
```
NFC_DEBUG="true" NFC_LISTEN_ADDRESS="0.0.0.0" NFC_LISTEN_PORT="6859" NFC_INFLUXDB_HOST="127.0.0.1" NFC_INFLUX_PORT="8086" NFC_INFLUXDB_TOKEN="JAD5kZ0n3GAQ3jdpe17NT5_NUg73GOvdjZjhxOMiJMx1cZyvLz-4DuR7K8xyRGlPcNQXLqrUTY20lWqbRiK--w==" NFC_INFLUXDB_BUCKET="nfCollector" NFC_INFLUXDB_ORG="MJMOHEBBI" NFC_IP_REPTATION_IPSUM="/opt/nfcollector/vendors/ipsum/ipsum.txt"  NFC_IP2L_ASN="/opt/nfcollector/vendors/ip2location/db/IP2LOCATION-LITE-ASN.IPV6.CSV/IP2LOCATION-LITE-ASN.IPV6.CSV" NFC_IP2L_IP="/opt/nfcollector/vendors/ip2location/db/IP2LOCATION-LITE-DB11.IPV6.BIN/IP2LOCATION-LITE-DB11.IPV6.BIN" NFC_IP2L_PROXY="/opt/nfcollector/vendors/ip2location/db/IP2PROXY-LITE-PX10.IPV6.CSV/IP2PROXY-LITE-PX10.IPV6.CSV" NFC_IP2L_LOCAL="/opt/nfcollector/vendors/ip2location/local-db/local.csv" go run cmd/collector/main.go
```


#remove
```
docker exec -it influxdb influx delete --org MJMOHEBBI --bucket nfCollector --start '2021-03-01T00:00:00.00Z' --stop '2021-05-29T00:00:00.00Z' --token VL-OzGDlxHlPjMUJM9nQeTWDQ5vcChicnXkVl_vowLud631Exc_seL62sLjq_9Pj5I5KO0i-5EfFdcspElV63A==
```



# build image

```
docker build --pull --rm -f "DockerFile" -t gonfcollector:beta "."
``` -->


# Go Netflow Collector (goNfCollector)
This repo will help you collect **Netflow** (version 1,5,6,7,9 and IPFIX) from network devices. It stores all the required information needed for further analysis in **InfluxDB** and visualize them using **Grafna**.

Currently we are using InfluxDB v2+ for stroring data. If You need older version, you can see [this repository](https://github.com/javadmohebbi/nfCollector).


### Features
- [X] **Supports almost all Netflow versions**: In order to decode Netflow we are using [tehmaze](https://github.com/tehmaze/netflow) go module. this module supports netflow version 1,5,6,7,9 & IPFIX.
- [X] **Container ready**: Just run a simple shell script to prepare your environment & run the containerized netflow collector
- [X] **IP Reputation check**: Check source & destination IPs for the reputation & potential threats.
    - [X] Currently we are using **IPSum** from [this repo](https://github.com/stamparm/ipsum)
    - [ ] [OpenIntelligence24.com](https://openIntelligence24.com) will be available soon. this will be a community based intelligence for checking IP, domains, ... reputatition.
- [ ] Machine Learning models & techniques to find threats like DDoS attacks through packet meta data
- [X] Get Geo Locations using IP2Location free lite database (IPv4 & IPv6)
- [X] Fetch AS Numebr & Name if possible from IP
- [X] Fetch Domain Name from IP if Possible (using **PTR** record)
- [ ] Define multiple data exporter:
  - [X] InfluxDB
  - [ ] Splunk (CEF)
  - [ ] Zabbix


# Quick Start
There are multiple ways to deploy "**netflow collector**" app & easiest ways is **all-in-one** deployment. This method will run `influxdb`, `grafana` & `gonfcollector` docker container using a shell script. No more further configuration are needed & everythings will be downloaded/configured using a `shell script`.
1. Downlaod the latest version:
  `wget https://download.openintelligence24.com/latest.sh`
2. Make this shell script executable
  `chmod +x latest.sh`
3. Run the downloaded shellscript.
  `./latest.sh`
   - You might be asked to enter your user's password during the execution.
   - At the end, it will let you know how to run the container.
   - **REQUIREMENTS**: `docker`, `docker-compose`, `wget` are required!

