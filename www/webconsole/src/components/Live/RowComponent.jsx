import React from 'react';


import {
    Box,
    Collapse,
    IconButton,
    Link,
    TableCell,
    TableRow,
    Typography
} from '@material-ui/core';

import KeyboardArrowUpIcon from '@material-ui/icons/KeyboardArrowUp'
import KeyboardArrowDownIcon from '@material-ui/icons/KeyboardArrowDown'
import NetflowVersionWidget from '../../widgets/NetflowVersion';


import ThreatLevelWidget from '../../widgets/ThreatLevel';
import DetailsComponent from './DetailsComponent';
import BytesWidget from './BytesWidget';
import PacketsWidget from './PacketsWidget';
import IPPortWidget from './IPPortWidget';
import TCPFlagsWidgets from './TCPFlagsWidgets';

import { Link as RouterLink } from 'react-router-dom';




const RowComponent = (props) => {
    const { row, isCapturing } = props;
    const [open, setOpen] = React.useState(false);

    return (
        <React.Fragment>
            <TableRow sx={{ '& > *': { borderBottom: 'unset' } }}>
                <TableCell>
                    <IconButton
                        aria-label="expand row"
                        size="small"
                        onClick={() => setOpen(!open)}
                        disabled={isCapturing}
                    >
                        {open ? <KeyboardArrowUpIcon /> : <KeyboardArrowDownIcon />}
                    </IconButton>
                </TableCell>
                <TableCell >
                    <Link
                        style={{ fontSize: '125%' }}
                        component={RouterLink}
                        to={`/devices/${row.device}/report`}
                        target='_blank'
                    >
                        {row.device}
                    </Link>
                </TableCell>
                <TableCell ><NetflowVersionWidget version={row.flowVersion} /></TableCell>
                <TableCell >{row.protoName}</TableCell>
                <TableCell >
                    <IPPortWidget ip={row.srcIP} proto={row.protoName}
                        port={row.srcPort} portName={row.srcPortName}
                    />
                    {
                        row.sCountryShort.toLowerCase() !== 'na' ?
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
                        row.dCountryShort.toLowerCase() !== 'na' ?
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
                <TableCell ><BytesWidget bytes={row.bytes} /> | <PacketsWidget packets={row.packets} /></TableCell>
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
                <TableCell style={{ paddingBottom: 0, paddingTop: 0, backgroundColor: 'gray' }} colSpan={8}>
                    <Collapse in={open} timeout="auto" unmountOnExit>
                        <Box sx={{ margin: 1 }}>
                            <DetailsComponent
                                row={row} />
                        </Box>
                    </Collapse>
                </TableCell>
            </TableRow>
        </React.Fragment>
    );
}


export default RowComponent;