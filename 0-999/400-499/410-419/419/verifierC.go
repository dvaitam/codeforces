package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

type pair struct{ u, v int }

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveC(n, p int, edges []pair) string {
	deg := make([]int, n+1)
	for _, e := range edges {
		deg[e.u]++
		deg[e.v]++
	}
	sortedDeg := make([]int, n)
	for i := 1; i <= n; i++ {
		sortedDeg[i-1] = deg[i]
	}
	sort.Ints(sortedDeg)
	var ans int64
	for i := 0; i < n; i++ {
		need := p - sortedDeg[i]
		j := sort.Search(n, func(j int) bool { return sortedDeg[j] >= need })
		if j < i+1 {
			j = i + 1
		}
		if j < n {
			ans += int64(n - j)
		}
	}
	sort.Slice(edges, func(i, j int) bool {
		if edges[i].u != edges[j].u {
			return edges[i].u < edges[j].u
		}
		return edges[i].v < edges[j].v
	})
	for i := 0; i < len(edges); {
		j := i + 1
		for j < len(edges) && edges[j] == edges[i] {
			j++
		}
		cnt := j - i
		u, v := edges[i].u, edges[i].v
		if deg[u]+deg[v] >= p && deg[u]+deg[v]-cnt < p {
			ans--
		}
		i = j
	}
	return fmt.Sprintf("%d", ans)
}

func generateCaseC(rng *rand.Rand) string {
	n := rng.Intn(6) + 2
	p := rng.Intn(n + 1)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, p)
	for i := 1; i <= n; i++ {
		x := rng.Intn(n) + 1
		for x == i {
			x = rng.Intn(n) + 1
		}
		y := rng.Intn(n) + 1
		for y == i || y == x {
			y = rng.Intn(n) + 1
		}
		fmt.Fprintf(&sb, "%d %d\n", x, y)
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseC(rng)
		scanner := bufio.NewScanner(strings.NewReader(tc))
		scanner.Split(bufio.ScanWords)
		var fields []string
		for scanner.Scan() {
			fields = append(fields, scanner.Text())
		}
		n, _ := strconv.Atoi(fields[0])
		pVal, _ := strconv.Atoi(fields[1])
		edges := make([]pair, n)
		idx := 2
		for j := 0; j < n; j++ {
			x, _ := strconv.Atoi(fields[idx])
			y, _ := strconv.Atoi(fields[idx+1])
			u, v := x, y
			if u > v {
				u, v = v, u
			}
			edges[j] = pair{u, v}
			idx += 2
		}
		expect := solveC(n, pVal, edges)
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
