package mongo_db

import (
	// "fmt"
	"github.com/GoranMandic91/euroleague_web_server/model"
	"github.com/GoranMandic91/euroleague_web_server/scraper"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
	"strconv"
	"github.com/spf13/viper"

	"fmt"
)

var (
	session *mgo.Session
	err     error
	host string = viper.GetString("development.database.host")
)

const (
	db_name           = "test"
	player_collection = "players"
	team_collection   = "teams"
	time_collection   = "time"
)

func InitializeDatabase() *mgo.Session {
	fmt.Println("MongoDB")

	session, _ = mgo.Dial("localhost")
	check(err)

	// viper.SetConfigName("config")
	// viper.AddConfigPath("/Users/goranmandic/Projects/work/src/github.com/GoranMandic91/euroleague_web_server/")
	// viper.ReadInConfig()
	// host := viper.GetString("development.database.host")
	// fmt.Println(host)

	// defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	fmt.Println("MongoDB session created")

	//get time of last db creation
	last := session.DB(db_name).C(time_collection)
	timeOfCreation := model.TimeOfCreation{}
	last.Find(bson.M{"id": "created_at"}).One(&timeOfCreation)

	db_valid := timeOfCreation.Date.Add(time.Hour * 24).Before(time.Now())

	if db_valid {

		//drop database because data are too old
		fmt.Println("Droping old database...")
		err = session.DB(db_name).DropDatabase()
		check(err)

		//set db time of creation
		current := session.DB(db_name).C(time_collection)
		currentTime := model.TimeOfCreation{"created_at", time.Now()}
		current.Insert(currentTime)

		//populate players collection
		fmt.Println("Creating players collection...")
		p := session.DB(db_name).C(player_collection)

		players := PopulateDbWithPlayers()

		bulkOfPlayers := p.Bulk()
		bulkOfPlayers.Unordered()
		bulkOfPlayers.Insert(players...)
		bulkOfPlayers.Run()
		fmt.Println("Player collection created successfully!")

		//populate teams collection
		fmt.Println("Creating teams collection...")
		t := session.DB(db_name).C(team_collection)

		teams := PopulateDbWithTeams()
		bulkOfTeams := t.Bulk()
		bulkOfTeams.Unordered()
		bulkOfTeams.Insert(teams...)
		bulkOfTeams.Run()
		fmt.Println("Team collection created successfully!")

		fmt.Println("MongoDB created successfully!")

	}
	return session
}

func PopulateDbWithPlayers() []interface{} {
	return scraper.GetAllPlayers()
}

func PopulateDbWithTeams() []interface{} {
	return scraper.GetAllTeams()
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func GetTeams(urlParams model.UrlParams) (model.ResponseData, error) {
	var teams []model.Team
	var filter = bson.M{}
	var query *mgo.Query

	s := session.Copy()
	defer s.Close()

	c := s.DB(db_name).C(team_collection)

	//filtering
	if urlParams.FilterBy != "" && urlParams.FilterValue != "" && urlParams.FilterOperator != "" {
		if urlParams.FilterBy == "id" {
			fv, _ := strconv.Atoi(urlParams.FilterValue)
			filter = bson.M{urlParams.FilterBy: bson.M{"$" + urlParams.FilterOperator: fv}}
		} else {
			filter = bson.M{urlParams.FilterBy: bson.M{"$" + urlParams.FilterOperator: urlParams.FilterValue}}
		}
	}

	//sorting
	if urlParams.OrderBy != "" {
		query = c.Find(filter).Sort(urlParams.OrderDirection + urlParams.OrderBy)
	} else {
		query = c.Find(filter)
	}
	num, _ := query.Count()

	//pagination
	if urlParams.Page != 0 {
		query = query.Skip((urlParams.Page - 1) * urlParams.PageSize).Limit(urlParams.PageSize)
	}

	err := query.All(&teams)
	direction := DirectionOfSort(urlParams.OrderBy, urlParams.OrderDirection)

	res := model.ResponseData{urlParams.Page, urlParams.PageSize, num, urlParams.OrderBy, direction, urlParams.FilterBy, urlParams.FilterValue, urlParams.FilterOperator, teams}
	return res, err
}

func GetTeamById(id int) (model.TeamInfo, error) {
	s := session.Copy()
	defer s.Close()

	c := s.DB(db_name).C(team_collection)

	var team model.TeamInfo
	err := c.Find(bson.M{"id": id}).One(&team)

	return team, err
}

func GetPlayersOfTeam(id int) ([]model.Player, error, error) {
	s := session.Copy()
	defer s.Close()

	t := s.DB(db_name).C(team_collection)
	c := s.DB(db_name).C(player_collection)

	var players []model.Player
	var team model.Team
	err1 := t.Find(bson.M{"id": id}).One(&team)
	err2 := c.Find(bson.M{"team": team.Name}).All(&players)

	return players, err1, err2
}

func GetPlayers(urlParams model.UrlParams) (model.ResponseData, error) {
	var players []model.Player
	var filter = bson.M{}
	var query *mgo.Query

	s := session.Copy()
	defer s.Close()

	c := s.DB(db_name).C(player_collection)

	//filtering
	if urlParams.FilterBy != "" && urlParams.FilterValue != "" && urlParams.FilterOperator != "" {
		if urlParams.FilterBy == "id" {
			fv, _ := strconv.Atoi(urlParams.FilterValue)
			filter = bson.M{urlParams.FilterBy: bson.M{"$" + urlParams.FilterOperator: fv}}
		} else {
			filter = bson.M{urlParams.FilterBy: bson.M{"$" + urlParams.FilterOperator: urlParams.FilterValue}}
		}
	}

	//sorting
	if urlParams.OrderBy != "" {
		query = c.Find(filter).Sort(urlParams.OrderDirection + urlParams.OrderBy)
	} else {
		query = c.Find(filter)
	}
	num, _ := query.Count()

	//pagination
	if urlParams.Page != 0 {
		query = query.Skip((urlParams.Page - 1) * urlParams.PageSize).Limit(urlParams.PageSize)
	}

	err := query.All(&players)
	direction := DirectionOfSort(urlParams.OrderBy, urlParams.OrderDirection)

	res := model.ResponseData{urlParams.Page, urlParams.PageSize, num, urlParams.OrderBy, direction, urlParams.FilterBy, urlParams.FilterValue, urlParams.FilterOperator, players}
	return res, err
}

func GetPlayerById(id int) (model.PlayerInfo, error) {
	s := session.Copy()
	defer s.Close()

	c := s.DB(db_name).C(player_collection)

	var player model.PlayerInfo
	err := c.Find(bson.M{"id": id}).One(&player)

	return player, err
}

func GetPlayersAndNations() ([]model.PlayerNation, error) {
	s := session.Copy()
	defer s.Close()

	c := s.DB(db_name).C(player_collection)

	var players []model.PlayerNation
	err := c.Find(bson.M{}).All(&players)

	return players, err
}

func DirectionOfSort(orderBy string, orderDirection string) string {
	if orderBy != "" {
		if orderDirection == "-" {
			return "desc"
		} else {
			return "asc"
		}
	}
	return ""
}
