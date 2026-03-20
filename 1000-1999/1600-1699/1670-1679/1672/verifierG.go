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

const vMOD int = 998244353

func vPow2(n int) []int {
	res := make([]int, n+1)
	res[0] = 1
	for i := 1; i <= n; i++ {
		res[i] = res[i-1] * 2 % vMOD
	}
	return res
}

type vDSU struct {
	parent []int
}

func newVDSU(n int) *vDSU {
	d := &vDSU{parent: make([]int, n)}
	for i := range d.parent {
		d.parent[i] = i
	}
	return d
}

func (d *vDSU) Find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.Find(d.parent[x])
	}
	return d.parent[x]
}

func (d *vDSU) Union(a, b int) {
	ra, rb := d.Find(a), d.Find(b)
	if ra != rb {
		d.parent[rb] = ra
	}
}

func solve(input string) string {
	rd := strings.NewReader(input)
	var r, c int
	fmt.Fscan(rd, &r, &c)
	grid := make([]string, r)
	for i := 0; i < r; i++ {
		fmt.Fscan(rd, &grid[i])
	}

	k := 0
	rowFixed := make([]int, r)
	colFixed := make([]int, c)
	rowQ := make([]int, r)
	colQ := make([]int, c)
	var edges [][2]int
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			ch := grid[i][j]
			if ch == '?' {
				k++
				rowQ[i]++
				colQ[j]++
				edges = append(edges, [2]int{i, j})
			} else if ch == '1' {
				rowFixed[i] ^= 1
				colFixed[j] ^= 1
			}
		}
	}

	pow := vPow2(k)

	if r%2 == 0 && c%2 == 0 {
		return fmt.Sprintf("%d", pow[k])
	}

	if r%2 == 0 && c%2 == 1 {
		ans := 0
		for p := 0; p < 2; p++ {
			cur := 1
			for i := 0; i < r; i++ {
				if rowQ[i] == 0 {
					if rowFixed[i] != p {
						cur = 0
						break
					}
				} else {
					cur = cur * pow[rowQ[i]-1] % vMOD
				}
			}
			ans = (ans + cur) % vMOD
		}
		return fmt.Sprintf("%d", ans)
	}

	if r%2 == 1 && c%2 == 0 {
		ans := 0
		for p := 0; p < 2; p++ {
			cur := 1
			for j := 0; j < c; j++ {
				if colQ[j] == 0 {
					if colFixed[j] != p {
						cur = 0
						break
					}
				} else {
					cur = cur * pow[colQ[j]-1] % vMOD
				}
			}
			ans = (ans + cur) % vMOD
		}
		return fmt.Sprintf("%d", ans)
	}

	// r and c are both odd
	// Use the same approach as the accepted solution:
	// DSU on bipartite graph (rows + columns), edges from '?' cells
	dsu := newVDSU(r + c)
	numEdges := 0
	for _, e := range edges {
		i, j := e[0], e[1]
		dsu.Union(i, r+j)
		numEdges++
	}
	comps := 0
	for i := 0; i < r+c; i++ {
		if dsu.Find(i) == i {
			comps++
		}
	}

	ans := 0
	for S := 0; S <= 1; S++ {
		req := make([]int, r+c)
		for i := 0; i < r; i++ {
			sum := rowFixed[i] // already computed as XOR of 1-bits
			req[i] = (S - sum%2 + 2) % 2
		}
		for j := 0; j < c; j++ {
			sum := colFixed[j]
			req[r+j] = (S - sum%2 + 2) % 2
		}

		compSum := make([]int, r+c)
		for i := 0; i < r+c; i++ {
			compSum[dsu.Find(i)] += req[i]
		}

		valid := true
		for i := 0; i < r+c; i++ {
			if dsu.Find(i) == i && compSum[i]%2 != 0 {
				valid = false
				break
			}
		}

		if valid {
			exp := numEdges - (r + c - comps)
			if exp >= 0 {
				ans = (ans + pow[exp]) % vMOD
			}
		}
	}
	return fmt.Sprintf("%d", ans)
}

func runExe(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
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

func genCase(rng *rand.Rand) string {
	r := rng.Intn(3) + 1
	c := rng.Intn(3) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", r, c))
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			v := rng.Intn(3)
			if v == 0 {
				sb.WriteByte('0')
			} else if v == 1 {
				sb.WriteByte('1')
			} else {
				sb.WriteByte('?')
			}
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		want := solve(input)
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
