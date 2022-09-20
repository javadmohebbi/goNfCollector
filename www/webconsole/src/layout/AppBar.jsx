import React, { useState } from 'react'
import { makeStyles } from '@material-ui/core/styles';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import IconButton from '@material-ui/core/IconButton';
import MenuIcon from '@material-ui/icons/Menu';
import Button from '@material-ui/core/Button';


import DrawerLayout from './Drawer'
import { Link } from 'react-router-dom';

const useStyles = makeStyles((theme) => ({
    root: {
        flexGrow: 1,
    },
    menuButton: {
        marginRight: theme.spacing(2),
    },
    title: {
        flexGrow: 1,
    },
}));

const AppBarLayout = () => {
    const classes = useStyles();


    const [isDrawerOpen, setIsDrawerOpen] = useState(false)

    const handleDrawerOpenToggle = () => {
        setIsDrawerOpen(!isDrawerOpen)
    }


    return (
        <div className={classes.root}>
            <AppBar position="static" color="primary">
                <Toolbar>
                    <IconButton edge="start" className={classes.menuButton} color="inherit" aria-label="menu"
                        onClick={(e) => { e.preventDefault(); setIsDrawerOpen(!isDrawerOpen) }}
                    >
                        <MenuIcon />
                    </IconButton>
                    <Typography variant="h6" className={classes.title}>
                        Netflow Analyzer
                    </Typography>

                    {/* <Button color="inherit" component={Link} to="/settings">Settings</Button> */}

                </Toolbar>

                <DrawerLayout drawerOpen={isDrawerOpen} onCloseDrawer={handleDrawerOpenToggle} />

            </AppBar>
        </div>
    );
}

export default AppBarLayout