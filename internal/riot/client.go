package riot

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/J4NN0/league-of-legends-fight-tactics/internal/log"
	"github.com/J4NN0/league-of-legends-fight-tactics/pkg/httpclient"
	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/api"
	"github.com/KnutZuidema/golio/datadragon"
	"github.com/sirupsen/logrus"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Docs: https://developer.riotgames.com/docs/lol#data-dragon_champions
const (
	dDragonLolAllChampionsURL = "https://ddragon.leagueoflegends.com/cdn/12.3.1/data/en_US/champion.json"
)

type Client interface {
	GetAllLoLChampions() ([]datadragon.ChampionDataExtended, error)
	GetLoLChampion(championName string) (datadragon.ChampionDataExtended, error)
}

type Concrete struct {
	log    log.Logger
	hc     *http.Client
	riotDD *golio.Client
}

func NewClient(log log.Logger, hc *http.Client, apiKey, region string) Client {
	client := golio.NewClient(apiKey,
		golio.WithRegion(api.Region(region)),
		golio.WithLogger(logrus.New().WithField("riot", region)),
	)
	return &Concrete{log: log, hc: hc, riotDD: client}
}

type dataDragonLoLAllChampionsResponse struct {
	Format  string                 `json:"format"`
	Version string                 `json:"version"`
	Data    map[string]interface{} `json:"data"`
}

func (c *Concrete) GetAllLoLChampions() ([]datadragon.ChampionDataExtended, error) {
	var ddAllChampionsResp dataDragonLoLAllChampionsResponse
	err := httpclient.Get(c.hc, dDragonLolAllChampionsURL, &ddAllChampionsResp)
	if err != nil {
		return []datadragon.ChampionDataExtended{}, err
	}

	ddChampions := make([]datadragon.ChampionDataExtended, len(ddAllChampionsResp.Data))
	for championName := range ddAllChampionsResp.Data {
		c.log.Printf("Fetching %s ...", championName)
		ddChampion, err := c.GetLoLChampion(championName)
		if err != nil {
			c.log.Warningf("Could not fetch champion %s: %v", championName, err)
		}
		ddChampions = append(ddChampions, ddChampion)
	}

	return ddChampions, nil
}

func (c *Concrete) GetLoLChampion(championName string) (datadragon.ChampionDataExtended, error) {
	sanitizedChampionName := sanitizeChampionName(championName)
	ddChampion, err := c.riotDD.DataDragon.GetChampion(sanitizedChampionName)
	if err != nil {
		return datadragon.ChampionDataExtended{}, fmt.Errorf("could not get champion from datadragon: %w", err)
	}
	return ddChampion, nil
}

// sanitizeChampionName Data Dragon APIs want champion name with first letter capitalized (e.g. TwistedFate, Jhin - plus Wukong’s internal name is monkeyking)
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
	case "wukong":
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
