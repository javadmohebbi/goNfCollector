import { UtlHttpRequest } from '../utils';


export const GetEthernetSummaryByDeviceByInterval = ({
    deviceId = 1,
    interval = '15m',

    signal
}) => {
    // AIP URL
    const url = `${process.env.REACT_APP_HTTP}/eth/get/device/${deviceId}/interval/${interval}`;

    // send request to server
    return UtlHttpRequest(
        url,
        'GET',
        false,
        false,
        false,
        signal
    );
};