package riot

import (
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"league-of-legends-fight-tactics/internal/log"
	"league-of-legends-fight-tactics/pkg/httpclient"
	"net/http"
	"strconv"
	"strings"
)

// Docs: https://developer.riotgames.com/docs/lol#data-dragon_champions
const (
	dDragonLolAllChampionsURL = "https://ddragon.leagueoflegends.com/cdn/12.3.1/data/en_US/champion.json"
	dDragonLolChampionBaseURL = "https://ddragon.leagueoflegends.com/cdn/12.3.1/data/en_US/champion"
)

type Client interface {
	GetAllLoLChampions() (championsData []DDragonChampionResponse, err error)
	GetLoLChampion(championName string) (championResponse DDragonChampionResponse, err error)
	getSpellDamage(spell spell) []float64
}

type Concrete struct {
	log log.Logger
	hc  *http.Client
}

func NewClient(log log.Logger, hc *http.Client) Client {
	return &Concrete{log: log, hc: hc}
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
	Hp                   float64 `json:"hp"`
	HpPerLevel           float64 `json:"hpperlevel"`
	Armor                float64 `json:"armor"`
	ArmorPerLevel        float64 `json:"armorperlevel"`
	SpellBlock           float64 `json:"spellblock"`
	SpellBlockPerLevel   float64 `json:"spellblockperlevel"`
	AttackDamage         float64 `json:"attackdamage"`
	AttackDamagePerLevel float64 `json:"attackdamageperlevel"`
	AttackSpeed          float64 `json:"attackspeed"`
	AttackSpeedPerLevel  float64 `json:"attackspeedperlevel"`
}

type spell struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	MaxRank    int       `json:"maxrank"`
	Cooldown   []float64 `json:"cooldown"`
	LevelTip   levelTip  `json:"leveltip"`
	EffectBurn []string  `json:"effectBurn"`
	Damage     []float64
}

type levelTip struct {
	Label []string `json:"label"`
}

func (c *Concrete) GetAllLoLChampions() (championsData []DDragonChampionResponse, err error) {
	var allChampionsResponse dDragonLoLAllChampionsResponse

	err = httpclient.Get(c.hc, dDragonLolAllChampionsURL, &allChampionsResponse)
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

func (c *Concrete) GetLoLChampion(championName string) (championResponse DDragonChampionResponse, err error) {
	sanitizedChampionName := sanitizeChampionName(championName)
	err = httpclient.Get(c.hc, getChampionURL(sanitizedChampionName), &championResponse)
	if err != nil {
		return DDragonChampionResponse{}, err
	}

	championResponse.DataName = sanitizedChampionName
	for i, s := range championResponse.Data[sanitizedChampionName].Spells {
		championResponse.Data[cases.Title(language.English).String(championResponse.DataName)].Spells[i].Damage = c.getSpellDamage(s)
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
		return cases.Title(language.English).String(championName)
	}
}

func getChampionURL(championName string) string {
	return fmt.Sprintf("%s/%s.json", dDragonLolChampionBaseURL, championName)
}

// setSpellsDamage
//
//	"ChampionName": {
//		"spells": [{
//			"id": "Q",
//	 	"leveltip":{
//	  		"label":[
//	  		"Damage", <------------------------ (at index 1)
//	  		"Attack Damage",
//	  		"Cooldown",
//	  		"@AbilityResourceName@ Cost",
//			]}
//			"effectBurn":[
//				null,
//				"95/130/165/200/235", <------------------------ (get damage at index 1)
//				"60/75/90/105/120",
//				...,
//	 	],
//			...,
//		}
//	}
func (c *Concrete) getSpellDamage(spell spell) []float64 {
	var spellDamages []float64

	for i, l := range spell.LevelTip.Label {
		if l == "Damage" {
			spellsDamagePerLevelString := strings.Split(spell.EffectBurn[i+1], "/")
			for _, damageLevel := range spellsDamagePerLevelString {
				spellDamagePerLevel, err := strconv.ParseFloat(damageLevel, 32)
				if err != nil {
					c.log.Warningf("Could not set %s spell damage: %v", spell.ID, err)
					return []float64{0, 0, 0, 0, 0}
				}
				spellDamages = append(spellDamages, spellDamagePerLevel)
			}
		}
	}

	if len(spellDamages) == 0 {
		spellDamages = []float64{0, 0, 0, 0, 0}
	}

	return spellDamages
}
