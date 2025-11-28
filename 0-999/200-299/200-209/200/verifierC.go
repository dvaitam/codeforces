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

const testcasesC = `20
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
BERLAND B 2:0
`

type Team struct {
	name     string
	points   int
	scored   int
	conceded int
}

func compute(lines []string) string {
	teams := make(map[string]*Team)
	games := make(map[string]int)
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) != 3 {
			continue
		}
		t1, t2 := parts[0], parts[1]
		sc := strings.Split(parts[2], ":")
		g1, _ := strconv.Atoi(sc[0])
		g2, _ := strconv.Atoi(sc[1])
		if teams[t1] == nil {
			teams[t1] = &Team{name: t1}
		}
		if teams[t2] == nil {
			teams[t2] = &Team{name: t2}
		}
		games[t1]++
		games[t2]++
		teams[t1].scored += g1
		teams[t1].conceded += g2
		teams[t2].scored += g2
		teams[t2].conceded += g1
		if g1 > g2 {
			teams[t1].points += 3
		} else if g1 < g2 {
			teams[t2].points += 3
		} else {
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
	orig := make(map[string]Team)
	for n, t := range teams {
		orig[n] = *t
	}
	limit := 1000
	for d := 1; d <= limit; d++ {
		for y := 0; y <= limit; y++ {
			x := y + d
			sim := make([]Team, 0, len(orig))
			for _, t := range orig {
				sim = append(sim, t)
			}
			for i := range sim {
				if sim[i].name == "BERLAND" {
					sim[i].points += 3
					sim[i].scored += x
					sim[i].conceded += y
				}
				if sim[i].name == opp {
					sim[i].scored += y
					sim[i].conceded += x
				}
			}
			sort.Slice(sim, func(i, j int) bool {
				a, b := sim[i], sim[j]
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
			for i, t := range sim {
				if t.name == "BERLAND" {
					rank = i
					break
				}
			}
			if rank >= 0 && rank < 2 {
				return fmt.Sprintf("%d:%d", x, y)
			}
		}
	}
	return "IMPOSSIBLE"
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	data := []byte(testcasesC)
	scanLines := bufio.NewScanner(bytes.NewReader(data))
	if !scanLines.Scan() {
		fmt.Println("bad file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(strings.TrimSpace(scanLines.Text()))
	cases := make([][]string, t)
	for i := 0; i < t; i++ {
		lines := make([]string, 5)
		for j := 0; j < 5; j++ {
			scanLines.Scan()
			lines[j] = scanLines.Text()
		}
		cases[i] = lines
	}
	expected := make([]string, t)
	for i, lns := range cases {
		expected[i] = compute(lns)
	}
	cmd := exec.Command(os.Args[1])
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("execution failed:", err)
		os.Exit(1)
	}
	outScan := bufio.NewScanner(bytes.NewReader(out))
	outScan.Split(bufio.ScanWords)
	for i := 0; i < t; i++ {
		if !outScan.Scan() {
			fmt.Printf("missing output for test %d\n", i+1)
			os.Exit(1)
		}
		got := outScan.Text()
		if got != expected[i] {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, expected[i], got)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
