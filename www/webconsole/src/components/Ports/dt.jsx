import React from 'react';
import { makeStyles } from '@material-ui/core/styles';

import WidgetDataTableComponent from '../../widgets/DataTable';

import moment from "moment";

// eslint-disable-next-line
import ReportButton from '../../widgets/ReportButton';

import EditButton from '../../widgets/EditButton';




const useStyles = makeStyles((theme) => ({
    root: {

    },
}))





const PortsDataTableComponent = ({
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

                URL="port/get/all"
                Columns={[
                    { key: "port", label: "Port", dataType: 'string', important: true },
                    { key: "port_proto", label: "Proto/Port", dataType: 'string', important: true },
                    { key: "port_info", label: "Info", dataType: 'string', reFormat: (val, row) => val === '' ? '-' : val },
                    { key: "updated_at", label: "LastActivity", dataType: 'dateTime', reFormat: (val, row) => moment(val).fromNow() },

                    {
                        key: "action_button",
                        label: "Operation",
                        dataType: "actionButton",
                        buttons: {
                            list: [
                                // {
                                //     url: row => `/ports/${row.port}/report`,
                                //     component: (btn, row) => (<ReportButton btn={btn} forObj={row.port} />),
                                // },
                                {
                                    url: row => () => {
                                        const prt = row.port.replace("/", " ");
                                        return `/ports/${row.port_id}/${prt}/edit`
                                    },
                                    component: (btn, row) => (<EditButton btn={btn} forObj={row.port} />),
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

export default PortsDataTableComponent;