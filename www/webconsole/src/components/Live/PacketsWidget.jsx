import { Chip } from '@material-ui/core';
import humanFormat from 'human-format';
import React from 'react';


const PacketsWidget = ({ packets }) => {
    return <Chip size="small" color="secondary"
        label={
            humanFormat(
                parseInt(packets),
            )
        }
    />
}

export default PacketsWidget;