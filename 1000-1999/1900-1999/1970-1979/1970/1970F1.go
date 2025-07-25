package main

import (
	"bufio"
	"fmt"
	"os"
)

type Player struct {
	x, y  int
	carry bool
}

type Quaffle struct {
	x, y    int
	carrier string // empty if on ground
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N, M int
	if _, err := fmt.Fscan(in, &N, &M); err != nil {
		return
	}

	redGoal := make([][]bool, N)
	blueGoal := make([][]bool, N)
	for i := range redGoal {
		redGoal[i] = make([]bool, M)
		blueGoal[i] = make([]bool, M)
	}

	players := make(map[string]*Player)
	var quaffle Quaffle

	for i := 0; i < N; i++ {
		for j := 0; j < M; j++ {
			var cell string
			fmt.Fscan(in, &cell)
			switch cell {
			case "..":
				// empty
			case ".Q":
				quaffle.x, quaffle.y = i, j
				quaffle.carrier = ""
			case "RG":
				redGoal[i][j] = true
			case "BG":
				blueGoal[i][j] = true
			default:
				// player code like R0 or B9
				players[cell] = &Player{x: i, y: j}
			}
		}
	}

	var T int
	fmt.Fscan(in, &T)

	redScore, blueScore := 0, 0
	midX, midY := N/2, M/2

	for t := 0; t < T; t++ {
		var entity, action string
		fmt.Fscan(in, &entity, &action)

		if action == "C" {
			var ball string
			fmt.Fscan(in, &ball)
			p, ok := players[entity]
			if ok && ball == ".Q" && quaffle.carrier == "" && quaffle.x == p.x && quaffle.y == p.y {
				quaffle.carrier = entity
				p.carry = true
			}
			continue
		}

		if action == "T" {
			p, ok := players[entity]
			if ok && p.carry {
				quaffle.x, quaffle.y = p.x, p.y
				quaffle.carrier = ""
				p.carry = false

				scored := false
				goalColor := byte(0)
				if redGoal[quaffle.x][quaffle.y] {
					scored = true
					goalColor = 'R'
				} else if blueGoal[quaffle.x][quaffle.y] {
					scored = true
					goalColor = 'B'
				}
				if scored {
					playerTeam := entity[0]
					scoringTeam := playerTeam
					if playerTeam == goalColor {
						if playerTeam == 'R' {
							scoringTeam = 'B'
						} else {
							scoringTeam = 'R'
						}
					}
					if scoringTeam == 'R' {
						redScore++
						fmt.Fprintf(out, "%d RED GOAL\n", t)
					} else {
						blueScore++
						fmt.Fprintf(out, "%d BLUE GOAL\n", t)
					}
					quaffle.x, quaffle.y = midX, midY
				}
			}
			continue
		}

		dx, dy := 0, 0
		switch action {
		case "U":
			dx = -1
		case "D":
			dx = 1
		case "L":
			dy = -1
		case "R":
			dy = 1
		}

		if entity == ".Q" {
			quaffle.x += dx
			quaffle.y += dy
		} else {
			p := players[entity]
			p.x += dx
			p.y += dy
			if p.carry {
				quaffle.x, quaffle.y = p.x, p.y
			}
		}
	}

	fmt.Fprintf(out, "FINAL SCORE: %d %d\n", redScore, blueScore)
}
