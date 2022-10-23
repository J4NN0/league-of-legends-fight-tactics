# League of Legends Fight Tactics

[![MIT Licence](https://badges.frapsoft.com/os/mit/mit.png?v=103)](https://opensource.org/licenses/mit-license.php)

Make two champions fight and get the best combination of spells based on the time (in seconds) it takes to kill the enemy champion.

Each champion data can be downloaded/updated from the official [Data Dragon League of Legends APIs](https://developer.riotgames.com/docs/lol#data-dragon_champions). Data Dragon is a way, used by Riot Games, of centralizing League of Legends game data and assets, including champions, items, runes, summoner spells, and profile icons. All of which can be used by third-party developers.

Once the data has been downloaded (custom data can be also provided as far as they follow the same `.yml` struct), it is possible to make two champions fight and get the best round of spells for which the first champion can kill the second as quickly as possible. The fight takes place only from a data point of view, i.e. no graphic combat is displayed on the screen.

The best spells output is saved in the corresponding `.loltactics` file where you can find the order of the spells, their name/id, damage and how much life the enemy has left in each spell round (until he reaches zero hp).

# Table of Contents

- [Setup](https://github.com/J4NN0/league-of-legends-fight-tactics#setup)
- [Usage](https://github.com/J4NN0/league-of-legends-fight-tactics#usage)
- [Champion Data](https://github.com/J4NN0/league-of-legends-fight-tactics#champion-data)

# Setup

1. Prerequisites

    - Install `golang` from: https://golang.org/doc/install

    - Install packages: 

          brew install make
          brew install golangci-lint

2. Clone repo

       git clone https://github.com/J4NN0/league-of-legends-fight-tactics.git
       cd league-of-legends-fight-tactics

3. Install

       make install-lol-tactics

# Usage

- Show help

      loltactics --help

- Download and/or upload champion data
   - Fetch specific champion data

         loltactics --download=CHAMPION_NAME

   - Fetch all champions data

         loltactics --download_all

- Fight tactics
   - Fight tactics between two (no more) champions (e.g. `lucian` vs `jhin`)

         loltactics --fight=lucian --fight=jhin

   - Generate all fights tactics

         loltactics --tactics

- Clean
  - Clean tactics file

        make clean-lol-fights

  - Clean champions data

        make clean-lol-champions

# Champion Data

Each league of legends champion is described by a `.yml` as follows:
```yml
version: 1.1.0
name: champion name
tags: Fighter, Tank
stats:
  hp: 500
  hp_per_level: 80
  armor: 20
  armor_per_level: 4
  spell_block: 30
  spell_block_per_level: 0.5
  attack_damage: 50
  attack_damage_per_level: 5
  attack_speed: 0.6
  attack_speed_per_level: 2
spells:
  - id: aa
    name: spell name
    damage: 5 
    cooldown: 0 
    cast: 2
  - id: q
    name: name1
    damage: 10
    max_rank: 5
    cooldown:
    - 15
    - 14
    - 13
    - 12
    - 11
    cast: 1
...
```

### Data overview

- `version`: data dragon [version](https://developer.riotgames.com/docs/lol#data-dragon_versions).
- `speels`: Contains the set of spells the champion can use in fight (e.g. `q`, `w`, `e`, `r`), including also auto-attack (i.e. `aa`).
- `cooldown`: Minimum length of time (in seconds) to wait after using an ability before it can be used again.
- `cast`: Length of time (in seconds) needed to summoning a spell.
