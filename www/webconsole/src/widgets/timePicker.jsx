import React from 'react'
import { makeStyles } from '@material-ui/core/styles';
import Popover from '@material-ui/core/Popover';
import Button from '@material-ui/core/Button';
import ButtonGroup from '@material-ui/core/ButtonGroup';
// import { CircularProgress } from '@material-ui/core';

const useStyles = makeStyles((theme) => ({
    typography: {
        padding: theme.spacing(2),
    },
    btGroupRoot: {
        textTransform: 'none',
        display: 'flex',
        '& > *': {
            margin: theme.spacing(1),
        },
    },
    noTransForm: {
        textTransform: 'none',
    },
}));

const minutes = [5, 15, 30, 45, 60]
const hours = [1, 2, 3, 4, 6, 12, 24]
const days = [1, 2, 3, 4, 5, 6]
const weeks = [1, 2, 3]
const months = [1, 3, 6, 9, 12]
const years = [1, 2, 3]

const getTimeRanges = () => {
    const mins = minutes.map(m => `${m}m`)
    const hrs = hours.map(m => `${m}h`)
    const dys = days.map(m => `${m}d`)
    const wks = weeks.map(m => `${m}w`)
    const mons = months.map(m => `${m}M`)
    const yrs = years.map(m => `${m}y`)

    return [
        ...mins, ...hrs, ...dys, ...wks, ...mons, ...yrs,
    ]
}

const TimePickerWidget = ({
    // call back for onTimeRangeChange
    onTimeRangeChange = undefined,

    // disabled state
    disabled = false,

    // busy state
    busy = false,

    defaultSelected = '15m'
}) => {


    // material ui classes
    const classes = useStyles();


    const [anchorEl, setAnchorEl] = React.useState(null);

    // selected value
    const [selected, setSelected] = React.useState(getTimeRanges()[getTimeRanges().indexOf(defaultSelected)])



    const handleClick = (event) => {
        setAnchorEl(event.currentTarget);
    };

    const handleClose = () => {
        setAnchorEl(null);
    };

    const handleTimeRangeButtonClick = (tr) => {
        setAnchorEl(null);
        setSelected(tr)


        if (typeof onTimeRangeChange !== 'undefined') {
            onTimeRangeChange(tr)
        }
    }

    const open = Boolean(anchorEl);
    const id = open ? 'simple-popover' : undefined;

    return (
        <div>
            <Button
                disabled={disabled || busy}
                className={classes.noTransForm}
                aria-describedby={id}
                variant="contained"
                color="primary"
                onClick={handleClick}>
                {selected}
            </Button>
            <Popover
                id={id}
                open={open}
                anchorEl={anchorEl}
                onClose={handleClose}
                anchorOrigin={{
                    vertical: 'bottom',
                    horizontal: 'center',
                }}
                transformOrigin={{
                    vertical: 'top',
                    horizontal: 'center',
                }}
            >
                {/* <Typography className={classes.typography}>{
                    getTimeRanges().map((t, i) => (
                        <span key={i}>{t}</span>)
                    )
                }</Typography> */}

                <div className={classes.btGroupRoot}>
                    {[
                        { lbl: "Minute", ch: "m", arr: minutes },
                        { lbl: "Hour", ch: "h", arr: hours },
                        { lbl: "Day", ch: "d", arr: days },
                        { lbl: "Week", ch: "w", arr: weeks },
                        { lbl: "Month", ch: "M", arr: months },
                        { lbl: "Year", ch: "y", arr: years }
                    ].map((items, i) => (
                        <ButtonGroup
                            key={i + items.ch}
                            orientation="vertical"
                            color="primary"
                            aria-label="vertical outlined primary button group"
                        >
                            {
                                items.arr.map((t, j) => (
                                    <Button
                                        variant={t + items.ch === selected ? 'contained' : 'outlined'}
                                        className={classes.noTransForm}
                                        key={j}
                                        onClick={(e) => { e.preventDefault(); handleTimeRangeButtonClick(t + items.ch) }}
                                    >
                                        {t + items.ch}
                                    </Button>
                                ))
                            }
                        </ButtonGroup>
                    ))}


                </div>
            </Popover>
        </div>
    )
}

export default TimePickerWidget