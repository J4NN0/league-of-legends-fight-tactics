package riot

import (
	"fmt"
	"league-of-legends-fight-tactics/internal/log"
	"league-of-legends-fight-tactics/pkg/httpclient"
	"net/http"
)

const dDragonLolAllChampionsUrl string = "https://ddragon.leagueoflegends.com/cdn/12.3.1/data/en_US/champion.json"
const dDragonLolChampionBaseUrl string = "https://ddragon.leagueoflegends.com/cdn/12.3.1/data/en_US/champion"

type ApiClient struct {
	log *log.Logger
	hc  *http.Client
}

func NewApiClient(log *log.Logger, hc *http.Client) *ApiClient {
	return &ApiClient{log: log, hc: hc}
}

type dDragonLoLAllChampionsResponse struct {
	Format  string                 `json:"format"`
	Version string                 `json:"version"`
	Data    map[string]interface{} `json:"data"`
}

type DDragonChampionResponse struct {
	Format  string                  `json:"format"`
	Version string                  `json:"version"`
	Data    map[string]championData `json:"data"`
}

type championData struct {
	Name  string   `json:"name"`
	Tags  []string `json:"tags"`
	Stats stats    `json:"stats"`
}

type stats struct {
	Hp                   float32 `json:"hp"`
	HpPerLevel           float32 `json:"hpperlevel"`
	Armor                float32 `json:"armor"`
	ArmorPerLevel        float32 `json:"armorperlevel"`
	SpellBlock           float32 `json:"spellblock"`
	SpellBlockPerLevel   float32 `json:"spellblockperlevel"`
	AttackDamage         float32 `json:"attackdamage"`
	AttackDamagePerLevel float32 `json:"attackdamageperlevel"`
	AttackSpeed          float32 `json:"attackspeed"`
	AttackSpeedPerLevel  float32 `json:"attackspeedperlevel"`
	Spells               []spell `json:"spells"`
}

type spell struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	MaxRank  int       `json:"maxrank"`
	Cooldown []float32 `json:"cooldown"`
}

func (c *ApiClient) FetchAllLoLChampions() (championsData []DDragonChampionResponse, err error) {
	var allChampionsResponse dDragonLoLAllChampionsResponse

	c.log.Printf("Fetching all league of legends champions ...\n")
	err = httpclient.Get(c.hc, dDragonLolAllChampionsUrl, &allChampionsResponse)
	if err != nil {
		return []DDragonChampionResponse{}, err
	}

	for championName := range allChampionsResponse.Data {
		c.log.Printf("Fetching %s ...\n", championName)
		championResponse, err := c.GetLoLChampion(championName)
		if err != nil {
			c.log.Warningf("Could not fetch champion %s: %v\n", championName, err)
		}
		championsData = append(championsData, championResponse)
	}

	return championsData, nil
}

func (c *ApiClient) GetLoLChampion(championName string) (championResponse DDragonChampionResponse, err error) {
	err = httpclient.Get(c.hc, getChampionUrl(championName), &championResponse)
	if err != nil {
		return DDragonChampionResponse{}, err
	}
	return championResponse, nil
}

func getChampionUrl(championName string) string {
	return fmt.Sprintf("%s/%s.json", dDragonLolChampionBaseUrl, championName)
}
