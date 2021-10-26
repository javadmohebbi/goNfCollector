
import React, { useCallback, useEffect, useState } from 'react';

import { makeStyles } from '@material-ui/core/styles';

import { Paper } from '@material-ui/core';

import CircularProgress from '@material-ui/core/CircularProgress';
import { HostsThreatReportByInterval } from '../../services/hosts';
import Alert from '@material-ui/lab/Alert';
import WidgetDataTableComponent from '../../widgets/DataTable';
import { Typography } from '@material-ui/core';

import Accordion from '@material-ui/core/Accordion';
import AccordionSummary from '@material-ui/core/AccordionSummary';
import AccordionDetails from '@material-ui/core/AccordionDetails';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';

import moment from "moment";
import ThreatLevelWidget from '../../widgets/ThreatLevel';



const useStyles = makeStyles((theme) => ({
    paperTitle: {
        marginBottom: theme.spacing(2),
        padding: theme.spacing(2),
        textAlign: 'center',
        background: "#353535",
    },
    paperTg: {
        marginBottom: theme.spacing(2)
    },
    tabs: {
        marginBottom: theme.spacing(2)
    },
    red: {
        backgroundColor: 'red',
        display: 'inline-table',
    },
    normal: {

    },
    item: {
        marginTop: theme.spacing(2),
        marginBottom: theme.spacing(2),
    },
}))


const controller = new AbortController();
const signal = controller.signal;



const HostReportThreatComponent = ({
    handleParentBusyState = () => { return },
    handleParentRefreshState = () => { return },
    interval = '15m',
    refresh = false,
    busy = false,

    // the device we want to filter
    // to show their summary
    host = '',
}) => {

    const classes = useStyles()

    const [hostThreats, setHostThreats] = useState([])
    const [selfBusy, setSelfBusy] = useState(true)

    // eslint-disable-next-line
    const [fetchError, setFetchError] = useState(false)

    const [accordPanel, setAccorPanel] = React.useState(0)

    const handleAccorPanelChange = (panel) => (event, i) => {
        setAccorPanel(i ? panel : false)
    }


    const getCallBack = useCallback(() => {
        setSelfBusy(true)


        handleParentBusyState(true)
        handleParentRefreshState(false)


        controller.abort()

        HostsThreatReportByInterval({ interval, host, signal }).then(async (json) => {
            if (json.error) {
                setFetchError(true)
            } else {
                const resp = await json.response.then((result) => result);
                console.log(resp);
                if (resp.threats.list !== null) {
                    setHostThreats(resp.threats.list)
                } else {
                    setHostThreats([])
                }
            }
            setSelfBusy(false)

            handleParentBusyState(false)
        })


        // eslint-disable-next-line
    }, [interval])


    useEffect(() => {
        if (refresh) {
            getCallBack()
        }
    }, [refresh, getCallBack])

    useEffect(() => {
        getCallBack()
    }, [interval, getCallBack])



    return (
        <div>
            <Paper variant="outlined" className={classes.paperTitle}>
                {(busy || selfBusy) ? <CircularProgress color="primary" size={15} /> : ''}
                {
                    !(busy || selfBusy) ? hostThreats.length === 0 ?
                        <Alert variant="filled" severity="success" >{`No threat detected on '${host}' in last '${interval}'`}</Alert>
                        :
                        <Alert variant="filled" severity="error" >{`Threat detected on '${host}' in last '${interval}'`}</Alert>
                        : ''
                }
                {
                    Array.isArray(hostThreats) ? hostThreats.map((obj, i) => (
                        <div key={obj.threat_id}>
                            <Accordion className={classes.item} key={i} expanded={accordPanel === i} onChange={handleAccorPanelChange(i)}>
                                <AccordionSummary className={classes.itemSum}
                                    expandIcon={<ExpandMoreIcon />}
                                    aria-controls={'acc-dev-content' + i}
                                    id={'acc-dev-' + i}
                                >
                                    <Typography>
                                        <ThreatLevelWidget label={'Threat: ' + obj.threat_source + ' | ' + obj.threat_kind + " > Level: "} level={obj.threat_reputation} />
                                    </Typography>
                                </AccordionSummary>
                                <AccordionDetails style={{ width: '100%' }}>

                                    <WidgetDataTableComponent
                                        refresh={refresh}
                                        URL={`flows/get/all/threat/${obj.threat_id}/interval/${interval}`}
                                        Columns={[
                                            {
                                                key: "src_host", label: "src", dataType: 'string',
                                                reFormat: (val, row) => <><span className={row.src_is_threat ? classes.red : classes.normal}>{val}</span>{`  ${row.src_port_name}`}</>,
                                            },
                                            {
                                                key: "dst_host", label: "dst", dataType: 'string',
                                                reFormat: (val, row) => <><span className={row.dst_is_threat ? classes.red : classes.normal}>{val}</span>{`  ${row.dst_port_name}`}</>,
                                            },
                                            {
                                                key: "next_hop_host", label: "next hop", dataType: 'string',
                                                reFormat: (val, row) => <><span className={row.next_hop_is_threat ? classes.red : classes.normal}>{val}</span></>,
                                            },
                                            {
                                                key: "protocol", label: "proto", dataType: 'string',
                                            },
                                            {
                                                key: "byte", label: "bytes", dataType: 'string',
                                            },
                                            {
                                                key: "packet", label: "packets", dataType: 'string',
                                            },
                                            {
                                                key: "created_at", label: "time", dataType: 'dateTime',
                                                reFormat: (val, row) => moment(val).fromNow()
                                            },
                                        ]}
                                        handleParentBusyState={handleParentBusyState}
                                        handleParentRefreshState={handleParentRefreshState}
                                        PerPage={5}
                                        Order="created_at"
                                        OrderType="desc"
                                    />
                                </AccordionDetails>
                            </Accordion>
                        </div>
                    ))
                        :
                        ''
                }
            </Paper>

        </div>
    )
}

export default HostReportThreatComponent;