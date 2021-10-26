import { Paper } from '@material-ui/core';
import React from 'react';
import { makeStyles } from '@material-ui/core/styles';

import CheckIcon from '@material-ui/icons/Check'
import CloseIcon from '@material-ui/icons/Close'


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

const ThreatUserAckWidget = ({ acked = false, closed = false, falsePositive = false, inline = false, labels = ["acked?", "closed?", "false+?"], showLabels = false }) => {

    const classes = useStyles()


    return (
        <div className={classes.root} style={{ display: inline ? 'inline-flex' : 'flex' }}>
            {[acked, closed, falsePositive].map((value, i) => (
                <Paper component="span" key={i} style={{
                    backgroundColor: value ? '#4FFF24' : '#F9FF24', color: 'black',
                    marginRight: '0.2em'
                }}>
                    {
                        showLabels ? labels[i] + '  ' : ''
                    }
                    {value ? <CheckIcon /> : <CloseIcon />} {value}
                </Paper>
            ))}

        </div>
    );
}

export default ThreatUserAckWidget;