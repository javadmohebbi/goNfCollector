import React from 'react'

import AppBarLayout from './AppBar';

import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';


const useStyles = makeStyles({
    root: {
        flexGrow: 1,
        maxWidth: '100%',
        overflowX: 'hidden',
        padding: '24px',
        height: 'calc(100vh - 70px)',
    },

    rootBTM: {
        // width: 500,
    },

});


const MainAdminLayout = (props) => {
    const classes = useStyles();


    return (
        <React.Fragment>
            <AppBarLayout />

            <div className={classes.root}>
                <Grid container spacing={2} >
                    <Grid item xs={12}>
                        {props.children}
                    </Grid>
                </Grid>
            </div>



        </React.Fragment>
    )
}

export default MainAdminLayout