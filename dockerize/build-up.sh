#!/bin/bash


# PROJECT DIR
PROJECT_DIR=$HOME/oi24/nfcollector

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
INFLUX_DIR=$PROJECT_DIR/INFLUX_DIR

# Grafana DIR
GRAFANA_DIR=$PROJECT_DIR/INFLUX_DIR

# downlaod latest version of binnary amd64
download_latest_version() {
    wget -O $PROJECT_DIR/bin/nfcollector https://download.openintelligence24.com/nf/bin/nfcollector
    chmod +x $PROJECT_DIR/bin/nfcollector

    wget -O $PROJECT_DIR/bin/nfupdater https://download.openintelligence24.com/nf/bin/nfupdater
    chmod +x $PROJECT_DIR/bin/nfupdater

    $PROJECT_DIR/bin/nfupdater -ipsum -ip2l -ip2l-asn -ip2l-proxy


}

# prepare needed directories & create if needed
prepare_dir()
{

    # project directory
    mkdir -pv $PROJECT_DIR

    # config files
    mkdir -pv $PROJECT_DIR/etc

    # vendors dir & files
    mkdir -pv $PROJECT_DIR/vendors/i2location
    mkdir -pv $PROJECT_DIR/vendors/i2location/local-db
    mkdir -pv $PROJECT_DIR/vendors/ipsum

    # influxDB directory
    mkdir -pv $INFLUX_DIR

    # grafana directory
    mkdir -pv $GRAFANA_DIR

}


prepare_dir
download_latest_version

