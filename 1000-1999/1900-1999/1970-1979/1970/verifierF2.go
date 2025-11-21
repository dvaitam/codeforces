package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const (
	goalNone = iota
	goalRed
	goalBlue
)

const refSource = "1000-1999/1900-1999/1970-1979/1970/1970F2.go"

type testCase struct {
	name  string
	input string
}

type player struct {
	id       string
	team     byte
	x, y     int
	carrying bool
	alive    bool
}

type ball struct {
	x, y    int
	carrier *player
}

type eventRecord struct {
	time  int
	kind  string
	label string
}

func (e eventRecord) String() string {
	if e.kind == "GOAL" {
		return fmt.Sprintf("%d %s GOAL", e.time, e.label)
	}
	return fmt.Sprintf("%d %s ELIMINATED", e.time, e.label)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF2.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()

	for idx, tc := range tests {
		expectedEvents, expRed, expBlue, err := simulate(tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to simulate test %d (%s): %v\n", idx+1, tc.name, err)
			os.Exit(1)
		}

		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}
		if err := checkOutput("reference", refOut, expectedEvents, expRed, expBlue, tc); err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		if err := checkOutput("candidate", candOut, expectedEvents, expRed, expBlue, tc); err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-1970F2-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref1970F2.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
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
	err := cmd.Run()
	return out.String(), err
}

func checkOutput(label, output string, expected []eventRecord, expRed, expBlue int, tc testCase) error {
	events, red, blue, err := parseOutput(output)
	if err != nil {
		return fmt.Errorf("%s output invalid: %v", label, err)
	}
	if len(events) != len(expected) {
		return fmt.Errorf("%s reported %d events but expected %d", label, len(events), len(expected))
	}
	for i := range expected {
		if events[i] != expected[i] {
			return fmt.Errorf("%s event %d mismatch: got %s expected %s", label, i+1, events[i].String(), expected[i].String())
		}
	}
	if red != expRed || blue != expBlue {
		return fmt.Errorf("%s final score mismatch: got %d %d expected %d %d", label, red, blue, expRed, expBlue)
	}
	return nil
}

func parseOutput(output string) ([]eventRecord, int, int, error) {
	scanner := bufio.NewScanner(strings.NewReader(output))
	var lines []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, 0, 0, fmt.Errorf("failed to read output: %v", err)
	}
	if len(lines) == 0 {
		return nil, 0, 0, fmt.Errorf("empty output")
	}
	finalLine := lines[len(lines)-1]
	eventLines := lines[:len(lines)-1]

	redScore, blueScore, err := parseFinalLine(finalLine)
	if err != nil {
		return nil, 0, 0, err
	}

	events := make([]eventRecord, 0, len(eventLines))
	for _, line := range eventLines {
		fields := strings.Fields(line)
		if len(fields) != 3 {
			return nil, 0, 0, fmt.Errorf("invalid event line %q", line)
		}
		t, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, 0, 0, fmt.Errorf("invalid time in %q", line)
		}
		switch fields[2] {
		case "GOAL":
			if fields[1] != "RED" && fields[1] != "BLUE" {
				return nil, 0, 0, fmt.Errorf("invalid team in %q", line)
			}
			events = append(events, eventRecord{time: t, kind: "GOAL", label: fields[1]})
		case "ELIMINATED":
			if len(fields[1]) != 2 || (fields[1][0] != 'R' && fields[1][0] != 'B') {
				return nil, 0, 0, fmt.Errorf("invalid player id in %q", line)
			}
			events = append(events, eventRecord{time: t, kind: "ELIMINATED", label: fields[1]})
		default:
			return nil, 0, 0, fmt.Errorf("invalid event type in %q", line)
		}
	}

	return events, redScore, blueScore, nil
}

func parseFinalLine(line string) (int, int, error) {
	fields := strings.Fields(line)
	if len(fields) != 4 || fields[0] != "FINAL" || fields[1] != "SCORE:" {
		return 0, 0, fmt.Errorf("invalid final line %q", line)
	}
	red, err := strconv.Atoi(fields[2])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid red score in %q", line)
	}
	blue, err := strconv.Atoi(fields[3])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid blue score in %q", line)
	}
	return red, blue, nil
}

func buildTests() []testCase {
	return []testCase{
		{
			name: "red_scores_own_goal",
			input: `3 3
RG .. BG
R0 .Q B0
.. .. ..
5
R0 R
R0 C .Q
R0 U
R0 L
R0 T
`,
		},
		{
			name: "double_elimination",
			input: `3 3
RG .. BG
R0 .Q B0
.. .B ..
3
R0 R
B0 L
.B U
`,
		},
		{
			name: "red_scores_correct_goal",
			input: `3 3
RG .. BG
R0 .Q B0
.. .. ..
5
R0 R
R0 C .Q
R0 U
R0 R
R0 T
`,
		},
	}
}

