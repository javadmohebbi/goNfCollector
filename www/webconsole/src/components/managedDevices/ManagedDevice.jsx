import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import { Typography, Paper, Grid, CircularProgress } from '@material-ui/core';
import ManagedDeviceDataTableComponent from './dt';
import { IconButton } from '@material-ui/core';
import RefreshIcon from '@material-ui/icons/Refresh'

const useStyles = makeStyles((theme) => ({
    root: {

    },
    loading: {
        marginLeft: theme.spacing(2)
    },
    h1: {
        fontSize: '20px',
        fontWeight: '700',
        textAlign: 'left',
        paddingRight: '10px'
    },
    paper: {
        padding: theme.spacing(2),
        textAlign: 'center',
        color: theme.palette.text.secondary,
    },
    dtHolder: {
        marginTop: theme.spacing(2),
        marginBottom: theme.spacing(2),
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
}))
const ManagedDeviceComponent = (props) => {

    const classes = useStyles();

    const [busy, setBusy] = React.useState(false);

    // eslint-disable-next-line
    const [refresh, setRefresh] = React.useState(false)

    const handleRefreshChange = (rf) => {
        setRefresh(rf)
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
                            Managed Devices
                            {
                                busy ? <CircularProgress className={classes.loading} color="primary" size={15} /> : ''
                            }
                        </Typography>
                    </Grid>
                    <Grid item xs={6} sm={6} md={2} className={classes.btnGrid}>
                        <div className={classes.btns}>
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

            <Grid container spacing={2} className={classes.dtHolder}>
                <Grid item xs={12} md={12} >
                    <ManagedDeviceDataTableComponent
                        busy={busy}
                        refresh={refresh}
                        handleParentBusyState={handleBusyState}
                        handleParentRefreshState={handleRefreshChange}
                    />
                </Grid>
            </Grid>

        </div>
    )
}

export default ManagedDeviceComponent;