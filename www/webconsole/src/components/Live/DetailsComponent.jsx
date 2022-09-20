import React from 'react';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import TCPFlagsWidgets from './TCPFlagsWidgets';
import BytesWidget from './BytesWidget';
import PacketsWidget from './PacketsWidget';
import { Chip, Typography } from '@material-ui/core';
import ThreatLevelWidget from '../../widgets/ThreatLevel';
import NetflowVersionWidget from '../../widgets/NetflowVersion';

const useStyles = makeStyles((theme) => ({
    root: {
        flexGrow: 1,
        width: '100%',
        padding: 5,
        margin: 5,
    },
}));

function DetailsComponent({ row }) {

    const classes = useStyles();

    return (
        <div className={classes.root}>
            <Grid container spacing={3}>
                <Grid item xs={12}>
                    <Grid container spacing={3}>
                        <Grid item>
                            <NetflowVersionWidget version={row.flowVersion} />
                        </Grid>
                        <Grid item>
                            <Typography variant="h6" component="span">
                                Device: {row.device}
                            </Typography>
                        </Grid>
                        <Grid item>
                            <Typography variant="h6" component="span">Proto: {row.protoName}</Typography>
                        </Grid>
                        <Grid item>
                            <Typography variant="h6" component="span">TCP Flags:
                                {
                                    row.protoName.toLowerCase() === 'tcp' ?
                                        <TCPFlagsWidgets
                                            fin={row.fin}
                                            syn={row.syn}
                                            rst={row.rst}
                                            psh={row.psh}
                                            ack={row.ack}
                                            urg={row.urg}
                                            ece={row.ece}
                                            cwr={row.cwr}
                                        /> : '-'
                                }

                            </Typography>
                        </Grid>
                        <Grid item>
                            <Typography variant="h6" component="span">Next Hop: {row.nextHop}</Typography>
                        </Grid>
                        <Grid item>
                            <Typography variant="h6" component="span">Bytes: <BytesWidget bytes={row.bytes} /></Typography>
                        </Grid>
                        <Grid item>
                            <Typography variant="h6" component="span">Packets: <PacketsWidget packets={row.packets} /></Typography>
                        </Grid>
                    </Grid>
                </Grid>
                <Grid item xs={6}>
                    <div>
                        <Grid container spacing={3}>
                            <Grid item xs={12}>
                                <Typography variant="h6">
                                    <Chip label="Source" style={{ fontSize: '120%' }} />
                                </Typography>
                            </Grid>
                            <Grid item xs={12}>
                                <div><Chip label={`IP: ${row.srcIP}`} style={{ fontSize: '110%', marginBottom: '10px' }} /></div>
                                <div><Chip label={`Port: ${row.srcPortName.includes(row.srcPort) ? row.srcPort : row.srcPortName}`} style={{ fontSize: '110%', marginBottom: '10px' }} /></div>
                                <div><Chip label={`Country: ${row.sCountryLong} (${row.sCountryShort})`} style={{ fontSize: '110%', marginBottom: '10px' }} /></div>
                                <div><Chip label={`State/Region: ${row.sState}`} style={{ fontSize: '110%', marginBottom: '10px' }} /></div>
                                <div><Chip label={`City: ${row.sCity}`} style={{ fontSize: '110%', marginBottom: '10px' }} /></div>
                                {
                                    row.isSrcThreat ?
                                        <div style={{ marginTop: '10px', marginBottom: '10px', fontSize: '110%' }}>
                                            <ThreatLevelWidget
                                                inline={true}
                                                level={row.srcThreatReputation}
                                                label={`Threat: ${row.srcThreatType} => ${row.srcThreatKind}`}
                                            />
                                        </div>
                                        : <></>
                                }
                            </Grid>
                        </Grid>
                    </div>
                </Grid>
                <Grid item xs={6}>
                    <div>
                        <Grid container spacing={3}>
                            <Grid item xs={12}>
                                <Typography variant="h6">
                                    <Chip label="Destination" style={{ fontSize: '120%' }} />
                                </Typography>
                            </Grid>
                            <Grid item xs={12}>
                                <div><Chip label={`IP: ${row.dstIP}`} style={{ fontSize: '110%', marginBottom: '10px' }} /></div>
                                <div><Chip label={`Port: ${row.dstPortName.includes(row.dstPort) ? row.dstPort : row.dstPortName}`} style={{ fontSize: '110%', marginBottom: '10px' }} /></div>
                                <div><Chip label={`Country: ${row.dCountryLong} (${row.dCountryShort})`} style={{ fontSize: '110%', marginBottom: '10px' }} /></div>
                                <div><Chip label={`State/Region: ${row.dState}`} style={{ fontSize: '110%', marginBottom: '10px' }} /></div>
                                <div><Chip label={`City: ${row.dCity}`} style={{ fontSize: '110%', marginBottom: '10px' }} /></div>
                                {
                                    row.isDstThreat ?
                                        <div style={{ marginTop: '10px', marginBottom: '10px', fontSize: '110%' }}>
                                            <ThreatLevelWidget
                                                inline={true}
                                                level={row.dstThreatReputation}
                                                label={`Threat: ${row.dstThreatType} => ${row.dstThreatKind}`}
                                            />
                                        </div>
                                        : <></>
                                }
                            </Grid>
                        </Grid>
                    </div>
                </Grid>
            </Grid>
        </div>
    );
}

export default DetailsComponent;