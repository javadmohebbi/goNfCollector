import { UtlHttpRequest } from '../utils';


export const GetTopProtoByDeviceByInterval = ({
    top = 15,
    deviceId = 1,
    interval = '15m',

    signal
}) => {
    // AIP URL
    const url = `${process.env.REACT_APP_HTTP}/protocol/get/top/${top}/device/${deviceId}/interval/${interval}`;

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