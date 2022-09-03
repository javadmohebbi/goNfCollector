import React from 'react'

import { HashRouter as Router, Route, Switch } from 'react-router-dom'
import DashboardComponent from './components/Dashbaord/Dashboard'
import GeoLocationsComponent from './components/GeoLocations/GeoLocations'
import HostsComponent from './components/Hosts/Hosts'
import MDEditComponent from './components/managedDevices/MDEdit'
import ManagedDeviceComponent from './components/managedDevices/ManagedDevice'
import MDReport from './components/managedDevices/MDReport'
import PortsComponent from './components/Ports/Ports'
import ProtocolsComponent from './components/Protocols/Protocols'
import ThreatsComponent from './components/Threats/Threats'

import MainAdminLayout from './layout/MainAdmin'
import HostReport from './components/Hosts/HostReport'
import HostEditComponent from './components/Hosts/HostEdit'
import PortEditComponent from './components/Ports/PortEdit'
import ProtocolEditComponent from './components/Protocols/ProtocolEdit'
import SettingsComponent from './components/Setting/Settings'

const MainRouter = () => {
    return (
        <Router>
            <MainAdminLayout>
                <Switch>
                    <Route exact path="/" component={DashboardComponent} name="dashboard" />
                    <Route exact path="/dashboard" component={DashboardComponent} name="dashboard" />

                    <Route exact path="/devices" component={ManagedDeviceComponent} name="devices" />
                    <Route exact path="/devices/:device/report" component={MDReport} name="devices_report" />
                    <Route exact path="/devices/:device/edit" component={MDEditComponent} name="devices_edit" />

                    <Route exact path="/hosts" component={HostsComponent} name="hosts" />
                    <Route exact path="/hosts/:host/report" component={HostReport} name="host_report" />
                    <Route exact path="/hosts/:hostId/:host/edit" component={HostEditComponent} name="host_edit" />

                    <Route exact path="/ports" component={PortsComponent} name="ports" />
                    <Route exact path="/ports/:portId/:port/edit" component={PortEditComponent} name="port_edit" />

                    <Route exact path="/protocols" component={ProtocolsComponent} name="protocols" />
                    <Route exact path="/protocols/:protoId/:protocol/edit" component={ProtocolEditComponent} name="proto_edit" />

                    <Route exact path="/geos" component={GeoLocationsComponent} name="geos" />

                    <Route exact path="/threats" component={ThreatsComponent} name="threats" />
                    <Route exact path="/threats/:host/report" component={HostReport} name="threat_host_report" />

                    {/* Settings */}
                    <Route exact path="/settings" component={SettingsComponent} name="settings" />
                </Switch>
            </MainAdminLayout>
        </Router>
    )
}


export default MainRouter