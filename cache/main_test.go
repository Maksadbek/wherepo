package cache

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/Maksadbek/wherepo/conf"
)

func TestMain(m *testing.M) {
	r := strings.NewReader(mockConf) // читает мок-данные из testdata.go
	app, err := conf.Read(r)         //
	if err != nil {
		panic(err)
	}
	err = Initialize(app)
	if err != nil {
		panic(err)
	}
	rc := pool.Get()
	defer rc.Close()

	for _, x := range FleetTest.Trackers {
		rc.Do("SADD", "fleet"+":"+FleetTest.FleetName, x)
	}

	// add mock user
	jusr, err := json.Marshal(testUsr[0])
	if err != nil {
		panic(err)
	}
	rc.Do(
		"SET",
		app.DS.Redis.UPrefix+":"+testUsr[0].Login,
		string(jusr),
	)

	retCode := m.Run()

	// clean up messed redis test zone
	/*
		for _, x := range FleetTest.Trackers {
			rc.Do("LPOP", FleetTest.FleetName)
			rc.Do("LPOP", x)
		}
	*/

	os.Exit(retCode)
}
