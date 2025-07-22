package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type logEntry struct {
	a, b string
	t    int
}

func solveA(n, d int, logs []logEntry) string {
	nameToID := make(map[string]int)
	var idToName []string
	pairTimes := make(map[int64][]int)
	outSet := make(map[int64]struct{})
	key := func(u, v int) int64 { return (int64(u) << 32) | int64(uint32(v)) }
	for _, lg := range logs {
		p1, ok := nameToID[lg.a]
		if !ok {
			p1 = len(idToName)
			nameToID[lg.a] = p1
			idToName = append(idToName, lg.a)
		}
		p2, ok := nameToID[lg.b]
		if !ok {
			p2 = len(idToName)
			nameToID[lg.b] = p2
			idToName = append(idToName, lg.b)
		}
		k12 := key(p1, p2)
		times := pairTimes[k12]
		idx := sort.Search(len(times), func(i int) bool { return times[i] >= lg.t })
		if idx == len(times) {
			times = append(times, lg.t)
		} else if times[idx] != lg.t {
			times = append(times, 0)
			copy(times[idx+1:], times[idx:])
			times[idx] = lg.t
		}
		pairTimes[k12] = times
		k21 := key(p2, p1)
		rev := pairTimes[k21]
		idx2 := sort.Search(len(rev), func(i int) bool { return rev[i] >= lg.t-d })
		if idx2 < len(rev) && rev[idx2] != lg.t {
			u, v := p1, p2
			if u > v {
				u, v = v, u
			}
			outSet[key(u, v)] = struct{}{}
		}
	}
	type pair struct{ u, v int }
	pairs := make([]pair, 0, len(outSet))
	for kk := range outSet {
		pairs = append(pairs, pair{u: int(kk >> 32), v: int(uint32(kk))})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].u != pairs[j].u {
			return pairs[i].u < pairs[j].u
		}
		return pairs[i].v < pairs[j].v
	})
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(pairs)))
	for _, p := range pairs {
		sb.WriteString(fmt.Sprintf("%s %s\n", idToName[p.u], idToName[p.v]))
	}
	return strings.TrimSpace(sb.String())
}

func run(bin, input string) (string, error) {
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

func randName(rng *rand.Rand) string {
	length := rng.Intn(3) + 1
	letters := []rune("abcde")
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rng.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(8) + 2
		d := rng.Intn(10) + 1
		logs := make([]logEntry, n)
		t := 0
		for j := 0; j < n; j++ {
			a := randName(rng)
			b := randName(rng)
			for b == a {
				b = randName(rng)
			}
			t += rng.Intn(3)
			logs[j] = logEntry{a: a, b: b, t: t}
			t += rng.Intn(3)
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, d))
		for _, lg := range logs {
			sb.WriteString(fmt.Sprintf("%s %s %d\n", lg.a, lg.b, lg.t))
		}
		input := sb.String()
		expected := solveA(n, d, logs)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
