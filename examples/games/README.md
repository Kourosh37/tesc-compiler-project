# TesLang Games

This folder contains runnable terminal games written in TesLang.

## Games

- `number_guess.tes`: guess a random number from 1 to 100.
- `dungeon_treasure_hunt.tes`: explore a dungeon, collect gold, survive events, and escape.
- `blackjack.tes`: simple blackjack-style card game.
- `battle_arena.tes`: turn-based arena combat with attacks and potions.
- `quiz_game.tes`: multiple-choice quiz using vectors and loops.
- `hangman_lite.tes`: numeric-letter hangman variant compatible with `scan()`.

## Run

Build the tools:

```powershell
.\scripts\windows\build-all.ps1
```

Run any game:

```powershell
.\bin\tes.exe .\examples\games\number_guess.tes
.\bin\tes.exe .\examples\games\dungeon_treasure_hunt.tes
.\bin\tes.exe .\examples\games\blackjack.tes
.\bin\tes.exe .\examples\games\battle_arena.tes
.\bin\tes.exe .\examples\games\quiz_game.tes
.\bin\tes.exe .\examples\games\hangman_lite.tes
```
