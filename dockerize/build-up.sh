#!/bin/bash -i


# PROJECT DIR
export PROJECT_DIR=$HOME/oi24/nfcollector

export NFC_LISTEN_ADDRESS="0.0.0.0"
export NFC_LISTEN_PORT="6859"
export NFC_INFLUXDB_HOST="influxdb"
export NFC_INFLUXDB_PORT="8086"
export NFC_INFLUXDB_TOKEN="TOKEN"
export NFC_INFLUXDB_BUCKET="nfCollector"
export NFC_INFLUXDB_ORG="OPENINTELLIGENCE"
export NFC_IP_REPTATION_IPSUM="$PROJECT_DIR/vendors/ipsum/ipsum.txt"
export NFC_IP2L_ASN="$PROJECT_DIR/vendors/ip2location/db/IP2LOCATION-LITE-ASN.IPV6.CSV/IP2LOCATION-LITE-ASN.IPV6.CSV"
export NFC_IP2L_IP="$PROJECT_DIR/vendors/ip2location/db/IP2LOCATION-LITE-DB11.IPV6.BIN/IP2LOCATION-LITE-DB11.IPV6.BIN"
export NFC_IP2L_PROXY="$PROJECT_DIR/vendors/ip2location/db/IP2PROXY-LITE-PX10.IPV6.CSV/IP2PROXY-LITE-PX10.IPV6.CSV"
export NFC_IP2L_LOCAL="$PROJECT_DIR/vendors/ip2location/local-db/local.csv"


# InfluxDB DIR
export INFLUX_DIR=$PROJECT_DIR/vendors/influxdb

# Grafana DIR
export GRAFANA_DIR=$PROJECT_DIR/vendors/grafana

# downlaod latest version of binnary amd64
download_latest_version() {

    echo ""
    echo "Downloading required files..."

    # download nfcollector
    wget -O $PROJECT_DIR/bin/nfcollector https://download.openintelligence24.com/nf/bin/nfcollector
    chmod +x $PROJECT_DIR/bin/nfcollector

    # download nfupdater
    wget -O $PROJECT_DIR/bin/nfupdater https://download.openintelligence24.com/nf/bin/nfupdater
    chmod +x $PROJECT_DIR/bin/nfupdater

    # download databases
    $PROJECT_DIR/bin/nfupdater -ipsum -ip2l -ip2l-asn -ip2l-proxy

    # download grafana dashboards & conf
    wget -O /tmp/grafana.tar.gz https://download.openintelligence24.com/vendors/grafana/grafana.tar.gz && tar -vxf /tmp/grafana.tar.gz -C $GRAFANA_DIR/

    wget -O ./docker-compose.yml https://download.openintelligence24.com/vendors/docker-compose/docker-compose.yml

    echo "...done!"
}

# prepare needed directories & create if needed
prepare_dir()
{

    echo ""
    echo "Creating needed directories..."

    # project directory
    mkdir -pv $PROJECT_DIR

    # config files
    mkdir -pv $PROJECT_DIR/etc

    # ibn
    mkdir -pv $PROJECT_DIR/bin

    # vendors dir & files
    mkdir -pv $PROJECT_DIR/vendors/ip2location
    mkdir -pv $PROJECT_DIR/vendors/ip2location/local-db
    mkdir -pv $PROJECT_DIR/vendors/ipsum

    # influxDB directory
    mkdir -pv $INFLUX_DIR

    # grafana directory
    mkdir -pv $GRAFANA_DIR

    echo "...done!"
}



# initializing influx db using a temporary
# influxDB container
get_influx_db_info() {

    echo ""
    echo "Renaming old configuraions..."
    unixnan=$(date +%s)
    mv $INFLUX_DIR/engine $INFLUX_DIR/engine.old.$unixnan -f > /dev/null 2>&1
    mv $INFLUX_DIR/influxd.bolt $INFLUX_DIR/influxd.bolt.old.$unixnan > /dev/null 2>&1
    echo "...done!"

    CONTAINERID=influxdb_tmp
    echo ""
    echo "Stop & remove container $CONTAINERID, if available..."
    docker stop $CONTAINERID > /dev/null 2>&1
    docker container rm $CONTAINERID > /dev/null 2>&1

    # docker network create --driver bridge tick-graf > /dev/null 2>&1

    echo "...done!"


    echo ""
    echo "Staring temporary InfluxDB container ($CONTAINERID)..."
    # create influxdb tmp image
    docker run -d \
    -v $INFLUX_DIR:/var/lib/influxdb2 \
    --name $CONTAINERID \
    influxdb:latest  > /dev/null 2>&1
    echo "...done!"

    NFC_INFLUXDB_TOKEN=`echo $CONTAINERID$(date)$USER | base64 -w 0`
    echo ""
    echo "InfluxDB Token is: $NFC_INFLUXDB_TOKEN"

    # wait until command finished; mean
    echo ""
    echo "Waiting for InfluxDB ($CONTAINERID) to be ready...!"
    docker exec -it $CONTAINERID wget http://localhost:8086 > /dev/null
    echo "...done!"

    echo ""
    echo ""
    echo "Initializing InfluxDB with this command...:"
    COMMAND_TO_RUN="docker exec -t $CONTAINERID influx setup --org $NFC_INFLUXDB_ORG --bucket $NFC_INFLUXDB_BUCKET --retention 7d --username admin --password influx_admin_secret --token $NFC_INFLUXDB_TOKEN --force "
    echo ""
    echo " >>> $COMMAND_TO_RUN"
    echo ""
    echo "...done!"

    echo ""
    echo "------------------------------------------"
    # RUN THE COMMAND
    $COMMAND_TO_RUN
    echo "------------------------------------------"

    echo ""
    echo "Stop & remove container $CONTAINERID"
    docker stop $CONTAINERID > /dev/null 2>&1
    docker container rm $CONTAINERID > /dev/null 2>&1
    echo "...done!"

}







# preparing direcoty we need
prepare_dir

# downloading latest version of needed files
download_latest_version

# config influx db using a temporary container
get_influx_db_info





docker-compose up -d