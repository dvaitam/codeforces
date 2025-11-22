package main

import (
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

const (
	randomArrays = 5
)

type operation struct {
	kind int // 0=pushback,1=pushfront,2=popback,3=popfront,4=min
	idx  int
}

type testPlan struct {
	n    int
	arrs [][]int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}

	candidate, candCleanup, err := prepareBinary(os.Args[1], "candidate2141C")
	if err != nil {
		fmt.Println("failed to prepare contestant binary:", err)
		return
	}
	if candCleanup != nil {
		defer candCleanup()
	}

	oracle, oracleCleanup, err := prepareOracle()
	if err != nil {
		fmt.Println("failed to prepare reference solution:", err)
		return
	}
	defer oracleCleanup()

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	plans := buildTestPlans(rng)
	totalTests := 0
	for i, plan := range plans {
		candOps, err := getOps(candidate, plan.n)
		if err != nil {
			fmt.Printf("test group %d (n=%d) failed: %v\n", i+1, plan.n, err)
			return
		}
		oracleOps, err := getOps(oracle, plan.n)
		if err != nil {
			fmt.Printf("oracle failed for n=%d: %v\n", plan.n, err)
			return
		}

		for j, arr := range plan.arrs {
			expect, err := execute(oracleOps, arr)
			if err != nil {
				fmt.Printf("oracle runtime error for n=%d array %d: %v\n", plan.n, j+1, err)
				return
			}
			got, err := execute(candOps, arr)
			if err != nil {
				fmt.Printf("runtime error for n=%d array %d: %v\n", plan.n, j+1, err)
				return
			}
			if expect != got {
				fmt.Printf("wrong answer for n=%d array %d: expected %d got %d\narray: %v\n", plan.n, j+1, expect, got, arr)
				return
			}
			totalTests++
		}
	}

	fmt.Printf("All %d tests passed.\n", totalTests)
}

func prepareBinary(path, prefix string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		abs, err := filepath.Abs(path)
		if err != nil {
			return "", nil, err
		}
		tmp := filepath.Join(os.TempDir(), fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano()))
		cmd := exec.Command("go", "build", "-o", tmp, abs)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", nil, fmt.Errorf("go build failed: %v: %s", err, out)
		}
		return tmp, func() { os.Remove(tmp) }, nil
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", nil, err
	}
	return abs, nil, nil
}

func prepareOracle() (string, func(), error) {
	dir := sourceDir()
	src := filepath.Join(dir, "2141C.go")
	tmp := filepath.Join(os.TempDir(), fmt.Sprintf("oracle2141C_%d", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", tmp, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", nil, fmt.Errorf("go build oracle failed: %v: %s", err, out)
	}
	return tmp, func() { os.Remove(tmp) }, nil
}

func getOps(bin string, n int) ([]operation, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(fmt.Sprintf("%d\n", n))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("execution failed: %v\n%s", err, string(out))
	}
	text := strings.TrimSpace(string(out))
	if text == "" {
		return nil, fmt.Errorf("empty output")
	}

	lines := splitNonEmpty(text)
	k, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return nil, fmt.Errorf("failed to parse command count: %v", err)
	}
	limit := n * (n + 2)
	if k < 1 || k > limit {
		return nil, fmt.Errorf("command count %d is outside [1, %d]", k, limit)
	}
	if len(lines)-1 != k {
		return nil, fmt.Errorf("expected %d commands, got %d", k, len(lines)-1)
	}

	ops := make([]operation, k)
	for i := 0; i < k; i++ {
		op, err := parseOp(lines[i+1], n)
		if err != nil {
			return nil, fmt.Errorf("invalid command %d: %v", i+1, err)
		}
		ops[i] = op
	}
	return ops, nil
}

func parseOp(line string, n int) (operation, error) {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return operation{}, fmt.Errorf("empty line")
	}
	switch fields[0] {
	case "pushback":
		if len(fields) != 2 {
			return operation{}, fmt.Errorf("invalid pushback syntax")
		}
		idx, err := parseIdx(fields[1], n)
		if err != nil {
			return operation{}, err
		}
		return operation{kind: 0, idx: idx}, nil
	case "pushfront":
		if len(fields) != 2 {
			return operation{}, fmt.Errorf("invalid pushfront syntax")
		}
		idx, err := parseIdx(fields[1], n)
		if err != nil {
			return operation{}, err
		}
		return operation{kind: 1, idx: idx}, nil
	case "popback":
		if len(fields) != 1 {
			return operation{}, fmt.Errorf("invalid popback syntax")
		}
		return operation{kind: 2}, nil
	case "popfront":
		if len(fields) != 1 {
			return operation{}, fmt.Errorf("invalid popfront syntax")
		}
		return operation{kind: 3}, nil
	case "min":
		if len(fields) != 1 {
			return operation{}, fmt.Errorf("invalid min syntax")
		}
		return operation{kind: 4}, nil
	default:
		return operation{}, fmt.Errorf("unknown command %q", fields[0])
	}
}

