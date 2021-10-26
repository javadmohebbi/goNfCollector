import React, { useCallback, useEffect, useState } from 'react';
import { EthernetsGetSummaryByDeviceByInterval } from '../../../services/ethernet';


import Chart from "react-apexcharts";
import { Card, CardContent, Typography } from '@material-ui/core';

import humanFormat from 'human-format'



const controller = new AbortController()
const signal = controller.signal


const EthernetSummaryWidget = (
    {
        handleParentBusyState = () => { return },
        handleParentRefreshState = () => { return },
        interval = '15m',
        refresh = false,
        deviceId,
        busy = false,
    }
) => {

    const [result, setResult] = useState({})
    const [selfBusy, setSelfBusy] = useState(false)



    const [chartData, setChartData] = useState({})

    // eslint-disable-next-line
    const [fetchError, setFetchError] = useState(false)

    const getTopProtocolsCallback = useCallback(() => {

        handleParentBusyState(true)
        handleParentRefreshState(false)

        setSelfBusy(true)


        controller.abort()
        EthernetsGetSummaryByDeviceByInterval({ deviceId, interval, signal }).then(async (json) => {
            if (json.error) {
                setFetchError(true)
            } else {
                const resp = await json.response.then((result) => result);
                if (resp !== null) {
                    setResult(resp)
                } else {
                    setResult([])
                }
            }
            setSelfBusy(false)

            handleParentBusyState(false)
        })


        // eslint-disable-next-line
    }, [])


    useEffect(() => {
        if (refresh) {
            getTopProtocolsCallback()
        }
    }, [refresh, getTopProtocolsCallback])

    useEffect(() => {
        getTopProtocolsCallback()
    }, [interval, getTopProtocolsCallback])






    useEffect(() => {

        if (result.length > 0) {


            const inEth = result.filter(d => d.ingres === true).map(d => d.data)
            const outEth = result.filter(d => d.outgres === true).map(d => d.data)


            // const inData = []
            // const getData = (arr) => {
            //     for (let i = 0; i < arr.length; i++) {
            //         const d = arr[i].map(d => [Date.parse(d._time), d.total_bytes])
            //         inData.push(d)
            //     }
            // }
            // getData(inEth)

            // console.log(inEth, outEth);

            const seriesIn = inEth.filter(d => d !== null).map(d => ({
                name: d[0].eth_key + " [in]",
                data: d.map(d => [Date.parse(d._time), d.total_bytes])
            }))
            const seriesOut = outEth.filter(d => d !== null).map(d => ({
                name: d[0].eth_key + " [out]",
                data: d.map(d => [Date.parse(d._time), d.total_bytes])
            }))

            // const series = {
            //     seriesIn, seriesOut
            // }
            // console.log(seriesIn);

            const newChartData = {
                theme: {
                    mode: 'dark',
                },
                toolbar: {
                    show: false
                },
                series: [
                    ...seriesOut, ...seriesIn

                ],
                chart: {
                    type: 'area',
                    height: 500,
                },
                annotations: {

                },

                fill: {
                    gradient: {
                        enabled: true,
                        opacityFrom: 0.55,
                        opacityTo: 0
                    }
                },
                stroke: {
                    curve: 'smooth',
                },
                noData: {
                    text: 'NO DATA',
                    align: 'center',
                    verticalAlign: 'middle',
                },
                xaxis: {
                    type: 'datetime',
                    labels: {
                        datetimeUTC: false
                    },
                    // min: Date.parse(result.In[0]._time),
                    tickAmount: 10,
                },
                yaxis: {
                    labels: {
                        formatter: function (value) {
                            return humanFormat(value || 0, { unit: "B" });
                        }
                    },
                },
                tooltip: {
                    x: {
                        format: 'dd MMM yyyy H:m'
                    }
                },
            }
            setChartData(newChartData)
        }
    }, [result])





    return (
        <div>
            <Card variant="outlined">
                <CardContent>


                    <Typography>
                        Ethernets Summary
                    </Typography>

                    {
                        !selfBusy
                            ?
                            <Chart
                                options={chartData}
                                series={chartData.series || []}

                                height="500"
                            />
                            :
                            '...'
                    }
                </CardContent>
            </Card>
        </div>
    );
}

export default EthernetSummaryWidget;