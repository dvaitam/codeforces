package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	n     int
	perms [][]int
}

type trade struct {
	player byte
	card   int
}

func main() {
	args := os.Args[1:]
	if len(args) >= 1 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierD.go /path/to/solution")
		os.Exit(1)
	}
	target := args[0]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build oracle: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := append(deterministicCases(), randomCases()...)
	input := buildInput(tests)

	oracleOut, err := runBinary(oracle, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle runtime error: %v\noutput:\n%s", err, oracleOut)
		os.Exit(1)
	}
	candOut, err := runBinary(target, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\noutput:\n%s", err, candOut)
		os.Exit(1)
	}

	possibility, err := parseOracleOutput(oracleOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse oracle output: %v\noutput:\n%s", err, oracleOut)
		os.Exit(1)
	}
	if err := validateCandidateOutput(candOut, tests, possibility); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	fmt.Println("All tests passed.")
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier location")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2028D-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleD")
	cmd := exec.Command("go", "build", "-o", outPath, "2028D.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("go build failed: %v\n%s", err, string(out))
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String() + stderr.String(), fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for p := 0; p < 3; p++ {
			for i := 1; i <= tc.n; i++ {
				if i > 1 {
					sb.WriteByte(' ')
				}
				sb.WriteString(strconv.Itoa(tc.perms[p][i]))
			}
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func deterministicCases() []testCase {
	return []testCase{
		makeCase([][]int{
			{0, 3, 1, 2},
			{0, 2, 3, 1},
			{0, 1, 2, 3},
		}),
		makeCase([][]int{
			{0, 1, 2, 3, 4},
			{0, 1, 2, 3, 4},
			{0, 1, 2, 3, 4},
		}),
		makeCase([][]int{
			{0, 4, 1, 3, 2},
			{0, 3, 4, 1, 2},
			{0, 2, 1, 4, 3},
		}),
	}
}

func makeCase(perms [][]int) testCase {
	n := len(perms[0]) - 1
	cp := make([][]int, 3)
	for i := 0; i < 3; i++ {
		cp[i] = append([]int(nil), perms[i]...)
	}
	return testCase{n: n, perms: cp}
}

func randomCases() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 80)
	for len(tests) < cap(tests) {
		n := rng.Intn(25) + 2
		perms := make([][]int, 3)
		for p := 0; p < 3; p++ {
			arr := rng.Perm(n)
			perms[p] = make([]int, n+1)
			for i := 0; i < n; i++ {
				perms[p][i+1] = arr[i] + 1
			}
		}
		tests = append(tests, testCase{n: n, perms: perms})
	}
	return tests
}

func parseOracleOutput(out string, t int) ([]bool, error) {
	tokens := strings.Fields(out)
	idx := 0
	res := make([]bool, t)
	for i := 0; i < t; i++ {
		if idx >= len(tokens) {
			return nil, fmt.Errorf("oracle output ended early on test %d", i+1)
		}
		token := strings.ToUpper(tokens[idx])
		idx++
		switch token {
		case "YES":
			res[i] = true
			if idx >= len(tokens) {
				return nil, fmt.Errorf("oracle output missing trade count on test %d", i+1)
			}
			k, err := strconv.Atoi(tokens[idx])
			if err != nil || k < 0 {
				return nil, fmt.Errorf("oracle trade count invalid on test %d: %v", i+1, tokens[idx])
			}
			idx++
			for step := 0; step < k; step++ {
				if idx+1 >= len(tokens) {
					return nil, fmt.Errorf("oracle output missing trade %d on test %d", step+1, i+1)
				}
				idx += 2
			}
		case "NO":
			res[i] = false
		default:
			return nil, fmt.Errorf("oracle output invalid token %q on test %d", token, i+1)
		}
	}
	if idx != len(tokens) {
		return nil, fmt.Errorf("oracle output has extra tokens")
	}
	return res, nil
}

func validateCandidateOutput(out string, tests []testCase, possible []bool) error {
	tokens := strings.Fields(out)
	idx := 0
	for i, tc := range tests {
		if idx >= len(tokens) {
			return fmt.Errorf("candidate output ended early on test %d", i+1)
		}
		token := strings.ToUpper(tokens[idx])
		idx++
		if token == "NO" {
			if possible[i] {
				return fmt.Errorf("test %d: candidate answered NO but solution exists", i+1)
			}
			continue
		}
		if token != "YES" {
			return fmt.Errorf("test %d: expected YES/NO token, got %q", i+1, tokens[idx-1])
		}
		if !possible[i] {
			return fmt.Errorf("test %d: candidate answered YES but problem is impossible", i+1)
		}
		if idx >= len(tokens) {
			return fmt.Errorf("test %d: missing trade count", i+1)
		}
		k, err := strconv.Atoi(tokens[idx])
		if err != nil || k < 0 {
			return fmt.Errorf("test %d: invalid trade count %q", i+1, tokens[idx])
		}
		idx++
		steps := make([]trade, 0, k)
		for step := 0; step < k; step++ {
			if idx >= len(tokens) {
				return fmt.Errorf("test %d: missing player for trade %d", i+1, step+1)
			}
			playerToken := strings.ToLower(tokens[idx])
			idx++
			if idx >= len(tokens) {
				return fmt.Errorf("test %d: missing card for trade %d", i+1, step+1)
			}
			cardVal, err := strconv.Atoi(tokens[idx])
			if err != nil {
				return fmt.Errorf("test %d: invalid card value %q on trade %d", i+1, tokens[idx], step+1)
			}
			idx++
			if len(playerToken) == 0 {
				return fmt.Errorf("test %d: empty player token on trade %d", i+1, step+1)
			}
			player := playerToken[0]
			switch player {
			case 'q', 'k', 'j':
			default:
				return fmt.Errorf("test %d: invalid player %q on trade %d", i+1, playerToken, step+1)
			}
			steps = append(steps, trade{player: player, card: cardVal})
		}
		if err := validateTrades(tc, steps); err != nil {
			return fmt.Errorf("test %d: %v", i+1, err)
		}
	}
	if idx != len(tokens) {
		return fmt.Errorf("candidate output has extra tokens")
	}
	return nil
}

func validateTrades(tc testCase, steps []trade) error {
	current := 1
	n := tc.n
	PlayerIndex := func(ch byte) int {
		switch ch {
		case 'q':
			return 0
		case 'k':
			return 1
		case 'j':
			return 2
		default:
			return -1
		}
	}

	for idx, step := range steps {
		if step.card < 1 || step.card > n {
			return fmt.Errorf("trade %d targets invalid card %d", idx+1, step.card)
		}
		if step.card <= current {
			return fmt.Errorf("trade %d does not increase card value (%d -> %d)", idx+1, current, step.card)
		}
		playerIdx := PlayerIndex(step.player)
		if playerIdx == -1 {
			return fmt.Errorf("trade %d has unknown player %c", idx+1, step.player)
		}
		currPref := tc.perms[playerIdx][current]
		nextPref := tc.perms[playerIdx][step.card]
		if currPref <= nextPref {
			return fmt.Errorf("trade %d invalid: player %c does not prefer card %d over %d", idx+1, step.player, current, step.card)
		}
		current = step.card
	}
	if current != n {
		return fmt.Errorf("sequence ends on card %d instead of %d", current, n)
	}
	return nil
}
