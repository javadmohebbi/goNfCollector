
import React, { useCallback, useEffect, useState } from 'react';


import { makeStyles } from '@material-ui/core/styles';

import { DeviceGetSummaryByInterval, DeviceGetSummaryByIntervalByDev } from '../../../services/devices';
import { Paper, Typography } from '@material-ui/core';

import CircularProgress from '@material-ui/core/CircularProgress';


import DeviceSummaryItem from './Item';


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

}))


const controller = new AbortController()
const signal = controller.signal



const DeviceSummaryComponents = ({
    handleParentBusyState = () => { return },
    handleParentRefreshState = () => { return },
    interval = '15m',
    refresh = false,
    busy = false,

    // the device we want to filter
    // to show their summary
    requestedDevices = '',

}) => {

    const classes = useStyles();

    const [deviceSummary, setDeviceSummary] = useState([])
    const [selfBusy, setSelfBusy] = useState(true)

    // eslint-disable-next-line
    const [fetchError, setFetchError] = useState(false)



    const getSummaryCallback = useCallback(() => {
        setSelfBusy(true)


        handleParentBusyState(true)
        handleParentRefreshState(false)


        controller.abort()
        if (requestedDevices === '') {
            DeviceGetSummaryByInterval({ interval, signal }).then(async (json) => {
                if (json.error) {
                    setFetchError(true)
                } else {
                    const resp = await json.response.then((result) => result);
                    if (resp !== null) {
                        setDeviceSummary(resp)
                    } else {
                        setDeviceSummary([])
                    }
                }
                setSelfBusy(false)

                handleParentBusyState(false)
            })
        } else {
            DeviceGetSummaryByIntervalByDev({ interval, device: requestedDevices, signal }).then(async (json) => {
                if (json.error) {
                    setFetchError(true)
                } else {
                    const resp = await json.response.then((result) => result);
                    if (resp !== null) {
                        setDeviceSummary(resp)
                    } else {
                        setDeviceSummary([])
                    }
                }
                setSelfBusy(false)

                handleParentBusyState(false)
            })
        }

        // eslint-disable-next-line
    }, [interval])


    useEffect(() => {
        if (refresh) {
            getSummaryCallback()
        }
    }, [refresh, getSummaryCallback])

    useEffect(() => {
        getSummaryCallback()
    }, [interval, getSummaryCallback])


    return (
        <div>
            <Paper variant="outlined" className={classes.paperTitle}>
                <Typography className={classes.paperTg}>
                    Netfow Exporter Summary {
                        selfBusy ? <CircularProgress color="primary" size={15} /> : ''
                    }
                </Typography>

                {
                    deviceSummary.map((dev, i) => (
                        <div key={i}>
                            {/* DEVICE SUMM */}
                            <DeviceSummaryItem device={dev} i={i} selfBusy={selfBusy}
                                interval={interval}
                                refresh={refresh}
                                handleParentRefreshState={handleParentRefreshState}
                                handleParentBusyState={handleParentBusyState}
                            />


                        </div>
                    ))
                }


            </Paper>

        </div >
    );
};

export default DeviceSummaryComponents;