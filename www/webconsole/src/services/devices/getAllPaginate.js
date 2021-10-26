import { UtlHttpRequest } from '../utils';


export const GetAllPaginate = ({
    page = 1,
    perPage = process.env.REACT_NORMAL_DT_PER_PAGE,
    filter = '',
    order = 'last_activity',
    orderType = 'desc',

    signal = null,
}) => {
    // AIP URL
    const url = `${process.env.REACT_APP_HTTP}/device/get/all`;

    // send request to server
    return UtlHttpRequest(
        url,
        'GET',
        false,
        false,
        {
            page: page,
            perPage: perPage,
            filter: filter,
            order: order,
            orderType: orderType,
        },
        false,
        signal
    );
};