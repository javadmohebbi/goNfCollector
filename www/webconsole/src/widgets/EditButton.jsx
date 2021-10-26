import React from 'react';
import Button from '@material-ui/core/Button';

import { withStyles } from '@material-ui/core/styles';
import { yellow } from '@material-ui/core/colors';
import EditIcon from '@material-ui/icons/Edit';
import { Tooltip } from '@material-ui/core';

const CustomButton = withStyles((theme) => ({
    root: {
        color: theme.palette.getContrastText(yellow[500]),
        backgroundColor: yellow[500],
        '&:hover': {
            backgroundColor: yellow[700],
        },
    },
}))(Button);


const EditButton = ({ btn, forObj = '', onClick = undefined }) => {


    return (
        <Tooltip title={`Edit ${forObj !== '' ? '"' + forObj + '"' : ''}`}>
            <CustomButton variant="contained" color="primary" onClick={onClick}>
                <EditIcon />
            </CustomButton>
        </Tooltip>
    );
}

export default EditButton;