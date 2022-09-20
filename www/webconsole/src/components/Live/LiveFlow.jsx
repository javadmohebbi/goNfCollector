import React, { useEffect, useState } from 'react';
import io from 'socket.io-client';
import { makeStyles } from '@material-ui/core/styles';
import { Button, Chip, CircularProgress, Grid, Paper, Table, TableBody, TableCell, TableContainer, TableHead, TablePagination, TableRow, Typography } from '@material-ui/core';

// import PlayCircleFilledWhiteIcon from '@material-ui/icons/PlayCircleFilledWhite';
// import StopCircleIcon from '@material-ui/icons/StopCircle';

import PlayCircleIcon from '@material-ui/icons/PlayArrow'
import StopCircleIcon from '@material-ui/icons/Stop'
import DownloadIcon from '@material-ui/icons/FileCopy'

import _ from 'lodash'
import RowComponent from './RowComponent';
import humanFormat from 'human-format';


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


const HeaderFooter = ({ recordsCount, sumBytes, sumPackets }) => {
    return <div style={{ marginTop: '15px', marginBottom: '15px' }}>
        <Chip style={{ fontSize: '120%', marginRight: '10px' }}
            color="primary"
            label={`Flow Counts: ${humanFormat(recordsCount)}`}
        />
        <Chip style={{ fontSize: '120%', marginRight: '10px' }}
            color="primary"
            label={`Total Bytes: ${humanFormat(sumBytes, { unit: 'B' })}`}
        />
        <Chip style={{ fontSize: '120%', marginRight: '10px' }}
            color="primary"
            label={`Total Packets: ${humanFormat(sumPackets)}`}
        />
    </div>
}


function LiveFlowComponent(props) {

    const [isConnected, setIsConnected] = useState(socket.connected);
    const [isFirstInit, setIsFirstInit] = useState(true);
    const [counter, setCounter] = useState(0);

    const [sumBytes, setSumBytes] = useState(0)
    const [sumPackets, setSumPackets] = useState(0)
    const [recordsCount, setRecordsCount] = useState(0)

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
        setRows([...newRows, ...rows]);
        setCounter(counter + 1);
        const byts = _.sumBy(newRows, (r) => parseInt(r.bytes))
        const pkts = _.sumBy(newRows, (r) => parseInt(r.packets))
        const rcds = newRows.length
        setSumBytes(byts + sumBytes)
        setSumPackets(pkts + sumPackets)
        setRecordsCount(rcds + recordsCount)
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
            setSumBytes(0)
            setSumPackets(0)
            setIsFirstInit(true)
            setRows([])
            setTableData(getPaginatedItems([], 0, rowsPerPage))
        }
    }

    const downloadFile = ({ data, fileName, fileType }) => {
        // Create a blob with the data we want to download as a file
        const blob = new Blob([data], { type: fileType })
        // Create an anchor element and dispatch a click event on it
        // to trigger a download
        const a = document.createElement('a')
        a.download = fileName
        a.href = window.URL.createObjectURL(blob)
        const clickEvt = new MouseEvent('click', {
            view: window,
            bubbles: true,
            cancelable: true,
        })
        a.dispatchEvent(clickEvt)
        a.remove()
    }

    const handleDownloadCsv = (e) => {
        e.preventDefault()

        let headers = []
        const records = [];

        for (let i = 0; i < rows.length; i++) {
            if (i === 0) {
                headers.push(Object.keys(rows[i]).join(','))
            }
            var vals = _.values(rows[i])
            vals.forEach((element, index, array) => {
                element = "\"" + element.toString() + "\"";
                array[index] = element;
            })
            // console.log(vals);
            records.push(vals.join(','))
        }


        downloadFile({
            // data: JSON.stringify(rows),
            data: [...headers, ...records].join('\r\n'),
            fileName: 'liveflow-exported.csv',
            fileType: 'text/csv',
        })
    }

    const handleDownloadJson = (e) => {
        e.preventDefault()
        downloadFile({
            data: JSON.stringify(rows),
            fileName: 'liveflow-exported.json',
            fileType: 'text/json',
        })
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
                        {
                            !isConnected && rows.length > 0 ?
                                <>
                                    <Button
                                        style={{ marginLeft: '10px' }}
                                        color="secondary"
                                        variant="contained"
                                        startIcon={<DownloadIcon />}
                                        onClick={handleDownloadCsv}
                                    >
                                        Download CSV
                                    </Button>
                                    <Button
                                        style={{ marginLeft: '10px' }}
                                        color="secondary"
                                        variant="contained"
                                        startIcon={<DownloadIcon />}
                                        onClick={handleDownloadJson}
                                    >
                                        Download JSON
                                    </Button>
                                </>

                                : <></>
                        }
                    </Grid>
                </Grid>
            </Paper>

            <Grid container spacing={2} className={classes.dtHolder}>
                <Grid item xs={12} md={12} >
                    <Paper className={classes.formPaper}>
                        <div>
                            {
                                !isConnected
                                    ?
                                    <p>Live flow is not running. Please click on "Play" button to see live flows!</p>
                                    :
                                    <></>
                            }
                            {
                                rows.length > 0 ?
                                    <HeaderFooter
                                        recordsCount={recordsCount}
                                        sumBytes={sumBytes}
                                        sumPackets={sumPackets}
                                    />
                                    : <></>
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
                                        <TableCell>TCP Flags</TableCell>
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    {
                                        tableData.data.map((row, i) => (
                                            <RowComponent
                                                isCapturing={isConnected}
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

                        {
                            rows.length > 0 ?
                                <HeaderFooter
                                    recordsCount={recordsCount}
                                    sumBytes={sumBytes}
                                    sumPackets={sumPackets}
                                />
                                : <></>
                        }

                        < TablePagination
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