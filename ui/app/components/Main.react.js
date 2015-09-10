var React = require('react');
var StatusStore = require('../stores/StatusStore');
var CarActions = require('../actions/StatusActions');
var UserActions = require('../actions/UserActions');
var Sidebar = require('./Sidebar.react');
var UserStore = require('../stores/UserStore');
var Status = require('./CarStatus.react');
var Mui  = require('material-ui');
var GoogleMap = require('react-google-maps').GoogleMap;
var Marker = require('react-google-maps').Marker;
var ThemeManager = new Mui.Styles.ThemeManager();
mui = require('material-ui')

var AppBar = Mui.AppBar,
    MenuItem= Mui.MenuItem, 
    IconButton= Mui.IconButton, 
    List  = Mui.List,
    Paper = Mui.Paper,
    LeftNav= Mui.LeftNav;

menuItems = [
    { 
        type: MenuItem.Types.SUBHEADER, 
        text: 'Resources' 
    },
    { 
        type: MenuItem.Types.LINK, 
        payload: 'https://github.com', 
        text: 'GitHub' 
    },
    { 
        text: 'Disabled', 
        disabled: true 
    },
    { 
        type: MenuItem.Types.LINK, 
        payload: 'https://www.google.com', 
        text: 'Disabled Link', 
        disabled: true 
    },
];


function getAllStatuses(){
    return StatusStore.getAll()
}

var StatusApp = React.createClass({
    childContextTypes: {
          muiTheme: React.PropTypes.object
    },
    getChildContext: function() {
        return {
            muiTheme: ThemeManager.getCurrentTheme()
        };
    },
    getInitialState: function(){
        this._bounds = new google.maps.LatLngBounds();
        var shape = {
            coords: [1, 1, 1, 20, 18, 20, 18 , 1],
            type: 'poly'
        };
        return {
            stats: {
                id: '',
                update: [],
            },
            isChildChecked: false
        }
    },

    componentDidMount: function(){
        StatusStore.addChangeListener(this._onChange);
        UserStore.addChangeListener(this._onAuth);
        var mapOptions = { zoom: 10 };
    },

    componentWillMount: function(){
       UserActions.Auth({
           login: "taxi",
           uid: "deadbeef",
           hash: "b5ea8985533defbf1d08d5ed2ac8fe9b",
           fleet: "436"
       });
    },
    componentWillUnmount: function(){
        StatusStore.removeChangeListener(this._onChange);
        UserStore.removeChangeListener(this._onAuth);
    },
    toggleLeftNav: function(){
        console.log(this.refs.leftNav);
        React.findDOMNode(this.refs.leftNav).toggle()
    },
    render: function(){
        var content = [];
        var markers = [];
        var update = this.state.stats.update;
        var checked = this.state.isChildChecked;
        update.forEach(function(group){
            content.push(<Sidebar key={group.groupName} stats={group}/>)
            group.data.forEach(function(vehicle){
                markers.push({
                    position:{
                        lat: vehicle.latitude,
                        lng: vehicle.longitude
                    },
                    key: vehicle.id
                });
            });
        });
        return (   
            <div>
                <LeftNav ref="leftNav" docked={false} menuItems={menuItems} />
                <AppBar
                    onLeftIconButtonTouchTap={this.toggleLeftNav}
                    title="Wherepo"
                    iconElementLeft={<IconButton></IconButton>}
                />
                <div style={{border: "solid 1px #d9d9d9", height: "100vh", float: "left", width:"69%"}} id={"map-canvas"}>
                    <GoogleMap containerProps={{style:{height:"100%"}}} ref="map" defaultZoom={3} 
                            defaultCenter={{lat: -25.363882, lng: 131.044922}}>
                            {markers.map(function(marker, index){
                                    return(<Marker {...marker} />);
                                })
                            }
                    </GoogleMap>
                </div>
                <div style={{width: "30%", border: "solid 1px #d9d9d9", height: "100vh", float: "left", overflow:"scroll"}}>
                    <List>
                        {content}
                    </List>
                </div>

            </div>
            )
    },
    _onChange: function(){
        this.setState({stats: getAllStatuses()});
    },
    _onAuth: function(){
        console.log("fuck");
        StatusStore.sendAjax();
        setInterval(function(){
            StatusStore.sendAjax();
        }, 5000);
    }
});

module.exports = StatusApp;
