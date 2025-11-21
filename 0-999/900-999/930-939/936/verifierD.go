package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const (
	refSource      = "0-999/900-999/930-939/936/936D.go"
	maxTransitions = 2_000_000
	verdictYes     = "yes"
	verdictNo      = "no"
)

type problemInput struct {
	n         int
	t         int
	totalObst int
	lanes     [2]laneData
}

type laneData struct {
	positions []int
	posIndex  map[int]int
}

type candidateOutput struct {
	verdict     string
	transitions []int
	shots       []shot
}

type shot struct {
	x    int
	lane int
}

func main() {
	if len(os.Args) != 2 {
		fail("usage: go run verifierD.go /path/to/candidate")
	}
	candidate := os.Args[1]

	inputData, err := io.ReadAll(os.Stdin)
	if err != nil {
		fail("failed to read input: %v", err)
	}

	prob, err := parseInput(inputData)
	if err != nil {
		fail("failed to parse input: %v", err)
	}

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	refOut, err := runProgram(exec.Command(refBin), inputData)
	if err != nil {
		fail("reference execution failed: %v", err)
	}
	refVerdict, err := parseVerdict(refOut)
	if err != nil {
		fail("failed to parse reference output: %v", err)
	}

	candOut, err := runProgram(commandFor(candidate), inputData)
	if err != nil {
		fail("candidate execution failed: %v", err)
	}
	plan, err := parseCandidateOutput(candOut, prob.n, prob.totalObst)
	if err != nil {
		fail("failed to parse candidate output: %v", err)
	}

	switch refVerdict {
	case verdictNo:
		if plan.verdict != verdictNo {
			fail("reference verdict is No but candidate printed Yes")
		}
	case verdictYes:
		if plan.verdict != verdictYes {
			fail("reference verdict is Yes but candidate printed No")
		}
		if err := validatePlan(prob, plan); err != nil {
			fail("invalid plan: %v", err)
		}
	default:
		fail("unknown reference verdict: %s", refVerdict)
	}

	fmt.Println("OK")
}

func parseInput(data []byte) (problemInput, error) {
	var prob problemInput
	reader := bufio.NewReader(bytes.NewReader(data))
	var n, m1, m2, t int
	if _, err := fmt.Fscan(reader, &n, &m1, &m2, &t); err != nil {
		return prob, err
	}
	prob.n = n
	prob.t = t
	prob.totalObst = m1 + m2
	for lane := 0; lane < 2; lane++ {
		var count int
		if lane == 0 {
			count = m1
		} else {
			count = m2
		}
		if count == 0 {
			prob.lanes[lane] = laneData{
				positions: nil,
				posIndex:  map[int]int{},
			}
			continue
		}
		pos := make([]int, count)
		for i := 0; i < count; i++ {
			if _, err := fmt.Fscan(reader, &pos[i]); err != nil {
				return prob, err
			}
		}
		index := make(map[int]int, count)
		for i, v := range pos {
			index[v] = i
		}
		prob.lanes[lane] = laneData{
			positions: pos,
			posIndex:  index,
		}
	}
	return prob, nil
}

func parseVerdict(out string) (string, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	token, err := readToken(reader)
	if err != nil {
		return "", err
	}
	token = strings.ToLower(token)
	if token != verdictYes && token != verdictNo {
		return "", fmt.Errorf("unexpected verdict token %q", token)
	}
	return token, nil
}

