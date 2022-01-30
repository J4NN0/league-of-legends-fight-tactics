# League of Legends Champion

A champion is described by a `.yml` file with the following struct:

```yml
name: "Name" # [string] champion name
stats:
  hp: 100 # [float32] hp amount
spells: # set of spells the champion can use in fight (including auto-attack)
  - id: aa # [string] id of the spell (auto-attack in this case)
    damage: 5 # [float32] damage of related spell 
    cooldown: 0 # [float32] cooldown (in seconds) of related spell 
    cast: 2 # [float32] cast (in seconds) of related spell
  - id: q
    damage: 10
    cooldown: 5
    cast: 1
  - id: w
    damage: 20
    cooldown: 5
    cast: 2
  - id: e
    damage: 30
    cooldown: 8
    cast: 3
  - id: r
    damage: 40
    cooldown: 10
    cast: 4
```

### Note

- `cooldown`: Minimum length of time to wait after using an ability before it can be used again.
- `cast`: Length of time needed to summoning a spell.
