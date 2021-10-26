import { UtlHttpRequest } from '../utils';


export const GetByDevice = ({
    device = 'uknown',

    signal
}) => {
    // AIP URL
    const url = `${process.env.REACT_APP_HTTP}/device/get/by/${device}`;

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