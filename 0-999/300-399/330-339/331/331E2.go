package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Edge struct {
	to int
	V  string
}

type State struct {
	u   int
	typ int
	Q   string
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	buf := make([]byte, 1024*1024)
	scanner.Buffer(buf, 1024*1024)

	nextInt := func() int {
		scanner.Scan()
		res := 0
		for _, b := range scanner.Bytes() {
			res = res*10 + int(b-'0')
		}
		return res
	}

	scanner.Scan()
	if len(scanner.Bytes()) == 0 {
		return
	}
	n := 0
	for _, b := range scanner.Bytes() {
		n = n*10 + int(b-'0')
	}
	m := nextInt()

	adj := make([][]Edge, n+1)
	for i := 0; i < m; i++ {
		u := nextInt()
		v := nextInt()
		k := nextInt()
		visions := make([]byte, k)
		for j := 0; j < k; j++ {
			visions[j] = byte(nextInt())
		}
		adj[u] = append(adj[u], Edge{to: v, V: string(visions)})
	}

	MOD := 1000000007

	curCount := make(map[State]int)
	curPath := make(map[State][]byte)

	for i := 1; i <= n; i++ {
		s := State{u: i, typ: 0, Q: string([]byte{byte(i)})}
		curCount[s] = 1
		curPath[s] = []byte{byte(i)}
	}

	var e1Path []byte
	ansE2 := make([]int, 2*n+1)

	for step := 1; step <= 2*n; step++ {
		nextCount := make(map[State]int)
		nextPath := make(map[State][]byte)

		for state, count := range curCount {
			path := curPath[state]
			if state.typ == 0 {
				for _, e := range adj[state.u] {
					T := state.Q + string([]byte{byte(e.to)})
					Ve := e.V

					var newState State
					valid := false

					if strings.HasPrefix(T, Ve) {
						newState = State{u: e.to, typ: 0, Q: T[len(Ve):]}
						valid = true
					} else if strings.HasPrefix(Ve, T) {
						newState = State{u: e.to, typ: 1, Q: Ve[len(T):]}
						valid = true
					}

					if valid {
						if newState.typ == 1 && len(newState.Q) > 2*n-step {
							continue
						}
						nextCount[newState] = (nextCount[newState] + count) % MOD
						if _, exists := nextPath[newState]; !exists {
							newP := make([]byte, len(path), len(path)+1)
							copy(newP, path)
							newP = append(newP, byte(e.to))
							nextPath[newState] = newP
						}
					}
				}
			} else {
				w := int(state.Q[0])
				for _, e := range adj[state.u] {
					if e.to == w {
						T := state.Q[1:] + e.V
						var newState State
						if len(T) == 0 {
							newState = State{u: w, typ: 0, Q: ""}
						} else {
							newState = State{u: w, typ: 1, Q: T}
						}

						if newState.typ == 1 && len(newState.Q) > 2*n-step {
							continue
						}

						nextCount[newState] = (nextCount[newState] + count) % MOD
						if _, exists := nextPath[newState]; !exists {
							newP := make([]byte, len(path), len(path)+1)
							copy(newP, path)
							newP = append(newP, byte(w))
							nextPath[newState] = newP
						}
						break
					}
				}
			}
		}

		curCount = nextCount
		curPath = nextPath

		ans := 0
		for state, count := range curCount {
			if state.typ == 0 && state.Q == "" {
				ans = (ans + count) % MOD
				if e1Path == nil && len(curPath[state]) <= 2*n {
					e1Path = curPath[state]
				}
			}
		}
		ansE2[step] = ans
	}

	if e1Path == nil {
		fmt.Println(0)
	} else {
		fmt.Println(len(e1Path))
		for i, v := range e1Path {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(v)
		}
		fmt.Println()
	}

	for i := 1; i <= 2*n; i++ {
		fmt.Println(ansE2[i])
	}
}
