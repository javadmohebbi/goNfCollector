import React from 'react';
import Button from '@material-ui/core/Button';

import { withStyles } from '@material-ui/core/styles';
import { lime } from '@material-ui/core/colors';
import MapIcon from '@material-ui/icons/Map';
import { Tooltip } from '@material-ui/core';

const CustomButton = withStyles((theme) => ({
    root: {
        color: theme.palette.getContrastText(lime[500]),
        backgroundColor: lime[500],
        '&:hover': {
            backgroundColor: lime[700],
        },
    },
}))(Button);


const MapButton = ({ btn, forObj = '', onClick = undefined }) => {



    return (
        <Tooltip title={`Locate on the map ${forObj !== '' ? '"' + forObj + '"' : ''}`}>
            <CustomButton variant="contained" color="primary" onClick={onClick}>
                <MapIcon />
            </CustomButton>
        </Tooltip>
    );
}

export default MapButton;