package main

import (
    "bytes"
    "fmt"
    "math"
    "math/rand"
    "os"
    "os/exec"
    "strconv"
    "strings"
)

func runExe(path string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

// We validate candidate outputs directly instead of comparing to a fixed reference
// because multiple valid operation sequences exist.

func buildCase(arr []int) []byte {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", len(arr)))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return []byte(sb.String())
}

func genRandomCase(rng *rand.Rand) []byte {
	n := rng.Intn(8) + 3 // 3..10
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(200) - 100
	}
	return buildCase(arr)
}

func genTests() [][]byte {
	rng := rand.New(rand.NewSource(44))
	tests := make([][]byte, 0, 100)
	tests = append(tests, buildCase([]int{1, 2, 3}))
	tests = append(tests, buildCase([]int{3, 2, 1}))
	for len(tests) < 100 {
		tests = append(tests, genRandomCase(rng))
	}
	return tests
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: go run verifierC.go /path/to/binary")
        return
    }
    bin := os.Args[1]

    tests := genTests()
    for i, tc := range tests {
        got, err := runExe(bin, tc)
        if err != nil {
            fmt.Printf("candidate runtime error on test %d: %v\n%s", i+1, err, got)
            os.Exit(1)
        }
        // Parse input array from tc
        lines := strings.Split(strings.TrimSpace(string(tc)), "\n")
        if len(lines) < 3 {
            fmt.Printf("internal error: malformed test input on %d\n", i+1)
            os.Exit(1)
        }
        n, _ := strconv.Atoi(strings.TrimSpace(lines[1]))
        arrStr := strings.Fields(lines[2])
        if len(arrStr) != n {
            fmt.Printf("internal error: expected %d numbers, got %d on test %d\n", n, len(arrStr), i+1)
            os.Exit(1)
        }
        a := make([]int64, n)
        for j := 0; j < n; j++ {
            v, _ := strconv.ParseInt(arrStr[j], 10, 64)
            a[j] = v
        }
        if !validate1635C(a, strings.TrimSpace(got)) {
            fmt.Printf("test %d failed\ninput:\n%sgot:\n%s\n", i+1, string(tc), got)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}

func nonDecreasing(a []int64) bool {
    for i := 0; i+1 < len(a); i++ {
        if a[i] > a[i+1] {
            return false
        }
    }
    return true
}

func validate1635C(a []int64, out string) bool {
    out = strings.TrimSpace(out)
    if out == "" {
        return false
    }
    fields := strings.Fields(out)
    // Case: -1
    if fields[0] == "-1" {
        // Accept -1 only when impossible by common criteria
        if nonDecreasing(a) {
            return false
        }
        n := len(a)
        if n < 2 {
            return true
        }
        if a[n-2] > a[n-1] {
            return true
        }
        // If last element negative and array not sorted, consider impossible
        if a[n-1] < 0 {
            return true
        }
        // Otherwise, likely solvable
        return false
    }
    // Otherwise, parse m and operations
    m, err := strconv.Atoi(fields[0])
    if err != nil || m < 0 {
        return false
    }
    // Extract following lines
    lines := strings.Split(out, "\n")
    if len(lines) < 1+m {
        return false
    }
    b := make([]int64, len(a))
    copy(b, a)
    for i := 0; i < m; i++ {
        parts := strings.Fields(lines[i+1])
        if len(parts) != 3 {
            return false
        }
        x, e1 := strconv.Atoi(parts[0])
        y, e2 := strconv.Atoi(parts[1])
        z, e3 := strconv.Atoi(parts[2])
        if e1 != nil || e2 != nil || e3 != nil {
            return false
        }
        n := len(b)
        if !(1 <= x && x < y && y < z && z <= n) {
            return false
        }
        val := b[y-1] - b[z-1]
        if math.Abs(float64(val)) >= 1e18 {
            return false
        }
        b[x-1] = val
    }
    return nonDecreasing(b)
}
