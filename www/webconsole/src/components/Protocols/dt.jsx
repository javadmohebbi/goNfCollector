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





const ProtocolsDataTableComponent = ({
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

                URL="protocol/get/all"
                Columns={[
                    { key: "protocol", label: "Protocol", dataType: 'string', important: true },
                    { key: "protocol_info", label: "Info", dataType: 'string', reFormat: (val, row) => val === '' ? '-' : val },
                    { key: "updated_at", label: "LastActivity", dataType: 'dateTime', reFormat: (val, row) => moment(val).fromNow() },

                    {
                        key: "action_button",
                        label: "Operation",
                        dataType: "actionButton",
                        buttons: {
                            list: [
                                // {
                                //     url: row => `/protocols/${row.protocol}/report`,
                                //     component: (btn, row) => (<ReportButton btn={btn} forObj={row.protocol} />),
                                // },
                                {
                                    url: row => `/protocols/${row.protocol_id}/${row.protocol}/edit`,
                                    component: (btn, row) => (<EditButton btn={btn} forObj={row.protocol} />),
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

export default ProtocolsDataTableComponent;