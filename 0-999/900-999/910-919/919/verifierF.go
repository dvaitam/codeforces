package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

var (
	handToID   [65536]int
	hands      [][4]int
	numHands   int
	status     []int8 // 0: Deal, 1: Win, 2: Loss
	degree     []int
	adjRev     [][]int
)

const (
	RES_DEAL = 0
	RES_WIN  = 1
	RES_LOSS = 2
)

func genHands(idx int, sum int, current [4]int) {
	if idx == 4 {
		for c := 0; c <= 8-sum; c++ {
			current[3] = c
			key := current[0] | (current[1] << 4) | (current[2] << 8) | (current[3] << 12)
			handToID[key] = numHands
			hands = append(hands, current)
			numHands++
		}
		return
	}
	for c := 0; c <= 8-sum; c++ {
		current[idx-1] = c
		genHands(idx+1, sum+c, current)
	}
}

func getID(cnt [4]int) int {
	key := cnt[0] | (cnt[1] << 4) | (cnt[2] << 8) | (cnt[3] << 12)
	return handToID[key]
}

func main() {
	for i := range handToID {
		handToID[i] = -1
	}
	genHands(1, 0, [4]int{})

	totalStates := numHands * numHands
	status = make([]int8, totalStates)
	degree = make([]int, totalStates)
	adjRev = make([][]int, totalStates)

	q := make([]int, 0, totalStates)

	for i := 0; i < numHands; i++ {
		for j := 0; j < numHands; j++ {
			u := i*numHands + j
			if i == 0 {
				status[u] = RES_WIN
				q = append(q, u)
				continue
			}
			if j == 0 {
				status[u] = RES_LOSS
				q = append(q, u)
				continue
			}

			hA := hands[i]
			hB := hands[j]
			countMoves := 0

			for valA := 1; valA <= 4; valA++ {
				if hA[valA-1] > 0 {
					for valB := 1; valB <= 4; valB++ {
						if hB[valB-1] > 0 {
							newVal := (valA + valB) % 5
							hNextA := hA
							hNextA[valA-1]--
							if newVal != 0 {
								hNextA[newVal-1]++
							}
							nextAID := getID(hNextA)
							v := j*numHands + nextAID
							adjRev[v] = append(adjRev[v], u)
							countMoves++
						}
					}
				}
			}
			degree[u] = countMoves
		}
	}

	head := 0
	for head < len(q) {
		v := q[head]
		head++
		res := status[v]
		for _, u := range adjRev[v] {
			if status[u] != RES_DEAL {
				continue
			}
			if res == RES_LOSS {
				status[u] = RES_WIN
				q = append(q, u)
			} else if res == RES_WIN {
				degree[u]--
				if degree[u] == 0 {
					status[u] = RES_LOSS
					q = append(q, u)
				}
			}
		}
	}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	if !scanner.Scan() {
		return
	}
	T, _ := strconv.Atoi(scanner.Text())

	for i := 0; i < T; i++ {
		scanner.Scan()
		f, _ := strconv.Atoi(scanner.Text())
		var cntA, cntB [4]int
		for k := 0; k < 8; k++ {
			scanner.Scan()
			val, _ := strconv.Atoi(scanner.Text())
			if val != 0 {
				cntA[val-1]++
			}
		}
		for k := 0; k < 8; k++ {
			scanner.Scan()
			val, _ := strconv.Atoi(scanner.Text())
			if val != 0 {
				cntB[val-1]++
			}
		}

		idA := getID(cntA)
		idB := getID(cntB)

		if idA == 0 {
			if f == 0 {
				fmt.Fprintln(writer, "Alice")
			} else {
				fmt.Fprintln(writer, "Alice")
			}
			continue
		}
		if idB == 0 {
			if f == 1 {
				fmt.Fprintln(writer, "Bob")
			} else {
				fmt.Fprintln(writer, "Bob")
			}
			continue
		}

		var res int8
		if f == 0 {
			res = status[idA*numHands+idB]
			if res == RES_WIN {
				fmt.Fprintln(writer, "Alice")
			} else if res == RES_LOSS {
				fmt.Fprintln(writer, "Bob")
			} else {
				fmt.Fprintln(writer, "Deal")
			}
		} else {
			res = status[idB*numHands+idA]
			if res == RES_WIN {
				fmt.Fprintln(writer, "Bob")
			} else if res == RES_LOSS {
				fmt.Fprintln(writer, "Alice")
			} else {
				fmt.Fprintln(writer, "Deal")
			}
		}
	}
}