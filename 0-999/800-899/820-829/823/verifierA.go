package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func buildRef() (string, error) {
	ref := "./refA.bin"
	cmd := exec.Command("go", "build", "-o", ref, "823A.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type Case struct{ input string }

func genCases() []Case {
	rng := rand.New(rand.NewSource(8230))
	cases := make([]Case, 100)
	for i := range cases {
		n := rng.Intn(18) + 3  // 3..20
		k := rng.Intn(n-2) + 2 // 2..n-1
		cases[i] = Case{fmt.Sprintf("%d %d\n", n, k)}
	}
	return cases
}

func runCase(bin, ref string, c Case) error {
	// parse input to retrieve n and k
	var n, k int
	fmt.Sscanf(c.input, "%d %d", &n, &k)

	expected, err := runBinary(ref, c.input)
	if err != nil {
		return fmt.Errorf("reference failed: %v", err)
	}
	expFields := strings.Fields(expected)
	if len(expFields) == 0 {
		return fmt.Errorf("reference produced no output")
	}
	expD, err := strconv.Atoi(expFields[0])
	if err != nil {
		return fmt.Errorf("invalid reference diameter: %v", err)
	}

	got, err := runBinary(bin, c.input)
	if err != nil {
		return err
	}
	lines := strings.FieldsFunc(got, func(r rune) bool { return r == '\n' || r == '\r' })
	if len(lines) == 0 {
		return fmt.Errorf("no output produced")
	}
	d, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return fmt.Errorf("invalid diameter: %v", err)
	}
	if d != expD {
		return fmt.Errorf("expected diameter %d, got %d", expD, d)
	}
	if len(lines)-1 != n-1 {
		return fmt.Errorf("expected %d edges, got %d", n-1, len(lines)-1)
	}

	adj := make([][]int, n+1)
	for i := 1; i < len(lines); i++ {
		parts := strings.Fields(lines[i])
		if len(parts) != 2 {
			return fmt.Errorf("invalid edge format on line %d", i+1)
		}
		u, err1 := strconv.Atoi(parts[0])
		v, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			return fmt.Errorf("invalid integers on line %d", i+1)
		}
		if u < 1 || u > n || v < 1 || v > n {
			return fmt.Errorf("edge out of range: %d %d", u, v)
		}
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	// BFS to check connectivity
	vis := make([]bool, n+1)
	queue := []int{1}
	vis[1] = true
	for q := 0; q < len(queue); q++ {
		v := queue[q]
		for _, to := range adj[v] {
			if !vis[to] {
				vis[to] = true
				queue = append(queue, to)
			}
		}
	}
	for i := 1; i <= n; i++ {
		if !vis[i] {
			return fmt.Errorf("graph is not connected")
		}
	}

	// count leaves
	leaves := 0
	for i := 1; i <= n; i++ {
		if len(adj[i]) == 1 {
			leaves++
		}
	}
	if leaves != k {
		return fmt.Errorf("expected %d leaves, got %d", k, leaves)
	}

	// helper BFS to compute diameter
	bfs := func(start int) (int, int) {
		dist := make([]int, n+1)
		for i := range dist {
			dist[i] = -1
		}
		q := []int{start}
		dist[start] = 0
		for qi := 0; qi < len(q); qi++ {
			v := q[qi]
			for _, to := range adj[v] {
				if dist[to] == -1 {
					dist[to] = dist[v] + 1
					q = append(q, to)
				}
			}
		}
		far := start
		for i := 1; i <= n; i++ {
			if dist[i] > dist[far] {
				far = i
			}
		}
		return far, dist[far]
	}

	far, _ := bfs(1)
	_, diam := bfs(far)
	if diam != d {
		return fmt.Errorf("reported diameter %d, computed %d", d, diam)
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if bin == "--" {
		if len(os.Args) < 3 {
			fmt.Println("usage: go run verifierA.go /path/to/binary")
			os.Exit(1)
		}
		bin = os.Args[2]
	}
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	cases := genCases()
	for i, c := range cases {
		if err := runCase(bin, ref, c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, c.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
