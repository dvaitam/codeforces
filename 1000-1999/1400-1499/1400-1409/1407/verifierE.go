package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type InputEdge struct {
	u, v, t int
}

// Reference solver logic
func solveRef(n, m int, inputEdges []InputEdge) (int, string) {
	if n == 1 {
		return 0, "0"
	}
	type Edge struct {
		to   int
		kind int
	}
	revGraph := make([][]Edge, n+1)
	for _, e := range inputEdges {
		revGraph[e.v] = append(revGraph[e.v], Edge{to: e.u, kind: e.t})
	}

	dist := make([]int, n+1)
	for i := range dist {
		dist[i] = -1
	}
	status := make([]int, n+1)
	colors := make([]byte, n+1)
	for i := range colors {
		colors[i] = '0'
	}

	queue := []int{n}
	dist[n] = 0

	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		d := dist[v]

		for _, e := range revGraph[v] {
			u := e.to
			t := e.kind
			if dist[u] != -1 {
				continue
			}
			if (status[u]>>t)&1 == 1 {
				continue
			}
			status[u] |= (1 << t)
			if status[u] == 3 {
				dist[u] = d + 1
				if t == 0 {
					colors[u] = '0'
				} else {
					colors[u] = '1'
				}
				queue = append(queue, u)
			} else {
				if t == 0 {
					colors[u] = '1'
				} else {
					colors[u] = '0'
				}
			}
		}
	}
	
	return dist[1], string(colors[1:])
}

func runCase(bin string, input string) error {
	// Parse input for Reference Solver
	var n, m int
	r := strings.NewReader(input)
	fmt.Fscan(r, &n, &m)
	edges := make([]InputEdge, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(r, &edges[i].u, &edges[i].v, &edges[i].t)
	}

	// Run Reference
	refDist, _ := solveRef(n, m, edges)

	// Run User Solution
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\nstderr: %s", err, errb.String())
	}
	
	userOutput := strings.TrimSpace(out.String())
	var userDist int
	var userColors string
	if _, err := fmt.Sscan(userOutput, &userDist, &userColors); err != nil {
		return fmt.Errorf("invalid output format")
	}

	// Verify User Validity
	if len(userColors) != n {
		return fmt.Errorf("color string length %d != n %d", len(userColors), n)
	}
	
	adj := make([][]int, n+1)
	for _, e := range edges {
		// u is 1-based. userColors is 0-based string (so index u-1).
		c := userColors[e.u-1]
		var color int
		if c == '0' { color = 0 } else { color = 1 }
		if color == e.t {
			adj[e.u] = append(adj[e.u], e.v)
		}
	}
	
	d := make([]int, n+1)
	for i := range d { d[i] = -1 }
	d[1] = 0
	q := []int{1}
	for len(q) > 0 {
		u := q[0]
		q = q[1:]
		if u == n { break }
		for _, v := range adj[u] {
			if d[v] == -1 {
				d[v] = d[u] + 1
				q = append(q, v)
			}
		}
	}
	realDist := d[n]
	
	if userDist != realDist {
		return fmt.Errorf("user claimed %d but actual shortest path is %d", userDist, realDist)
	}
	
	// Verify Optimality
	if refDist == -1 {
		if userDist != -1 {
			return fmt.Errorf("user found path %d but optimal is blocked (-1)", userDist)
		}
	} else {
		if userDist == -1 {
			return fmt.Errorf("user claims -1 but optimal is %d", refDist)
		}
		if userDist < refDist {
			return fmt.Errorf("user found %d, optimal is %d", userDist, refDist)
		}
	}
	
	return nil
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(30) + 2
	m := rng.Intn(60) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < m; i++ {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		t := rng.Intn(2)
		fmt.Fprintf(&sb, "%d %d %d\n", u, v, t)
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		if err := runCase(bin, input); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}