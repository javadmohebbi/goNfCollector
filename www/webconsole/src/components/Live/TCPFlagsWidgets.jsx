import { Chip } from '@material-ui/core';
import React from 'react';

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

export default TCPFlagsWidgets;