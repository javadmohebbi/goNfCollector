import React from 'react';
import { makeStyles } from '@material-ui/core/styles';

import WidgetDataTableComponent from '../../widgets/DataTable';

import moment from "moment";

import ReportButton from '../../widgets/ReportButton';
import ThreatLevelWidget from '../../widgets/ThreatLevel';
import ThreatUserAckWidget from '../../widgets/ThreatUserAck';




const useStyles = makeStyles((theme) => ({
    root: {

    },
}))





const ThreatsDataTableComponent = ({
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

                URL="threat/get/all"
                Columns={[
                    { key: "threat_source", label: "Source", dataType: 'string', important: true },
                    { key: "threat_kind", label: "Kind", dataType: 'string', important: true },
                    // { key: "threat_counter", label: "Count", dataType: 'string', reFormat: (val, row) => <ThreatLevelWidget level={row.threat_counter} justCounter /> },
                    { key: "threat_reputation", label: "Level", dataType: 'string', reFormat: (val, row) => <ThreatLevelWidget level={val} /> },
                    { key: "threat_host", label: "RelatedHost", dataType: 'string', reFormat: (val, row) => row.threat_host_info !== '' ? val + ' (' + row.threat_host_info + ')' : val },
                    // { key: "threat_acked", label: "Acked|Closed|False+", dataType: 'string', reFormat: (val, row) => <ThreatUserAckWidget acked={row.threat_acked} closed={row.threat_closed} falsePositive={row.threat_false_positive} /> },
                    { key: "created_at", label: "FirstSeen", dataType: 'dateTime', reFormat: (val, row) => moment(val).fromNow() },
                    { key: "updated_at", label: "LastActivity", dataType: 'dateTime', reFormat: (val, row) => moment(val).fromNow() },

                    {
                        key: "action_button",
                        label: "Operation",
                        dataType: "actionButton",
                        buttons: {
                            list: [
                                {
                                    url: row => `/threats/${row.threat_host}/report`,
                                    component: (btn, row) => (<ReportButton btn={btn} forObj={row.host} />),
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

export default ThreatsDataTableComponent;