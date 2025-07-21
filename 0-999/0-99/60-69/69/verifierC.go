package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type comp struct {
	name string
	req  []int
}

type purchase struct {
	hero int
	item int
}

type testCaseC struct {
	k, n, m, q int
	basic      []string
	comps      []comp
	ops        []purchase
	expected   string
}

func solveCase(tc testCaseC) string {
	k, n, m := tc.k, tc.n, tc.m
	basicNames := tc.basic
	compNames := make([]string, m)
	compReq := make([][]int, m)
	for i, c := range tc.comps {
		compNames[i] = c.name
		compReq[i] = c.req
	}
	basicCnt := make([][]int, k)
	compCnt := make([][]int, k)
	for i := 0; i < k; i++ {
		basicCnt[i] = make([]int, n)
		compCnt[i] = make([]int, m)
	}
	for _, op := range tc.ops {
		h := op.hero
		bi := op.item
		basicCnt[h][bi]++
		for j := 0; j < m; j++ {
			can := true
			for idx, need := range compReq[j] {
				if need > 0 && basicCnt[h][idx] < need {
					can = false
					break
				}
			}
			if can {
				for idx, need := range compReq[j] {
					if need > 0 {
						basicCnt[h][idx] -= need
					}
				}
				compCnt[h][j]++
				break
			}
		}
	}
	var out strings.Builder
	type pair struct {
		name string
		cnt  int
	}
	for i := 0; i < k; i++ {
		var lst []pair
		for bi := 0; bi < n; bi++ {
			if basicCnt[i][bi] > 0 {
				lst = append(lst, pair{basicNames[bi], basicCnt[i][bi]})
			}
		}
		for cj := 0; cj < m; cj++ {
			if compCnt[i][cj] > 0 {
				lst = append(lst, pair{compNames[cj], compCnt[i][cj]})
			}
		}
		sort.Slice(lst, func(a, b int) bool { return lst[a].name < lst[b].name })
		fmt.Fprintln(&out, len(lst))
		for _, p := range lst {
			fmt.Fprintf(&out, "%s %d\n", p.name, p.cnt)
		}
	}
	return strings.TrimSpace(out.String())
}

func generateCase(rng *rand.Rand) testCaseC {
	for {
		k := rng.Intn(3) + 1
		n := rng.Intn(3) + 1
		m := rng.Intn(3)
		q := rng.Intn(10) + 1
		basic := make([]string, n)
		for i := 0; i < n; i++ {
			basic[i] = fmt.Sprintf("b%d", i)
		}
		comps := make([]comp, m)
		for i := 0; i < m; i++ {
			req := make([]int, n)
			cnt := rng.Intn(n) + 1
			used := make(map[int]bool)
			for j := 0; j < cnt; j++ {
				idx := rng.Intn(n)
				for used[idx] {
					idx = rng.Intn(n)
				}
				used[idx] = true
				req[idx] = rng.Intn(2) + 1
			}
			comps[i] = comp{name: fmt.Sprintf("c%d", i), req: req}
		}
		ops := make([]purchase, q)
		valid := true
		basicCnt := make([][]int, k)
		for i := 0; i < k; i++ {
			basicCnt[i] = make([]int, n)
		}
		for step := 0; step < q && valid; step++ {
			hero := rng.Intn(k)
			item := rng.Intn(n)
			ops[step] = purchase{hero, item}
			basicCnt[hero][item]++
			craftable := 0
			for j := 0; j < m; j++ {
				can := true
				for idx, need := range comps[j].req {
					if need > 0 && basicCnt[hero][idx] < need {
						can = false
						break
					}
				}
				if can {
					craftable++
				}
			}
			if craftable > 1 {
				valid = false
				break
			}
			if craftable == 1 {
				for j := 0; j < m; j++ {
					can := true
					for idx, need := range comps[j].req {
						if need > 0 && basicCnt[hero][idx] < need {
							can = false
							break
						}
					}
					if can {
						for idx, need := range comps[j].req {
							if need > 0 {
								basicCnt[hero][idx] -= need
							}
						}
						break
					}
				}
			}
		}
		if !valid {
			continue
		}
		exp := solveCase(testCaseC{k: k, n: n, m: m, q: q, basic: basic, comps: comps, ops: ops})
		return testCaseC{k: k, n: n, m: m, q: q, basic: basic, comps: comps, ops: ops, expected: exp}
	}
}

func run(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(3))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d %d\n", tc.k, tc.n, tc.m, tc.q)
		for _, name := range tc.basic {
			fmt.Fprintf(&sb, "%s\n", name)
		}
		for _, cp := range tc.comps {
			var parts []string
			for idx, need := range cp.req {
				if need > 0 {
					parts = append(parts, fmt.Sprintf("%s %d", tc.basic[idx], need))
				}
			}
			fmt.Fprintf(&sb, "%s: %s\n", cp.name, strings.Join(parts, ", "))
		}
		for _, op := range tc.ops {
			fmt.Fprintf(&sb, "%d %s\n", op.hero+1, tc.basic[op.item])
		}
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected:\n%s\ngot:\n%s\n", i+1, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
