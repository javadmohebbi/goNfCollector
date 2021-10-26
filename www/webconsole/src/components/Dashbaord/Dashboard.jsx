import React from 'react'

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import CircularProgress from '@material-ui/core/CircularProgress';

import TimePickerWidget from '../../widgets/timePicker'
import { IconButton, Typography } from '@material-ui/core';
import RefreshIcon from '@material-ui/icons/Refresh'

import DeviceSummaryComponents from './DeviceSummary/DeviceSummary';



const useStyles = makeStyles((theme) => ({
    top: {
        flexGrow: 1,
    },
    paper: {
        padding: theme.spacing(2),
        textAlign: 'center',
        color: theme.palette.text.secondary,
    },
    h1: {
        fontSize: '20px',
        fontWeight: '700',
        textAlign: 'left',
        paddingRight: '10px'
    },
    title: {
        display: 'flex',
        alignItems: 'center',
    },
    btnGrid: {
        textAlign: 'right',
    },
    btns: {
        display: 'flex',
        justifyContent: 'flex-end',
        alignItems: 'center',
    },
    dashboardBody: {
        flexGrow: 1,
        marginTop: '24px',
        marginBottom: '24px'
    },
}));





const DashboardComponent = ({ pageTitle = "Dashboard" }) => {

    const classes = useStyles();

    const [busy, setBusy] = React.useState(false);

    const [timeRange, setTimeRange] = React.useState('15m')

    const [refresh, setRefresh] = React.useState(false)



    const handleRefreshChange = (rf) => {
        setRefresh(rf)
    }

    const handleOnTimeRangeChange = (tr) => {
        setTimeRange(tr)
    }

    const handleBusyState = (b) => {
        setBusy(b)
    }


    return (
        <div className={classes.root}>
            <Paper className={classes.paper}>
                <Grid container spacing={1}
                    direction="row"
                    justify="flex-start"
                    alignItems="center"
                >
                    <Grid item xs={6} sm={6} md={10} className={classes.title}>

                        <Typography
                            variant="h1"
                            color="primary"
                            className={classes.h1}
                        >
                            {pageTitle}
                        </Typography>
                        {
                            busy ? <CircularProgress color="primary" size={15} /> : ''
                        }
                    </Grid>
                    <Grid item xs={6} sm={6} md={2} className={classes.btnGrid}>
                        <div className={classes.btns}>
                            <TimePickerWidget busy={busy} onTimeRangeChange={handleOnTimeRangeChange} />
                            <IconButton
                                onClick={(e) => { e.preventDefault(); setRefresh(true) }}
                                aria-label="refresh" component="span" disabled={busy}
                            >
                                {
                                    busy
                                        ? <CircularProgress color="primary" size={15} />
                                        : <RefreshIcon fontSize="small" />
                                }
                            </IconButton>
                        </div>
                    </Grid>
                </Grid>
            </Paper>

            <div>
                <Grid container spacing={2} className={classes.dashboardBody} >

                    <Grid item xs={12}>
                        {/* Device SUMMARY */}
                        <DeviceSummaryComponents
                            refresh={refresh}
                            interval={timeRange}
                            busy={busy}
                            handleParentBusyState={handleBusyState}
                            handleParentRefreshState={handleRefreshChange}
                        />
                    </Grid>
                </Grid>
            </div>
        </div>
    )
}

export default DashboardComponent