import { UtlHttpRequest } from '../utils';


export const GetById = ({
    id = 'uknown',

    signal
}) => {
    // AIP URL
    const url = `${process.env.REACT_APP_HTTP}/protocol/get/by/id/${id}`;

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