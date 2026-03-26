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

// ── embedded solver (CF-accepted 309E) ──────────────────────────────

type Segment struct {
	id int
	l  int
	r  int
}

func solve(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var n int
	fmt.Fscan(reader, &n)

	segments := make([]Segment, n)
	for i := 0; i < n; i++ {
		segments[i].id = i
		fmt.Fscan(reader, &segments[i].l, &segments[i].r)
	}

	adj := make([][]int, n)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if maxInt(segments[i].l, segments[j].l) <= minInt(segments[i].r, segments[j].r) {
				adj[i] = append(adj[i], j)
				adj[j] = append(adj[j], i)
			}
		}
	}

	orderSegs := make([]Segment, n)
	copy(orderSegs, segments)
	sort.Slice(orderSegs, func(i, j int) bool {
		if orderSegs[i].r != orderSegs[j].r {
			return orderSegs[i].r < orderSegs[j].r
		}
		return orderSegs[i].l < orderSegs[j].l
	})

	order := make([]int, n)
	for i := 0; i < n; i++ {
		order[i] = orderSegs[i].id
	}

	ans := make([]int, n)
	for i := 0; i < n; i++ {
		ans[i] = i
	}

	if n == 1 {
		return "1\n"
	}

	limit := make([]int, n)
	placed := make([]bool, n)
	cnt := make([]int, n)
	curAns := make([]int, n)
	freq := make([]int, n)

	check := func(M int) bool {
		for i := 0; i < n; i++ {
			limit[i] = n - 1
			placed[i] = false
		}
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				cnt[j] = 0
			}
			for j := 0; j < n; j++ {
				if !placed[j] {
					cnt[limit[j]]++
				}
			}
			for x := 1; x < n; x++ {
				cnt[x] += cnt[x-1]
			}
			if i > 0 && cnt[i-1] > 0 {
				return false
			}
			reqX := n - 1
			for x := i; x < n; x++ {
				allowed := x - i + 1 - cnt[x]
				if allowed == 0 {
					if x < reqX {
						reqX = x
					}
				} else if allowed < 0 {
					return false
				}
			}
			bestJ := -1
			for _, j := range order {
				if placed[j] {
					continue
				}
				if limit[j] > reqX {
					continue
				}
				possible := true
				if i+M < n {
					for x := 0; x < n; x++ {
						freq[x] = 0
					}
					for _, k := range adj[j] {
						if !placed[k] {
							freq[limit[k]]++
						}
					}
					c := 0
					for x := n - 1; x >= i+M; x-- {
						allowed := x - i + 1 - cnt[x]
						if limit[j] <= x {
							if c > allowed {
								possible = false
								break
							}
						} else {
							if c > allowed-1 {
								possible = false
								break
							}
						}
						c += freq[x]
					}
				}
				if possible {
					bestJ = j
					break
				}
			}
			if bestJ == -1 {
				return false
			}
			placed[bestJ] = true
			curAns[i] = bestJ
			for _, k := range adj[bestJ] {
				if !placed[k] && limit[k] > i+M {
					limit[k] = i + M
				}
			}
		}
		return true
	}

	low := 1
	high := n - 1
	for low <= high {
		mid := (low + high) / 2
		if check(mid) {
			copy(ans, curAns)
			high = mid - 1
		} else {
			low = mid + 1
		}
	}

	var sb strings.Builder
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(ans[i] + 1))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ── test generation ─────────────────────────────────────────────────

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(8) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		l := rng.Intn(1000)
		r := l + rng.Intn(1000)
		if i > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(fmt.Sprintf("%d %d", l, r))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCandidate(exe, input string) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		in := generateCase(rng)
		exp := strings.TrimSpace(solve(in))
		got, err := runCandidate(exe, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d: expected %q got %q\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
