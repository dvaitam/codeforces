package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

var parent []int
var rankArr []int

func find(u int) int {
	if parent[u] != u {
		parent[u] = find(parent[u])
	}
	return parent[u]
}

func union(u, v int) {
	ru, rv := find(u), find(v)
	if ru == rv {
		return
	}
	if rankArr[ru] < rankArr[rv] {
		parent[ru] = rv
	} else if rankArr[ru] > rankArr[rv] {
		parent[rv] = ru
	} else {
		parent[rv] = ru
		rankArr[ru]++
	}
}

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveD(r *bufio.Reader) string {
	var n int
	if _, err := fmt.Fscan(r, &n); err != nil {
		return ""
	}
	vals := make([]int, n)
	maxA := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(r, &vals[i])
		if vals[i] > maxA {
			maxA = vals[i]
		}
	}
	mv := make([]int, maxA+1)
	for i := 0; i <= maxA; i++ {
		mv[i] = -1
	}
	for i, v := range vals {
		mv[v] = i
	}
	parent = make([]int, n)
	rankArr = make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
	}
	uMax := int(math.Sqrt(float64(maxA)))
	for u := 2; u <= uMax; u++ {
		vStart := 1
		if u%2 == 1 {
			vStart = 2
		}
		for v := vStart; v < u; v += 2 {
			if gcd(u, v) != 1 {
				continue
			}
			uu, vv := u*u, v*v
			x := uu - vv
			y := 2 * u * v
			z := uu + vv
			if x <= maxA && y <= maxA {
				ix, iy := mv[x], mv[y]
				if ix >= 0 && iy >= 0 {
					union(ix, iy)
				}
			}
			if x <= maxA && z <= maxA {
				ix, iz := mv[x], mv[z]
				if ix >= 0 && iz >= 0 {
					union(ix, iz)
				}
			}
			if y <= maxA && z <= maxA {
				iy, iz := mv[y], mv[z]
				if iy >= 0 && iz >= 0 {
					union(iy, iz)
				}
			}
		}
	}
	seen := make(map[int]bool)
	comps := 0
	for i := 0; i < n; i++ {
		r := find(i)
		if !seen[r] {
			seen[r] = true
			comps++
		}
	}
	return fmt.Sprintf("%d", comps)
}

func generateCaseD(rng *rand.Rand) string {
	n := rng.Intn(10) + 1
	vals := make([]int, n)
	for i := 0; i < n; i++ {
		vals[i] = rng.Intn(40) + 1
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i, v := range vals {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseD(rng)
		expect := solveD(bufio.NewReader(strings.NewReader(tc)))
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
