package main

import (
	"bufio"
	"fmt"
	"os"
)

const base = 9

// encode counts into unique integer
func encode(c [5]int) int {
	res := 0
	mul := 1
	for i := 0; i < 5; i++ {
		res += c[i] * mul
		mul *= base
	}
	return res
}

var configs [][5]int
var idOf [base * base * base * base * base]int

func gen(pos, rem int, cur [5]int) {
	if pos == 4 {
		cur[4] = rem
		id := len(configs)
		configs = append(configs, cur)
		idOf[encode(cur)] = id
		return
	}
	for i := 0; i <= rem; i++ {
		cur[pos] = i
		gen(pos+1, rem-i, cur)
	}
}

func main() {
	gen(0, 8, [5]int{})
	ncfg := len(configs)
	total := ncfg * ncfg * 2

	type edge struct{ to int }

	pre := make([][]int, total)
	deg := make([]int, total)
	status := make([]int8, total) // 0-unknown 1-win 2-lose

	stateID := func(a, b, t int) int { return (a*ncfg+b)*2 + t }

	isZero := func(c [5]int) bool {
		return c[1] == 0 && c[2] == 0 && c[3] == 0 && c[4] == 0
	}

	// build graph
	for idA := 0; idA < ncfg; idA++ {
		for idB := 0; idB < ncfg; idB++ {
			ca := configs[idA]
			cb := configs[idB]
			for turn := 0; turn < 2; turn++ {
				s := stateID(idA, idB, turn)
				if isZero(ca) {
					if turn == 0 {
						status[s] = 1
					} else {
						status[s] = 2
					}
					continue
				}
				if isZero(cb) {
					if turn == 1 {
						status[s] = 1
					} else {
						status[s] = 2
					}
					continue
				}
				if turn == 0 { // Alice moves
					for a := 1; a <= 4; a++ {
						if ca[a] == 0 {
							continue
						}
						for b := 1; b <= 4; b++ {
							if cb[b] == 0 {
								continue
							}
							c := (a + b) % 5
							na := ca
							nb := cb
							na[a]--
							na[c]++
							nb[b]--
							nb[c]++
							ida := idOf[encode(na)]
							idb := idOf[encode(nb)]
							t := stateID(ida, idb, 1)
							pre[t] = append(pre[t], s)
							deg[s]++
						}
					}
				} else { // Bob moves
					for a := 1; a <= 4; a++ {
						if cb[a] == 0 {
							continue
						}
						for b := 1; b <= 4; b++ {
							if ca[b] == 0 {
								continue
							}
							c := (a + b) % 5
							na := ca
							nb := cb
							nb[a]--
							nb[c]++
							na[b]--
							na[c]++
							ida := idOf[encode(na)]
							idb := idOf[encode(nb)]
							t := stateID(ida, idb, 0)
							pre[t] = append(pre[t], s)
							deg[s]++
						}
					}
				}
			}
		}
	}

	// BFS propagation
	q := make([]int, 0)
	for i := 0; i < total; i++ {
		if status[i] != 0 {
			q = append(q, i)
		}
	}
	head := 0
	for head < len(q) {
		s := q[head]
		head++
		for _, p := range pre[s] {
			if status[p] != 0 {
				continue
			}
			if status[s] == 2 {
				status[p] = 1
				q = append(q, p)
			} else if status[s] == 1 {
				deg[p]--
				if deg[p] == 0 {
					status[p] = 2
					q = append(q, p)
				}
			}
		}
	}

	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var f int
		fmt.Fscan(reader, &f)
		var ca [5]int
		var cb [5]int
		for i := 0; i < 8; i++ {
			var x int
			fmt.Fscan(reader, &x)
			ca[x]++
		}
		for i := 0; i < 8; i++ {
			var x int
			fmt.Fscan(reader, &x)
			cb[x]++
		}
		idA := idOf[encode(ca)]
		idB := idOf[encode(cb)]
		s := stateID(idA, idB, f)
		switch status[s] {
		case 1:
			if f == 0 {
				fmt.Fprintln(writer, "Alice")
			} else {
				fmt.Fprintln(writer, "Bob")
			}
		case 2:
			if f == 0 {
				fmt.Fprintln(writer, "Bob")
			} else {
				fmt.Fprintln(writer, "Alice")
			}
		default:
			fmt.Fprintln(writer, "Deal")
		}
	}
}
