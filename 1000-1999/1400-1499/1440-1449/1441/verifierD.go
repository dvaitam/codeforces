package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type edge struct{ u, v int }
type testCase struct {
	n      int
	colors []int
	edges  []edge
}

func (tc testCase) Input() string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(strconv.Itoa(tc.n) + "\n")
	for i, c := range tc.colors {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(c))
	}
	sb.WriteByte('\n')
	for _, e := range tc.edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
	}
	return sb.String()
}

// Embedded correct solver for 1441D
func solve1441D(input string) string {
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	scanInt := func() int {
		scanner.Scan()
		res, _ := strconv.Atoi(scanner.Text())
		return res
	}

	var out strings.Builder
	w := bufio.NewWriter(&out)

	t := scanInt()
	for i := 0; i < t; i++ {
		n := scanInt()
		color := make([]int, n+1)
		for j := 1; j <= n; j++ {
			color[j] = scanInt()
		}
		adj := make([][]int, n+1)
		for j := 0; j < n-1; j++ {
			u := scanInt()
			v := scanInt()
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}

		ans := 0
		const negInf = -1000000000

		var dfs func(u, p int) (int, int)
		dfs = func(u, p int) (int, int) {
			dp0, dp1 := negInf, negInf
			if color[u] == 1 {
				dp0 = 0
			} else if color[u] == 2 {
				dp1 = 0
			}

			for _, v := range adj[u] {
				if v == p {
					continue
				}
				v0, v1 := dfs(v, u)

				if color[u] == 1 {
					tmpV1 := v1
					if tmpV1 != negInf {
						tmpV1++
					}
					bestV := v0
					if tmpV1 > bestV {
						bestV = tmpV1
					}
					if dp0+bestV > ans {
						ans = dp0 + bestV
					}
					if bestV > dp0 {
						dp0 = bestV
					}
				} else if color[u] == 2 {
					tmpV0 := v0
					if tmpV0 != negInf {
						tmpV0++
					}
					bestV := v1
					if tmpV0 > bestV {
						bestV = tmpV0
					}
					if dp1+bestV > ans {
						ans = dp1 + bestV
					}
					if bestV > dp1 {
						dp1 = bestV
					}
				} else {
					if dp0 != negInf && v1 != negInf {
						if dp0+v1+1 > ans {
							ans = dp0 + v1 + 1
						}
					}
					if dp1 != negInf && v0 != negInf {
						if dp1+v0+1 > ans {
							ans = dp1 + v0 + 1
						}
					}
					if dp0 != negInf && v0 != negInf {
						if dp0+v0 > ans {
							ans = dp0 + v0
						}
					}
					if dp1 != negInf && v1 != negInf {
						if dp1+v1 > ans {
							ans = dp1 + v1
						}
					}
					if v0 > dp0 {
						dp0 = v0
					}
					if v1 > dp1 {
						dp1 = v1
					}
				}
			}
			return dp0, dp1
		}

		dfs(1, 0)
		fmt.Fprintln(w, (ans+3)/2)
	}
	w.Flush()
	return strings.TrimSpace(out.String())
}

func genCase(rng *rand.Rand) testCase {
	n := rng.Intn(6) + 1
	colors := make([]int, n)
	for i := 0; i < n; i++ {
		colors[i] = rng.Intn(3)
	}
	edges := make([]edge, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges = append(edges, edge{p, i})
	}
	return testCase{n, colors, edges}
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out, stderr strings.Builder
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
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{{n: 1, colors: []int{0}, edges: []edge{}}}
	for i := 0; i < 100; i++ {
		cases = append(cases, genCase(rng))
	}
	for i, tc := range cases {
		inp := tc.Input()
		exp := solve1441D(inp)
		got, err := runExe(bin, inp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, inp)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, inp)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
