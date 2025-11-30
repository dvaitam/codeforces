package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

// Embedded copy of testcasesC.txt so the verifier is self-contained.
const testcasesRaw = `20
A B 5:5
A C 2:2
BERLAND C 4:1
B C 4:0
BERLAND B 4:5
BERLAND B 4:3
BERLAND C 4:2
A B 0:0
B C 2:3
A C 2:3
A C 1:1
BERLAND A 0:1
A B 2:1
BERLAND C 1:4
B C 4:2
BERLAND A 5:4
A B 2:4
A C 2:2
BERLAND B 3:1
B C 3:5
BERLAND A 3:4
A B 4:2
B C 5:3
BERLAND B 3:2
A C 4:5
BERLAND B 2:5
BERLAND A 1:4
A B 2:3
B C 2:2
A C 5:4
A B 3:4
BERLAND A 2:5
BERLAND B 4:0
A C 2:5
B C 0:1
BERLAND B 2:4
A C 1:5
A B 0:4
B C 1:2
BERLAND A 1:1
A B 2:1
BERLAND C 1:5
B C 0:0
BERLAND B 0:0
A C 0:0
A C 1:5
B C 1:4
BERLAND B 5:0
A B 3:4
BERLAND A 0:1
A B 2:2
A C 3:0
BERLAND C 2:3
B C 4:4
BERLAND B 5:0
A B 0:5
BERLAND A 5:2
B C 0:0
BERLAND C 3:1
A C 4:4
B C 2:2
BERLAND C 4:3
BERLAND A 5:0
A B 5:4
A C 1:5
BERLAND C 1:0
B C 3:5
A C 1:4
BERLAND B 5:0
A B 1:1
A B 4:1
B C 4:4
BERLAND B 5:2
BERLAND A 2:5
A C 3:2
A B 3:1
A C 0:4
B C 5:0
BERLAND B 1:0
BERLAND A 0:0
B C 4:5
A B 3:3
A C 2:4
BERLAND B 5:3
BERLAND C 1:5
BERLAND C 4:4
BERLAND B 0:3
A B 4:4
B C 1:0
A C 5:3
A C 2:5
BERLAND C 2:2
A B 0:5
B C 3:0
BERLAND A 0:2
A B 5:3
BERLAND C 3:1
B C 4:4
A C 0:0
BERLAND B 2:0`

type team struct {
	name     string
	points   int
	scored   int
	conceded int
}

func parseMatch(line string) (string, string, int, int, error) {
	fields := strings.Fields(strings.TrimSpace(line))
	if len(fields) != 3 {
		return "", "", 0, 0, fmt.Errorf("bad line %q", line)
	}
	sides := strings.Split(fields[2], ":")
	if len(sides) != 2 {
		return "", "", 0, 0, fmt.Errorf("bad score %q", fields[2])
	}
	left, err := strconv.Atoi(sides[0])
	if err != nil {
		return "", "", 0, 0, fmt.Errorf("parse score: %w", err)
	}
	right, err := strconv.Atoi(sides[1])
	if err != nil {
		return "", "", 0, 0, fmt.Errorf("parse score: %w", err)
	}
	return fields[0], fields[1], left, right, nil
}

func solveCase(lines []string) (string, error) {
	teams := make(map[string]*team)
	games := make(map[string]int)
	for _, line := range lines {
		t1, t2, g1, g2, err := parseMatch(line)
		if err != nil {
			return "", err
		}
		if teams[t1] == nil {
			teams[t1] = &team{name: t1}
		}
		if teams[t2] == nil {
			teams[t2] = &team{name: t2}
		}
		games[t1]++
		games[t2]++
		teams[t1].scored += g1
		teams[t1].conceded += g2
		teams[t2].scored += g2
		teams[t2].conceded += g1
		switch {
		case g1 > g2:
			teams[t1].points += 3
		case g1 < g2:
			teams[t2].points += 3
		default:
			teams[t1].points++
			teams[t2].points++
		}
	}
	opp := ""
	for name, cnt := range games {
		if name != "BERLAND" && cnt == 2 {
			opp = name
			break
		}
	}
	orig := make(map[string]team, len(teams))
	for name, t := range teams {
		orig[name] = *t
	}
	found := false
	bestX, bestY := 0, 0
	for diff := 1; diff <= 1000 && !found; diff++ {
		for y := 0; y <= 1000; y++ {
			x := y + diff
			snapshot := make([]team, 0, len(orig))
			for _, t := range orig {
				snapshot = append(snapshot, t)
			}
			for i := range snapshot {
				if snapshot[i].name == "BERLAND" {
					snapshot[i].points += 3
					snapshot[i].scored += x
					snapshot[i].conceded += y
				}
				if snapshot[i].name == opp {
					snapshot[i].scored += y
					snapshot[i].conceded += x
				}
			}
			sort.Slice(snapshot, func(i, j int) bool {
				a, b := snapshot[i], snapshot[j]
				if a.points != b.points {
					return a.points > b.points
				}
				da := a.scored - a.conceded
				db := b.scored - b.conceded
				if da != db {
					return da > db
				}
				if a.scored != b.scored {
					return a.scored > b.scored
				}
				return a.name < b.name
			})
			rank := -1
			for i, t := range snapshot {
				if t.name == "BERLAND" {
					rank = i
					break
				}
			}
			if rank >= 0 && rank < 2 {
				bestX, bestY = x, y
				found = true
				break
			}
		}
	}
	if !found {
		return "IMPOSSIBLE", nil
	}
	return fmt.Sprintf("%d:%d", bestX, bestY), nil
}

func loadTestcases() ([][]string, error) {
	scan := bufio.NewScanner(strings.NewReader(testcasesRaw))
	if !scan.Scan() {
		return nil, fmt.Errorf("empty testcase data")
	}
	t, err := strconv.Atoi(strings.TrimSpace(scan.Text()))
	if err != nil {
		return nil, fmt.Errorf("parse t: %w", err)
	}
	cases := make([][]string, 0, t)
	for i := 0; i < t; i++ {
		lines := make([]string, 0, 5)
		for j := 0; j < 5; j++ {
			if !scan.Scan() {
				return nil, fmt.Errorf("test %d: missing line %d", i+1, j+1)
			}
			lines = append(lines, scan.Text())
		}
		cases = append(cases, lines)
	}
	return cases, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	testcases, err := loadTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	bin := os.Args[1]
	for idx, lines := range testcases {
		expect, err := solveCase(lines)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: failed to solve: %v\n", idx+1, err)
			os.Exit(1)
		}
		input := strings.Join(lines, "\n") + "\n"
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx+1, err, stderr.String())
			os.Exit(1)
		}
		fields := strings.Fields(stdout.String())
		if len(fields) == 0 {
			fmt.Printf("test %d: no output\n", idx+1)
			os.Exit(1)
		}
		if fields[0] != expect {
			fmt.Printf("test %d failed: expected %s got %s\n", idx+1, expect, fields[0])
			os.Exit(1)
		}
		if len(fields) > 1 {
			fmt.Printf("test %d: extra output detected\n", idx+1)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
