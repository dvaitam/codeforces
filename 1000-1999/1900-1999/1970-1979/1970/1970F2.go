package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const (
	GoalNone = iota
	GoalRed
	GoalBlue
)

type Player struct {
	id       string
	team     byte
	x, y     int
	carrying bool
	alive    bool
}

type Ball struct {
	x, y    int
	carrier *Player
}

func moveXY(x, y *int, act string) {
	switch act {
	case "U":
		*x = *x - 1
	case "D":
		*x = *x + 1
	case "L":
		*y = *y - 1
	case "R":
		*y = *y + 1
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N, M int
	if _, err := fmt.Fscan(in, &N, &M); err != nil {
		return
	}

	board := make([][]int, N)
	for i := range board {
		board[i] = make([]int, M)
	}

	players := make(map[string]*Player)
	var quaffle Ball
	var bludger Ball
	bludgerExist := false

	for i := 0; i < N; i++ {
		for j := 0; j < M; j++ {
			var s string
			fmt.Fscan(in, &s)
			switch s {
			case "..":
			case "RG":
				board[i][j] = GoalRed
			case "BG":
				board[i][j] = GoalBlue
			case ".Q":
				quaffle.x, quaffle.y = i, j
			case ".B":
				bludger.x, bludger.y = i, j
				bludgerExist = true
			default:
				if len(s) == 2 {
					players[s] = &Player{id: s, team: s[0], x: i, y: j, alive: true}
				}
			}
		}
	}

	midX, midY := (N-1)/2, (M-1)/2

	var T int
	fmt.Fscan(in, &T)

	redScore, blueScore := 0, 0

	for t := 0; t < T; t++ {
		var entity, act string
		fmt.Fscan(in, &entity, &act)
		var target string
		if act == "C" {
			fmt.Fscan(in, &target)
		}

		events := []string{}

		switch entity {
		case ".Q":
			if quaffle.carrier == nil {
				moveXY(&quaffle.x, &quaffle.y, act)
			}
		case ".B":
			if bludgerExist {
				moveXY(&bludger.x, &bludger.y, act)
				elim := []string{}
				for id, p := range players {
					if p.alive && p.x == bludger.x && p.y == bludger.y {
						elim = append(elim, id)
					}
				}
				if len(elim) > 0 {
					sort.Strings(elim)
					for _, id := range elim {
						p := players[id]
						p.alive = false
						if p.carrying {
							p.carrying = false
							quaffle.carrier = nil
							quaffle.x, quaffle.y = p.x, p.y
						}
						events = append(events, fmt.Sprintf("%d %s ELIMINATED", t, id))
					}
				}
			}
		default:
			p := players[entity]
			if p == nil || !p.alive {
				break
			}
			switch act {
			case "U", "D", "L", "R":
				moveXY(&p.x, &p.y, act)
				if p.carrying {
					quaffle.x, quaffle.y = p.x, p.y
				}
				if bludgerExist && p.x == bludger.x && p.y == bludger.y {
					p.alive = false
					if p.carrying {
						p.carrying = false
						quaffle.carrier = nil
						quaffle.x, quaffle.y = p.x, p.y
					}
					events = append(events, fmt.Sprintf("%d %s ELIMINATED", t, p.id))
				}
			case "C":
				if target == ".Q" && quaffle.carrier == nil && p.x == quaffle.x && p.y == quaffle.y {
					quaffle.carrier = p
					p.carrying = true
				}
			case "T":
				if p.carrying {
					p.carrying = false
					quaffle.carrier = nil
					quaffle.x, quaffle.y = p.x, p.y
					goal := board[p.x][p.y]
					if goal == GoalRed {
						events = append(events, fmt.Sprintf("%d BLUE GOAL", t))
						blueScore++
						quaffle.x, quaffle.y = midX, midY
					} else if goal == GoalBlue {
						events = append(events, fmt.Sprintf("%d RED GOAL", t))
						redScore++
						quaffle.x, quaffle.y = midX, midY
					}
				}
			}
		}

		for _, e := range events {
			fmt.Fprintln(out, e)
		}
	}

	fmt.Fprintf(out, "FINAL SCORE: %d %d\n", redScore, blueScore)
}
