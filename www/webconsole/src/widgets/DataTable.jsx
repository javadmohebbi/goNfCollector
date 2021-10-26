import React, { useEffect, useState } from 'react';
import PropTypes from 'prop-types';
import { Link } from 'react-router-dom';
import { makeStyles, withStyles, useTheme } from '@material-ui/core/styles';
import { FormControl, InputAdornment, Paper, Table, TableBody, TableCell, TableContainer, TableHead, TablePagination, TableRow } from '@material-ui/core';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import InputLabel from '@material-ui/core/InputLabel';
import SearchIcon from '@material-ui/icons/Search'
import Alert from '@material-ui/lab/Alert';

import IconButton from '@material-ui/core/IconButton';
import FirstPageIcon from '@material-ui/icons/FirstPage';
import KeyboardArrowLeft from '@material-ui/icons/KeyboardArrowLeft';
import KeyboardArrowRight from '@material-ui/icons/KeyboardArrowRight';
import LastPageIcon from '@material-ui/icons/LastPage'

import axios from 'axios';

const useStyles = makeStyles((theme) => ({
    root: {
        width: '100%'
    },
    theadcell: {
        fontWeight: '900',
        // background: "#f14624",
        // color: "black",
    },
    fetchError: {
        color: 'red'
    }
}))


var timerId;
var CancelToken = axios.CancelToken;
var cancel;



const StyledTableRow = withStyles((theme) => ({
    root: {
        '&:nth-of-type(odd)': {
            backgroundColor: '#7B7B7B',
        },
        '&:nth-of-type(odd):hover': {
            backgroundColor: '#505050',
        },
        '&:hover': {
            backgroundColor: '#252525',
        },

    },
}))(TableRow);


const ReFormat = (callback, v, row) => {
    if (typeof callback !== 'undefined') {
        return callback(v, row)
    } else {
        return v
    }
}


const useStyles1 = makeStyles((theme) => ({
    root: {
        flexShrink: 0,
        marginLeft: theme.spacing(2.5),
    },
}));

function TablePaginationActions(props) {
    const classes = useStyles1();
    const theme = useTheme();
    const { count, page, rowsPerPage, onPageChange } = props;

    const handleFirstPageButtonClick = (event) => {
        onPageChange(event, 0);
    };

    const handleBackButtonClick = (event) => {
        onPageChange(event, page - 1);
    };

    const handleNextButtonClick = (event) => {
        onPageChange(event, page + 1);
    };

    const handleLastPageButtonClick = (event) => {
        onPageChange(event, Math.max(0, Math.ceil(count / rowsPerPage) - 1));
    };

    return (
        <div className={classes.root}>
            <IconButton
                onClick={handleFirstPageButtonClick}
                disabled={page === 0}
                aria-label="first page"
            >
                {theme.direction === 'rtl' ? <LastPageIcon /> : <FirstPageIcon />}
            </IconButton>
            <IconButton onClick={handleBackButtonClick} disabled={page === 0} aria-label="previous page">
                {theme.direction === 'rtl' ? <KeyboardArrowRight /> : <KeyboardArrowLeft />}
            </IconButton>
            <IconButton
                onClick={handleNextButtonClick}
                disabled={page >= Math.ceil(count / rowsPerPage) - 1}
                aria-label="next page"
            >
                {theme.direction === 'rtl' ? <KeyboardArrowLeft /> : <KeyboardArrowRight />}
            </IconButton>
            <IconButton
                onClick={handleLastPageButtonClick}
                disabled={page >= Math.ceil(count / rowsPerPage) - 1}
                aria-label="last page"
            >
                {theme.direction === 'rtl' ? <FirstPageIcon /> : <LastPageIcon />}
            </IconButton>
        </div>
    );
}
TablePaginationActions.propTypes = {
    count: PropTypes.number.isRequired,
    onPageChange: PropTypes.func.isRequired,
    page: PropTypes.number.isRequired,
    rowsPerPage: PropTypes.number.isRequired,
};



const ActionButtons = ({ buttons, row }) => (
    <>
        {
            buttons.map((btn, i) => (
                <React.Fragment key={i}>
                    {
                        typeof btn.url !== 'undefined'
                            ?
                            <Link to={btn.url(row)}>
                                {btn.component(btn, row)}
                            </Link>
                            :
                            btn.component(btn, row)
                    }

                </React.Fragment>
            ))
        }
    </>
)


