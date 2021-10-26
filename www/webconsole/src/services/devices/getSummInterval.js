import { UtlHttpRequest } from '../utils';


export const GetSummaryByInterval = ({
    interval = '15m',

    signal
}) => {
    // AIP URL
    const url = `${process.env.REACT_APP_HTTP}/device/get/summary/interval/${interval}`;

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