func simulate(input string) ([]eventRecord, int, int, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return nil, 0, 0, err
	}

	board := make([][]int, n)
	for i := range board {
		board[i] = make([]int, m)
	}

	players := make(map[string]*player)
	var quaffle ball
	var bludger ball
	bludgerExist := false

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			var cell string
			if _, err := fmt.Fscan(reader, &cell); err != nil {
				return nil, 0, 0, err
			}
			switch cell {
			case "..":
			case "RG":
				board[i][j] = goalRed
			case "BG":
				board[i][j] = goalBlue
			case ".Q":
				quaffle.x, quaffle.y = i, j
			case ".B":
				bludger.x, bludger.y = i, j
				bludgerExist = true
			default:
				if len(cell) == 2 {
					players[cell] = &player{id: cell, team: cell[0], x: i, y: j, alive: true}
				} else {
					return nil, 0, 0, fmt.Errorf("invalid cell %q", cell)
				}
			}
		}
	}

	midX, midY := (n-1)/2, (m-1)/2

	var steps int
	if _, err := fmt.Fscan(reader, &steps); err != nil {
		return nil, 0, 0, err
	}

	events := []eventRecord{}
	redScore, blueScore := 0, 0

	for t := 0; t < steps; t++ {
		var entity, action string
		if _, err := fmt.Fscan(reader, &entity, &action); err != nil {
			return nil, 0, 0, err
		}
		target := ""
		if action == "C" {
			if _, err := fmt.Fscan(reader, &target); err != nil {
				return nil, 0, 0, err
			}
		}

		switch entity {
		case ".Q":
			if quaffle.carrier == nil && isMove(action) {
				move(&quaffle.x, &quaffle.y, action)
			}
		case ".B":
			if bludgerExist && isMove(action) {
				move(&bludger.x, &bludger.y, action)
				eliminated := []*player{}
				for _, p := range players {
					if p.alive && p.x == bludger.x && p.y == bludger.y {
						eliminated = append(eliminated, p)
					}
				}
				if len(eliminated) > 0 {
					sort.Slice(eliminated, func(i, j int) bool {
						return eliminated[i].id < eliminated[j].id
					})
					for _, p := range eliminated {
						eliminatePlayer(p, &quaffle, t, &events)
					}
				}
			}
		default:
			p := players[entity]
			if p == nil || !p.alive {
				continue
			}
			switch action {
			case "U", "D", "L", "R":
				move(&p.x, &p.y, action)
				if p.carrying {
					quaffle.x, quaffle.y = p.x, p.y
				}
				if bludgerExist && p.x == bludger.x && p.y == bludger.y {
					eliminatePlayer(p, &quaffle, t, &events)
				}
			case "C":
				if target == ".Q" && quaffle.carrier == nil && p.x == quaffle.x && p.y == quaffle.y {
					quaffle.carrier = p
					p.carrying = true
				}
			case "T":
				if p.carrying {
					p.carrying = false
					quaffle.carrier = nil
					quaffle.x, quaffle.y = p.x, p.y
					switch board[p.x][p.y] {
					case goalRed:
						events = append(events, eventRecord{time: t, kind: "GOAL", label: "BLUE"})
						blueScore++
						quaffle.x, quaffle.y = midX, midY
					case goalBlue:
						events = append(events, eventRecord{time: t, kind: "GOAL", label: "RED"})
						redScore++
						quaffle.x, quaffle.y = midX, midY
					}
				}
			}
		}
	}

	return events, redScore, blueScore, nil
}

func move(x, y *int, action string) {
	switch action {
	case "U":
		*x = *x - 1
	case "D":
		*x = *x + 1
	case "L":
		*y = *y - 1
	case "R":
		*y = *y + 1
	}
}

func isMove(action string) bool {
	return action == "U" || action == "D" || action == "L" || action == "R"
}

func eliminatePlayer(p *player, quaffle *ball, t int, events *[]eventRecord) {
	if !p.alive {
		return
	}
	p.alive = false
	if p.carrying {
		p.carrying = false
		quaffle.carrier = nil
		quaffle.x, quaffle.y = p.x, p.y
	}
	*events = append(*events, eventRecord{time: t, kind: "ELIMINATED", label: p.id})
}
