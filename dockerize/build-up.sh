#!/bin/bash

# This shell-script will install the latest version of
# nfcollector to your system & has some system requirements including
# docker, docker-compose, wget. If not satisfied, will help you install
# them using a guide printed in the terminal
# Check more about it's usage at: https://github.com/javadmohebbi/goNfCollector




# change the color of echo output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color


# RANDOM QUERY STRING FOR PREVENTING CACHE
RAND_STR=$(date | md5sum | cut -d" " -f1)


## This VARIABLES are exported because of
## nfupdater commad. This command needs some of them
## Don't panic, this env vars are temparary & will no longer exist if
## you re-open the Terminal
## ;-)


# PROJECT DIR
export PROJECT_DIR=$HOME/oi24/nfcollector

# USER ID
export USER_ID=$(id -u)

# Listen address
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
    wget -O $PROJECT_DIR/bin/nfcollector https://download.openintelligence24.com/nf/bin/nfcollector?rnd=$RAND_STR
    chmod +x $PROJECT_DIR/bin/nfcollector

    # download nfupdater
    wget -O $PROJECT_DIR/bin/nfupdater https://download.openintelligence24.com/nf/bin/nfupdater?rnd=$RAND_STR
    chmod +x $PROJECT_DIR/bin/nfupdater

    # download databases
    $PROJECT_DIR/bin/nfupdater -ipsum -ip2l -ip2l-asn -ip2l-proxy

    # download grafana dashboards & conf
    wget -O /tmp/grafana.tar.gz https://download.openintelligence24.com/vendors/grafana/grafana.tar.gz?rnd=$RAND_STR && tar -vxf /tmp/grafana.tar.gz -C $GRAFANA_DIR/
    # chmod +r $GRAFANA_DIR/ -Rv
    echo -e "${YELLOW} Chaning permissions & owner of grafana directory. Maybe you need to enter 'sudo' password ${NC}"
    sudo chmod -Rv a+w $GRAFANA_DIR/
    sudo chown $USER:root $GRAFANA_DIR/ -Rv

    wget -O ./docker-compose.yml https://download.openintelligence24.com/vendors/docker-compose/docker-compose.yml?rnd=$RAND_STR


    # DOWNLAOD TEMPLATES
    wget -O $PROJECT_DIR/etc/collector.yml https://download.openintelligence24.com/nf/etc/nfcol-bash.yml?rnd=$RAND_STR
    wget -O $PROJECT_DIR/etc/location.yml https://download.openintelligence24.com/nf/etc/nfloc-bash.yml?rnd=$RAND_STR


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

    echo -e "${GREEN}...done!${NC}"
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
    echo ""
    echo -e "${NC}= = = = = = = = = = = = = = = = = = = = = = = = = = = = = = "
    echo ""
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


    echo ""
    echo -e "${YELLOW}- - - - - - - - - - - - - - - - "
    echo ""
    echo -e " ${YELLOW}To start containers:"
    echo -e " ${NC}\$ ${GREEN}cd ${PROJECT_DIR} && docker-compose up -d${NC}"
    echo -e " ${YELLOW}To stop containers:"
    echo -e " ${NC}\$ ${GREEN}cd ${PROJECT_DIR} && docker-compose down${NC}"
    echo ""
    echo -e "${YELLOW}- - - - - - - - - - - - - - - - "
    echo ""
    echo ""

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



