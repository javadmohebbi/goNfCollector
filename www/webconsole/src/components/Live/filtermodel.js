
export const fltModel = {
    isFilterEnable: false,
    device: '',
    ip: '',
    port: '',
    srcOrDst: 'both', // src, dst, both
    proto: '',
    country: '', // long form
    region: '',
    city: '',
    flags: {
        filtered: false, // to enable or disable this filter
        list: {
            fin: false,
            syn: false,
            rst: false,
            psh: false,
            ack: false,
            urg: false,
            ece: false,
            cwr: false,
        }
    },
    threat: {
        filtered: false, // to enable or disable this filter
        isThreat: false
    },
    flowVersion: 0, // version number; 0 = disable this filter
}