import React from 'react';
import Button from '@material-ui/core/Button';

import { useHistory } from 'react-router-dom'

import { withStyles } from '@material-ui/core/styles';
import { grey } from '@material-ui/core/colors';
import ArrowBackIcon from '@material-ui/icons/ArrowBack';
import { Tooltip } from '@material-ui/core';

const CustomButton = withStyles((theme) => ({
    root: {
        color: theme.palette.getContrastText(grey[500]),
        backgroundColor: grey[500],
        '&:hover': {
            backgroundColor: grey[700],
        },
    },
}))(Button);


const BackButton = ({ btn, forObj = '', onClick = undefined, className = {}, backURL = "/", goback = false }) => {


    const history = useHistory();

    return (
        <Tooltip title={goback ? 'Back' : `Back to ${forObj !== '' ? '"' + forObj + '"' : ''}`}>
            <CustomButton variant="contained" color="primary" onClick={
                e => {
                    e.preventDefault();
                    if (goback) {
                        history.goBack()
                    } else {
                        history.push(backURL);
                    }
                }
            }
                className={className}
            >
                <ArrowBackIcon />
            </CustomButton>
        </Tooltip>
    );
}

export default BackButton;