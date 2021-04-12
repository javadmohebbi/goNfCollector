#!/bin/bash

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color


# PROJECT DIR
export PROJECT_DIR=$HOME/oi24/nfcollector

# export PROJECT_DIR=/opt/openintelligence24/nfcollector


export USER_ID=$(id -u)

export NFC_LISTEN_ADDRESS="0.0.0.0"

# total number of used cpu
export NFC_CPU_NUM="0"


export NFC_LISTEN_PORT="6859"
export NFC_INFLUXDB_HOST="127.0.0.1"
export NFC_INFLUXDB_PORT="8086"
export NFC_INFLUXDB_TOKEN="5vqt0q0b4g_lZwNgp7-8GgPq5Nxf3YY37xbVZP_ypeK_G3dwdNlTrAkcKN_Q6QzbmG-Th96lT_65Kp0j2UD1HA=="
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
    echo -e "${YELLOW} Downloading required files...${NC}"

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
    # chmod +r $GRAFANA_DIR/ -Rv
    echo -e "${YELLOW} Chaning permissions & owner of grafana directory. Maybe you need to enter 'sudo' password ${NC}"
    sudo chmod -Rv a+w $GRAFANA_DIR/
    sudo chown $USER:root $GRAFANA_DIR/ -Rv

    wget -O ./docker-compose.yml https://download.openintelligence24.com/vendors/docker-compose/docker-compose.yml


    # download templates
    # wget -O $PROJECT_DIR/etc/collector.yml https://download.openintelligence24.com/nf/etc/nfcol-bash.yml
    # wget -O $PROJECT_DIR/etc/ip2location.yml https://download.openintelligence24.com/nf/etc/nfloc-bash.yml

    echo -e "${GREEN}...done!${NC}"
}

# prepare needed directories & create if needed
prepare_dir()
{

    echo ""
    echo -e "${YELLOW} Creating needed directories...${NC}"

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

    echo -e "${GRREN}...done!${NC}"
}



# initializing influx db using a temporary
# influxDB container
get_influx_db_info() {

    echo ""
    echo -e "${YELLOW} Renaming old configuraions...${NC}"
    unixnan=$(date +%s)
    mv $INFLUX_DIR/engine $INFLUX_DIR/engine.old.$unixnan -f > /dev/null 2>&1
    mv $INFLUX_DIR/influxd.bolt $INFLUX_DIR/influxd.bolt.old.$unixnan > /dev/null 2>&1
    echo -e "${GREEN}...done!${NC}"

    CONTAINERID=influxdb_tmp
    echo ""
    echo "Stop & remove container $CONTAINERID, if available..."
    docker stop $CONTAINERID > /dev/null 2>&1
    docker container rm $CONTAINERID > /dev/null 2>&1

    docker network create --driver bridge tick-graf > /dev/null 2>&1

    echo -e "${GREEN}...done!${NC}"


    echo ""
    echo -e "${YELLOW} Staring temporary InfluxDB container ($CONTAINERID)...${NC}"
    # create influxdb tmp image
    docker run -d \
    -v $INFLUX_DIR:/var/lib/influxdb2 \
    --name $CONTAINERID \
    influxdb:latest  > /dev/null 2>&1
    echo -e "${GREEN}...done!${NC}"

    # NFC_INFLUXDB_TOKEN=`echo $CONTAINERID$(date)$USER | base64 -w 0`
    echo ""
    echo "InfluxDB Token is: $NFC_INFLUXDB_TOKEN"

    # wait until command finished; mean
    echo ""
    echo -e "${YELLOW} Waiting for InfluxDB (${CONTAINERID}) to be ready...!${NC}"
    docker exec -it $CONTAINERID wget http://localhost:8086 > /dev/null
    echo -e "${GREEN}...done!${NC}"

    echo ""
    echo ""
    echo -e "${YELLOW} Initializing InfluxDB with this command...:${NC}"
    COMMAND_TO_RUN="docker exec -t $CONTAINERID influx setup --org $NFC_INFLUXDB_ORG --bucket $NFC_INFLUXDB_BUCKET --retention 7d --username admin --password influx_admin_secret --token $NFC_INFLUXDB_TOKEN --force "
    echo ""
    echo -e " >>> command to run: ${YELLOW} $COMMAND_TO_RUN ${NC}"
    echo ""
    echo -e "${GREEN}...done!${NC}"

    echo ""
    echo "------------------------------------------"
    # RUN THE COMMAND
    $COMMAND_TO_RUN
    echo "------------------------------------------"

    echo ""
    echo -e "${YELLOW} Stop & remove container $CONTAINERID ${NC}"
    docker stop $CONTAINERID > /dev/null 2>&1
    docker container rm $CONTAINERID > /dev/null 2>&1
    echo -e "${GREEN}...done!${NC}"

}



