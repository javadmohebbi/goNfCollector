import React from 'react';


import { makeStyles } from '@material-ui/core/styles';

import Accordion from '@material-ui/core/Accordion';
import AccordionSummary from '@material-ui/core/AccordionSummary';
import AccordionDetails from '@material-ui/core/AccordionDetails';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import Grid from '@material-ui/core/Grid';
import { Typography } from '@material-ui/core';

import SimpleTextWidget from '../../../widgets/SimpleText';
import TopPortWidget from '../TopPort/TopPort';
import TopProtocolWidget from '../TopProto/TopProtocol';
import TopHostWidget from '../TopHost/TopHost';
import EthernetSummaryWidget from '../Ethernet/EthernetSummary';
import TopCountryWidget from '../TopGeo/TopCountry';


const colorsToRows = [
    { bgcolor: '#09F1AD', color: 'black' },
    { bgcolor: '#E656FF', color: 'black' },
    { bgcolor: '#56FF99', color: 'black' },
    { bgcolor: '#A5FF56', color: 'black' },
    { bgcolor: '#FF569F', color: 'black' },
    { bgcolor: '#FFBA56', color: 'black' },
]


const useStyles = makeStyles((theme) => ({
    item: {
        marginTop: theme.spacing(2),
        marginBottom: theme.spacing(2),
    },
    itemSum: {
        '& div': { justifyContent: 'center' },
    },
    itemTitle: {
        fontSize: 20,
    }
}))



