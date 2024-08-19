package main

import (
  "os"
  "log"
  mgo "gopkg.in/mgo.v2"
  . "github.com/AnastasiyaGapochkina01/go-api/model"
)

type MongoConfig struct {
	mongoHost string
	mongoPort string
	mongoDb   string
	username  string
	password  string
	msgCol    string
}

var mongoConfig = MongoConfig{
	mongoHost: getEnv("MONGO_HOST", "localhost"),
	mongoPort: getEnv("MONGO_PORT", "27017"),
	mongoDb:   getEnv("MONGO_DB", "wods_db"),
	username:  getEnv("MONGO_USER", "apiadmin"),
	password:  getEnv("MONGO_PASS", "rapiadmin"),
	msgCol:    getEnv("MONGO_MSG_COLLECTION", "wods_db"),
}


func putWod(wod *Wod) error {
	sessionCopy := Session.Copy()
	defer sessionCopy.Close()

	// put
	var coll = sessionCopy.DB(mongoConfig.mongoDb).C(mongoConfig.msgCol)
	err := coll.Insert(wod)
	if err != nil {
		log.Printf("ERROR: fail put wod, %s", err.Error())
	}

	return err
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getWod(id string) (*Wod, error) {
	sessionCopy := Session.Copy()
	defer sessionCopy.Close()

	// get
	var coll = sessionCopy.DB(mongoConfig.mongoDb).C(mongoConfig.msgColl)
	wod := &Wod{}
	err := coll.Find(bson.M{"Id": id}).One(wod)
	if err != nil {
		log.Printf("ERROR: no wod with id, %s", id)
	}

	log.Printf("INFO: found wod, %+v", wod)

	return wod, err
}

func updateWod(id string, status string) error {
	sessionCopy := Session.Copy()
	defer sessionCopy.Close()

	// update msg status
	var coll = sessionCopy.DB(mongoConfig.mongoDb).C(mongoConfig.msgColl)
	err := coll.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"status": status}})
	if err != nil {
		log.Printf("ERROR: fail update state id, %s", id)
	}

	return err
}


func getWods(receiver string, status string) []Wod {
	sessionCopy := Session.Copy()
	defer sessionCopy.Close()

	// get msgs
	var coll = sessionCopy.DB(mongoConfig.mongoDb).C(mongoConfig.msgColl)
	var wods []Wod
	err := coll.Find(bson.M{"status": status}).Select(bson.M{"id": 1, "attr": 1}).All(&wods)
	if err != nil {
		log.Printf("ERROR: fail get wods, %s", err.Error())
	}

	log.Printf("INFO: found wods, %+v", msgs)

	return msgs
}

var Session *mgo.Session

func initSession() {
	// db info
	info := &mgo.DialInfo{
		Addrs:    []string{mongoConfig.mongoHost},
		Database: mongoConfig.mongoDb,
		Username: mongoConfig.username,
		Password: mongoConfig.password,
	}
        log.Printf("INFO connecting mongo with info, %+v", info)
	// connect to mongo
	s, err := mgo.DialWithInfo(info)
	if err != nil {
		log.Printf("ERROR connecting mongo, %s ", err.Error())
		return
	}
	s.SetMode(mgo.Monotonic, true)

	Session = s
}

func clearSession() {
	Session.Close()
}


