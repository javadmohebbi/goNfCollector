import { UtlHttpRequest } from '../utils';


export const GetSummaryByIntervalByDevice = ({
    interval = '15m',
    device = 'uknown',

    signal
}) => {
    // AIP URL
    const url = `${process.env.REACT_APP_HTTP}/device/get/summary/interval/${interval}/by/${device}`;

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