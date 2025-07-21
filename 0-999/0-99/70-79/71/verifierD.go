package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type pair struct{ r, s int }

func cardData(card string, crank, suit string) (int, int) {
	if card[0] == 'J' && card[1] >= '1' && card[1] <= '2' {
		return -1, int(card[1] - '1')
	}
	r := strings.IndexByte(crank, card[0])
	s := strings.IndexByte(suit, card[1])
	return r, s
}

func solveD(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var out bytes.Buffer

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return ""
	}
	crank := "A23456789TJQK"
	suit := "CDHS"
	inDeck := make([][]bool, 13)
	for i := range inDeck {
		inDeck[i] = make([]bool, 4)
		for j := 0; j < 4; j++ {
			inDeck[i][j] = true
		}
	}
	card := make([][]string, n)
	var jokers []int
	for i := 0; i < n; i++ {
		card[i] = make([]string, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(in, &card[i][j])
			r, s := cardData(card[i][j], crank, suit)
			if r >= 0 {
				inDeck[r][s] = false
			} else {
				jokers = append(jokers, s)
			}
		}
	}
	var placeJokers []int
	jreplace := [2][]pair{}
	place := func(y, x int) bool {
		used := make([]bool, 13)
		var localJokers []int
		for dy := 0; dy < 3; dy++ {
			for dx := 0; dx < 3; dx++ {
				cy, cx := y+dy, x+dx
				r, s := cardData(card[cy][cx], crank, suit)
				if r < 0 {
					localJokers = append(localJokers, s)
				} else {
					if used[r] {
						return false
					}
					used[r] = true
				}
			}
		}
		for _, v := range localJokers {
			placeJokers = append(placeJokers, v)
		}
		for i := 0; i < 13; i++ {
			if !used[i] {
				for j := 0; j < 4; j++ {
					if inDeck[i][j] && len(localJokers) > 0 {
						t := localJokers[len(localJokers)-1]
						jreplace[t] = append(jreplace[t], pair{i, j})
						if len(localJokers) > 1 {
							localJokers = localJokers[:len(localJokers)-1]
							break
						}
					}
				}
			}
		}
		return true
	}
	for y1 := 0; y1+3 <= n; y1++ {
		for x1 := 0; x1+3 <= m; x1++ {
			for y2 := 0; y2+3 <= n; y2++ {
				for x2 := 0; x2+3 <= m; x2++ {
					if y2+3 <= y1 || y1+3 <= y2 || x1+3 <= x2 || x2+3 <= x1 {
						placeJokers = placeJokers[:0]
						jreplace = [2][]pair{{}, {}}
						if !place(y1, x1) {
							continue
						}
						if !place(y2, x2) {
							continue
						}
						var rep [2]pair
						if len(placeJokers) == 1 {
							t := placeJokers[0]
							if len(jreplace[t]) == 0 {
								continue
							}
							rep[t] = jreplace[t][0]
						} else if len(placeJokers) == 2 {
							if len(jreplace[0]) == 0 || len(jreplace[1]) == 0 {
								continue
							}
							rep[0] = jreplace[0][0]
							rep[1] = jreplace[1][0]
							if rep[0] == rep[1] && len(jreplace[0]) > 1 {
								rep[0] = jreplace[0][1]
							}
							if rep[0] == rep[1] && len(jreplace[1]) > 1 {
								rep[1] = jreplace[1][1]
							}
						}
						if len(jokers) > len(placeJokers) {
							if len(placeJokers) == 1 {
								other := 1 - placeJokers[0]
								for i := 0; i < 13; i++ {
									for j := 0; j < 4; j++ {
										if inDeck[i][j] && (i != rep[placeJokers[0]].r || j != rep[placeJokers[0]].s) {
											rep[other] = pair{i, j}
										}
									}
								}
							} else {
								cnt := 0
								for i := 0; i < 13; i++ {
									for j := 0; j < 4; j++ {
										if inDeck[i][j] {
											rep[cnt] = pair{i, j}
											cnt = (cnt + 1) % 2
										}
									}
								}
							}
						}
						fmt.Fprintln(&out, "Solution exists.")
						if len(jokers) == 1 {
							t := jokers[0]
							p := rep[t]
							fmt.Fprintf(&out, "Replace J%d with %c%c.\n", t+1, crank[p.r], suit[p.s])
						} else if len(jokers) == 2 {
							fmt.Fprintf(&out, "Replace J1 with %c%c and J2 with %c%c.\n", crank[rep[0].r], suit[rep[0].s], crank[rep[1].r], suit[rep[1].s])
						} else {
							fmt.Fprintln(&out, "There are no jokers.")
						}
						fmt.Fprintf(&out, "Put the first square to (%d, %d).\n", y1+1, x1+1)
						fmt.Fprintf(&out, "Put the second square to (%d, %d).\n", y2+1, x2+1)
						return out.String()
					}
				}
			}
		}
	}
	fmt.Fprintln(&out, "No solution.")
	return out.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(3) + 3
	m := rng.Intn(3) + 3
	ranks := "A23456789TJQK"
	suits := "CDHS"
	deck := make([]string, 0, 54)
	for i := 0; i < len(ranks); i++ {
		for j := 0; j < len(suits); j++ {
			deck = append(deck, string(ranks[i])+string(suits[j]))
		}
	}
	deck = append(deck, "J1", "J2")
	for i := len(deck) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		deck[i], deck[j] = deck[j], deck[i]
	}
	total := n * m
	cards := deck[:total]
	var in strings.Builder
	fmt.Fprintf(&in, "%d %d\n", n, m)
	idx := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				in.WriteByte(' ')
			}
			in.WriteString(cards[idx])
			idx++
		}
		in.WriteByte('\n')
	}
	expect := solveD(in.String())
	return in.String(), expect
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, buf.String())
	}
	if strings.TrimSpace(buf.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected\n%s\ngot\n%s", expected, buf.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
