import { UtlHttpRequest } from '../utils';


export const GetTopPortByDeviceByInterval = ({
    top = 15,
    deviceId = 1,
    direction = 'src',
    interval = '15m',

    signal
}) => {
    // AIP URL
    const url = `${process.env.REACT_APP_HTTP}/port/get/top/${top}/device/${deviceId}/direction/${direction}/interval/${interval}`;

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