replace_collector_template() {

    conf=${PROJECT_DIR}/etc/collector.yml

    echo ""
    echo -e "${YELLOW} Preparing ${conf} file...${NC}"

    PWD_ESCP=$(echo $INFLUX_DIR | sed 's_/_\\/_g')
    sed -i "s/_INFLUX_DIR_/$PWD_ESCP/g"  $conf

    PWD_ESCP=$(echo $GRAFANA_DIR | sed 's_/_\\/_g')
    sed -i "s/_GRAFANA_DIR_/$PWD_ESCP/g"  $conf

    PWD_ESCP=$(echo $PROJECT_DIR | sed 's_/_\\/_g')
    sed -i "s/_PROJECT_DIR_/$PWD_ESCP/g"  $conf

    sed -i "s/_NFC_CPU_NUM_/$NFC_CPU_NUM/g"  $conf

    sed -i "s/_NFC_LISTEN_ADDRESS_/$NFC_LISTEN_ADDRESS/g"  $conf
    sed -i "s/_NFC_LISTEN_PORT_/$NFC_LISTEN_PORT/g" . $conf
    sed -i "s/_NFC_INFLUXDB_HOST_/$NFC_INFLUXDB_HOST/g"  $conf
    sed -i "s/_NFC_INFLUXDB_PORT_/$NFC_INFLUXDB_PORT/g"  $conf
    sed -i "s/_NFC_INFLUXDB_TOKEN_/$NFC_INFLUXDB_TOKEN/g"  $conf
    sed -i "s/_NFC_INFLUXDB_BUCKET_/$NFC_INFLUXDB_BUCKET/g"  $conf
    sed -i "s/_NFC_INFLUXDB_ORG_/$NFC_INFLUXDB_ORG/g" $conf

    echo -e "${GRREN}...done!${NC}"
}


replace_location_template() {

    conf=${PROJECT_DIR}/etc/location.yml

    echo ""
    echo -e "${YELLOW} Preparing ${conf} file...${NC}"

    PWD_ESCP=$(echo $NFC_IP2L_ASN | sed 's_/_\\/_g')
    sed -i "s/_NFC_IP2L_ASN_/$PWD_ESCP/g"  $conf

    PWD_ESCP=$(echo $NFC_IP2L_IP | sed 's_/_\\/_g')
    sed -i "s/_NFC_IP2L_IP_/$PWD_ESCP/g"  $conf

    PWD_ESCP=$(echo $NFC_IP2L_PROXY | sed 's_/_\\/_g')
    sed -i "s/_NFC_IP2L_PROXY_/$PWD_ESCP/g"  $conf

    PWD_ESCP=$(echo $NFC_IP2L_LOCAL | sed 's_/_\\/_g')
    sed -i "s/_NFC_IP2L_LOCAL_/$PWD_ESCP/g"  $conf


    echo -e "${GRREN}...done!${NC}"
}



requirement_check() {

    echo ""
    echo -e "${YELLOW} Cheking system requirements...${NC}"

    #check if  docker command is installed
    commad_path=$(command -v docker)
    if [ -z  "$commad_path" ] ; then
        # OK
        echo -e "${RED} >>> DOCKER is not installed. Please install docker and try again.${NC}"
        echo -e "${RED} >>> Docker installation guide: https://docs.docker.com/engine/install"
        exit 1
    else
        echo -e "${GREEN} Found 'docker' in '$commad_path'"
    fi


    #check if  docker-compose command is installed
    commad_path=$(command -v docker-compose)
    if [ -z  "$commad_path" ] ; then
        # OK
        echo -e "${RED} >>> DOCKER-COMPOSE is not installed. Please install docker-compose and try again.${NC}"
        echo -e "${RED} >>> Docker-Compose installation guide: https://docs.docker.com/compose/install/"
        exit 1
    else
        echo -e "${GREEN} Found 'docker-compose' in '$commad_path'"
    fi


    #check if wget command is installed
    commad_path=$(command -v wget)
    if [ -z  "$commad_path" ] ; then
        # OK
        echo -e "${RED} >>> WGET is not installed. Please install wget and try again.${NC}"
        echo -e "${RED} >>> To install wget:"
        echo -e "${RED} >>> \tDebian/Ubuntu/LinuxMint: ${YELLOW}sudo apt-get install wget"
        echo -e "${RED} >>> \tRedhat/CentOS/Fedora: ${YELLOW}sudo yum install wget"
        echo -e "${NC}"
        exit 1
    else
        echo -e "${GREEN} Found 'wget' in '$commad_path'"
    fi


    echo -e "${YELLOW} ...requirements s atisfied!${NC}"
    echo ""
}


