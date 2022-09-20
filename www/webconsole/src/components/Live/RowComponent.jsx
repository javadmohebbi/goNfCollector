import React from 'react';


import {
    Box,
    Chip,
    Collapse,
    IconButton,
    TableCell,
    TableRow,
    Link,
    Typography
} from '@material-ui/core';

import KeyboardArrowUpIcon from '@material-ui/icons/KeyboardArrowUp'
import KeyboardArrowDownIcon from '@material-ui/icons/KeyboardArrowDown'
import NetflowVersionWidget from '../../widgets/NetflowVersion';
import humanFormat from 'human-format';

import { Link as RouterLink } from 'react-router-dom';
import ThreatLevelWidget from '../../widgets/ThreatLevel';


const IPPortWidget = ({ ip, proto, port, portName }) => {
    const prt = proto.toLowerCase()
    const prtProto = `${prt}/${port}`

    if (prtProto === portName) {
        return <>
            <Link component={RouterLink} to={`/hosts/${ip}/report`} target='_blank'>
                {ip}
            </Link>
            <Chip size="small" label={port} />
        </>
    } else {
        return <>
            <Link component={RouterLink} to={`/hosts/${ip}/report`} target='_blank'>
                {ip}
            </Link>
            <Chip size="small" label={portName} />
        </>
    }
}

const BytesWidget = ({ bytes }) => {
    return <Chip size="small" color="secondary"
        label={
            humanFormat(
                parseInt(bytes),
                { unit: 'B' }
            )
        }
    />
}

const PacketWidget = ({ packets }) => {
    return <Chip size="small" color="secondary"
        label={
            humanFormat(
                parseInt(packets),
            )
        }
    />
}

const TCPFlagsWidgets = ({ fin, syn, rst, psh, ack, urg, ece, cwr }) => {
    return <>
        <Chip label="FIN" size="small" color={fin ? 'primary' : 'secondary'} />
        <Chip label="SYN" size="small" color={syn ? 'primary' : 'secondary'} />
        <Chip label="RST" size="small" color={rst ? 'primary' : 'secondary'} />
        <Chip label="PSH" size="small" color={psh ? 'primary' : 'secondary'} />
        <Chip label="ACK" size="small" color={ack ? 'primary' : 'secondary'} />
        <Chip label="URG" size="small" color={urg ? 'primary' : 'secondary'} />
        <Chip label="ECE" size="small" color={ece ? 'primary' : 'secondary'} />
        <Chip label="CWR" size="small" color={cwr ? 'primary' : 'secondary'} />
    </>
}

const RowComponent = (props) => {
    const { row, keyID } = props;
    const [open, setOpen] = React.useState(false);

    console.log(row);

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
                <TableCell >
                    <IPPortWidget ip={row.srcIP} proto={row.protoName}
                        port={row.srcPort} portName={row.srcPortName}
                    />
                    {
                        row.sCountryShort.toLowerCase() !== 'NA' ?
                            <Typography component={'p'} style={{ fontSize: '80%' }}>
                                {`${row.sCity} (${row.sCountryShort})`}
                            </Typography>
                            : <></>
                    }
                    {
                        row.isSrcThreat ?
                            <div style={{ marginTop: '10px', marginBottom: '10px' }}>
                                <ThreatLevelWidget
                                    inline={true}
                                    level={row.srcThreatReputation}
                                    label={`(${row.srcThreatType})`}
                                />
                            </div>
                            : <></>
                    }
                </TableCell>
                <TableCell >
                    <IPPortWidget ip={row.dstIP} proto={row.protoName}
                        port={row.dstPort} portName={row.dstPortName}
                    />
                    {
                        row.dCountryShort.toLowerCase() !== 'NA' ?
                            <Typography component={'p'} style={{ fontSize: '80%' }}>
                                {`${row.dCity} (${row.dCountryShort})`}
                            </Typography>
                            : <></>
                    }
                    {
                        row.isDstThreat ?
                            <div style={{ marginTop: '10px', marginBottom: '10px' }}>
                                <ThreatLevelWidget
                                    inline={true}
                                    level={row.dstThreatReputation}
                                    label={`(${row.dstThreatType})`}
                                />
                            </div>
                            : <></>
                    }
                </TableCell>
                <TableCell ><BytesWidget bytes={row.bytes} /> | <PacketWidget packets={row.packets} /></TableCell>
                <TableCell >
                    {
                        row.protoName.toLowerCase() === 'tcp' ?
                            <>
                                <TCPFlagsWidgets
                                    fin={row.fin}
                                    syn={row.syn}
                                    rst={row.rst}
                                    psh={row.psh}
                                    ack={row.ack}
                                    urg={row.urg}
                                    ece={row.ece}
                                    cwr={row.cwr}
                                />
                            </>
                            :
                            '-'
                    }
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


export default RowComponent;