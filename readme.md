# Smogon Stat Tool

A command line tool written in Go that is able to parse Pokemon Showdown usage statistics from the human-readable text files at https://www.smogon.com/stats/ and print information information to the user's console. 

# Motivation

While battling a human opponent and making moves on a strict time limit, it is useful to know which moves a pokemon they are using commonly utilises or in generations without team preview, which unrevealed teammates might accompany them. This data is available in .txt files hosted on Smogon (the most prominent online community for competitive 1v1 Pokemon battling), but it can be prohibotively slow to find the relevant data. Smogon Stat Tool is intended to be a command line utility that allows users to quickly look up this information with a few keystrokes.

# Usage

To 'install':
1. Install Go, if your machine does not have it already
2. Clone this repository to your local machine
3. Open a terminal in that folder
4. Use go run . <flags> <mode> <arguments> to run the tool

Currently only one mode is implemented, for moveset statistics. You can specify this mode using any of the following:
    m, mv, move, moves
The moveset mode accepts between three and five arguments. The first three are always the name of the pokemon, the tier to query stats for as a two-letter abbreviation, and either an integer 0-3 representing which of the four ratings stats can be weighted by to use or an integer skill rating to find the closest available stats for. 
Optionally, you can specify the generation to get stats for as a fourth argument, in the format 'gen8'. SST will default to the current generation if one is not specified. 
An override URL to get stats from can be provided as a fourth or fifth argument; in that case SST will check only there for a valid txt file to pull stats from and ignore the provided tier, format and generation. 

One flag is currently implemented, `--forcePathGen`. If provided, SST will always walk the stats folder tree to find the stats file it needs rather than guessing the URL first and walking the folder tree only as a backup. 

### Examples
    go run . m jirachi ou 1700 gen3
This command would display a list of moves commonly used by Jirachi in the gen3ou tier, weighted by the closest available skill value to 1700 (likely 1760).  

    go run . --forcePathGen Dodrio nu 2
This command would display a list of moves commonly used by Dodrio in the gen8nu tier, weighted by the second-highest available skill value (likely 1630). It would skip attempting to guess the correct path for the file we want and instead immediately walk the folder tree from https://www.smogon.com/stats/ to find it. 

# TBA

Several improvements and expansions are planned, including:
 - Cache usage data that's already been received
 - Teammate prediction and suggestion
 - Speed stat estimation
 - 'Insights' mode comparing usage between skill ratings for a Pokemon