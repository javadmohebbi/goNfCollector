import React from 'react';
import { Link as RouterLink } from 'react-router-dom';

import {
    Chip,
    Link,
} from '@material-ui/core';

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

export default IPPortWidget;