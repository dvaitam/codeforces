package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Embedded oracle solver for 1361C
func solveOracle(input string) string {
	content := []byte(input)
	pos := 0

	nextInt := func() int32 {
		for pos < len(content) && (content[pos] < '0' || content[pos] > '9') {
			pos++
		}
		if pos >= len(content) {
			return 0
		}
		var res int32 = 0
		for pos < len(content) && content[pos] >= '0' && content[pos] <= '9' {
			res = res*10 + int32(content[pos]-'0')
			pos++
		}
		return res
	}

	n := int(nextInt())
	if n == 0 {
		return ""
	}

	U := make([]int32, n)
	V := make([]int32, n)
	for i := 0; i < n; i++ {
		U[i] = nextInt()
		V[i] = nextInt()
	}

	findRoot := func(parent []int32, x int32) int32 {
		root := x
		for parent[root] != root {
			root = parent[root]
		}
		curr := x
		for parent[curr] != root {
			nxt := parent[curr]
			parent[curr] = root
			curr = nxt
		}
		return root
	}

	deg := make([]int32, 1<<20)
	parent := make([]int32, 1<<20)

	var bestK int
	for k := 20; k >= 0; k-- {
		mask := int32((1 << k) - 1)

		for i := 0; i < n; i++ {
			u := U[i] & mask
			v := V[i] & mask
			parent[u] = u
			parent[v] = v
			deg[u] = 0
			deg[v] = 0
		}

		for i := 0; i < n; i++ {
			u := U[i] & mask
			v := V[i] & mask
			deg[u]++
			deg[v]++
		}

		valid := true
		for i := 0; i < n; i++ {
			u := U[i] & mask
			v := V[i] & mask
			if deg[u]%2 != 0 || deg[v]%2 != 0 {
				valid = false
				break
			}
		}

		if !valid {
			continue
		}

		for i := 0; i < n; i++ {
			u := U[i] & mask
			v := V[i] & mask
			rootU := findRoot(parent, u)
			rootV := findRoot(parent, v)
			if rootU != rootV {
				parent[rootU] = rootV
			}
		}

		firstRoot := findRoot(parent, U[0]&mask)
		for i := 1; i < n; i++ {
			if findRoot(parent, U[i]&mask) != firstRoot {
				valid = false
				break
			}
		}

		if valid {
			bestK = k
			break
		}
	}

	var sb strings.Builder
	sb.WriteString(strconv.Itoa(bestK))
	sb.WriteByte('\n')

	mask := int32((1 << bestK) - 1)
	head := make([]int32, 1<<20)
	for i := range head {
		head[i] = -1
	}

	next := make([]int32, 2*n)
	to := make([]int32, 2*n)
	pf := make([]int32, 2*n)
	pt := make([]int32, 2*n)
	eid := make([]int32, 2*n)

	edgeCount := int32(0)
	addEdge := func(u, v, id, pFrom, pTo int32) {
		to[edgeCount] = v
		eid[edgeCount] = id
		pf[edgeCount] = pFrom
		pt[edgeCount] = pTo
		next[edgeCount] = head[u]
		head[u] = edgeCount
		edgeCount++
	}

	for i := int32(0); i < int32(n); i++ {
		u := U[i] & mask
		v := V[i] & mask
		addEdge(u, v, i, 2*i+1, 2*i+2)
		addEdge(v, u, i, 2*i+2, 2*i+1)
	}

	usedEdge := make([]bool, n)
	path := make([]int32, 0, 2*n)

	type Frame struct {
		u  int32
		pf int32
		pt int32
	}
	stack := make([]Frame, 0, 2*n+1)
	stack = append(stack, Frame{u: U[0] & mask, pf: -1, pt: -1})

	for len(stack) > 0 {
		u := stack[len(stack)-1].u

		e := head[u]
		for e != -1 && usedEdge[eid[e]] {
			e = next[e]
		}
		head[u] = e

		if e != -1 {
			usedEdge[eid[e]] = true
			stack = append(stack, Frame{u: to[e], pf: pf[e], pt: pt[e]})
		} else {
			f := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if f.pf != -1 {
				path = append(path, f.pt, f.pf)
			}
		}
	}

	for i := len(path) - 1; i >= 0; i-- {
		sb.WriteString(strconv.Itoa(int(path[i])))
		sb.WriteByte(' ')
	}
	sb.WriteByte('\n')
	return sb.String()
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		a := rng.Intn(1024)
		b := rng.Intn(1024)
		sb.WriteString(fmt.Sprintf("%d %d\n", a, b))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const cases = 100
	for i := 1; i <= cases; i++ {
		input := genCase(rng)
		expect := strings.TrimSpace(solveOracle(input))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed\nexpected:\n%s\n\ngot:\n%s\n", i, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}
