package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// ---------- embedded solver from accepted solution ----------

type Pair struct {
	x, y int
}

func solveF(input string) string {
	scanner := strings.NewReader(input)
	var n, m int
	fmt.Fscan(scanner, &n, &m)

	adj := make([][]bool, n)
	adj_mask := make([]uint32, n)
	for i := 0; i < n; i++ {
		adj[i] = make([]bool, n)
	}

	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(scanner, &u, &v)
		u--
		v--
		adj[u][v] = true
		adj[v][u] = true
		adj_mask[u] |= 1 << v
		adj_mask[v] |= 1 << u
	}

	dp_path := make([][14][14]bool, 1<<n)
	for i := 0; i < n; i++ {
		dp_path[1<<i][i][i] = true
	}

	for mask := 1; mask < (1 << n); mask++ {
		for i := 0; i < n; i++ {
			if (mask & (1 << i)) == 0 {
				continue
			}
			for j := 0; j < n; j++ {
				if (mask & (1 << j)) == 0 {
					continue
				}
				if dp_path[mask][i][j] {
					for k := 0; k < n; k++ {
						if (mask&(1<<k)) == 0 && adj[j][k] {
							dp_path[mask|(1<<k)][i][k] = true
							dp_path[mask|(1<<k)][k][i] = true
						}
					}
				}
			}
		}
	}

	valid_pairs := make([][]Pair, 1<<n)
	for mask := 1; mask < (1 << n); mask++ {
		if bits.OnesCount(uint(mask)) <= 1 {
			continue
		}
		for i := 0; i < n; i++ {
			if (mask & (1 << i)) == 0 {
				continue
			}
			for j := i + 1; j < n; j++ {
				if (mask&(1<<j)) != 0 && dp_path[mask][i][j] {
					valid_pairs[mask] = append(valid_pairs[mask], Pair{i, j})
				}
			}
		}
	}

	dp := make([]int, 1<<n)
	for i := range dp {
		dp[i] = 1e9
	}
	dp[1] = 0

	parent := make([]int, 1<<n)
	parent_T := make([]int, 1<<n)

	for mask := 1; mask < (1 << n); mask++ {
		if dp[mask] == 1e9 {
			continue
		}
		comp := ((1 << n) - 1) ^ mask
		for T := comp; T > 0; T = (T - 1) & comp {
			can_attach := false
			if bits.OnesCount(uint(T)) == 1 {
				x := bits.TrailingZeros(uint(T))
				if bits.OnesCount(uint(adj_mask[x]&uint32(mask))) >= 2 {
					can_attach = true
				}
			} else {
				for _, p := range valid_pairs[T] {
					if (adj_mask[p.x]&uint32(mask)) != 0 && (adj_mask[p.y]&uint32(mask)) != 0 {
						can_attach = true
						break
					}
				}
			}

			if can_attach {
				if dp[mask]+1 < dp[mask|T] {
					dp[mask|T] = dp[mask] + 1
					parent[mask|T] = mask
					parent_T[mask|T] = T
				}
			}
		}
	}

	curr := (1 << n) - 1
	var edges [][2]int

	for curr != 1 {
		p := parent[curr]
		T := parent_T[curr]

		if bits.OnesCount(uint(T)) == 1 {
			x := bits.TrailingZeros(uint(T))
			count := 0
			for u := 0; u < n; u++ {
				if (p&(1<<u)) != 0 && adj[x][u] {
					edges = append(edges, [2]int{x, u})
					count++
					if count == 2 {
						break
					}
				}
			}
		} else {
			var best_p Pair
			for _, pair := range valid_pairs[T] {
				if (adj_mask[pair.x]&uint32(p)) != 0 && (adj_mask[pair.y]&uint32(p)) != 0 {
					best_p = pair
					break
				}
			}

			for u := 0; u < n; u++ {
				if (p&(1<<u)) != 0 && adj[best_p.x][u] {
					edges = append(edges, [2]int{best_p.x, u})
					break
				}
			}

			for w := 0; w < n; w++ {
				if (p&(1<<w)) != 0 && adj[best_p.y][w] {
					edges = append(edges, [2]int{best_p.y, w})
					break
				}
			}

			curr_v := best_p.y
			curr_T := T
			for curr_v != best_p.x {
				curr_T ^= (1 << curr_v)
				for prev := 0; prev < n; prev++ {
					if (curr_T&(1<<prev)) != 0 && adj[curr_v][prev] && dp_path[curr_T][best_p.x][prev] {
						edges = append(edges, [2]int{curr_v, prev})
						curr_v = prev
						break
					}
				}
			}
		}
		curr = p
	}

	var sb strings.Builder
	fmt.Fprintln(&sb, len(edges))
	for _, edge := range edges {
		fmt.Fprintf(&sb, "%d %d\n", edge[0]+1, edge[1]+1)
	}
	return sb.String()
}

// ---------- verifier logic ----------

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func genTests() []string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []string
	for len(tests) < 100 {
		n := rng.Intn(4) + 2
		m := rng.Intn(4) + n - 1
		edgeSet := make(map[[2]int]bool)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		for len(edgeSet) < m {
			u := rng.Intn(n) + 1
			v := rng.Intn(n) + 1
			if u == v {
				continue
			}
			if u > v {
				u, v = v, u
			}
			key := [2]int{u, v}
			if edgeSet[key] {
				continue
			}
			edgeSet[key] = true
			fmt.Fprintf(&sb, "%d %d\n", u, v)
		}
		tests = append(tests, sb.String())
	}
	return tests
}

// parseEdgeCount extracts just the edge count from solver output
func parseEdgeCount(s string) (int, bool) {
	fields := strings.Fields(s)
	if len(fields) == 0 {
		return 0, false
	}
	v, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, false
	}
	return v, true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	tests := genTests()
	for i, input := range tests {
		expStr := solveF(input)
		gotStr, err := runExe(bin, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		expCount, ok1 := parseEdgeCount(expStr)
		gotCount, ok2 := parseEdgeCount(gotStr)
		if !ok1 || !ok2 {
			fmt.Printf("Test %d failed: could not parse output\nInput:\n%sExpected:\n%sGot:\n%s\n", i+1, input, expStr, gotStr)
			os.Exit(1)
		}
		if expCount != gotCount {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s\n", i+1, input, expStr, gotStr)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