func parseCandidateOutput(out string, n int, maxShots int) (candidateOutput, error) {
	var plan candidateOutput
	reader := bufio.NewReader(strings.NewReader(out))
	token, err := readToken(reader)
	if err != nil {
		return plan, err
	}
	token = strings.ToLower(token)
	if token != verdictYes && token != verdictNo {
		return plan, fmt.Errorf("expected Yes/No, got %q", token)
	}
	plan.verdict = token
	if token == verdictNo {
		if extra, err := readToken(reader); err != io.EOF {
			if err == nil {
				return plan, fmt.Errorf("unexpected extra token %q after verdict No", extra)
			}
			return plan, err
		}
		return plan, nil
	}
	kStr, err := readToken(reader)
	if err != nil {
		return plan, fmt.Errorf("missing transition count: %v", err)
	}
	k, err := strconv.Atoi(kStr)
	if err != nil || k < 0 {
		return plan, fmt.Errorf("invalid transition count %q", kStr)
	}
	if k > maxTransitions {
		return plan, fmt.Errorf("too many transitions (%d)", k)
	}
	plan.transitions = make([]int, k)
	for i := 0; i < k; i++ {
		tok, err := readToken(reader)
		if err != nil {
			return plan, fmt.Errorf("not enough transition coordinates")
		}
		val64, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return plan, fmt.Errorf("invalid transition coordinate %q", tok)
		}
		if val64 < 0 || val64 > int64(n)+1 {
			return plan, fmt.Errorf("transition coordinate %d out of range [0,%d]", val64, n+1)
		}
		val := int(val64)
		if i > 0 && val <= plan.transitions[i-1] {
			return plan, fmt.Errorf("transition coordinates must be strictly increasing")
		}
		plan.transitions[i] = val
	}

	sStr, err := readToken(reader)
	if err != nil {
		return plan, fmt.Errorf("missing shots count: %v", err)
	}
	s, err := strconv.Atoi(sStr)
	if err != nil || s < 0 {
		return plan, fmt.Errorf("invalid shots count %q", sStr)
	}
	if s > maxShots {
		return plan, fmt.Errorf("shots count %d exceeds total obstacles %d", s, maxShots)
	}
	plan.shots = make([]shot, s)
	for i := 0; i < s; i++ {
		xTok, err := readToken(reader)
		if err != nil {
			return plan, fmt.Errorf("not enough shot coordinates")
		}
		yTok, err := readToken(reader)
		if err != nil {
			return plan, fmt.Errorf("missing lane for shot %d", i+1)
		}
		xVal, err := strconv.ParseInt(xTok, 10, 64)
		if err != nil {
			return plan, fmt.Errorf("invalid shot x %q", xTok)
		}
		yVal, err := strconv.Atoi(yTok)
		if err != nil {
			return plan, fmt.Errorf("invalid shot lane %q", yTok)
		}
		if xVal < 1 || xVal > int64(n) {
			return plan, fmt.Errorf("shot x=%d out of range [1,%d]", xVal, n)
		}
		if yVal != 1 && yVal != 2 {
			return plan, fmt.Errorf("shot lane must be 1 or 2, got %d", yVal)
		}
		if i > 0 && xVal < int64(plan.shots[i-1].x) {
			return plan, fmt.Errorf("shot timestamps must be non-decreasing")
		}
		plan.shots[i] = shot{
			x:    int(xVal),
			lane: yVal,
		}
	}
	if extra, err := readToken(reader); err != io.EOF {
		if err == nil {
			return plan, fmt.Errorf("unexpected extra token %q", extra)
		}
		return plan, err
	}
	return plan, nil
}

