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

func run(bin, input string) (string, error) {
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

type participant struct {
	name  string
	score int
}

func expectedTable(ps []participant) []string {
	sort.Slice(ps, func(i, j int) bool {
		if ps[i].score != ps[j].score {
			return ps[i].score > ps[j].score
		}
		return ps[i].name < ps[j].name
	})
	var res []string
	n := len(ps)
	i := 0
	for i < n {
		j := i + 1
		for j < n && ps[j].score == ps[i].score {
			j++
		}
		place := ""
		if j-i == 1 {
			place = fmt.Sprintf("%d", i+1)
		} else {
			place = fmt.Sprintf("%d-%d", i+1, j)
		}
		for k := i; k < j; k++ {
			res = append(res, fmt.Sprintf("%s %s", place, ps[k].name))
		}
		i = j
	}
	return res
}

func randName(rng *rand.Rand) string {
	l := rng.Intn(8) + 3
	b := make([]byte, l)
	for i := range b {
		b[i] = byte('a' + rng.Intn(26))
	}
	return string(b)
}

func generateCase(rng *rand.Rand) (string, []string) {
	n := rng.Intn(10) + 1
	ps := make([]participant, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	used := make(map[string]bool)
	for i := 0; i < n; i++ {
		name := randName(rng)
		for used[name] {
			name = randName(rng)
		}
		used[name] = true
		score := rng.Intn(1001)
		ps[i] = participant{name, score}
		sb.WriteString(fmt.Sprintf("%s %d\n", name, score))
	}
	exp := expectedTable(ps)
	return sb.String(), exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expectedLines := generateCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		outLines := strings.Split(strings.TrimSpace(out), "\n")
		if len(outLines) != len(expectedLines) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d lines got %d\ninput:\n%s", i+1, len(expectedLines), len(outLines), input)
			os.Exit(1)
		}
		for j := range expectedLines {
			if strings.TrimSpace(outLines[j]) != expectedLines[j] {
				fmt.Fprintf(os.Stderr, "case %d failed: line %d expected '%s' got '%s'\ninput:\n%s", i+1, j+1, expectedLines[j], outLines[j], input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
