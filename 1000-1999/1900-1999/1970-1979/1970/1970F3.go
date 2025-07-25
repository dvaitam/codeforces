package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Player struct {
	team  byte
	x, y  int
	alive bool
	carry bool
}

type Ball struct {
	x, y  int
	carry string
}

type Coord struct{ x, y int }

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N, M int
	if _, err := fmt.Fscan(in, &N, &M); err != nil {
		return
	}

	redGoals := make(map[Coord]bool)
	blueGoals := make(map[Coord]bool)
	players := make(map[string]*Player)
	var quaffle Ball
	var bludger Ball
	var snitch Ball
	hasBludger := false
	hasSnitch := false

	for i := 0; i < N; i++ {
		for j := 0; j < M; j++ {
			var s string
			fmt.Fscan(in, &s)
			c := Coord{i, j}
			switch {
			case s == "..":
			case s == "RG":
				redGoals[c] = true
			case s == "BG":
				blueGoals[c] = true
			case s == ".Q":
				quaffle.x, quaffle.y = i, j
			case s == ".B":
				hasBludger = true
				bludger.x, bludger.y = i, j
			case s == ".S":
				hasSnitch = true
				snitch.x, snitch.y = i, j
			default:
				// player
				if len(s) == 2 {
					p := &Player{team: s[0], x: i, y: j, alive: true}
					players[s] = p
				}
			}
		}
	}

	var T int
	fmt.Fscan(in, &T)

	center := Coord{(N - 1) / 2, (M - 1) / 2}
	scoreR, scoreB := 0, 0

	gameOver := false
	for t := 0; t < T && !gameOver; t++ {
		var ent, act string
		fmt.Fscan(in, &ent, &act)
		var target string
		if act == "C" {
			fmt.Fscan(in, &target)
		}

		if ent[0] == 'R' || ent[0] == 'B' { // player action
			p := players[ent]
			if p == nil || !p.alive {
				continue
			}
			switch act {
			case "U":
				p.x--
			case "D":
				p.x++
			case "L":
				p.y--
			case "R":
				p.y++
			case "C":
				if target == ".Q" {
					p.carry = true
					quaffle.carry = ent
					quaffle.x, quaffle.y = p.x, p.y
				} else if target == ".S" {
					if p.team == 'R' {
						scoreR += 10
						fmt.Fprintf(out, "%d RED CATCH GOLDEN SNITCH\n", t)
					} else {
						scoreB += 10
						fmt.Fprintf(out, "%d BLUE CATCH GOLDEN SNITCH\n", t)
					}
					gameOver = true
				}
			case "T":
				if p.carry {
					p.carry = false
					quaffle.carry = ""
					quaffle.x, quaffle.y = p.x, p.y
					c := Coord{p.x, p.y}
					scored := false
					if p.team == 'R' {
						if blueGoals[c] {
							scoreR++
							fmt.Fprintf(out, "%d RED GOAL\n", t)
							scored = true
						} else if redGoals[c] {
							scoreB++
							fmt.Fprintf(out, "%d BLUE GOAL\n", t)
							scored = true
						}
					} else {
						if redGoals[c] {
							scoreB++
							fmt.Fprintf(out, "%d BLUE GOAL\n", t)
							scored = true
						} else if blueGoals[c] {
							scoreR++
							fmt.Fprintf(out, "%d RED GOAL\n", t)
							scored = true
						}
					}
					if scored {
						quaffle.x, quaffle.y = center.x, center.y
						quaffle.carry = ""
					}
				}
			}
			if p.carry {
				quaffle.x, quaffle.y = p.x, p.y
			}
		} else { // ball action
			switch ent {
			case ".Q":
				if quaffle.carry == "" {
					switch act {
					case "U":
						quaffle.x--
					case "D":
						quaffle.x++
					case "L":
						quaffle.y--
					case "R":
						quaffle.y++
					}
				}
			case ".B":
				if hasBludger {
					switch act {
					case "U":
						bludger.x--
					case "D":
						bludger.x++
					case "L":
						bludger.y--
					case "R":
						bludger.y++
					}
				}
			case ".S":
				if hasSnitch {
					switch act {
					case "U":
						snitch.x--
					case "D":
						snitch.x++
					case "L":
						snitch.y--
					case "R":
						snitch.y++
					}
				}
			}
		}

		// check eliminations after movement
		if act == "U" || act == "D" || act == "L" || act == "R" {
			eliminated := make([]string, 0)
			if hasBludger {
				bx, by := bludger.x, bludger.y
				for id, pl := range players {
					if pl.alive && pl.x == bx && pl.y == by {
						eliminated = append(eliminated, id)
					}
				}
			}
			if len(eliminated) > 0 {
				sort.Strings(eliminated)
				for _, id := range eliminated {
					pl := players[id]
					pl.alive = false
					if pl.carry {
						pl.carry = false
						quaffle.carry = ""
						quaffle.x, quaffle.y = pl.x, pl.y
					}
					fmt.Fprintf(out, "%d %s ELIMINATED\n", t, id)
				}
			}
		}
	}

	fmt.Fprintf(out, "FINAL SCORE: %d %d\n", scoreR, scoreB)
}
