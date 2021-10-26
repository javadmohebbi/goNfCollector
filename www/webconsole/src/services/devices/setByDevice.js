import { UtlHttpRequest } from '../utils';


export const SetByDevice = ({
    device = 'uknown',
    name = '',
    info = '',

    signal
}) => {
    // AIP URL
    const url = `${process.env.REACT_APP_HTTP}/device/set/by/${device}`;

    // send request to server
    return UtlHttpRequest(
        url,
        'POST',
        false,
        {
            device: device,
            name: name,
            info: info,
        },
        false,
        signal
    );
};