# add cron job
cron_jobs() {
    echo ""
    echo -e "${YELLOW} Create or update crontab for '${USER}' and add scheduled job...${NC}"

    sh_path=$(echo "$PROJECT_DIR/bin/updateIPSum.sh")

    command_to_run=$(echo "NFC_IP_REPTATION_IPSUM=$NFC_IP_REPTATION_IPSUM NFC_IP2L_ASN=$NFC_IP2L_ASN NFC_IP2L_IP=$NFC_IP2L_IP NFC_IP2L_PROXY=$NFC_IP2L_PROXY NFC_IP2L_LOCAL=$NFC_IP2L_LOCAL $PROJECT_DIR/bin/nfupdater -ipsum")
    schedule_to=$(echo "0 7 * * * ")

    echo "#!/bin/bash" > $sh_path
    echo $command_to_run >> $sh_path
    chmod +x -v $sh_path


    echo -e "${YELLOW} this command will be added to crontab if not available in 'crontab -l': ${NC}'nfupdater -ipsum'"

    if echo $(crontab -l) | grep -q "updateIPSum" ; then
        echo -e "${YELLOW} >>> Detected ${NC}'updateIPSum' ${YELLOW}and no further action is required!"
    else
        crontab -l | echo "${schedule_to}${sh_path}" | crontab -
        echo -e "${YELLOW} >>> ${NC}'${schedule_to}${sh_path}' ${YELLOW}added to crontab!"
    fi

    echo -e "${GRREN}...done!${NC}"

}





# print info about this shell-script
# and let user confirm the execution
echo ""
echo -e "${YELLOW} This shell-script will download & install ${NC}'nfcollector' ${YELLOW}to your system."
echo -e "${YELLOW} - Please check your server is connected to the Internet."
echo -e "${YELLOW} - The project directory is ${NC}'${PROJECT_DIR}'"
echo -e "${YELLOW} - These 'Docker' containers will be downloaded during installation:"
echo -e "${YELLOW} \t- ${NC}'influxdb:2.0'"
echo -e "${YELLOW} \t- ${NC}'grafana/grafana'"
echo -e "${YELLOW} \t- ${NC}'javadmohebbi/gonfcollector'"
echo -e "${YELLOW} - Also this script will download all the other required files to ${NC}${PROJECT_DIR}${YELLOW}. Including:"
echo -e "${YELLOW} \t- ${NC}'https://download.openintelligence24.com/nf/bin/nfcollector'"
echo -e "${YELLOW} \t- ${NC}'https://download.openintelligence24.com/nf/bin/nfupdater'"
echo -e "${YELLOW} \t- ${NC}'https://download.openintelligence24.com/vendors/grafana/grafana.tar.gz'"
echo -e "${YELLOW} \t- ${NC}'https://download.openintelligence24.com/vendors/docker-compose/docker-compose.yml'"
echo -e "${YELLOW} \t- ${NC}'nfupdater -ipsum -ip2l -ip2l-asn -ip2l-proxy'${YELLOW} command will download ${NC}'ip2location lite DBs'${YELLOW} and ${NC}'ipsum IP reputation'${YELLOW} from the internet & palce them in ${NC}'${PROJECT_DIR}.'"
echo -e "${YELLOW} - And a daily cron job will be created to download ${NC}'IPSUM' ${YELLOW}daily at ${NC}'7AM'"
echo ""
echo -e "${YELLOW} For more info visit: ${NC}'https://github.com/javadmohebbi/goNfCollector'"
echo ""

read -r -p "Do you want to continue? [y/N] " resp
if [[ "$resp" =~ ^([yY][eE][sS]|[yY])$ ]]
then
    echo -e "${YELLOW} Start working...${NC}"
else
    echo -e "${YELLOW} Thank you! shell-script will be exited now${NC}"
    exit 0
fi





# check the required packages
requirement_check


# preparing direcoty we need
prepare_dir

# downloading latest version of needed files
download_latest_version

# config influx db using a temporary container
get_influx_db_info


#rename compose & other templates
replace_compose_template

replace_collector_template
replace_location_template

# add cron jobs
cron_jobs

# print info
print_info

