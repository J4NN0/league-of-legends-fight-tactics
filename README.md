# League of Legends Fight Tactics

[![MIT Licence](https://badges.frapsoft.com/os/mit/mit.png?v=103)](https://opensource.org/licenses/mit-license.php)

Make two champions fight and get the best combination of spells based on the time (in seconds) it takes to kill the enemy champion.

Each champion data can be downloaded/updated using the official [Data Dragon League of Legends APIs](https://developer.riotgames.com/docs/lol#data-dragon_champions). Data Dragon is a way, used by Riot Games, of centralizing League of Legends game data and assets, including champions, items, runes, summoner spells, and profile icons. All of which can be used by third-party developers.

Once the data has been downloaded, it is possible to make two champions fight (virtually) and get the best spell output for which the first champion can kill the second as quickly as possible. The fight takes place only from a data point of view, i.e. no graphic combat is displayed on the screen.

The best spells output is saved in the corresponding `.loltactics` file where you can find the order of the spells, their name/id, damage and how much life the enemy has left in each spell round (until he reaches zero hp).

# Table of Contents

- [Usage](https://github.com/J4NN0/league-of-legends-fight-tactics#usage)
- [Champion Data](https://github.com/J4NN0/league-of-legends-fight-tactics#champion-data)

# Usage

1. Prerequisites

    - Install `golang` from: https://golang.org/doc/install

    - Install packages: 

          brew install make
          brew install golangci-lint

2. Clone repo

       git clone https://github.com/J4NN0/league-of-legends-fight-tactics.git
       cd league-of-legends-fight-tactics

3. Run 

   1. Download and/or upload champion data
      - Fetch specific champion data
   
            make run-fetch f=CHAMPION_NAME
   
      - Fetch all champions data

            make run-fetchall

   2. Fight tactics
      - Fight tactics between two specific champions (`c1` vs `c2`)

            make run c1=CHAMPION_NAME c2=CHAMPION_NAME
   
      - Generate all fights tactics

            make run-all

# Champion Data

Each league of legends champion is described by a `.yml` file with the following struct:

```yml
version: 1.1.0
name: champion name
tags: Fighter, Tank
stats:
  hp: 0
  hp_per_level: 580
  armor: 38
  armor_per_level: 3.25
  spell_block: 32
  spell_block_per_level: 1.25
  attack_damage: 60
  attack_damage_per_level: 5
  attack_speed: 0.651
  attack_speed_per_level: 2.5
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

### Note

- `version`: data dragon [version](https://developer.riotgames.com/docs/lol#data-dragon_versions).
- `speels`: Contains the set of spells the champion can use in fight (e.g. `q`, `w`, `e`, `r`), including also auto-attack (i.e. `aa`).
- `cooldown`: Minimum length of time (in seconds) to wait after using an ability before it can be used again.
- `cast`: Length of time (in seconds) needed to summoning a spell.
