package main

import (
    "bufio"
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "strconv"
    "strings"
    "time"
)

// Compute the best achievable sum via brute force over split positions.
// For a bitmask M over positions [0..n-1], where a set bit at r indicates
// we "take" a[r] and also end a zero-run just before r, the achieved sum is:
//   sum(arr[r] for r where bit r is 1) + sum((lenGap)^2 over zero-runs)
// where zero-runs are stretches between chosen positions and at ends.
func bestSum(arr []int) int {
    n := len(arr)
    if n == 0 {
        return 0
    }
    maxSum := -1
    for M := 0; M < (1 << n); M++ {
        sum := 0
        l := 0
        for r := 0; r < n; r++ {
            if (M>>r)&1 == 1 {
                // Close gap [l, r)
                gap := r - l
                sum += gap * gap
                sum += arr[r]
                l = r + 1
            }
        }
        // trailing gap
        gap := n - l
        sum += gap * gap
        if sum > maxSum {
            maxSum = sum
        }
    }
    return maxSum
}

// Apply candidate operations to the array and return resulting sum.
// Each operation sets a[l..r] to mex of that subarray.
func applyAndSum(a []int, ops [][2]int) (int, error) {
    n := len(a)
    b := append([]int(nil), a...)
    for _, op := range ops {
        l := op[0]
        r := op[1]
        if l < 0 || r < 0 || l >= n || r >= n || l > r {
            return 0, fmt.Errorf("invalid op indices: %d %d", l+1, r+1)
        }
        // Compute mex within b[l..r]
        size := r - l + 1
        seen := make([]bool, size+2)
        for i := l; i <= r; i++ {
            v := b[i]
            if 0 <= v && v <= size+1 {
                if v < len(seen) {
                    seen[v] = true
                }
            }
        }
        mex := 0
        for mex < len(seen) && seen[mex] {
            mex++
        }
        for i := l; i <= r; i++ {
            b[i] = mex
        }
    }
    s := 0
    for _, v := range b {
        s += v
    }
    return s, nil
}

func genCase(rng *rand.Rand) (string, []int, int) {
    n := rng.Intn(6) + 1
    arr := make([]int, n)
    var in bytes.Buffer
    fmt.Fprintf(&in, "%d\n", n)
    for i := 0; i < n; i++ {
        // Random values including zeros and small positives
        arr[i] = rng.Intn(10)
        if i > 0 {
            in.WriteByte(' ')
        }
        fmt.Fprintf(&in, "%d", arr[i])
    }
    in.WriteByte('\n')
    bs := bestSum(arr)
    return in.String(), arr, bs
}

func run(bin, input string) (string, error) {
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    err := cmd.Run()
    return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))
    for i := 0; i < 100; i++ {
        in, arr, want := genCase(rng)
        raw, err := run(bin, in)
        if err != nil {
            fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
            os.Exit(1)
        }
        // Parse candidate output
        scanner := bufio.NewScanner(strings.NewReader(raw))
        scanner.Split(bufio.ScanLines)
        if !scanner.Scan() {
            fmt.Fprintf(os.Stderr, "test %d failed: empty output\ninput:\n%s\n", i+1, in)
            os.Exit(1)
        }
        first := strings.Fields(scanner.Text())
        if len(first) < 2 {
            fmt.Fprintf(os.Stderr, "test %d failed: malformed first line: %q\n", i+1, scanner.Text())
            os.Exit(1)
        }
        sumGot, err1 := strconv.Atoi(first[0])
        m, err2 := strconv.Atoi(first[1])
        if err1 != nil || err2 != nil {
            fmt.Fprintf(os.Stderr, "test %d failed: non-integer header: %q\n", i+1, scanner.Text())
            os.Exit(1)
        }
        if m < 0 || m > 500000 {
            fmt.Fprintf(os.Stderr, "test %d failed: invalid m=%d\n", i+1, m)
            os.Exit(1)
        }
        var ops [][2]int
        for j := 0; j < m; j++ {
            if !scanner.Scan() {
                fmt.Fprintf(os.Stderr, "test %d failed: insufficient operations lines (got %d)\n", i+1, j)
                os.Exit(1)
            }
            f := strings.Fields(scanner.Text())
            if len(f) < 2 {
                fmt.Fprintf(os.Stderr, "test %d failed: malformed op line: %q\n", i+1, scanner.Text())
                os.Exit(1)
            }
            l, er1 := strconv.Atoi(f[0])
            r, er2 := strconv.Atoi(f[1])
            if er1 != nil || er2 != nil {
                fmt.Fprintf(os.Stderr, "test %d failed: non-integer op: %q\n", i+1, scanner.Text())
                os.Exit(1)
            }
            // Convert to 0-based
            l--
            r--
            ops = append(ops, [2]int{l, r})
        }
        // Validate sum first
        if sumGot != want {
            fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%s\nexpected sum:\n%d\ngot header sum:\n%d\n", i+1, in, want, sumGot)
            os.Exit(1)
        }
        // Simulate operations
        finalSum, err := applyAndSum(arr, ops)
        if err != nil {
            fmt.Fprintf(os.Stderr, "test %d failed: %v\n", i+1, err)
            os.Exit(1)
        }
        if finalSum != want {
            fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%s\nexpected final sum:\n%d\nobtained final sum:\n%d\noutput:\n%s\n", i+1, in, want, finalSum, strings.TrimSpace(raw))
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}
