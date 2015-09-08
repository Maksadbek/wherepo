package cache

import (
	"encoding/json"
	"sort"
	"strconv"
	"time"

	"github.com/Maksadbek/wherepo/conf"
	"github.com/Maksadbek/wherepo/logger"
	"github.com/garyburd/redigo/redis"
)

var (
	config conf.App // config
	pool   *redis.Pool
)

// structure for fleet
type Fleet struct {
	Id      string           `json:"id"`           // unique id of fleet
	Update  map[string][]Pos `json:"update"`       // and its tracker's info
	LastReq int64            `json:"last_request"` // current unix time
}

func (f Fleet) MarshalJSON() ([]byte, error) {
	m := struct {
		ID     string `json:"id"`
		Update []struct {
			GroupName string `json:"groupName"`
			Data      []Pos  `json:"data"`
		} `json:"update"`
	}{
		ID: f.Id,
	}
	keys := make([]string, 0)
	for groupName := range f.Update {
		keys = append(keys, groupName)
	}

	sort.Strings(keys)
	for _, groupName := range keys {
		data := struct {
			GroupName string `json:"groupName"`
			Data      []Pos  `json:"data"`
		}{
			GroupName: groupName,
		}
		for _, position := range f.Update[groupName] {
			data.Data = append(data.Data, position)
		}
		m.Update = append(m.Update, data)
	}
	return json.Marshal(m)
}

// structure for fleet
type FleetTracker struct {
	Fleet    string
	Trackers []string
}

func newPool(server string) *redis.Pool {
	return &redis.Pool{
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},
	}
}
func Initialize(c conf.App) (err error) {
	m := make(map[string]interface{})
	m["config"] = c
	logger.Log.Info("Rcache initialization")
	config = c
	// create redis pool
	pool = newPool(c.DS.Redis.Host)
	return
}

// GetTrackers can be used to get array of tracker of particular fleet
// start and stop are range values of list, default is 0,200, can be set from config
func GetTrackers(fleet string, start, stop int) (trackers []string, err error) {
	rc := pool.Get()
	defer rc.Close()
	logger.FuncLog("rcache.GetTrackers", conf.InfoListOfTrackers, nil, nil)
	// get list of trackers ids from cache
	v, err := redis.Strings(rc.Do("SMEMBERS", config.DS.Redis.FPrefix+":"+fleet)) //
	if err != nil {
		logger.FuncLog("cache.GetTrackers", conf.ErrGetListOfTrackers, nil, err)
		return
	}

	// range over list and append it to the slice
	for _, val := range v {
		trackers = append(trackers, val)
	}
	return
}

// PushRedis can be used to save fleet var into redis
func PushRedis(fleet Fleet) (err error) {
	logger.FuncLog("rcache.PushRedis", conf.InfoPushFleet, nil, nil)
	rc := pool.Get()
	defer rc.Close()
	// get list of trackers ids from cache
	// range over map of Pos and push them
	for _, x := range fleet.Update {
		for _, pos := range x {
			jpos, err := json.Marshal(pos)
			if err != nil {
				logger.FuncLog("cache.PushRedis", conf.ErrGetListOfTrackers, nil, err)
				return err
			}
			rc.Do("RPUSH", config.DS.Redis.TPrefix+":"+strconv.Itoa(pos.Id), jpos) // prefix can be set from conf
		}
	}
	return
}

// AddFleetTrackers can be used to save list of trackers to redis
func AddFleetTrackers(ftracker []FleetTracker) error {
	rc := pool.Get()
	defer rc.Close()
	for _, tracker := range ftracker {
		// range over tracker data
		rc.Do("DEL", "fleet"+":"+tracker.Fleet)
		for _, x := range tracker.Trackers {
			// add tracker to list
			rc.Do("SADD", "fleet"+":"+tracker.Fleet, x)
		}
	}
	return nil
}
