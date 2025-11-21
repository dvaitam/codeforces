package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const referenceSolution = "0-999/200-299/250-259/250/250B.go"

type zeroRun struct {
	start  int
	length int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := manualTests()
	for len(tests) < 150 {
		tests = append(tests, randomTestInput(rng))
	}

	for idx, input := range tests {
		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		if err := compareOutputs(refOut, candOut); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n", idx+1, err, input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReferenceBinary() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-250B-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref250B.bin")
	cmd := exec.Command("go", "build", "-o", binPath, referenceSolution)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	cleanup := func() {
		_ = os.RemoveAll(dir)
	}
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
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func manualTests() []string {
	tests := []string{
		"6\na56f:d3:0:0124:01:f19a:1000:00\na56f:00d3:0000:0124:0001::\na56f::0124:0001:0000:1234:0ff0\na56f:0000::0000:0001:0000:1234:0ff0\n::\n0ea::4d:f4:6:0\n",
		singleAddress("::"),
		singleAddress("0:0:0:0:0:0:0:0"),
		singleAddress("1::"),
		singleAddress("::1"),
		singleAddress("abcd::"),
		singleAddress("f:00f:0:1:02:3:4:5"),
		buildInput([]string{"::", "0::", "1:2:3:4:5:6:7:8"}),
	}
	return tests
}

func singleAddress(addr string) string {
	return buildInput([]string{addr})
}

func buildInput(addrs []string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(addrs)))
	for _, addr := range addrs {
		sb.WriteString(addr)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func randomTestInput(rng *rand.Rand) string {
	n := rng.Intn(100) + 1
	addrs := make([]string, n)
	for i := 0; i < n; i++ {
		addrs[i] = randomShortIPv6(rng)
	}
	return buildInput(addrs)
}

func randomShortIPv6(rng *rand.Rand) string {
	vals := make([]int, 8)
	for i := 0; i < 8; i++ {
		vals[i] = rng.Intn(1 << 16)
	}
	if rng.Intn(2) == 0 {
		start := rng.Intn(8)
		length := rng.Intn(8-start) + 1
		for i := start; i < start+length; i++ {
			vals[i] = 0
		}
	}
	blocks := make([]string, 8)
	for i, v := range vals {
		blocks[i] = fmt.Sprintf("%04x", v)
	}
	shortBlocks := make([]string, 8)
	for i, block := range blocks {
		leading := 0
		for leading < len(block) && block[leading] == '0' {
			leading++
		}
		maxRemove := leading
		if maxRemove == len(block) {
			maxRemove = len(block) - 1
		}
		remove := 0
		if maxRemove > 0 {
			remove = rng.Intn(maxRemove + 1)
		}
		short := block[remove:]
		if len(short) == 0 {
			short = "0"
		}
		shortBlocks[i] = short
	}

	runs := collectZeroRuns(vals)
	useCompression := len(runs) > 0 && rng.Intn(2) == 0
	if useCompression {
		run := runs[rng.Intn(len(runs))]
		left := strings.Join(shortBlocks[:run.start], ":")
		right := strings.Join(shortBlocks[run.start+run.length:], ":")
		switch {
		case left == "" && right == "":
			return "::"
		case left == "":
			return "::" + right
		case right == "":
			return left + "::"
		default:
			return left + "::" + right
		}
	}
	return strings.Join(shortBlocks, ":")
}

func collectZeroRuns(vals []int) []zeroRun {
	var runs []zeroRun
	for i := 0; i < len(vals); {
		if vals[i] == 0 {
			j := i
			for j < len(vals) && vals[j] == 0 {
				j++
			}
			runs = append(runs, zeroRun{start: i, length: j - i})
			i = j
		} else {
			i++
		}
	}
	return runs
}

func compareOutputs(expect, got string) error {
	expLines := normalizeLines(expect)
	gotLines := normalizeLines(got)
	if len(expLines) != len(gotLines) {
		return fmt.Errorf("expected %d lines got %d", len(expLines), len(gotLines))
	}
	for i := range expLines {
		if expLines[i] != gotLines[i] {
			return fmt.Errorf("line %d mismatch: expected %q got %q", i+1, expLines[i], gotLines[i])
		}
	}
	return nil
}

func normalizeLines(out string) []string {
	out = strings.ReplaceAll(out, "\r\n", "\n")
	out = strings.TrimRight(out, "\n")
	if len(out) == 0 {
		return nil
	}
	return strings.Split(out, "\n")
}
