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

// ---------- embedded solver for 1236D ----------

// solveEmbedded implements the correct solution for CF 1236D.
// The doll starts at (1,1) facing right. At each cell it may move forward
// or turn right once then move forward. We simulate the walk and check
// whether every non-obstacle cell is visited exactly once.
func solveEmbedded(input string) string {
	var n, m, k int
	r := strings.NewReader(input)
	fmt.Fscan(r, &n, &m, &k)

	obs := make(map[[2]int]bool)
	for i := 0; i < k; i++ {
		var x, y int
		fmt.Fscan(r, &x, &y)
		obs[[2]int{x, y}] = true
	}

	if obs[[2]int{1, 1}] {
		return "No"
	}

	total := n*m - k

	// direction vectors: 0=right, 1=down, 2=left, 3=up
	dx := [4]int{0, 1, 0, -1}
	dy := [4]int{1, 0, -1, 0}

	visited := make(map[[2]int]bool)
	visited[[2]int{1, 1}] = true
	count := 1
	x, y := 1, 1
	dir := 0

	for count < total {
		// Try to move forward
		nx, ny := x+dx[dir], y+dy[dir]
		if nx >= 1 && nx <= n && ny >= 1 && ny <= m && !obs[[2]int{nx, ny}] && !visited[[2]int{nx, ny}] {
			x, y = nx, ny
			visited[[2]int{x, y}] = true
			count++
			continue
		}
		// Try to turn right and then move forward
		dir = (dir + 1) % 4
		nx, ny = x+dx[dir], y+dy[dir]
		if nx >= 1 && nx <= n && ny >= 1 && ny <= m && !obs[[2]int{nx, ny}] && !visited[[2]int{nx, ny}] {
			x, y = nx, ny
			visited[[2]int{x, y}] = true
			count++
			continue
		}
		// Cannot move
		break
	}

	if count == total {
		return "Yes"
	}
	return "No"
}

// ---------- verifier infrastructure ----------

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "cand*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func genCase(r *rand.Rand) string {
	n := r.Intn(10) + 1
	m := r.Intn(10) + 1
	maxK := n*m - 1
	if maxK < 0 {
		maxK = 0
	}
	k := r.Intn(maxK + 1)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
	used := make(map[[2]int]bool)
	used[[2]int{1, 1}] = true
	for i := 0; i < k; {
		x := r.Intn(n) + 1
		y := r.Intn(m) + 1
		if used[[2]int{x, y}] {
			continue
		}
		used[[2]int{x, y}] = true
		sb.WriteString(fmt.Sprintf("%d %d\n", x, y))
		i++
	}
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	cand, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input := genCase(rng)
		expect := solveEmbedded(input)
		got, err := run(cand, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
