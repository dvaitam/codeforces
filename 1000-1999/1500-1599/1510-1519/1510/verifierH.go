package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

const testCount = 120

type segment struct {
	L, R int
}

type interval struct {
	L, R int
}

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "1510H.go")
	tmp, err := os.CreateTemp("", "oracle1510H")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	tmp.Close()
	os.Remove(path)
	cmd := exec.Command("go", "build", "-o", path, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build oracle: %v\n%s", err, string(out))
	}
	return path, nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func genLaminar(r *rand.Rand, n int) []segment {
	type node struct {
		seg      segment
		children []node
	}
	var build func(depth, maxDepth int, L, R int, needed *int) node
	build = func(depth, maxDepth int, L, R int, needed *int) node {
		if *needed == 0 {
			return node{seg: segment{L, R}}
		}
		k := 0
		if depth < maxDepth {
			k = r.Intn(3)
		}
		children := make([]node, 0, k)
		left := L + 1
		span := R - L - 1
		if span <= 0 {
			k = 0
		}
		lengths := make([]int, k)
		remaining := span - k
		for i := 0; i < k; i++ {
			lenPart := 1
			if i == k-1 {
				lenPart += remaining
			} else {
				extra := 0
				if remaining > 0 {
					extra = r.Intn(remaining + 1)
				}
				lenPart += extra
				remaining -= extra
			}
			lengths[i] = lenPart
		}
		for i := 0; i < k; i++ {
			nodeL := left
			nodeR := left + lengths[i]
			left = nodeR + 1
			if nodeR > R-1 {
				nodeR = R - 1
			}
			child := build(depth+1, maxDepth, nodeL, nodeR, needed)
			children = append(children, child)
			*needed--
			if *needed <= 0 {
				break
			}
		}
		return node{seg: segment{L, R}, children: children}
	}
	maxDepth := 4
	needed := n - 1
	root := build(0, maxDepth, 0, 1000000000, &needed)
	res := []segment{}
	var flatten func(nd node)
	flatten = func(nd node) {
		res = append(res, nd.seg)
		sort.Slice(nd.children, func(i, j int) bool {
			return nd.children[i].seg.L < nd.children[j].seg.L
		})
		for _, ch := range nd.children {
			flatten(ch)
		}
	}
	flatten(root)
	if len(res) > n {
		res = res[:n]
	}
	return res
}

func genCase(r *rand.Rand) (string, []segment) {
	n := r.Intn(30) + 1
	if r.Intn(5) == 0 {
		n = r.Intn(80) + 1
	}
	segs := genLaminar(r, n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(segs))
	for _, seg := range segs {
		fmt.Fprintf(&sb, "%d %d\n", seg.L, seg.R)
	}
	return sb.String(), segs
}

func parseOutput(out string, n int) (int, []interval, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	totalLine, err := reader.ReadString('\n')
	if err != nil {
		return 0, nil, fmt.Errorf("missing total length line")
	}
	totalLine = strings.TrimSpace(totalLine)
	total, err := strconv.Atoi(totalLine)
	if err != nil {
		return 0, nil, fmt.Errorf("invalid total length")
	}
	res := make([]interval, n)
	for i := 0; i < n; i++ {
		line, err := reader.ReadString('\n')
		if err != nil && len(line) == 0 {
			return 0, nil, fmt.Errorf("missing interval line %d", i+1)
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return 0, nil, fmt.Errorf("line %d should contain 2 integers", i+1)
		}
		l, err := strconv.Atoi(fields[0])
		if err != nil {
			return 0, nil, fmt.Errorf("invalid l on line %d", i+1)
		}
		r, err := strconv.Atoi(fields[1])
		if err != nil {
			return 0, nil, fmt.Errorf("invalid r on line %d", i+1)
		}
		if l >= r {
			return 0, nil, fmt.Errorf("line %d has non-positive length interval", i+1)
		}
		res[i] = interval{l, r}
	}
	if extra, err := reader.ReadString('\n'); err == nil && len(strings.TrimSpace(extra)) > 0 {
		return 0, nil, fmt.Errorf("extra output detected")
	}
	return total, res, nil
}

func validateIntervals(segs []segment, intervals []interval) error {
	if len(segs) != len(intervals) {
		return fmt.Errorf("interval count mismatch")
	}
	total := 0
	active := []interval{}
	for i, seg := range segs {
		iv := intervals[i]
		if iv.L < seg.L || iv.R > seg.R {
			return fmt.Errorf("interval %d not within segment", i+1)
		}
		total += iv.R - iv.L
		active = append(active, iv)
	}
	sort.Slice(active, func(i, j int) bool {
		if active[i].L == active[j].L {
			return active[i].R < active[j].R
		}
		return active[i].L < active[j].L
	})
	for i := 1; i < len(active); i++ {
		if active[i].L < active[i-1].R {
			return fmt.Errorf("intervals overlap")
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	r := rand.New(rand.NewSource(1))
	for t := 0; t < testCount; t++ {
		input, segs := genCase(r)
		expectStr, err := runProgram(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		expectTotal, expectIntervals, err := parseOutput(expectStr, len(segs))
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid oracle output on test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		gotStr, err := runProgram(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		gotTotal, gotIntervals, err := parseOutput(gotStr, len(segs))
		if err != nil {
			fmt.Printf("test %d failed\ninput:\n%s\nerror: %v\n", t+1, input, err)
			os.Exit(1)
		}
		if err := validateIntervals(segs, gotIntervals); err != nil {
			fmt.Printf("test %d failed\ninput:\n%s\nerror: %v\n", t+1, input, err)
			os.Exit(1)
		}
		if gotTotal != expectTotal {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected total: %d\ngot: %d\n", t+1, input, expectTotal, gotTotal)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", testCount)
}
