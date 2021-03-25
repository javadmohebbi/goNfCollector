# How to deploy TICK inside docker
1. Create `tick` network
   - ```docker network create --driver bridge tick-net```

2. Create needed direcories for `data` & `config`
   - ```mkdir /opt/nfcollector/vendors/influxdb/data -pv```
   - ```mkdir /opt/nfcollector/vendors/influxdb/config -pv```

3. Get & run the container
   -
    ```
    docker run -d \
    --network tick-net \
    -p 8086:8086 \
    -p 8082:8082 \
    -p 8089:8089 \
    -v /opt/nfcollector/vendors/influxdb/data:/var/lib/influxdb2/engine \
    -v /opt/nfcollector/vendors/influxdb/config:/etc/influxdb2 \
    --name influxdb \
    influxdb:latest
    ```