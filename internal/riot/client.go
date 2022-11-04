package riot

import (
	"fmt"
	"net/http"

	"github.com/J4NN0/league-of-legends-fight-tactics/internal/log"
	"github.com/J4NN0/league-of-legends-fight-tactics/pkg/httpclient"
	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/api"
	"github.com/KnutZuidema/golio/datadragon"
	"github.com/sirupsen/logrus"
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
		golio.WithLogger(logrus.New().WithField("foo", "bar")),
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
	ddChampion, err := c.riotDD.DataDragon.GetChampion(championName)
	if err != nil {
		return datadragon.ChampionDataExtended{}, fmt.Errorf("could not get champion from datadragon: %w", err)
	}
	return ddChampion, nil
}
