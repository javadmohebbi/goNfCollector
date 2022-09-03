import React, { useState } from 'react';
import { makeStyles } from '@material-ui/core/styles';

import WidgetDataTableComponent from '../../widgets/DataTable';

import moment from "moment";
import ReactCountryFlag from "react-country-flag"

// eslint-disable-next-line
import ReportButton from '../../widgets/ReportButton';

import { Grid } from '@material-ui/core';
import MapButton from '../../widgets/MapButton';
import Modal from 'react-modal';
import CloseButton from '../../widgets/CloseButton';
import MapMarkerWidget from '../../widgets/MapMarkerWidget';


Modal.setAppElement('#root');


const useStyles = makeStyles((theme) => ({
    root: {

    },
}))











const GeoLocationsDataTableComponent = ({
    busy = false,
    refresh = false,
    handleParentBusyState = () => { return },
    handleParentRefreshState = () => { return },

}) => {

    const classes = useStyles();


    const [isMapModelOpen, setIsMapModelOpen] = useState(false)
    const [modalMapInfo, setModalMapInfo] = useState({})


    const handleModalOpen = (
        position = [0, 0],
        popupTitle = '',
    ) => {
        setModalMapInfo({ position: position, popupTitle: popupTitle })
        setIsMapModelOpen(true)
    }


    return (
        <div className={classes.root}>
            <WidgetDataTableComponent
                refresh={refresh}

                URL="geo/get/all"
                Columns={[
                    {
                        key: "country_short", label: "Country Short", dataType: 'string', important: true,
                        reFormat: (val, row) => <> {val} {
                            row.country_long === 'NA'
                                ? ''
                                : <ReactCountryFlag
                                    svg
                                    countryCode={val}
                                    style={{
                                        width: '1.5em',
                                        height: '1.5em',
                                    }}
                                    title={row.country_long}
                                />
                        }</>
                    },
                    { key: "country_long", label: "Country Long", dataType: 'string', important: true },
                    { key: "region", label: "Region", dataType: 'string' },
                    { key: "city", label: "City", dataType: 'string' },
                    {
                        key: "latitude", label: "Map", dataType: 'string',
                        reFormat: (val, row) => (
                            <>
                                {
                                    row.latitude === 0 && row.longitude === 0
                                        ? '-'
                                        :
                                        <>
                                            <MapButton forObj={row.country_short + ' > ' + row.city}
                                                onClick={(e) => {
                                                    handleModalOpen([row.latitude, row.longitude], [row.city + ' (' + row.country_long + ')'])
                                                }}
                                            />
                                        </>
                                }
                            </>
                        )
                    },
                    { key: "updated_at", label: "LastActivity", dataType: 'dateTime', reFormat: (val, row) => moment(val).fromNow() },

                    // {
                    //     key: "action_button",
                    //     label: "Operation",
                    //     dataType: "actionButton",
                    //     buttons: {
                    //         list: [
                    //             {
                    //                 url: row => `/geos/${row.geo_id}/report`,
                    //                 component: (btn, row) => (<ReportButton btn={btn} forObj={row.country_long + ' > ' + row.city} />),
                    //             },
                    //         ]
                    //     }
                    // }
                ]}

                handleParentBusyState={handleParentBusyState}
                handleParentRefreshState={handleParentRefreshState}
                PerPage={15}
                Order="updated_at"
                OrderType="desc"
            />

            <Modal
                isOpen={isMapModelOpen}
            >

                <Grid container spacing={2}>
                    <Grid item xs={12} style={{ alignItems: 'flex-end', color: 'black', fontWeight: 'bold' }} >
                        <CloseButton onClick={(e) => setIsMapModelOpen(false)} />
                        <span style={{ margin: 'auto 20px' }}> {modalMapInfo.popupTitle} </span>
                    </Grid>
                    <Grid item xs={12}>
                        <MapMarkerWidget
                            popupTitle={modalMapInfo.popupTitle}
                            position={modalMapInfo.position}
                            zoom={7}
                        />
                    </Grid>
                </Grid>



            </Modal>
        </div >
    )
}

export default GeoLocationsDataTableComponent;