func parseIdx(token string, n int) (int, error) {
	if !strings.HasPrefix(token, "a[") || !strings.HasSuffix(token, "]") {
		return 0, fmt.Errorf("invalid index token %q", token)
	}
	numStr := token[2 : len(token)-1]
	idx, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, fmt.Errorf("invalid index %q: %v", numStr, err)
	}
	if idx < 0 || idx >= n {
		return 0, fmt.Errorf("index %d out of bounds", idx)
	}
	return idx, nil
}

func execute(ops []operation, arr []int) (int64, error) {
	deq := make([]int, 0)
	freq := make(map[int]int)
	var sum int64
	minVal := 0
	hasMin := false

	addVal := func(x int) {
		freq[x]++
		if !hasMin || x < minVal {
			minVal = x
			hasMin = true
		}
	}
	removeVal := func(x int) {
		freq[x]--
		if freq[x] == 0 {
			delete(freq, x)
			if x == minVal {
				hasMin = false
				for v := range freq {
					if !hasMin || v < minVal {
						minVal = v
						hasMin = true
					}
				}
			}
		}
	}

	for i, op := range ops {
		switch op.kind {
		case 0: // pushback
			v := arr[op.idx]
			deq = append(deq, v)
			addVal(v)
		case 1: // pushfront
			v := arr[op.idx]
			deq = append([]int{v}, deq...)
			addVal(v)
		case 2: // popback
			if len(deq) == 0 {
				return 0, fmt.Errorf("popback on empty at command %d", i+1)
			}
			last := deq[len(deq)-1]
			deq = deq[:len(deq)-1]
			removeVal(last)
		case 3: // popfront
			if len(deq) == 0 {
				return 0, fmt.Errorf("popfront on empty at command %d", i+1)
			}
			first := deq[0]
			deq = deq[1:]
			removeVal(first)
		case 4: // min
			if len(deq) == 0 || !hasMin {
				return 0, fmt.Errorf("min on empty at command %d", i+1)
			}
			sum += int64(minVal)
		default:
			return 0, fmt.Errorf("unknown op kind %d", op.kind)
		}
	}
	return sum, nil
}

func buildTestPlans(rng *rand.Rand) []testPlan {
	nVals := []int{1, 2, 3, 5, 10, 50, 200, 500}
	plans := make([]testPlan, 0, len(nVals))
	for _, n := range nVals {
		plan := testPlan{n: n}
		plan.arrs = append(plan.arrs, makeConstArray(n, 1))
		plan.arrs = append(plan.arrs, increasingArray(n))
		plan.arrs = append(plan.arrs, decreasingArray(n))
		plan.arrs = append(plan.arrs, alternatingArray(n))
		for i := 0; i < randomArrays; i++ {
			plan.arrs = append(plan.arrs, randomArray(n, rng))
		}
		plans = append(plans, plan)
	}
	return plans
}

func makeConstArray(n int, v int) []int {
	arr := make([]int, n)
	for i := range arr {
		arr[i] = v
	}
	return arr
}

func increasingArray(n int) []int {
	arr := make([]int, n)
	for i := range arr {
		arr[i] = i + 1
	}
	return arr
}

func decreasingArray(n int) []int {
	arr := make([]int, n)
	for i := range arr {
		arr[i] = n - i
	}
	return arr
}

func alternatingArray(n int) []int {
	arr := make([]int, n)
	for i := range arr {
		if i%2 == 0 {
			arr[i] = 0
		} else {
			arr[i] = 1000
		}
	}
	return arr
}

func randomArray(n int, rng *rand.Rand) []int {
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(1_000_000_000)
	}
	return arr
}

func splitNonEmpty(text string) []string {
	raw := strings.Split(text, "\n")
	res := make([]string, 0, len(raw))
	for _, line := range raw {
		trim := strings.TrimSpace(line)
		if trim != "" {
			res = append(res, trim)
		}
	}
	return res
}

func sourceDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "."
	}
	return filepath.Dir(file)
}