const DeviceSummaryItem = ({
    device,
    i,
    selfBusy,
    handleParentBusyState = () => { return },
    handleParentRefreshState = () => { return },
    interval = '15m',
    refresh = false,
}) => {


    const classes = useStyles();

    const [accordPanel, setAccorPanel] = React.useState(0)

    const handleAccorPanelChange = (panel) => (event, i) => {
        setAccorPanel(i ? panel : false)
    }




    return (
        <Accordion className={classes.item} key={i} expanded={accordPanel === i} onChange={handleAccorPanelChange(i)}>
            <AccordionSummary className={classes.itemSum}
                expandIcon={<ExpandMoreIcon />}
                aria-controls={'acc-dev-content' + i}
                id={'acc-dev-' + i}
            >
                <Typography style={{ color: colorsToRows[i % colorsToRows.length].bgcolor }} className={classes.itemTitle}>
                    {device.device_name === '' ? device.device : device.device_name + ` (${device.device})`}
                </Typography>
            </AccordionSummary>
            <AccordionDetails>
                <Grid container justifycontent="center" spacing={2}>
                    <Grid item xs={12}>
                        <Grid container justifycontent="center" spacing={2}>
                            <Grid item xs={12} sm={12} md={6} lg={4}>
                                <SimpleTextWidget
                                    bgcolor={colorsToRows[i % colorsToRows.length].bgcolor}
                                    color={colorsToRows[i % colorsToRows.length].color}
                                    label="Flow Count"
                                    belongsTo={
                                        device.device_name === '' ? device.device : device.device_name
                                    }
                                    unit="decimal"
                                    humarReadable={{
                                        unit: "",
                                    }}
                                    value={device.flow_count}
                                    busy={selfBusy}
                                />
                            </Grid>
                            <Grid item xs={12} sm={12} md={6} lg={4}>
                                <SimpleTextWidget
                                    bgcolor={colorsToRows[i % colorsToRows.length].bgcolor}
                                    color={colorsToRows[i % colorsToRows.length].color}
                                    label="Total Bytes"
                                    belongsTo={
                                        device.device_name === '' ? device.device : device.device_name
                                    }
                                    unit="decimal"
                                    humarReadable={{
                                        unit: "B",
                                    }}
                                    value={device.total_bytes}
                                    busy={selfBusy}
                                />
                            </Grid>
                            <Grid item xs={12} sm={12} md={6} lg={4}>
                                <SimpleTextWidget
                                    bgcolor={colorsToRows[i % colorsToRows.length].bgcolor}
                                    color={colorsToRows[i % colorsToRows.length].color}
                                    label="Total Packets"
                                    belongsTo={
                                        device.device_name === '' ? device.device : device.device_name
                                    }
                                    unit="decimal"
                                    humarReadable={{
                                        unit: "",
                                    }}
                                    value={device.total_packets}
                                    busy={selfBusy}
                                />
                            </Grid>
                        </Grid>
                    </Grid>

                    {
                        selfBusy ? '' :
                            <>
                                {/* SRC AND DST PORTS */}
                                <Grid item xs={12}>
                                    <Grid container justifycontent="center" spacing={2}>
                                        <Grid item xs={12} sm={12} md={6}>
                                            <TopPortWidget
                                                interval={interval}
                                                top={'10'}
                                                deviceId={device.device_id}
                                                direction='src'
                                                refresh={refresh}
                                                handleParentRefreshState={handleParentRefreshState}
                                                handleParentBusyState={handleParentBusyState}
                                            />
                                        </Grid>
                                        <Grid item xs={12} sm={12} md={6}>
                                            <TopPortWidget
                                                interval={interval}
                                                top={'10'}
                                                deviceId={device.device_id}
                                                direction='dst'
                                                refresh={refresh}
                                                handleParentRefreshState={handleParentRefreshState}
                                                handleParentBusyState={handleParentBusyState}
                                            />
                                        </Grid>
                                    </Grid>
                                </Grid>

                                {/* SRC AND DST Hosts */}
                                <Grid item xs={12}>
                                    <Grid container justifycontent="center" spacing={2}>
                                        <Grid item xs={12} sm={12} md={6}>
                                            <TopHostWidget
                                                interval={interval}
                                                top={'10'}
                                                deviceId={device.device_id}
                                                direction='src'
                                                refresh={refresh}
                                                handleParentRefreshState={handleParentRefreshState}
                                                handleParentBusyState={handleParentBusyState}
                                            />
                                        </Grid>
                                        <Grid item xs={12} sm={12} md={6}>
                                            <TopHostWidget
                                                interval={interval}
                                                top={'10'}
                                                deviceId={device.device_id}
                                                direction='dst'
                                                refresh={refresh}
                                                handleParentRefreshState={handleParentRefreshState}
                                                handleParentBusyState={handleParentBusyState}
                                            />
                                        </Grid>
                                    </Grid>
                                </Grid>


                                {/* SRC AND DST Country */}
                                <Grid item xs={12}>
                                    <Grid container justifycontent="center" spacing={2}>
                                        <Grid item xs={12} sm={12} md={6}>
                                            <TopCountryWidget
                                                interval={interval}
                                                top={'10'}
                                                deviceId={device.device_id}
                                                direction='src'
                                                refresh={refresh}
                                                handleParentRefreshState={handleParentRefreshState}
                                                handleParentBusyState={handleParentBusyState}
                                            />
                                        </Grid>
                                        <Grid item xs={12} sm={12} md={6}>
                                            <TopCountryWidget
                                                interval={interval}
                                                top={'10'}
                                                deviceId={device.device_id}
                                                direction='dst'
                                                refresh={refresh}
                                                handleParentRefreshState={handleParentRefreshState}
                                                handleParentBusyState={handleParentBusyState}
                                            />
                                        </Grid>
                                    </Grid>
                                </Grid>



                                {/* TOP PROTOCOL  */}
                                <Grid item xs={12} md={6}>
                                    {/* SRC AND DST PORTS */}
                                    <Grid container justifycontent="center" spacing={2}>
                                        <Grid item xs={12}>
                                            <TopProtocolWidget
                                                interval={interval}
                                                top={'10'}
                                                deviceId={device.device_id}
                                                refresh={refresh}
                                                handleParentRefreshState={handleParentRefreshState}
                                                handleParentBusyState={handleParentBusyState}
                                            />
                                        </Grid>
                                    </Grid>
                                </Grid>


                                {/* TOP Ethernet  */}
                                <Grid item xs={12} md={6}>
                                    {/* SRC AND DST PORTS */}
                                    <Grid container justifycontent="center" spacing={2}>
                                        <Grid item xs={12}>
                                            <EthernetSummaryWidget
                                                interval={interval}
                                                deviceId={device.device_id}
                                                refresh={refresh}
                                                handleParentRefreshState={handleParentRefreshState}
                                                handleParentBusyState={handleParentBusyState}
                                            />
                                        </Grid>
                                    </Grid>
                                </Grid>
                            </>
                    }
                </Grid>


            </AccordionDetails>
        </Accordion>
    )

}



export default DeviceSummaryItem