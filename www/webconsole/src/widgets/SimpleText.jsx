import React from 'react';
import Paper from '@material-ui/core/Paper';

import { makeStyles } from '@material-ui/core/styles';
import { CircularProgress, Typography } from '@material-ui/core';

import humanFormat from 'human-format'


const useStyles = makeStyles((theme) => ({
    paper: {
        padding: theme.spacing(2),
        textAlign: 'center',
    },
    label: {
        fontSize: 14,
    },
    value: {
        fontSize: 18,
        fontWeight: 900,
        marginTop: theme.spacing(1)
    }
}))

const SimpleTextWidget = (props) => {

    const classes = useStyles()


    return (
        <div>
            <Paper className={classes.paper}
                style={{ backgroundColor: props.bgcolor || 'grey', color: props.color || 'white' }}
            >
                <Typography className={classes.label} variant='h4'>
                    {props.label}
                    {/* {
                        typeof props.belongsTo !== 'undefined' &&
                            props.belongsTo !== ''
                            ?
                            <strong>
                                {' | ' + props.belongsTo}
                            </strong>
                            :
                            ''
                    } */}
                </Typography>
                <Typography variant='h3' className={classes.value}>
                    {
                        props.busy
                            ?
                            <CircularProgress color="primary" style={{ color: props.color || 'white' }} size={15} />
                            :
                            <>
                                {props.unit === 'decimal' ? humanFormat(props.value, { unit: props.humarReadable.unit }) : ''}
                            </>
                    }
                </Typography>
            </Paper>
        </div>
    );
};

export default SimpleTextWidget;