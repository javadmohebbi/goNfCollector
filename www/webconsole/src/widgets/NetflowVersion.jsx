import { Paper } from '@material-ui/core';
import React from 'react';
import { makeStyles } from '@material-ui/core/styles';

const colors = [
    { version: 1, bgcolor: '#FF9D57', color: 'black', label: 'v1' },
    { version: 5, bgcolor: '#FF57AA', color: 'black', label: 'v5' },
    { version: 6, bgcolor: '#FFCF57', color: 'black', label: 'v6' },
    { version: 7, bgcolor: '#DCFF57', color: 'black', label: 'v7' },
    { version: 9, bgcolor: '#57FEFF', color: 'black', label: 'v9' },
    { version: 10, bgcolor: '#57FF79', color: 'black', label: 'ipfix' },
]


const getVersion = version => {
    for (let i = 0; i < colors.length; i++) {
        if (version === colors[i].version) {
            return colors[i]
        }
    }
    return { bgcolor: 'black', color: 'white', label: 'unknown!' }
}



const useStyles = makeStyles((theme) => ({
    root: {
        display: 'flex',
        justifyContent: 'flex-start',
    },
    paper: {
        paddingTop: theme.spacing(1),
        paddingBottom: theme.spacing(1),
        paddingRight: theme.spacing(2),
        paddingLeft: theme.spacing(2)
    }
}))

const NetflowVersionWidget = ({ version }) => {

    const classes = useStyles()
    const v = getVersion(version)

    return (
        <div className={classes.root}>
            <Paper className={classes.paper} style={{ color: v.color, backgroundColor: v.bgcolor, }}>
                {v.label}
            </Paper >
        </div>
    );
}

export default NetflowVersionWidget;