print_info(){

    echo -e "${NC}Information you need to know:"
    echo -e "\n\t${GREEN} Project Directory: ${YELLOW}${PROJECT_DIR}"

    echo -e "\n\t${GREEN} InfluxDB:"
    echo -e "\t\t${NC} address:${YELLOW}${NFC_INFLUXDB_HOST}:${NFC_INFLUXDB_PORT}"
    echo -e "\t\t${NC} token:${YELLOW}${NFC_INFLUXDB_TOKEN}"
    echo -e "\t\t${NC} Web UI Credentials:"
    echo -e "\t\t\t${NC} username: ${YELLOW}admin"
    echo -e "\t\t\t${NC} password: ${YELLOW}influx_admin_secret"
    echo -e "\n\t${GREEN} Grafana:"
    echo -e "\t\t${NC} address:${YELLOW}127.0.0.1:3000"
    echo -e "\t\t${NC} Web UI Credentials:"
    echo -e "\t\t\t${NC} username: ${YELLOW}admin"
    echo -e "\t\t\t${NC} password: ${YELLOW}secret"

    echo -e "\n\t${GREEN} nfcollector:"
    echo -e "\t\t${NC} address:${YELLOW}${NFC_LISTEN_ADDRESS}:${NFC_LISTEN_PORT}(udp)"


    echo -e "- - - - - - - - - - - - - - - - "
    echo -e " ${YELLOW}To start containers:"
    echo -e " ${GREEN}cd ${PROJECT_DIR} && docker-compose up -d${NC}"
    echo -e " ${YELLOW}To stop containers:"
    echo -e " ${GREEN}cd ${PROJECT_DIR} && docker-compose down${NC}"
    echo -e "- - - - - - - - - - - - - - - - "

    echo -e "${NC}"
}


replace_compose_template() {
    echo ""
    echo -e "${YELLOW} Preparing docker-compose.yml file...${NC}"

    PWD_ESCP=$(echo $INFLUX_DIR | sed 's_/_\\/_g')
    sed -i "s/_INFLUX_DIR_/$PWD_ESCP/g" ./docker-compose.yml

    PWD_ESCP=$(echo $GRAFANA_DIR | sed 's_/_\\/_g')
    sed -i "s/_GRAFANA_DIR_/$PWD_ESCP/g" ./docker-compose.yml

    PWD_ESCP=$(echo $PROJECT_DIR | sed 's_/_\\/_g')
    sed -i "s/_PROJECT_DIR_/$PWD_ESCP/g" ./docker-compose.yml

    sed -i "s/_NFC_CPU_NUM_/$NFC_CPU_NUM/g" ./docker-compose.yml

    sed -i "s/_NFC_LISTEN_ADDRESS_/$NFC_LISTEN_ADDRESS/g" ./docker-compose.yml
    sed -i "s/_NFC_LISTEN_PORT_/$NFC_LISTEN_PORT/g" ./docker-compose.yml
    sed -i "s/_NFC_INFLUXDB_HOST_/$NFC_INFLUXDB_HOST/g" ./docker-compose.yml
    sed -i "s/_NFC_INFLUXDB_PORT_/$NFC_INFLUXDB_PORT/g" ./docker-compose.yml
    sed -i "s/_NFC_INFLUXDB_TOKEN_/$NFC_INFLUXDB_TOKEN/g" ./docker-compose.yml
    sed -i "s/_NFC_INFLUXDB_BUCKET_/$NFC_INFLUXDB_BUCKET/g" ./docker-compose.yml
    sed -i "s/_NFC_INFLUXDB_ORG_/$NFC_INFLUXDB_ORG/g" ./docker-compose.yml


    mv ./docker-compose.yml $PROJECT_DIR/docker-compose.yml -v

    echo -e "${GRREN}...done!${NC}"
}



# preparing direcoty we need
prepare_dir

# downloading latest version of needed files
download_latest_version

# config influx db using a temporary container
get_influx_db_info


#rename compose
replace_compose_template



# print info
print_info

