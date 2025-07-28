package main

import (
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

func run(bin, input string) (string, error) {
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

type entry struct {
	sum int
	id  int
}

func bestAvg(n int, pairs [][2]int) (float64, map[int]int) {
	teamVals := make(map[int][]int)
	teams := make([]int, 0)
	for _, p := range pairs {
		m := p[0]
		k := p[1]
		if _, ok := teamVals[m]; !ok {
			teams = append(teams, m)
		}
		teamVals[m] = append(teamVals[m], k)
	}
	best := 0.0
	sums := make(map[int]int)
	for t := 1; t <= 20; t++ {
		entries := make([]entry, 0, len(teams))
		for _, id := range teams {
			vals := teamVals[id]
			s := 0
			for _, v := range vals {
				if v >= t {
					s += t
				} else {
					s += v
				}
			}
			if s > 0 {
				entries = append(entries, entry{s, id})
			}
		}
		if len(entries) < t {
			continue
		}
		sort.Slice(entries, func(i, j int) bool { return entries[i].sum > entries[j].sum })
		total := 0
		ids := make([]int, t)
		for i := 0; i < t; i++ {
			total += entries[i].sum
			ids[i] = entries[i].id
		}
		avg := float64(total) / float64(t)
		if avg > best {
			best = avg
			sums = make(map[int]int)
			for i := 0; i < t; i++ {
				sums[ids[i]] = 1
			}
		}
	}
	return best, sums
}

func computeAvg(t int, selected []int, pairs [][2]int) float64 {
	idSet := make(map[int]bool)
	for _, id := range selected {
		idSet[id] = true
	}
	total := 0
	for _, p := range pairs {
		if !idSet[p[0]] {
			continue
		}
		if t <= p[1] {
			total += t
		} else {
			total += p[1]
		}
	}
	return float64(total) / float64(t)
}

func parseOutput(output string) (int, []int, error) {
	parts := strings.Fields(output)
	if len(parts) == 0 {
		return 0, nil, fmt.Errorf("empty output")
	}
	t, err := strconv.Atoi(parts[0])
	if err != nil || t < 0 {
		return 0, nil, fmt.Errorf("invalid t")
	}
	if len(parts)-1 != t {
		return 0, nil, fmt.Errorf("expected %d ids, got %d", t, len(parts)-1)
	}
	ids := make([]int, t)
	for i := 0; i < t; i++ {
		v, err := strconv.Atoi(parts[i+1])
		if err != nil {
			return 0, nil, fmt.Errorf("invalid id %q", parts[i+1])
		}
		ids[i] = v
	}
	return t, ids, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		n := rand.Intn(8) + 2
		pairs := make([][2]int, n)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for j := 0; j < n; j++ {
			m := rand.Intn(8) + 1
			k := rand.Intn(5) + 1
			pairs[j] = [2]int{m, k}
			fmt.Fprintf(&sb, "%d %d\n", m, k)
		}
		input := sb.String()
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		tOut, ids, err := parseOutput(out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\noutput:%s\n", i+1, err, out)
			os.Exit(1)
		}
		best, _ := bestAvg(n, pairs)
		avg := computeAvg(tOut, ids, pairs)
		if absFloat(best-avg) > 1e-6 {
			fmt.Fprintf(os.Stderr, "case %d failed: suboptimal average\n", i+1)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func absFloat(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