func validatePlan(prob problemInput, plan candidateOutput) error {
	if len(plan.transitions) > maxTransitions {
		return fmt.Errorf("transition count exceeds limit")
	}
	if len(plan.shots) > prob.totalObst {
		return fmt.Errorf("shots exceed number of obstacles")
	}
	eventCap := prob.totalObst + len(plan.transitions) + len(plan.shots) + 2
	eventSet := make(map[int]struct{}, eventCap)
	eventSet[0] = struct{}{}
	eventSet[prob.n+1] = struct{}{}
	for _, lane := range prob.lanes {
		for _, pos := range lane.positions {
			eventSet[pos] = struct{}{}
		}
	}
	for _, tr := range plan.transitions {
		eventSet[tr] = struct{}{}
	}
	for _, sh := range plan.shots {
		eventSet[sh.x] = struct{}{}
	}
	events := make([]int, 0, len(eventSet))
	for x := range eventSet {
		events = append(events, x)
	}
	sort.Ints(events)

	alive := [2][]bool{}
	for lane := 0; lane < 2; lane++ {
		if len(prob.lanes[lane].positions) == 0 {
			alive[lane] = nil
			continue
		}
		alive[lane] = make([]bool, len(prob.lanes[lane].positions))
		for i := range alive[lane] {
			alive[lane][i] = true
		}
	}

	lane := 1
	transIdx := 0
	shotIdx := 0
	nextReady := int64(prob.t)

	for _, x := range events {
		if shotIdx < len(plan.shots) && plan.shots[shotIdx].x < x {
			return fmt.Errorf("shots must be listed in chronological order")
		}
		if transIdx < len(plan.transitions) && plan.transitions[transIdx] < x {
			return fmt.Errorf("transitions must be strictly increasing")
		}
		if x >= 1 && x <= prob.n {
			if idx, ok := prob.lanes[lane-1].posIndex[x]; ok && alive[lane-1][idx] {
				return fmt.Errorf("collision at x=%d lane=%d", x, lane)
			}
		}
		for transIdx < len(plan.transitions) && plan.transitions[transIdx] == x {
			newLane := 3 - lane
			if x >= 1 && x <= prob.n {
				if idx, ok := prob.lanes[newLane-1].posIndex[x]; ok && alive[newLane-1][idx] {
					return fmt.Errorf("transition at x=%d enters obstacle on lane %d", x, newLane)
				}
			}
			lane = newLane
			transIdx++
		}
		for shotIdx < len(plan.shots) && plan.shots[shotIdx].x == x {
			sh := plan.shots[shotIdx]
			if sh.lane != lane {
				return fmt.Errorf("shot at x=%d declared lane %d but tank in lane %d", sh.x, sh.lane, lane)
			}
			if int64(sh.x) < nextReady {
				return fmt.Errorf("shot at x=%d fired before gun reloaded (ready at %d)", sh.x, nextReady)
			}
			if err := destroyNearest(&prob, &alive, lane-1, sh.x); err != nil {
				return fmt.Errorf("shot at x=%d failed: %v", sh.x, err)
			}
			nextReady = int64(sh.x) + int64(prob.t)
			shotIdx++
		}
	}

	if transIdx != len(plan.transitions) {
		return fmt.Errorf("unprocessed transitions remain")
	}
	if shotIdx != len(plan.shots) {
		return fmt.Errorf("unprocessed shots remain")
	}
	return nil
}

func destroyNearest(prob *problemInput, alive *[2][]bool, lane int, shotX int) error {
	positions := prob.lanes[lane].positions
	if len(positions) == 0 {
		return fmt.Errorf("no obstacles on lane")
	}
	target := shotX + 1
	idx := sort.Search(len(positions), func(i int) bool {
		return positions[i] >= target
	})
	for idx < len(positions) {
		if alive[lane][idx] {
			alive[lane][idx] = false
			return nil
		}
		idx++
	}
	return fmt.Errorf("no obstacle ahead to destroy")
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "936D-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func runProgram(cmd *exec.Cmd, input []byte) (string, error) {
	cmd.Stdin = bytes.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := stderr.String()
		if msg == "" {
			msg = stdout.String()
		}
		return "", fmt.Errorf("%v\n%s", err, msg)
	}
	return stdout.String(), nil
}

func readToken(r *bufio.Reader) (string, error) {
	var sb strings.Builder
	for {
		b, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				return "", io.EOF
			}
			return "", err
		}
		if !isSpace(b) {
			sb.WriteByte(b)
			break
		}
	}
	for {
		b, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				return sb.String(), nil
			}
			return "", err
		}
		if isSpace(b) {
			break
		}
		sb.WriteByte(b)
	}
	return sb.String(), nil
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\n' || b == '\r' || b == '\t' || b == '\v' || b == '\f'
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
