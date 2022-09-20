import { Chip } from '@material-ui/core';
import humanFormat from 'human-format';
import React from 'react';

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


export default BytesWidget;