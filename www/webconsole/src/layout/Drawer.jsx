import React from 'react'
import Drawer from '@material-ui/core/Drawer';
import List from '@material-ui/core/List';
import Divider from '@material-ui/core/Divider';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import { makeStyles } from '@material-ui/core/styles';
import { Link } from 'react-router-dom'
import { useLocation } from 'react-router-dom'

import DashboardIcon from '@material-ui/icons/Dashboard';
import RouterIcon from '@material-ui/icons/Router';
// eslint-disable-next-line
import ViewListIcon from '@material-ui/icons/ViewList';
import ComputerIcon from '@material-ui/icons/Computer';
import WebIcon from '@material-ui/icons/Web';
import MapIcon from '@material-ui/icons/Map';
import WarningIcon from '@material-ui/icons/Warning';

const useStyles = makeStyles({
    list: {
        width: 250,
    },

});


const DrawerLayout = (props) => {
    const location = useLocation();


    const classes = useStyles();

    const { drawerOpen } = props || false
    const { onCloseDrawer } = props

    const toggleDrawer = (open) => (event) => {
        if (event.type === 'keydown' && (event.key === 'Tab' || event.key === 'Shift')) {
            return;
        }

        onCloseDrawer()
    };

    return (
        <React.Fragment>
            <Drawer open={drawerOpen} anchor={'left'} onClose={toggleDrawer(false)}>
                <div
                    className={classes.list}
                    role="presentation"
                    onClick={toggleDrawer(false)}
                    onKeyDown={toggleDrawer(false)}
                >
                    <List>
                        {[
                            { name: 'Dashboard', Icon: () => (<DashboardIcon />), to: '/dashboard' },
                            { name: 'Managed Devices', Icon: () => (<RouterIcon />), to: '/devices' },
                            // { name: 'Raw Filter', Icon: () => (<ViewListIcon />), to: '/raw/filter' },
                        ].map((item, index) => (
                            <ListItem selected={location.pathname === item.to} component={Link} to={item.to} button key={item.name}>
                                <ListItemIcon>{<item.Icon />}</ListItemIcon>
                                <ListItemText primary={item.name} />
                            </ListItem>
                        ))}
                    </List>

                    <Divider />



                    <Divider />


                    {/* HOSTS */}
                    <List>
                        {[
                            { name: 'Hosts', Icon: () => (<ComputerIcon />), to: '/hosts' },
                            { name: 'Ports', Icon: () => (<WebIcon />), to: '/ports' },
                            { name: 'Protocols', Icon: () => (<WebIcon />), to: '/protocols' },
                            // { name: 'Geo Locations', Icon: () => (<MapIcon />), to: '/geos' },
                            { name: 'Threats', Icon: () => (<WarningIcon />), to: '/threats' },
                        ].map((item, index) => (
                            <ListItem selected={location.pathname === item.to} component={Link} to={item.to} button key={item.name}>
                                <ListItemIcon>
                                    {
                                        typeof item.Icon !== 'undefined' ?
                                            <item.Icon />
                                            : ''
                                    }
                                </ListItemIcon>
                                <ListItemText primary={item.name} />
                            </ListItem>
                        ))}
                    </List>



                </div>
            </Drawer>
        </React.Fragment >
    )
}

export default DrawerLayout