# League of Legends Fight Tactics

[![MIT Licence](https://badges.frapsoft.com/os/mit/mit.png?v=103)](https://opensource.org/licenses/mit-license.php)

Make two champions fight and get the best combination of spells based on the time (in seconds) taken to kill the enemy champion.

Each champion data can be downloaded/updated from the official [Data Dragon League of Legends APIs](https://developer.riotgames.com/docs/lol#data-dragon_champions). Data Dragon is a way, used by Riot Games, of centralizing League of Legends game data and assets, including champions, items, runes, summoner spells, and profile icons. All of which can be used by third-party developers.

Once the data has been downloaded (custom data can be also provided as far as they follow the same `.yml` struct), it is possible to make two champions fight and get the best round of spells for which the first champion can kill the second as quickly as possible. The fight takes place only from a data point of view, i.e. no graphic combat is displayed on the screen.

The best spells output is saved in the corresponding `.loltactics` file where you can find the order of the spells, their name/id, damage and how much life the enemy has left in each spell round (until he reaches zero hp).

### DISCLAIMER

Unfortunately the **data quality** is not the best and apparently this is a [problem known to riot and its community](https://riot-api-libraries.readthedocs.io/en/latest/ddragon.html#common-issues):
> The data in ddragon is inaccurate, especially champion spell data and item stats. This is an unfortunate situation that is surprisingly difficult to solve. If you want to know why, you can ask on Discord. There is no perfect or even close to perfect source for champion spell data, despite significant effort.

The current best resource should be [League Wikia](https://leagueoflegends.fandom.com/wiki/League_of_Legends_Wiki). Since there is no official API, it is not easy (and mostly not sustainable/feasible over time) to download the data from the previously mentioned site (as it would need web scraping).

In conclusion, this tool will perform at its best if the data quality is medium/good. If you are interested in the outcome of the fight between two champions - and do not want to rely on the data download from Data Dragon League of Legends - you can manually edit the relevant `.yml` file and use the tool as shown below.

Last but not least, take a look at the resources listed below - they might be helpful.

# Table of Contents

- [Setup](https://github.com/J4NN0/league-of-legends-fight-tactics#setup)
- [Usage](https://github.com/J4NN0/league-of-legends-fight-tactics#usage)
- [Champion Data](https://github.com/J4NN0/league-of-legends-fight-tactics#champion-data)
- [Import Package](https://github.com/J4NN0/league-of-legends-fight-tactics#import-package)
- [Resources](https://github.com/J4NN0/league-of-legends-fight-tactics#resources)

# Setup

1. Prerequisites

    - Install `golang` from: https://golang.org/doc/install

    - Install packages: 

          brew install make
          brew install mockery
          brew install golangci-lint

2. Clone repo

       git clone https://github.com/J4NN0/league-of-legends-fight-tactics.git
       cd league-of-legends-fight-tactics

3. Install CLI

       make install-lol-tactics
   
4. Set up environment variables

    | Variable     | Description                                                | Optional |
     --------------|------------------------------------------------------------|----------|
    | RIOT_API_KEY | Riot Developer [API Key](https://developer.riotgames.com). | Yes      |
    | LOL_REGION   | League of Legends region code.                             | No       | 

    Valid `LOL_REGION`:

    ```
    Brasil              = "br1"
    Europe North East   = "eun1"
    Europe West         = "euw1"
    Japan               = "jp1"
    Korea               = "kr"
    Latin America North = "la1"
    Latin America South = "la2"
    North America       = "na1"
    Oceania             = "oc1"
    Turkey              = "tr1"
    Russia              = "ru"
    PBE                 = "pbe1"
    ```

    Before running (either with CLI or `make`), add environment variables above and then source them however you like:

       cp .env.sample .env

# Usage

- Shows list of commands or help for one command

      loltactics help, h

- Download and/or upload champion(s) data
   - Fetch specific champion data (e.g. `jhin`)

         loltactics download, d jhin

   - Fetch all champions data

         loltactics download_all, da, a

- Fight tactics
   - Fight tactics between two (neither less nor more) champions (e.g. `lucian` vs `jhin`)

         loltactics fight, f lucian jhin

   - Generate all fights tactics

         loltactics tactics, t

- Clean
  - Clean tactics file

        make clean-lol-fights

  - Clean champions data

        make clean-lol-champions

# Champion Data

Each league of legends champion is described by a `.yml` as follows:
```yml
id: Chogath
name: Cho'Gath
title: the Terror of the Void
tags: Tank, Mage
passive:
  name: Carnivore
  description: Whenever Cho'Gath kills a unit, he recovers Health and Mana. The values restored increase with Cho'Gath's level.
stats:
  health_points: 644
  attack_damage: 69
  attack_speed: 0
spells:
  - id: aa
    name: Auto Attack
    max_rank: 1
    damage:
      - 69
    cooldown:
      - 0
    cast: 0
  - id: Rupture
    name: Rupture
    max_rank: 5
    damage:
      - 80
      - 135
      - 190
      - 245
      - 300
    cooldown:
      - 6
      - 6
      - 6
      - 6
      - 6
    cast: 0
...
```

### Data overview

- `id`: riot champions internal name (where `name` is the "public" champion's name).
- `speels`: Contains the set of spells the champion can use in fight (e.g. `q`, `w`, `e`, `r`), including also auto-attack (i.e. `aa`).
- `cooldown`: Minimum length of time (in seconds) to wait after using an ability before it can be used again.
- `cast`: Length of time (in seconds) needed to summoning a spell.

# Import Package

You can import tactics tool as external lib and use it as you prefer.

```go
package main

import (
    "fmt"

    "github.com/J4NN0/league-of-legends-fight-tactics/pkg/logger"
    "github.com/J4NN0/league-of-legends-fight-tactics/pkg/lol"
)

func main() {
    log := logger.New("lol-tactics")
    lolTactics := lol.NewTactics(log)
    
    lolChampion1, err := lolTactics.ReadChampion("LOL_CHAMPION_1")
    if err != nil {
        fmt.Printf("Could not load champion: %v\n", err)
        return
    }
    
    lolChampion2, err := lolTactics.ReadChampion("LOL_CHAMPION_2")
    if err != nil {
        fmt.Printf("Could not load champion: %v\n", err)
        return
    }
    
    fightTactic := lolTactics.Fight(lolChampion1, lolChampion2)
    fmt.Printf("Enemy defeated: %v\n", fightTactic)
}
```

# Resources

- [Data Dragon](https://developer.riotgames.com/docs/lol#data-dragon_champions)
- [GitHub Community Dragon](https://github.com/CommunityDragon)
- [Community Dragon](https://raw.communitydragon.org)
- [Raw CDragon](https://raw.communitydragon.org/latest/plugins/rcp-be-lol-game-data/global/default/v1/)
- [Riot API Libraries](https://riot-api-libraries.readthedocs.io/en/latest/libraries.html)
- [League Wikia](https://leagueoflegends.fandom.com/wiki/League_of_Legends_Wiki)
