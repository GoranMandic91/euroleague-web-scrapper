>Instalation
1) start mongodb with `mongod --dbpath data`
2) run from root `go run main.go`
3) visit http://localhost:8081


>API endpoints:
- `/teams` - get all euroleageu teams
- `/teams/{id}` - get info about specific team
- `/players` - get all euroleague players
- `/players/{id}` - get info about specific player
- `/teams/{id}/players` - get all team players
- `/stats/country-players` - get stats about players by country

