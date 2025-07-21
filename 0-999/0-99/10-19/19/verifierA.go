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

type team struct {
	name           string
	points, gd, gs int
}

func runCandidate(bin, input string) (string, error) {
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

func solveCase(n int, names []string, matches []match) string {
	teams := make([]*team, n)
	idx := make(map[string]int, n)
	for i, name := range names {
		teams[i] = &team{name: name}
		idx[name] = i
	}
	for _, m := range matches {
		t1 := teams[idx[m.name1]]
		t2 := teams[idx[m.name2]]
		t1.gs += m.s1
		t1.gd += m.s1 - m.s2
		t2.gs += m.s2
		t2.gd += m.s2 - m.s1
		if m.s1 > m.s2 {
			t1.points += 3
		} else if m.s1 < m.s2 {
			t2.points += 3
		} else {
			t1.points++
			t2.points++
		}
	}
	sort.Slice(teams, func(i, j int) bool {
		a, b := teams[i], teams[j]
		if a.points != b.points {
			return a.points > b.points
		}
		if a.gd != b.gd {
			return a.gd > b.gd
		}
		if a.gs != b.gs {
			return a.gs > b.gs
		}
		return false
	})
	m := n / 2
	selected := make([]string, m)
	for i := 0; i < m; i++ {
		selected[i] = teams[i].name
	}
	sort.Strings(selected)
	return strings.Join(selected, "\n")
}

type match struct {
	name1, name2 string
	s1, s2       int
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6)*2 + 2 // even between 2 and 12
	names := make([]string, n)
	for i := 0; i < n; i++ {
		names[i] = fmt.Sprintf("team%d", i+1)
	}
	total := n * (n - 1) / 2
	matches := make([]match, 0, total)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			s1 := rng.Intn(6)
			s2 := rng.Intn(6)
			matches = append(matches, match{names[i], names[j], s1, s2})
		}
	}
	// Build input
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for _, name := range names {
		fmt.Fprintln(&sb, name)
	}
	for _, m := range matches {
		fmt.Fprintf(&sb, "%s-%s %d:%d\n", m.name1, m.name2, m.s1, m.s2)
	}
	expect := solveCase(n, names, matches)
	return sb.String(), expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
