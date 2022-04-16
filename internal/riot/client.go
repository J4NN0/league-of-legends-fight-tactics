package riot

import (
	"fmt"
	"league-of-legends-fight-tactics/pkg/httpclient"
	"net/http"
	"strconv"
	"strings"
)

// Docs: https://developer.riotgames.com/docs/lol#data-dragon_champions
const dDragonLolAllChampionsUrl string = "https://ddragon.leagueoflegends.com/cdn/12.3.1/data/en_US/champion.json"
const dDragonLolChampionBaseUrl string = "https://ddragon.leagueoflegends.com/cdn/12.3.1/data/en_US/champion"

type ApiClient struct {
	log Logger
	hc  *http.Client
}

type Logger interface {
	Printf(fmt string, args ...interface{})
	Warningf(fmt string, args ...interface{})
	Fatalf(fmt string, args ...interface{})
}

func NewApiClient(log Logger, hc *http.Client) *ApiClient {
	return &ApiClient{log: log, hc: hc}
}

type dDragonLoLAllChampionsResponse struct {
	Format  string                 `json:"format"`
	Version string                 `json:"version"`
	Data    map[string]interface{} `json:"data"`
}

type DDragonChampionResponse struct {
	Format   string `json:"format"`
	Version  string `json:"version"`
	DataName string
	Data     map[string]championData `json:"data"`
}

type championData struct {
	Name   string   `json:"name"`
	Tags   []string `json:"tags"`
	Stats  stats    `json:"stats"`
	Spells []spell  `json:"spells"`
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
}

type spell struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	MaxRank    int       `json:"maxrank"`
	Cooldown   []float32 `json:"cooldown"`
	Leveltip   leveltip  `json:"leveltip"`
	EffectBurn []string  `json:"effectBurn"`
	Damage     []float32
}

type leveltip struct {
	Label []string `json:"label"`
}

func (c *ApiClient) GetAllLoLChampions() (championsData []DDragonChampionResponse, err error) {
	var allChampionsResponse dDragonLoLAllChampionsResponse

	err = httpclient.Get(c.hc, dDragonLolAllChampionsUrl, &allChampionsResponse)
	if err != nil {
		return []DDragonChampionResponse{}, err
	}

	for championName := range allChampionsResponse.Data {
		c.log.Printf("Fetching %s ...", championName)
		championResponse, err := c.GetLoLChampion(championName)
		if err != nil {
			c.log.Warningf("Could not fetch champion %s: %v", championName, err)
		}
		championsData = append(championsData, championResponse)
	}

	return championsData, nil
}

func (c *ApiClient) GetLoLChampion(championName string) (championResponse DDragonChampionResponse, err error) {
	sanitizedChampionName := sanitizeChampionName(championName)
	err = httpclient.Get(c.hc, getChampionUrl(sanitizedChampionName), &championResponse)
	if err != nil {
		return DDragonChampionResponse{}, err
	}

	championResponse.DataName = sanitizedChampionName
	for i, s := range championResponse.Data[sanitizedChampionName].Spells {
		championResponse.Data[strings.Title(championResponse.DataName)].Spells[i].Damage = c.getSpellDamage(s)
	}

	return championResponse, nil
}

// sanitizeChampionName endpoint dDragonLolChampionBaseUrl wants champion name with first letter capitalized (e.g. TwistedFate, Jhin)
func sanitizeChampionName(championName string) string {
	switch strings.ToLower(championName) {
	case "aurelionsol":
		return "AurelionSol"
	case "drmundo":
		return "DrMundo"
	case "jarvaniv":
		return "JarvanIV"
	case "kogmaw":
		return "KogMaw"
	case "leesin":
		return "LeeSin"
	case "masteryi":
		return "MasterYi"
	case "missfortune":
		return "MissFortune"
	case "monkeyking":
		return "MonkeyKing"
	case "reksai":
		return "RekSai"
	case "tahmkench":
		return "TahmKench"
	case "twistedfate":
		return "TwistedFate"
	case "xinzhao":
		return "XinZhao"
	default:
		return strings.Title(strings.ToLower(championName))
	}
}

func getChampionUrl(championName string) string {
	return fmt.Sprintf("%s/%s.json", dDragonLolChampionBaseUrl, championName)
}

// setSpellsDamage
// "ChampionName": {
// 	"spells": [{
//		"id": "Q",
//  	"leveltip":{
//   		"label":[
//   		"Damage", <------------------------ (at index 1)
//   		"Attack Damage",
//   		"Cooldown",
//   		"@AbilityResourceName@ Cost",
//		]}
//		"effectBurn":[
//			null,
//			"95/130/165/200/235", <------------------------ (get damage at index 1)
//			"60/75/90/105/120",
//			...,
//  	],
//		...,
//	}
// }
func (c *ApiClient) getSpellDamage(spell spell) []float32 {
	var spellDamages []float32

	for i, l := range spell.Leveltip.Label {
		if l == "Damage" {
			spellsDamagePerLevelString := strings.Split(spell.EffectBurn[i+1], "/")
			for _, damageLevel := range spellsDamagePerLevelString {
				spellDamagePerLevel, err := strconv.ParseFloat(damageLevel, 32)
				if err != nil {
					c.log.Warningf("Could not set %s spell damage: %v", spell.ID, err)
					return []float32{0, 0, 0, 0, 0}
				}
				spellDamages = append(spellDamages, float32(spellDamagePerLevel))
			}
		}
	}

	if len(spellDamages) == 0 {
		spellDamages = []float32{0, 0, 0, 0, 0}
	}

	return spellDamages
}
