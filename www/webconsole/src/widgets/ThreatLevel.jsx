import { Paper } from '@material-ui/core';
import React from 'react';
import { makeStyles } from '@material-ui/core/styles';

const colors = [
    { level: 'Very Low', bgcolor: '#253BFF', color: 'white' },
    { level: 'Low', bgcolor: '#5A6AFF', color: 'black' },
    { level: 'Medium', bgcolor: '#FFC25A', color: 'black' },
    { level: 'High', bgcolor: '#FF6625', color: 'black' },
    { level: 'Very High', bgcolor: '#F61A2D', color: 'white' },
]


const getLevel = level => {
    if (level > 8) {
        return colors[4]
    }
    if (level <= 8 && level > 6) {
        return colors[3]
    }
    if (level <= 6 && level > 4) {
        return colors[2]
    }
    if (level <= 4 && level > 2) {
        return colors[1]
    }
    return colors[0]
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

const ThreatLevelWidget = ({ level, justCounter = false, inline = false, label = '' }) => {

    const classes = useStyles()
    const v = getLevel(level)

    return (
        <div className={classes.root} style={{ display: inline ? 'inline' : 'unset' }}>
            <Paper className={classes.paper} style={{ color: v.color, backgroundColor: v.bgcolor, display: inline ? 'inline' : 'unset' }}>
                {
                    label !== '' ? label + ' ' : ''
                }
                {
                    justCounter ? level : level + '-' + v.level
                }

            </Paper >
        </div>
    );
}

export default ThreatLevelWidget;