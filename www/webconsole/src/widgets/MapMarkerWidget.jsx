import React from 'react';

import { MapContainer, TileLayer, Marker, Popup } from 'react-leaflet'



const MapMarkerWidget = (
    {
        position = [0, 0],
        zoom = 10,
        scrollWheelZoom = false,
        tileLayer = {
            attribution: '&copy; <a href="http://osm.org/copyright">OpenStreetMap</a> contributors',
            url: 'https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png',
        },
        popupTitle = '',
    }
) => {


    return (
        <MapContainer style={{ height: 'calc(100vh - 200px)' }} center={position} zoom={zoom} scrollWheelZoom={false} >
            <TileLayer
                attribution='&copy; <a href="http://osm.org/copyright">OpenStreetMap</a> contributors'
                url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
            />
            <Marker position={position}>
                <Popup>
                    {popupTitle}
                </Popup>
            </Marker>
        </MapContainer>

    );
}

export default MapMarkerWidget;