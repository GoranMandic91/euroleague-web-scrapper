# Installation
1) Start mongodb `mongod --dbpath data`
2) Run from root `go run main.go`
3) Visit http://localhost:8081


# API endpoints
- `/teams` - get all euroleague teams
- `/teams/{id}` - get info about specific team
- `/players` - get all euroleague players
- `/players/{id}` - get info about specific player
- `/teams/{id}/players` - get all team players
- `/stats/country-players` - get stats about players by 
country
- API supports filtering and paging, you just need to add filter params in query.

# Examples
> `GET: /teams?filter_by=id&filter_value=6&filter_operator=gte`
>
> `GET: /players?order_by=name&order_direction=asc&page=1&page_size=10`


