import React, { useEffect, useState } from 'react';
import io from 'socket.io-client';
import { makeStyles } from '@material-ui/core/styles';
import { Box, Button, Chip, CircularProgress, Collapse, Grid, IconButton, Paper, Table, TableBody, TableCell, TableContainer, TableHead, TablePagination, TableRow, Typography } from '@material-ui/core';

// import PlayCircleFilledWhiteIcon from '@material-ui/icons/PlayCircleFilledWhite';
// import StopCircleIcon from '@material-ui/icons/StopCircle';

import PlayCircleIcon from '@material-ui/icons/PlayArrow'
import StopCircleIcon from '@material-ui/icons/Stop'
import KeyboardArrowUpIcon from '@material-ui/icons/KeyboardArrowUp'
import KeyboardArrowDownIcon from '@material-ui/icons/KeyboardArrowDown'
import NetflowVersionWidget from '../../widgets/NetflowVersion';
import _ from 'lodash'


const useStyles = makeStyles((theme) => ({
    root: {

    },
    paper: {
        padding: theme.spacing(2),
        textAlign: 'center',
        color: theme.palette.text.secondary,
    },
    loading: {
        marginLeft: theme.spacing(2)
    },
    backButton: {
        marginRight: theme.spacing(2),
    },
    h1: {
        fontSize: '20px',
        fontWeight: '700',
        textAlign: 'left',
        paddingRight: '10px',
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
    formPaper: {
        padding: theme.spacing(2),
    },
    dtHolder: {
        marginTop: theme.spacing(2),
        marginBottom: theme.spacing(2),
    },
}))


function Row(props) {
    const { row, keyID } = props;
    const [open, setOpen] = React.useState(false);

    return (
        <React.Fragment>
            <TableRow sx={{ '& > *': { borderBottom: 'unset' } }}>
                <TableCell>
                    <IconButton
                        aria-label="expand row"
                        size="small"
                        onClick={() => setOpen(!open)}
                    >
                        {open ? <KeyboardArrowUpIcon /> : <KeyboardArrowDownIcon />}
                    </IconButton>
                </TableCell>
                <TableCell >{keyID}{row.device}</TableCell>
                <TableCell ><NetflowVersionWidget version={row.flowVersion} /></TableCell>
                <TableCell >{row.protoName}</TableCell>
                <TableCell >{row.srcIP} : {row.srcPort}</TableCell>
                <TableCell >{row.dstIP} : {row.dstPort}</TableCell>
                <TableCell >{row.bytes} | {row.packets}</TableCell>
                <TableCell >
                    <Chip label="FIN" size="small" color={row.fin ? 'primary' : 'secondary'} />
                    <Chip label="SYN" size="small" color={row.syn ? 'primary' : 'secondary'} />
                    <Chip label="RST" size="small" color={row.rst ? 'primary' : 'secondary'} />
                    <Chip label="PSH" size="small" color={row.psh ? 'primary' : 'secondary'} />
                    <Chip label="ACK" size="small" color={row.ack ? 'primary' : 'secondary'} />
                    <Chip label="URG" size="small" color={row.urg ? 'primary' : 'secondary'} />
                    <Chip label="ECE" size="small" color={row.ece ? 'primary' : 'secondary'} />
                    <Chip label="CWR" size="small" color={row.cwr ? 'primary' : 'secondary'} />
                </TableCell>
            </TableRow>
            <TableRow>
                <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={6}>
                    <Collapse in={open} timeout="auto" unmountOnExit>
                        <Box sx={{ margin: 1 }}>
                            {/* <Typography variant="h6" gutterBottom component="div">
                  History
                </Typography> */}
                            {/* <Table size="small" aria-label="purchases">
                  <TableHead>
                    <TableRow>
                      <TableCell>Date</TableCell>
                      <TableCell>Customer</TableCell>
                      <TableCell align="right">Amount</TableCell>
                      <TableCell align="right">Total price ($)</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {row.history.map((historyRow) => (
                      <TableRow key={historyRow.date}>
                        <TableCell component="th" scope="row">
                          {historyRow.date}
                        </TableCell>
                        <TableCell>{historyRow.customerId}</TableCell>
                        <TableCell align="right">{historyRow.amount}</TableCell>
                        <TableCell align="right">
                          {Math.round(historyRow.amount * row.price * 100) / 100}
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table> */}
                            Details
                        </Box>
                    </Collapse>
                </TableCell>
            </TableRow>
        </React.Fragment>
    );
}


const socket = io(
    process.env.REACT_APP_WS,
    {
        transports: ["websocket"],
        autoConnect: false
    }
);


const getPaginatedItems = (items, page, pageSize) => {
    var pg = page || 1,
        pgSize = pageSize || 100,
        offset = (pg - 1) * pgSize,
        pagedItems = _.drop(items, offset).slice(0, pgSize);
    return {
        page: pg - 1,
        pageSize: pgSize,
        total: items.length,
        total_pages: Math.ceil(items.length / pgSize),
        data: pagedItems
    };
}

function LiveFlowComponent(props) {

    const [isConnected, setIsConnected] = useState(socket.connected);
    const [isFirstInit, setIsFirstInit] = useState(true);
    const [counter, setCounter] = useState(0);

    // eslint-disable-next-line
    const [page, setPage] = React.useState(1);
    const [rowsPerPage, setRowsPerPage] = React.useState(25);
    const [tableData, setTableData] = React.useState({
        page: 0,
        pageSize: 25,
        total: 0,
        total_pages: 0,
        data: []
    })
    const [rows, setRows] = useState([])
    const [newRows, setNewRows] = useState([])
    const classes = useStyles();


    // useEffect(() => {
    //     setRows([...newRows, ...rows])
    // },[newRows, rows])

    useEffect(() => {
        if (isFirstInit && rows.length > 0) {
            // console.log(isFirstInit, rows);
            setTableData(getPaginatedItems(rows, 0, 25))
            setIsFirstInit(false)
        }
    }, [isFirstInit, rows])


    useEffect(() => {
        setRows([...newRows, ...rows])
        setCounter(counter + 1)
        // eslint-disable-next-line
    }, [newRows])

    useEffect(() => {
        // console.log(tableData);
    }, [tableData])

    useEffect(() => {
        setTableData(getPaginatedItems(rows, page, rowsPerPage))
    }, [rows, page, rowsPerPage])

    // useEffect(() => {
    //     console.log(tableData);
    // }, [tableData])

    // const handleChangePage = (event, newPage) => {
    //     setPage(newPage);
    //     setTableData(getPaginatedItems(rows, newPage, rowsPerPage))
    // };
    const handleChangePage = (event, newPage) => {
        setPage(newPage + 1);
        console.log(page, newPage + 1);
        // setTableData(getPaginatedItems(rows, newPage, rowsPerPage))
    };


    const handleChangeRowsPerPage = (event) => {
        const nrpp = parseInt(event.target.value, 25)
        setRowsPerPage(nrpp);
        setPage(0);
        // setTableData(getPaginatedItems(rows, 0, nrpp))
    };




    useEffect(() => {
        socket.on('connect', () => {
            setIsConnected(true);
            socket.emit("join", "liveflow");
        });
        socket.on('disconnect', () => {
            setIsConnected(false);
        });
        socket.on('live-flow', (data) => {
            const flowData = JSON.parse(data)
            if (flowData.hasOwnProperty('payload')) {
                // console.log(flowData.payload);
                // console.log(rows);
                // const dt = [...flowData.payload, ...rows]
                // setRows(dt)
                // console.log("dt", dt);
                setNewRows(flowData.payload)
            }
        })
        return () => {
            socket.off('connect');
            socket.off('disconnect');
            socket.off('live-flow');
            socket.disconnect();
            setIsFirstInit(true)
        };
        // eslint-disable-next-line
    }, [])

    // eslint-disable-next-line
    const [busy, setBusy] = React.useState(false);

    const handleTogglePlayStop = (e) => {
        e.preventDefault()
        if (isConnected) {
            socket.disconnect()
        } else {
            socket.connect()
            setIsFirstInit(true)
            setRows([])
            setTableData(getPaginatedItems([], 0, rowsPerPage))
        }
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
                        {/* <BackButton
                            forObj="Dashboard"
                            backURL="/"
                            className={classes.backButton}
                        /> */}
                        <Typography
                            variant="h1"
                            color="primary"
                            className={classes.h1}
                        >
                            Live Flow
                            {
                                busy ? <CircularProgress className={classes.loading} color="primary" size={15} /> : ''
                            }

                        </Typography>
                        {
                            !isConnected ?
                                <Button
                                    color="primary"
                                    variant="contained"
                                    startIcon={<PlayCircleIcon />}
                                    onClick={handleTogglePlayStop}
                                >
                                    Click to play Live Flow
                                </Button>
                                :

                                <Button
                                    color="secondary"
                                    variant="contained"
                                    startIcon={<StopCircleIcon />}
                                    onClick={handleTogglePlayStop}
                                >
                                    Stop
                                </Button>
                        }
                    </Grid>
                </Grid>
            </Paper>

            <Grid container spacing={2} className={classes.dtHolder}>
                <Grid item xs={12} md={12} >
                    <Paper className={classes.formPaper}>
                        <div>
                            {
                                !isConnected ? <p>Live flow is not started. Please click on "Play" button to see live flows!</p> : ''
                            }
                        </div>

                        <TableContainer component={Paper}>
                            <Table size="small" aria-label="Live FLow Table">
                                <TableHead>
                                    <TableRow>
                                        <TableCell />
                                        <TableCell>Device</TableCell>
                                        <TableCell>Version</TableCell>
                                        <TableCell>Proto</TableCell>
                                        <TableCell>Source</TableCell>
                                        <TableCell>Destination</TableCell>
                                        <TableCell>Bytes | Packets</TableCell>
                                        <TableCell>TCPFlags</TableCell>
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    {
                                        tableData.data.map((row, i) => (
                                            <Row
                                                key={"r-" + i}
                                                // keyID={i + counter + ' '}
                                                keyID={' '}
                                                row={row}
                                            />
                                        ))
                                    }
                                </TableBody>
                            </Table>

                        </TableContainer>

                        <TablePagination
                            component="div"
                            count={rows.length}
                            page={tableData.page}
                            // onPageChange={handleChangePage}
                            onChangePage={handleChangePage}
                            rowsPerPage={tableData.pageSize}
                            onRowsPerPageChange={handleChangeRowsPerPage}
                            rowsPerPageOptions={[5, 10, 25, 50, 100, 150, 200, 400]}
                        />
                        {/* <p><b>{tableData.page}</b></p> */}
                    </Paper>
                </Grid>
            </Grid>
        </div>
    );
}

export default LiveFlowComponent;