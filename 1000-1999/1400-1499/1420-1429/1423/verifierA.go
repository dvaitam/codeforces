package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

// refSolve is the embedded reference solver from cf_latest_1423_A.go
func refSolve(input string) string {
	type Cost struct {
		city int
		cost int
	}

	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)
	buf := make([]byte, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)

	if !scanner.Scan() {
		return ""
	}

	n := 0
	for _, v := range scanner.Bytes() {
		n = n*10 + int(v-'0')
	}

	if n%2 != 0 {
		return "-1"
	}

	pref := make([][]int, n+1)
	rank := make([][]int, n+1)

	for i := 1; i <= n; i++ {
		pref[i] = make([]int, n-1)
		rank[i] = make([]int, n+1)

		costs := make([]Cost, n-1)
		idx := 0
		for j := 1; j <= n; j++ {
			if i == j {
				continue
			}
			scanner.Scan()
			c := 0
			for _, v := range scanner.Bytes() {
				c = c*10 + int(v-'0')
			}
			costs[idx] = Cost{city: j, cost: c}
			idx++
		}

		sort.Slice(costs, func(a, b int) bool {
			return costs[a].cost < costs[b].cost
		})

		for j := 0; j < n-1; j++ {
			pref[i][j] = costs[j].city
			rank[i][costs[j].city] = j
		}
	}

	head := make([]int, n+1)
	tail := make([]int, n+1)
	next_val := make([][]int, n+1)
	prev_val := make([][]int, n+1)
	match := make([]int, n+1)

	for i := 1; i <= n; i++ {
		head[i] = 0
		tail[i] = n - 2
		next_val[i] = make([]int, n-1)
		prev_val[i] = make([]int, n-1)
		for j := 0; j < n-1; j++ {
			next_val[i][j] = j + 1
			prev_val[i][j] = j - 1
		}
		next_val[i][n-2] = -1
	}

	remove := func(i, idx int) {
		p := prev_val[i][idx]
		nx := next_val[i][idx]
		if p != -1 {
			next_val[i][p] = nx
		} else {
			head[i] = nx
		}
		if nx != -1 {
			prev_val[i][nx] = p
		} else {
			tail[i] = p
		}
	}

	Q := make([]int, 0, n*n)
	for i := 1; i <= n; i++ {
		Q = append(Q, i)
	}

	processQ := func() bool {
		for len(Q) > 0 {
			x := Q[0]
			Q = Q[1:]

			if head[x] == -1 {
				return false
			}

			y := pref[x][head[x]]
			idx_y := rank[y][x]

			curr := next_val[y][idx_y]
			for curr != -1 {
				w := pref[y][curr]
				remove(w, rank[w][y])
				if head[w] == -1 {
					return false
				}
				if match[y] == w {
					match[y] = 0
					Q = append(Q, w)
				}
				curr = next_val[y][curr]
			}

			tail[y] = idx_y
			next_val[y][idx_y] = -1
			match[y] = x
		}
		return true
	}

	if !processQ() {
		return "-1"
	}

	in_path := make([]int, n+1)
	for i := 0; i <= n; i++ {
		in_path[i] = -1
	}

	for {
		found_cycle := false
		for i := 1; i <= n; i++ {
			if head[i] != -1 && head[i] != tail[i] {
				found_cycle = true

				seq := []int{}
				p := i
				for {
					if in_path[p] != -1 {
						cycle := seq[in_path[p]:]
						for _, u := range cycle {
							sec := pref[u][next_val[u][head[u]]]
							nxt := pref[sec][tail[sec]]

							remove(nxt, rank[nxt][sec])
							if head[nxt] == -1 {
								return "-1"
							}

							remove(sec, rank[sec][nxt])
							if head[sec] == -1 {
								return "-1"
							}

							match[sec] = pref[sec][tail[sec]]
							Q = append(Q, nxt)
						}
						for _, node := range seq {
							in_path[node] = -1
						}
						break
					}
					in_path[p] = len(seq)
					seq = append(seq, p)

					sec := pref[p][next_val[p][head[p]]]
					nxt := pref[sec][tail[sec]]
					p = nxt
				}
				break
			}
		}

		if !found_cycle {
			break
		}

		if !processQ() {
			return "-1"
		}
	}

	for i := 1; i <= n; i++ {
		if head[i] == -1 || head[i] != tail[i] {
			return "-1"
		}
	}

	ans := make([]int, n+1)
	for i := 1; i <= n; i++ {
		ans[i] = pref[i][head[i]]
	}

	var sb strings.Builder
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteString(" ")
		}
		sb.WriteString(fmt.Sprintf("%d", ans[i]))
	}
	return sb.String()
}

func runProgram(bin, input string) (string, error) {
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

func genCase(rng *rand.Rand) string {
	n := rng.Intn(4)*2 + 2 // even 2..8
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		first := true
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			if !first {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", rng.Intn(10)+1))
			first = false
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		expect := refSolve(input)
		got, err := runProgram(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected: %s\ngot: %s\ninput:\n%s", i+1, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
