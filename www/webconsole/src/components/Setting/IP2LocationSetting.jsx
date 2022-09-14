import React, { useEffect, useState } from 'react';

// prepare socket io client
// import io from 'socket.io-client';

// const socket = io("http://127.0.0.1:9999/", { transports: ['websocket'] });
import useSocket from 'use-socket.io-client';

var timerId = -1

const IP2LocationSettings = () => {

    const [socket] = useSocket("ws://127.0.0.1:9999", {
        autoConnect: false,
        transports: ['websocket'],
    })

    const [data, setData] = useState("NODATA")



    useEffect(() => {

        socket.connect()

        return () => {
            console.log("disconnected!");
            socket.disconnect()
            socket.close()
        }
        // eslint-disable-next-line
    }, [])


    // useEffect(() => {
    socket.on('connect', () => {

        clearInterval(timerId)

        console.log("on connect");
        timerId = setInterval(() => {
            socket.emit("getIP2L", "update");
        }, 1000);


    })


    socket.on('pong', function (data) {
        console.log("PONG::", data);
        setData(data)
    })

    socket.on('ping', function (data) {
        console.log("PING::", data);
    })



    return (
        <div>
            {data}
        </div>
    )

};

export default IP2LocationSettings;