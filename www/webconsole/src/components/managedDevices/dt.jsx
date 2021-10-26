import React from 'react';
import { makeStyles } from '@material-ui/core/styles';

import WidgetDataTableComponent from '../../widgets/DataTable';

import moment from "moment";
import NetflowVersionWidget from '../../widgets/NetflowVersion';
import ReportButton from '../../widgets/ReportButton';
import EditButton from '../../widgets/EditButton';




const useStyles = makeStyles((theme) => ({
    root: {

    },
}))





const ManagedDeviceDataTableComponent = ({
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

                URL="device/get/all"
                Columns={[
                    { key: "device", label: "Device IP", dataType: 'string', important: true },
                    { key: "device_name", label: "Name", dataType: 'string', reFormat: (val, row) => val === '' ? '-' : val },
                    { key: "device_info", label: "Info", dataType: 'string', reFormat: (val, row) => val === '' ? '-' : val },
                    { key: "last_activity", label: "LastActivity", dataType: 'dateTime', reFormat: (val, row) => moment(val).fromNow() },
                    { key: "flow_version", label: "Netflow Version", dataType: 'string', reFormat: (val, row) => <NetflowVersionWidget version={val} /> },
                    {
                        key: "action_button",
                        label: "Operation",
                        dataType: "actionButton",
                        buttons: {
                            list: [
                                {
                                    url: row => `/devices/${row.device}/report`,
                                    component: (btn, row) => (<ReportButton btn={btn} forObj={row.device} />),
                                },
                                {
                                    url: row => `/devices/${row.device}/edit`,
                                    component: (btn, row) => (<EditButton btn={btn} forObj={row.device} />),
                                },
                                // {
                                //     url: row => `/devices/${row.device}/delete`,
                                //     component: (btn, row) => (<DeleteButton btn={btn} forObj={row.device} />),
                                // },
                            ]
                        }
                    }
                ]}
                Order="last_activity"
                handleParentBusyState={handleParentBusyState}
                handleParentRefreshState={handleParentRefreshState}
                PerPage={15}
            />
        </div>
    )
}

export default ManagedDeviceDataTableComponent;