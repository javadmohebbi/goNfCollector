import React, { useState } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import { Button, Paper, TextField } from '@material-ui/core';
import ClearIcon from '@material-ui/icons/Clear';
import SearchIcon from '@material-ui/icons/Search';
import _ from 'lodash'
import { fltModel } from './filtermodel';



const useStyles = makeStyles((theme) => ({
    root: {

    },
    paper: {
        padding: theme.spacing(2),
        paddingBottom: 0,
        textAlign: 'center',
        color: theme.palette.text.secondary,
    },
    formRoot: {
        '& .MuiTextField-root': {
            padding: theme.spacing(1),
        },
    },
    formButton: {
        marginTop: theme.spacing(2),
        marginRight: theme.spacing(1),
        marginLeft: theme.spacing(1),
    },
}))

function FilterFormComponent({
    filter,
    callback,
}) {


    const classes = useStyles();
    const [flt, setFlt] = useState({ ...filter })

    const handleFilterClick = (e) => {
        e.preventDefault()
        callback({ ...flt }, true)
    }
    const handleFilterClear = (e) => {
        e.preventDefault()
        setFlt({ ...fltModel })
        callback({ ...fltModel }, false)
    }

    return (
        <div className={classes.root}>
            <Paper className={classes.paper}>
                <form
                    className={classes.formRoot}
                    noValidate
                    autoComplete="off"
                    onSubmit={(e) => {
                        e.preventDefault()
                    }}
                >

                    <TextField
                        id="device"
                        label="Filter Device"
                        variant="filled"
                        value={flt.device}
                        onChange={
                            (e) => {
                                setFlt({ ...flt, device: e.target.value })
                            }
                        }
                    />

                    <TextField
                        id="ip"
                        label="Filter IP"
                        variant="filled"
                        value={flt.ip}
                        onChange={
                            (e) => {
                                setFlt({ ...flt, ip: e.target.value })
                            }
                        }
                    />

                    <TextField
                        id="port"
                        label="Filter Port"
                        variant="filled"
                        value={flt.port}
                        onChange={
                            (e) => {
                                setFlt({ ...flt, port: e.target.value })
                            }
                        }
                    />

                    <TextField
                        id="proto"
                        label="Filter Protocol"
                        variant="filled"
                        value={flt.proto}
                        onChange={
                            (e) => {
                                setFlt({ ...flt, proto: e.target.value })
                            }
                        }
                    />

                    <TextField
                        id="country"
                        label="Filter Country"
                        variant="filled"
                        value={flt.country}
                        onChange={
                            (e) => {
                                setFlt({ ...flt, country: e.target.value })
                            }
                        }
                    />





                </form><br />

                <Button
                    variant="contained"
                    color="primary"
                    className={classes.formButton}
                    startIcon={<SearchIcon />}
                    onClick={handleFilterClick}
                >
                    Click to Filter
                </Button>

                {
                    !_.isEqual(flt, fltModel) ?
                        <Button
                            variant="contained"
                            color="secondary"
                            className={classes.formButton}
                            startIcon={<ClearIcon />}
                            onClick={handleFilterClear}
                        >
                            Clear Filters
                        </Button> :
                        <></>
                }
                <br />
                By clicking Filter or Clear Filter button, below table will be cleared and be filled with new data!

            </Paper>
        </div>
    );
}

export default FilterFormComponent;