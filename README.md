# League of Legends Fight Tactics

# Usage

1. Prerequisites

    - Install `golang` from: `https://golang.org/doc/install`

    - Install packages: 

           brew install make
           brew install golangci-lint

2. Clone repo

       git clone https://github.com/J4NN0/league-of-legends-fight-tactics.git
       cd league-of-legends-fight-tactics

3. Run 

   - Fight tactics between `champion1` and `champion2`

         make run c1=<championName> c2=<championName>
   
   - Generate all fights tactics

         make run-all