const WidgetDataTableComponent = ({
    busy = false,
    refresh = false,

    handleParentBusyState = () => { return },
    handleParentRefreshState = () => { return },
    Page = 1,
    PerPage = process.env.REACT_NORMAL_DT_PER_PAGE,
    Order = '',
    OrderType = '',
    Filter = '',
    URL = '',
    Columns = [],
    NeedFilter = true,
}) => {

    const classes = useStyles();


    const [data, setData] = useState([]);
    const [loading, setLoading] = useState(false);
    const [totalRows, setTotalRows] = useState(0);

    const [perPage, setPerPage] = useState(PerPage || process.env.REACT_NORMAL_DT_PER_PAGE);
    const [filter, setFilter] = useState(Filter || '')

    // eslint-disable-next-line
    const [order, setOrder] = useState(Order || '')

    // eslint-disable-next-line
    const [orderType, setOrderType] = useState(OrderType || 'desc')
    const [thePage, setThePage] = useState(Page || 1);

    // eslint-disable-next-line
    const [hasFetchError, setHasFetchError] = useState(false)


    const GetDataFromServer = (page, perPage, order, orderType, filter) => {
        clearTimeout(timerId);
        if (typeof cancel === 'function') {
            cancel()
        }
        timerId = setTimeout(() => {
            fetchData(page, perPage, order, orderType, filter)
        }, filter === '' ? 500 : 1000)
    }


    const getURL = (`${process.env.REACT_APP_HTTP}/${URL}`)

    const fetchData = async (page, perPage, order = '', orderType = '', filter = '') => {
        handleParentBusyState(true)
        setLoading(true);

        // eslint-disable-next-line
        const response = await axios.get(
            getURL,
            {
                params: {
                    page: page,
                    perPage: perPage,
                    order: order,
                    orderType: orderType,
                    filter: filter
                },
                headers: {
                    'Content-Type': 'application/json',
                    "FromPage": window.location,
                },
                cancelToken: new CancelToken(function executor(c) {
                    // An executor function receives a cancel function as a parameter
                    cancel = c;
                }),
            }
        ).then(async (response) => {
            // const src = (`?page=${page}&perPage=${perPage}&order=${order}&orderType=${orderType}&filter=${filter}`)
            // history.push({
            //     pathname: history.location.pathname,
            //     search: src
            // })
            setThePage(page)


            // const resp = await response.then((result) => result);
            // console.log(resp);
            // console.log(response);

            setData(response.data.Result);
            setTotalRows(response.data.Pagination.total);
            setLoading(false);
            handleParentBusyState(false)
            setHasFetchError(false)
        })
            .catch((error) => {
                if (typeof error.response === 'undefined') {
                    setHasFetchError(true)
                } else {
                    if (error.response.status === 404) {
                        setHasFetchError(false)
                    } else {
                        setHasFetchError(true)
                    }
                }
                setData([])
                setLoading(false);
                handleParentBusyState(false)
            });
    };

    // eslint-disable-next-line
    const forceReaload = () => {
        GetDataFromServer(thePage, perPage, order, orderType, filter)
    }


    const handlePageChange = (event, newPage) => {

        setThePage(newPage + 1)
        // fetchData(page);
        GetDataFromServer(newPage + 1, perPage, order, orderType, filter)
    };

    // const handlePerRowsChange = async (newPerPage, page) => {
    //     GetDataFromServer(thePage, newPerPage, order, orderType, filter)
    //     setPerPage(0);
    // };
    const handlePerRowsChange = event => {
        console.log(event);
        const newPerPage = parseInt(event.target.value, 10)
        setPerPage(newPerPage);
        setThePage(1);
        GetDataFromServer(thePage, newPerPage, order, orderType, filter)
    };


    const handleFilterChange = (e) => {
        setFilter(e.target.value)
        setThePage(1)
    }

    useEffect(() => {
        // fetchData(thePage);
        GetDataFromServer(thePage, perPage, order, orderType, filter)
        // eslint-disable-next-line
    }, []);

    useEffect(() => {
        // fetchData(thePage)
        GetDataFromServer(thePage, perPage, order, orderType, filter)
        // eslint-disable-next-line
    }, [filter])


    useEffect(() => {
        // fetchData(thePage)
        handleParentRefreshState(false)
        GetDataFromServer(thePage, perPage, order, orderType, filter)
        // eslint-disable-next-line
    }, [refresh])




    // useEffect(() => {
    //     DeviceGetAllPaginate()
    // }, [])
    const cells = Columns.map(c => c.key)





    return (
        <div className={classes.root} >


            <Paper>
                <TableContainer>
                    <Table size="small">
                        {/* HEADER */}
                        <TableHead>

                            {
                                NeedFilter
                                    ?
                                    <TableRow>
                                        <TableCell colSpan={cells.length}>
                                            <FormControl fullWidth variant="outlined">
                                                <InputLabel htmlFor="filterTable">Filter Table</InputLabel>
                                                <OutlinedInput
                                                    id="filterTable"
                                                    value={filter}
                                                    onChange={handleFilterChange}
                                                    startAdornment={<InputAdornment position="start"><SearchIcon /></InputAdornment>}
                                                    labelWidth={(("Filter Table".length - 4) * 10) + 1}
                                                />
                                            </FormControl>
                                        </TableCell>
                                    </TableRow>
                                    : <></>
                            }
                            <TableRow>
                                {Columns.map(col => (
                                    <TableCell className={classes.theadcell} key={col.key} align="left">
                                        {col.label}
                                    </TableCell>
                                ))}
                            </TableRow>
                        </TableHead>


                        {/* BODY */}
                        <TableBody style={{
                            opacity: loading ? '0.5' : '1.0'
                        }}>

                            {
                                !loading ? data.map((v, i) => (
                                    <StyledTableRow key={i}>
                                        {
                                            cells.map((c, j) => (
                                                <TableCell key={j} >
                                                    {
                                                        typeof Columns[j].dataType !== 'undefined'
                                                            && Columns[j].dataType === 'actionButton'
                                                            ?
                                                            <ActionButtons buttons={Columns[j].buttons.list} row={data[i]} />
                                                            :
                                                            <span style={{
                                                                fontWeight: Columns[j].important ? '900' : '500',
                                                                fontSize: Columns[j].important ? '110%' : '100%',

                                                            }}>
                                                                {ReFormat(Columns[j].reFormat || undefined, v[c], v)}
                                                            </span>
                                                    }

                                                </TableCell>
                                            ))
                                        }

                                    </StyledTableRow>
                                ))
                                    :
                                    ''
                            }
                            {
                                loading ?
                                    <>
                                        {
                                            Array.from(Array(perPage).keys()).map((r, ri) => (
                                                <TableRow key={ri}>
                                                    {cells.map((c, j) => (
                                                        <TableCell key={j} >
                                                            . . .
                                                        </TableCell>
                                                    ))}
                                                </TableRow>
                                            ))
                                        }
                                    </>
                                    :
                                    <></>
                            }
                            {
                                !loading && data.length === 0 ?
                                    <TableRow>
                                        <TableCell colSpan={cells.length} align="center">
                                            {
                                                hasFetchError
                                                    ?
                                                    // <span className={classes.fetchError}>Error: Could not get data from server!</span>
                                                    <Alert variant="filled" severity="error" >Error: Could not get data from server!</Alert>
                                                    :
                                                    <Alert variant="filled" severity="info" >There is no records! {filter !== '' ? `Please change your filter "${filter}" and try again.` : ''}</Alert>
                                            }

                                        </TableCell>
                                    </TableRow>
                                    :
                                    <></>
                            }
                        </TableBody>


                    </Table>
                    <TablePagination

                        // begin custom
                        ActionsComponent={() => <TablePaginationActions
                            count={totalRows}
                            page={thePage - 1}
                            rowsPerPage={perPage}
                            onPageChange={handlePageChange}
                        />}
                        SelectProps={{
                            inputProps: { 'aria-label': 'rows per page' },
                            native: true,
                        }}
                        // end custom


                        component="div"
                        count={totalRows}
                        page={thePage - 1}
                        // onChangePage={handlePageChange}
                        rowsPerPage={perPage}
                        onChangeRowsPerPage={handlePerRowsChange}
                        rowsPerPageOptions={[5, 10, 15, 25, 50, 100, 200]}




                    />
                </TableContainer>
            </Paper>
        </div >
    )
}

export default WidgetDataTableComponent;