//go:generate mockery --case underscore --dir . --name Client --output ./mocks

package riot

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/J4NN0/league-of-legends-fight-tactics/pkg/logger"
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

var wg sync.WaitGroup

const (
	wgN = 30
)

type Client interface {
	GetAllLoLChampions() ([]datadragon.ChampionDataExtended, error)
	GetLoLChampion(championName string) (datadragon.ChampionDataExtended, error)
}

type Concrete struct {
	log    logger.Logger
	hc     *http.Client
	riotDD *golio.Client
}

func NewClient(log logger.Logger, hc *http.Client, apiKey, region string) Client {
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
	err := c.httpGet(dDragonLolAllChampionsURL, &ddAllChampionsResp)
	if err != nil {
		return []datadragon.ChampionDataExtended{}, err
	}

	errChan := make(chan error)
	go func() {
		for e := range errChan {
			c.log.Warningf("%v", e)
		}
	}()

	ddChampions := make([]datadragon.ChampionDataExtended, 0, len(ddAllChampionsResp.Data))
	ddChampionRespChan := make(chan datadragon.ChampionDataExtended)
	go func() {
		for ddChampion := range ddChampionRespChan {
			ddChampions = append(ddChampions, ddChampion)
		}
	}()

	championNameChan := make(chan string, wgN)

	wg.Add(wgN)
	for i := 0; i < wgN; i++ {
		go c.getDDChampionData(championNameChan, ddChampionRespChan, errChan)
	}

	for championName := range ddAllChampionsResp.Data {
		championNameChan <- championName
	}

	close(championNameChan)
	wg.Wait()

	close(errChan)
	close(ddChampionRespChan)

	return ddChampions, nil
}

func (c *Concrete) getDDChampionData(championNameChan chan string, ddChampionRespChan chan datadragon.ChampionDataExtended, errChan chan error) {
	defer wg.Done()
	for championName := range championNameChan {
		c.log.Printf("Fetching %s ...", championName)
		ddChampion, err := c.GetLoLChampion(championName)
		if err != nil {
			errChan <- fmt.Errorf("could not fetch champion %s: %v", championName, err)
		}
		ddChampionRespChan <- ddChampion
	}
}

func (c *Concrete) httpGet(url string, destination interface{}) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create API request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.hc.Do(req)
	if err != nil {
		return fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("got HTTP status code %d: %s", resp.StatusCode, resp.Body)
	}

	err = json.Unmarshal(respBody, &destination)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return nil
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
