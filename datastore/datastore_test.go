package datastore

import (
	"strings"
	"testing"

	"bitbucket.org/maksadbek/go-mon-service/conf"
)

func TestGetTrackers(t *testing.T) {
	mockConf := `[ds]
    [ds.mysql]
        dsn = "root:toor@tcp(localhost:3306)/maxtrack"
        interval = 1
    [ds.redis]
		host = ":6379"
		fprefix = "fleet"
        tprefix = "tracker"
	[srv]
		port = ":1234"
	[log]
		path = "info.log"
	`
	c := strings.NewReader(mockConf)
	app, err := conf.Read(c)
	if err != nil {
		t.Error(err)
	}

	// mysql setup
	err = Initialize(app.DS)
	if err != nil {
		t.Error(err)
	}
	_, err = GetTrackers("202")
}

func TestUsrTrackersPartialCars(t *testing.T) {
	usr, err := UsrTrackers("Kamilka")
	if err != nil {
		t.Error(err)
	}
	if usr.Trackers[0] == "0" {
		t.Errorf("want %s, got %s", "0", usr.Trackers[0])
	}
}

func TestUsrTrackersAllCars(t *testing.T) {
	usr, err := UsrTrackers("newmax")
	if err != nil {
		t.Error(err)
	}
	if usr.Trackers[0] != "0" {
		t.Errorf("want %s, got %s", "0", usr.Trackers[0])
	}
}

func TestCacheFleetTrackers(t *testing.T){
    fleetTrackers, err := CacheFleetTrackers()
    if err != nil {
            t.Error(err)
    }
    t.Logf("%+v\n", fleetTrackers)
}
