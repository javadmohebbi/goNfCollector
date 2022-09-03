import { UtlHttpRequest } from '../utils';


export const SetById = ({
    id = 0,
    info = '',

    signal
}) => {
    // AIP URL
    const url = `${process.env.REACT_APP_HTTP}/port/set/by/id/${id}`;

    // send request to server
    return UtlHttpRequest(
        url,
        'POST',
        false,
        {
            id: parseInt(id),
            info: info,
        },
        false,
        signal
    );
};