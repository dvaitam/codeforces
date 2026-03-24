package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const solverSrc = `package main

import (
	"io"
	"os"
	"sort"
	"strconv"
)

func nextInt(data []byte, idx *int) int {
	n := len(data)
	i := *idx
	for i < n {
		c := data[i]
		if c >= '0' && c <= '9' {
			break
		}
		i++
	}
	val := 0
	for i < n {
		c := data[i]
		if c < '0' || c > '9' {
			break
		}
		val = val*10 + int(c-'0')
		i++
	}
	*idx = i
	return val
}

func main() {
	data, _ := io.ReadAll(os.Stdin)
	idx := 0

	n := nextInt(data, &idx)
	m := nextInt(data, &idx)
	k := nextInt(data, &idx)

	adj := make([][]int, n+1)
	prev := make([]int, 1, 2*m+1)
	curr := make([]int, 1, 2*m+1)
	curr[0] = 1

	mul := n + 1
	pairID := make(map[int]int, 2*m+1)

	for i := 0; i < m; i++ {
		u := nextInt(data, &idx)
		v := nextInt(data, &idx)

		id := len(prev)
		prev = append(prev, u)
		curr = append(curr, v)
		adj[u] = append(adj[u], id)
		pairID[u*mul+v] = id

		id = len(prev)
		prev = append(prev, v)
		curr = append(curr, u)
		adj[v] = append(adj[v], id)
		pairID[v*mul+u] = id
	}

	forb := make([][]int, len(prev))
	for i := 0; i < k; i++ {
		a := nextInt(data, &idx)
		b := nextInt(data, &idx)
		c := nextInt(data, &idx)
		if id, ok := pairID[a*mul+b]; ok {
			forb[id] = append(forb[id], c)
		}
	}

	for i := 1; i < len(forb); i++ {
		if len(forb[i]) > 1 {
			sort.Ints(forb[i])
		}
	}

	parent := make([]int, len(prev))
	for i := range parent {
		parent[i] = -2
	}
	parent[0] = -1

	queue := make([]int, 1, len(prev))
	queue[0] = 0
	end := -1

	for head := 0; head < len(queue); head++ {
		s := queue[head]
		if curr[s] == n {
			end = s
			break
		}

		fl := forb[s]
		b := curr[s]

		if len(fl) == 0 {
			for _, nid := range adj[b] {
				if parent[nid] == -2 {
					parent[nid] = s
					queue = append(queue, nid)
				}
			}
		} else if len(fl) <= 8 {
			for _, nid := range adj[b] {
				c := curr[nid]
				blocked := false
				for _, x := range fl {
					if x == c {
						blocked = true
						break
					}
				}
				if blocked {
					continue
				}
				if parent[nid] == -2 {
					parent[nid] = s
					queue = append(queue, nid)
				}
			}
		} else {
			for _, nid := range adj[b] {
				c := curr[nid]
				l, r := 0, len(fl)
				for l < r {
					mid := (l + r) >> 1
					if fl[mid] < c {
						l = mid + 1
					} else {
						r = mid
					}
				}
				if l < len(fl) && fl[l] == c {
					continue
				}
				if parent[nid] == -2 {
					parent[nid] = s
					queue = append(queue, nid)
				}
			}
		}
	}

	if end == -1 {
		os.Stdout.Write([]byte("-1\n"))
		return
	}

	chain := make([]int, 0)
	for s := end; s != -1; s = parent[s] {
		chain = append(chain, s)
	}

	out := make([]byte, 0, len(chain)*8)
	out = strconv.AppendInt(out, int64(len(chain)-1), 10)
	out = append(out, '\n')
	for i := len(chain) - 1; i >= 0; i-- {
		if i != len(chain)-1 {
			out = append(out, ' ')
		}
		out = strconv.AppendInt(out, int64(curr[chain[i]]), 10)
	}
	out = append(out, '\n')
	os.Stdout.Write(out)
}
`

func buildSolver() (string, func(), error) {
	dir, err := os.MkdirTemp("", "verE59")
	if err != nil {
		return "", nil, err
	}
	cleanup := func() { os.RemoveAll(dir) }
	src := filepath.Join(dir, "solver.go")
	if err := os.WriteFile(src, []byte(solverSrc), 0644); err != nil {
		cleanup()
		return "", nil, err
	}
	bin := filepath.Join(dir, "solver")
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("build solver: %v\n%s", err, out)
	}
	return bin, cleanup, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stderr = &out
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCaseE(rng *rand.Rand) string {
	n := rng.Intn(49) + 2 // max n = 50
	edges := make([][2]int, 0)
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			if rng.Intn(2) == 0 {
				edges = append(edges, [2]int{i, j})
			}
		}
	}
	if len(edges) == 0 {
		edges = append(edges, [2]int{1, n})
	}
	forb := make([][3]int, 0)
	maxForb := rng.Intn(3)
	for i := 0; i < maxForb; i++ {
		a := rng.Intn(n) + 1
		b := rng.Intn(n) + 1
		c := rng.Intn(n) + 1
		for a == b {
			b = rng.Intn(n) + 1
		}
		for c == a || c == b {
			c = rng.Intn(n) + 1
		}
		forb = append(forb, [3]int{a, b, c})
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, len(edges), len(forb)))
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	for _, t := range forb {
		fmt.Fprintf(&sb, "%d %d %d\n", t[0], t[1], t[2])
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	ref, cleanup, err := buildSolver()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 5; i++ {
		in := generateCaseE(rng)
		exp, err := run(ref, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d ref failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		// For this problem, compare only the first line (distance) since paths may differ
		expLines := strings.SplitN(exp, "\n", 2)
		outLines := strings.SplitN(out, "\n", 2)
		if expLines[0] != outLines[0] {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
