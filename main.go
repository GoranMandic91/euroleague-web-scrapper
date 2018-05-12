package main

import (
	"encoding/json"
	"fmt"
	"github.com/GoranMandic91/euroleague_web_server/map-sort"
	"github.com/GoranMandic91/euroleague_web_server/model"
	"github.com/GoranMandic91/euroleague_web_server/mongo_db"
	"github.com/julienschmidt/httprouter"
	//"github.com/gorilla/handlers"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {

	router := httprouter.New()

	start := time.Now()

	initDB()

	fmt.Printf("MongoDB initialized for %s\n", time.Since(start))

	router.GET("/teams", getTeams)

	router.GET("/teams/:id", getTeamById)

	router.GET("/teams/:id/players", getPlayersOfTeam)

	router.GET("/players", getPlayers)

	router.GET("/players/:id", getPlayerById)

	router.GET("/stats/country-players", getStats)

	log.Fatal(http.ListenAndServe(":8081", router))


}

func initDB() {
	mongo_db.InitializeDatabase()
}

func getTeams(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	urlParams := SetUrlParams(r)
	teams, err := mongo_db.GetTeams(urlParams)

	if err != nil {
		ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
		log.Println("Failed get all books: ", err)
		return
	}

	respBody, err := json.MarshalIndent(teams, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	ResponseWithJSON(w, respBody, http.StatusOK)

}
func getTeamById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	ids := ps.ByName("id")
	id, _ := strconv.Atoi(ids)

	team, err := mongo_db.GetTeamById(id)

	if err != nil {
		ErrorWithJSON(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		log.Println("Failed get team with given id: ", err)
		return
	}

	respBody, err := json.MarshalIndent(team, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	ResponseWithJSON(w, respBody, http.StatusOK)
}

func getPlayersOfTeam(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ids := ps.ByName("id")
	id, _ := strconv.Atoi(ids)
	players, err, _ := mongo_db.GetPlayersOfTeam(id)

	if err != nil {
		ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
		log.Println("Failed get player of team with given id: ", err)
		return
	}

	respBody, err := json.MarshalIndent(players, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	ResponseWithJSON(w, respBody, http.StatusOK)
}

func getPlayers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	urlParams := SetUrlParams(r)
	players, err := mongo_db.GetPlayers(urlParams)

	if err != nil {
		ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
		log.Println("Failed get all players: ", err)
		return
	}

	respBody, err := json.MarshalIndent(players, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	ResponseWithJSON(w, respBody, http.StatusOK)

}

func getPlayerById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ids := ps.ByName("id")
	id, _ := strconv.Atoi(ids)
	player, err := mongo_db.GetPlayerById(id)

	if err != nil {
		ErrorWithJSON(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		log.Println("Failed get player with given id: ", err)
		return
	}

	respBody, err := json.MarshalIndent(player, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	ResponseWithJSON(w, respBody, http.StatusOK)
}

func getStats(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	mapOfNations := make(map[string]int)
	stats := []model.Statistic{}
	players, err := mongo_db.GetPlayersAndNations()

	if err != nil {
		ErrorWithJSON(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		log.Println("Failed get all players: ", err)
		return
	}

	for _, s := range players {
		if val, ok := mapOfNations[s.Nationality]; ok {
			mapOfNations[s.Nationality] = val + 1
		} else {
			mapOfNations[s.Nationality] = 1
		}
	}
	sortMap := mapsort.NewMapSorter(mapOfNations)
	sortMap.Sort()
	for i, key := range sortMap.Keys {
		stats = append(stats, model.Statistic{key, sortMap.Values[i]})
	}
	respBody, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	ResponseWithJSON(w, respBody, http.StatusOK)

}

func ErrorWithJSON(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, "message: { %s }", message)
}

func ResponseWithJSON(w http.ResponseWriter, json []byte, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(json)
}

func SetUrlParams(r *http.Request) model.UrlParams {
	var page, pageSize int

	p := r.URL.Query().Get("page")
	if p != "" {
		page, _ = strconv.Atoi(p)
	} else {
		page = 1
	}
	ps := r.URL.Query().Get("page_size")
	if ps != "" {
		pageSize, _ = strconv.Atoi(ps)
	} else {
		// default page size
		pageSize = 100
	}
	orderBy := r.URL.Query().Get("order_by")
	orderDirection := r.URL.Query().Get("order_direction")
	if orderDirection == "desc" {
		orderDirection = "-"
	} else {
		orderDirection = ""
	}
	filterBy := r.URL.Query().Get("filter_by")
	filterValue := r.URL.Query().Get("filter_value")
	filterOperator := r.URL.Query().Get("filter_operator")

	urlParams := model.UrlParams{page, pageSize, orderBy, orderDirection, filterBy, filterValue, filterOperator}
	return urlParams
}
