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

// Embedded solver for computing optimal k (minimum path cover via Hopcroft-Karp).
func solveOptimalK(n, m int, rows []string) int {
	bcnt := (n + 63) >> 6
	cols := make([][]uint64, m)
	for j := 0; j < m; j++ {
		cols[j] = make([]uint64, bcnt)
	}
	for i := 0; i < n; i++ {
		word := i >> 6
		bit := uint(i & 63)
		mask := uint64(1) << bit
		for j := 0; j < m; j++ {
			if rows[i][j] == '1' {
				cols[j][word] |= mask
			}
		}
	}

	equalCols := func(a, b int) bool {
		for i := 0; i < bcnt; i++ {
			if cols[a][i] != cols[b][i] {
				return false
			}
		}
		return true
	}

	adj := make([][]int, m)
	for u := 0; u < m; u++ {
		for v := 0; v < m; v++ {
			if u == v {
				continue
			}
			subset := true
			eq := true
			for k := 0; k < bcnt; k++ {
				if cols[u][k]&^cols[v][k] != 0 {
					subset = false
					break
				}
				if cols[u][k] != cols[v][k] {
					eq = false
				}
			}
			if subset && (!eq || u < v) {
				adj[u] = append(adj[u], v)
			}
		}
	}

	_ = equalCols

	// Hopcroft-Karp
	pairU := make([]int, m)
	pairV := make([]int, m)
	dist := make([]int, m)
	for i := 0; i < m; i++ {
		pairU[i] = -1
		pairV[i] = -1
	}

	bfs := func() bool {
		q := make([]int, 0, m)
		for u := 0; u < m; u++ {
			if pairU[u] == -1 {
				dist[u] = 0
				q = append(q, u)
			} else {
				dist[u] = -1
			}
		}
		found := false
		for head := 0; head < len(q); head++ {
			u := q[head]
			for _, v := range adj[u] {
				pu := pairV[v]
				if pu == -1 {
					found = true
				} else if dist[pu] == -1 {
					dist[pu] = dist[u] + 1
					q = append(q, pu)
				}
			}
		}
		return found
	}

	var dfs func(int) bool
	dfs = func(u int) bool {
		for _, v := range adj[u] {
			pu := pairV[v]
			if pu == -1 || (dist[pu] == dist[u]+1 && dfs(pu)) {
				pairU[u] = v
				pairV[v] = u
				return true
			}
		}
		dist[u] = -1
		return false
	}

	for bfs() {
		for u := 0; u < m; u++ {
			if pairU[u] == -1 {
				dfs(u)
			}
		}
	}

	// Count paths = unmatched vertices (starts of chains)
	matching := 0
	for u := 0; u < m; u++ {
		if pairU[u] != -1 {
			matching++
		}
	}
	return m - matching
}

func genTest(seed int64) (int, int, []string, string) {
	r := rand.New(rand.NewSource(seed))
	n := r.Intn(3) + 1
	m := r.Intn(3) + 1
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d\n", n, m)
	rows := make([]string, n)
	for i := 0; i < n; i++ {
		var row strings.Builder
		for j := 0; j < m; j++ {
			if r.Intn(2) == 1 {
				row.WriteByte('1')
			} else {
				row.WriteByte('0')
			}
		}
		rows[i] = row.String()
		b.WriteString(rows[i])
		b.WriteByte('\n')
	}
	return n, m, rows, b.String()
}

func runProg(path string, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func validateOutput(n, m int, rows []string, output string, optimalK int) error {
	fields := strings.Fields(output)
	pos := 0
	nextInt := func() (int, error) {
		if pos >= len(fields) {
			return 0, fmt.Errorf("unexpected end of output")
		}
		v, err := strconv.Atoi(fields[pos])
		pos++
		return v, err
	}

	k, err := nextInt()
	if err != nil {
		return fmt.Errorf("parse k: %v", err)
	}
	if k != optimalK {
		return fmt.Errorf("expected k=%d got k=%d", optimalK, k)
	}

	group := make([]int, m)
	for j := 0; j < m; j++ {
		g, err := nextInt()
		if err != nil {
			return fmt.Errorf("parse group[%d]: %v", j, err)
		}
		if g < 1 || g > k {
			return fmt.Errorf("group[%d]=%d out of range [1,%d]", j, g, k)
		}
		group[j] = g
	}

	threshold := make([]int, m)
	for j := 0; j < m; j++ {
		t, err := nextInt()
		if err != nil {
			return fmt.Errorf("parse threshold[%d]: %v", j, err)
		}
		if t < 1 || t > 1000000000 {
			return fmt.Errorf("threshold[%d]=%d out of range", j, t)
		}
		threshold[j] = t
	}

	devLevel := make([][]int, n)
	for i := 0; i < n; i++ {
		devLevel[i] = make([]int, k)
		for g := 0; g < k; g++ {
			v, err := nextInt()
			if err != nil {
				return fmt.Errorf("parse devLevel[%d][%d]: %v", i, g, err)
			}
			if v < 1 || v > 1000000000 {
				return fmt.Errorf("devLevel[%d][%d]=%d out of range", i, g, v)
			}
			devLevel[i][g] = v
		}
	}

	// Verify access table
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			hasAccess := devLevel[i][group[j]-1] >= threshold[j]
			shouldHave := rows[i][j] == '1'
			if hasAccess != shouldHave {
				return fmt.Errorf("dev %d doc %d: hasAccess=%v shouldHave=%v", i+1, j+1, hasAccess, shouldHave)
			}
		}
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	for i := 0; i < 100; i++ {
		n, m, rows, input := genTest(int64(i))
		optimalK := solveOptimalK(n, m, rows)

		gotOut, err := runProg(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate crashed on test %d: %v\n", i, err)
			os.Exit(1)
		}

		if err := validateOutput(n, m, rows, gotOut, optimalK); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\ninput:\n%s\noutput:\n%s\n", i, err, input, gotOut)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
