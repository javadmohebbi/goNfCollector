import React from 'react';
import Button from '@material-ui/core/Button';

import { withStyles } from '@material-ui/core/styles';
import { red } from '@material-ui/core/colors';
import DeleteForeverIcon from '@material-ui/icons/DeleteForever';
import { Tooltip } from '@material-ui/core';


const CustomButton = withStyles((theme) => ({
    root: {
        color: theme.palette.getContrastText(red[500]),
        backgroundColor: red[500],
        '&:hover': {
            backgroundColor: red[700],
        },
    },
}))(Button);


const DeleteButton = ({ btn, forObj = '', onClick = undefined }) => {


    return (
        <Tooltip title={`Delete ${forObj !== '' ? '"' + forObj + '"' : ''}`}>
            <CustomButton variant="contained" color="primary" onClick={onClick}>
                <DeleteForeverIcon />
            </CustomButton>
        </Tooltip>
    );

}

export default DeleteButton;