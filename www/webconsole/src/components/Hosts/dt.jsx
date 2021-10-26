import React from 'react';
import { makeStyles } from '@material-ui/core/styles';

import WidgetDataTableComponent from '../../widgets/DataTable';

import moment from "moment";

import ReportButton from '../../widgets/ReportButton';
import EditButton from '../../widgets/EditButton';




const useStyles = makeStyles((theme) => ({
    root: {

    },
}))





const HostsDataTableComponent = ({
    busy = false,
    refresh = false,
    handleParentBusyState = () => { return },
    handleParentRefreshState = () => { return },

}) => {

    const classes = useStyles();


    return (
        <div className={classes.root}>
            <WidgetDataTableComponent
                refresh={refresh}

                URL="host/get/all"
                Columns={[
                    { key: "host", label: "Host IP Address", dataType: 'string', important: true },
                    { key: "host_info", label: "Info", dataType: 'string', reFormat: (val, row) => val === '' ? '-' : val },
                    { key: "updated_at", label: "LastActivity", dataType: 'dateTime', reFormat: (val, row) => moment(val).fromNow() },

                    {
                        key: "action_button",
                        label: "Operation",
                        dataType: "actionButton",
                        buttons: {
                            list: [
                                {
                                    url: row => `/hosts/${row.host}/report`,
                                    component: (btn, row) => (<ReportButton btn={btn} forObj={row.host} />),
                                },
                                {
                                    url: row => `/hosts/${row.host}/edit`,
                                    component: (btn, row) => (<EditButton btn={btn} forObj={row.host} />),
                                },
                            ]
                        }
                    }
                ]}

                handleParentBusyState={handleParentBusyState}
                handleParentRefreshState={handleParentRefreshState}
                PerPage={15}
                Order="updated_at"
                OrderType="desc"
            />
        </div>
    )
}

export default HostsDataTableComponent;