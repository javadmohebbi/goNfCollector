import { Box, Grid, Paper, Tab, Tabs, Typography } from '@material-ui/core';
import CircularProgress from '@material-ui/core/CircularProgress';
import { makeStyles } from '@material-ui/core/styles';
import React from 'react';
import IP2LocationSettings from './IP2LocationSetting';



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
    tabPanel: {
        padding: theme.spacing(3),
    }
}));


// interface TabPanelProps {
//     children?: React.ReactNode;
//     index: number;
//     value: number;
// }

function TabPanel(props) {
    const { children, value, index, ...other } = props;

    return (
        <div
            role="tabpanel"
            hidden={value !== index}
            id={`simple-tabpanel-${index}`}
            aria-labelledby={`simple-tab-${index}`}
            {...other}
        >
            {value === index && (
                <Box sx={{ p: 3, marginTop: '15px' }}>
                    <Typography>{children}</Typography>
                </Box>
            )}
        </div>
    );
}


function a11yProps(index) {
    return {
        id: `simple-tab-${index}`,
        'aria-controls': `simple-tabpanel-${index}`,
    };
}


const SettingsComponent = () => {
    const classes = useStyles();

    const [value, setValue] = React.useState(0);

    const handleChange = (event, newValue) => {
        setValue(newValue);
    };

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
                <Paper className={classes.bodyPaper}>

                    <Tabs value={value} onChange={handleChange}
                        variant="scrollable"
                        scrollButtons='on'
                        // allowScrollButtonsMobile
                        aria-label="basic tabs example">
                        <Tab label="GeoLocation Database" {...a11yProps(0)} />
                        <Tab label="Threat Sources Database" {...a11yProps(1)} />
                        <Tab label="Netflow Configuration" {...a11yProps(2)} />
                    </Tabs>
                    <TabPanel value={value} index={0} className={classes.tabPanel}>
                        <IP2LocationSettings />
                    </TabPanel>
                    <TabPanel value={value} index={1} className={classes.tabPanel}>
                        Threat
                    </TabPanel>
                    <TabPanel value={value} index={2} className={classes.tabPanel}>
                        Netflow
                    </TabPanel>
                </Paper>
            </div>

        </div>
    );
};

export default SettingsComponent;