import { Grid, Paper, Typography } from '@material-ui/core';
import CircularProgress from '@material-ui/core/CircularProgress';
import { makeStyles } from '@material-ui/core/styles';
import React from 'react';



const useStyles = makeStyles((theme) => ({
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
    paper: {
        padding: theme.spacing(2),
        textAlign: 'center',
        color: theme.palette.text.secondary,
    },
    bodyDiv: {
        marginTop: theme.spacing(2),
    },
    bodyPaper: {
        padding: theme.spacing(2),
        color: theme.palette.text.secondary,
    },
}));



const SettingsComponent = () => {
    const classes = useStyles();

    const [busy, setBusy] = React.useState(false);
    const [refresh, setRefresh] = React.useState(false)

    return (
        <div className={classes.root}>
            <Paper className={classes.paper}>
                <Grid container spacing={1}
                    direction="row"
                    justify="flex-start"
                    alignItems="center">
                    <Grid item xs={6} sm={6} md={10} className={classes.title}>

                        <Typography
                            variant="h1"
                            color="primary"
                            className={classes.h1}
                        >
                            Settings
                        </Typography>
                        {
                            busy ? <CircularProgress color="primary" size={15} /> : ''
                        }
                    </Grid>
                </Grid>
            </Paper>


            <div className={classes.bodyDiv}>
                <Paper className={classes.bodyPaper}>Body</Paper>
            </div>

        </div>
    );
};

export default SettingsComponent;