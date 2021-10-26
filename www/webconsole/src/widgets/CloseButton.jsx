import React from 'react';
import Button from '@material-ui/core/Button';

import { withStyles } from '@material-ui/core/styles';
import { amber } from '@material-ui/core/colors';
import CloseIcon from '@material-ui/icons/Close';
import { Tooltip } from '@material-ui/core';

const CustomButton = withStyles((theme) => ({
    root: {
        color: theme.palette.getContrastText(amber[500]),
        backgroundColor: amber[500],
        '&:hover': {
            backgroundColor: amber[700],
        },
    },
}))(Button);


const CloseButton = ({ btn, forObj = '', onClick = undefined }) => {



    return (
        <Tooltip title={`Close`}>
            <CustomButton variant="contained" color="primary" onClick={onClick}>
                <CloseIcon />
            </CustomButton>
        </Tooltip>
    );
}

export default CloseButton;