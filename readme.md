# Smogon Stat Tool

A command line tool written in Go that is able to parse Pokemon Showdown usage statistics from the human-readable text files at https://www.smogon.com/stats/ and print information information to the user's console. 

# Motivation

While battling a human opponent and making moves on a strict time limit, it is useful to know which moves a pokemon they are using commonly utilises or in generations without team preview, which unrevealed teammates might accompany them. This data is available in .txt files hosted on Smogon (the most prominent online community for competitive 1v1 Pokemon battling), but it can be prohibotively slow to find the relevant data. Smogon Stat Tool is intended to be a command line utility that allows users to quickly look up this information with a few keystrokes.

# Usage

To 'install':
1. Install Go, if your machine does not have it already
2. Clone this repository to your local machine
3. Open a terminal in that folder
4. Use go run . <arguments> to run the tool

Currently only one mode is implemented, for moveset statistics. To use this mode provide three arguments: the name of a pokemon, the tier to query stats for, and your current skill rating.

# TBA

Many improvements and expansions are planned, including:
 - Predict file paths and walk directories only as a fallback
 - Cache usage data that's already been received
 - Revised argument parsing in place of the mix of args and flags used currently
 - Teammate prediction and suggestion
 - Speed stat estimation