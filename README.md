# League of Legends Fight Tactics

# Table of Contents

- [Usage](https://github.com/J4NN0/league-of-legends-fight-tactics#usage)
- [League of Legends Champion Data](https://github.com/J4NN0/league-of-legends-fight-tactics#league-of-legends-champion-data)

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

# League of Legends Champion Data

A  League of Legends champion is described by a `.yml` file with the following struct:

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

- `speels`: Contains the set of spells the champion can use in fight (e.g. `q`, `w`, `e`, `r`), including also auto-attack (`aa`).
- `cooldown`: Minimum length of time (in seconds) to wait after using an ability before it can be used again.
- `cast`: Length of time (in seconds) needed to summoning a spell.
