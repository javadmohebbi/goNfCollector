import React, { useEffect, useState } from 'react';
import io from 'socket.io-client';

const socket = io(
    process.env.REACT_APP_WS,
    {
        transports: ["websocket"],
        autoConnect: false
    }
);


const IP2LocationSettings = () => {

    const [isConnected, setIsConnected] = useState(socket.connected);
    const [pong, setPong] = useState(0);

    useEffect(() => {

        socket.connect()

        socket.on('connect', () => {
            setIsConnected(true);
            socket.emit("join", "ip2l");
        });

        socket.on('disconnect', () => {
            socket.connect()
            setIsConnected(false);
        });

        socket.on('pong', (data) => {
            setPong(data)
        });

        return () => {
            socket.off('connect');
            socket.off('disconnect');
            socket.off('pong');
            socket.disconnect();
        };
    }, [])


    return (
        <div>
            IP2location Settings

            <p>Connected: {'' + isConnected}</p>
            <p>Pong: {'' + pong}</p>
        </div>
    )

};

export default IP2LocationSettings;