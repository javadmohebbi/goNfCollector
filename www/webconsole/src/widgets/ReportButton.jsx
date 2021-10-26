import React from 'react';
import Button from '@material-ui/core/Button';

import { withStyles } from '@material-ui/core/styles';
import { orange } from '@material-ui/core/colors';
import PieChartIcon from '@material-ui/icons/PieChart';
import { Tooltip } from '@material-ui/core';


const CustomButton = withStyles((theme) => ({
    root: {
        color: theme.palette.getContrastText(orange[500]),
        backgroundColor: orange[500],
        '&:hover': {
            backgroundColor: orange[700],
        },
    },
}))(Button);



const ReportButton = ({ btn, forObj = '', onClick = undefined }) => {

    return (
        <Tooltip title={`Show report ${forObj !== '' ? 'on "' + forObj + '"' : ''}`}>
            <CustomButton variant="contained" color="primary" onClick={onClick}>
                <PieChartIcon />
            </CustomButton>
        </Tooltip>
    );
}

export default ReportButton;