package scraper

import (
	"github.com/GoranMandic91/euroleague_web_server/model"
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"net/http"
	"strings"
	"sync"
)

func GetAllTeams() []interface{} {

	var arrayOfTeamNodes []*html.Node
	var allTeams []interface{}
	var wg sync.WaitGroup
	queueOfTeams := make(chan model.TeamInfo, 5)

	url := "http://www.euroleague.net/competition/teams"
	arrayOfTeamNodes = GetTeamUrl(url)

	wg.Add(len(arrayOfTeamNodes))
	for i := 0; i < len(arrayOfTeamNodes); i++ {
		go func(i int) {
			nameNode, coachNode, presidentNode, websiteNode := GetTeamData("http://www.euroleague.net" + scrape.Attr(arrayOfTeamNodes[i], "href"))
			name := scrape.Text(nameNode)
			coach := PrettyName(scrape.Text(coachNode))
			president := scrape.Text(presidentNode)
			website := scrape.Attr(websiteNode, "href")
			queueOfTeams <- model.TeamInfo{i, name, coach, president, website}
		}(i)
	}

	go func() {
		for t := range queueOfTeams {
			allTeams = append(allTeams, t)
			wg.Done()
		}
	}()
	wg.Wait()

	return allTeams
}

func GetTeamUrl(url string) []*html.Node {
	resp, err := http.Get(url)
	check(err)

	root, err := html.Parse(resp.Body)
	check(err)

	matcher := func(n *html.Node) bool {
		return n.DataAtom == atom.A && n.Parent != nil && scrape.Attr(n.Parent, "class") == "RoasterName"
	}

	teams := scrape.FindAll(root, matcher)

	return teams
}
func GetTeamData(url string) (*html.Node, *html.Node, *html.Node, *html.Node) {

	resp, err := http.Get(url)
	check(err)

	root, err := html.Parse(resp.Body)
	check(err)

	matchName := func(n *html.Node) bool {
		return n.DataAtom == atom.Div && n.Parent != nil && scrape.Attr(n.Parent, "id") == "team-title" && scrape.Attr(n, "class") == "title"
	}

	matchCoach := func(n *html.Node) bool {
		return n.DataAtom == atom.A && n.Parent != nil && scrape.Attr(n.Parent, "class") == "name" && n.Parent.Parent != nil && scrape.Attr(n.Parent.Parent, "class") == "item coach"
	}

	matchPresident := func(n *html.Node) bool {
		return n.DataAtom == atom.Dd && n.Parent != nil && n.Parent.DataAtom == atom.Dl && n.PrevSibling != nil && n.PrevSibling.PrevSibling != nil && scrape.Attr(n.PrevSibling.PrevSibling, "class") == "president"
	}

	matchWebsite := func(n *html.Node) bool {
		return n.DataAtom == atom.A && n.Parent != nil && scrape.Attr(n.Parent, "class") == "website"
	}

	name, _ := scrape.Find(root, matchName)
	coach, _ := scrape.Find(root, matchCoach)
	president, _ := scrape.Find(root, matchPresident)
	website, _ := scrape.Find(root, matchWebsite)

	return name, coach, president, website
}

func GetAllPlayers() []interface{} {

	var arrayOfPlayerNodes []*html.Node
	var allPlayers []interface{}
	var wg sync.WaitGroup

	gueueOfLetters := make(chan bool, 5)
	queueOfPlayers := make(chan model.PlayerInfo, 5)

	for letter := 1; letter <= 26; letter++ {
		url := "http://www.euroleague.net/competition/players?letter=" + string('A'-1+letter)
		go func() {
			players := GetPlayerUrl(url, gueueOfLetters)
			arrayOfPlayerNodes = append(arrayOfPlayerNodes, players...)
		}()
	}

	for i := 0; i < 26; i++ {
		<-gueueOfLetters
	}

	wg.Add(len(arrayOfPlayerNodes))
	for i := 0; i < len(arrayOfPlayerNodes); i++ {
		go func(i int) {
			var height string
 			teamNode, nationNode, heightNode, bornNode, shirtNode := GetPlayerData("http://www.euroleague.net" + scrape.Attr(arrayOfPlayerNodes[i], "href"))
			team := scrape.Text(teamNode)
			nation := scrape.Text(nationNode)
			if heightNode !=nil {
				height = scrape.Text(heightNode)
			}else {
				height= "Height:  "
			}
			born := scrape.Text(bornNode)
			shirt := scrape.Text(shirtNode)

			queueOfPlayers <- model.PlayerInfo{i, scrape.Text(arrayOfPlayerNodes[i]), team, nation[13:], height[8:], born[6:], shirt}
		}(i)
	}

	go func() {
		for t := range queueOfPlayers {
			allPlayers = append(allPlayers, t)
			wg.Done()
		}
	}()

	wg.Wait()

	return allPlayers

}

func GetPlayerUrl(url string, ch chan bool) []*html.Node {
	resp, err := http.Get(url)
	check(err)

	root, err := html.Parse(resp.Body)
	check(err)

	matcher := func(n *html.Node) bool {
		return n.DataAtom == atom.A && n.Parent != nil && n.NextSibling != nil && n.NextSibling.DataAtom == atom.A && scrape.Attr(n.Parent, "class") == "item"
	}

	players := scrape.FindAll(root, matcher)
	ch <- true
	return players
}

func GetPlayerData(url string) (*html.Node, *html.Node, *html.Node, *html.Node, *html.Node) {
	resp, err := http.Get(url)
	check(err)

	root, err := html.Parse(resp.Body)
	check(err)

	matchTeam := func(n *html.Node) bool {
		return n.DataAtom == atom.A && n.Parent != nil && scrape.Attr(n.Parent, "class") == "club"
	}

	matchNation := func(n *html.Node) bool {
		return n.DataAtom == atom.Span && n.Parent != nil && scrape.Attr(n.Parent, "class") == "summary-second" && strings.Contains(scrape.Text(n), "Nationality:")
	}

	matchHeight := func(n *html.Node) bool {
		return n.DataAtom == atom.Span && n.Parent != nil && scrape.Attr(n.Parent, "class") == "summary-second" && strings.Contains(scrape.Text(n), "Height:")
	}

	matchBorn := func(n *html.Node) bool {
		return n.DataAtom == atom.Span && n.Parent != nil && scrape.Attr(n.Parent, "class") == "summary-second" && strings.Contains(scrape.Text(n), "Born:")
	}

	matchShirt := func(n *html.Node) bool {
		return n.DataAtom == atom.Span && scrape.Attr(n, "class") == "dorsal"
	}

	team, _ := scrape.Find(root, matchTeam)
	nation, _ := scrape.Find(root, matchNation)
	height, _ := scrape.Find(root, matchHeight)
	born, _ := scrape.Find(root, matchBorn)
	shirt, _ := scrape.Find(root, matchShirt)

	return team, nation, height, born, shirt
}

func PrettyName(name string) string {
	s := strings.SplitN(name, ", ", 2)
	name = strings.Title(strings.ToLower(s[1])) + " " + strings.Title(strings.ToLower(s[0]))
	return name
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
