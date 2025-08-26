package main

import (
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
    "time"
)

func baseDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Dir(file)
}

func prepareBinary(path, tag string) (string, error) {
	if strings.HasSuffix(path, ".go") {
		bin := filepath.Join(os.TempDir(), tag+"_"+fmt.Sprint(time.Now().UnixNano()))
		cmd := exec.Command("go", "build", "-o", bin, path)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build %s: %v\n%s", path, err, out)
		}
		return bin, nil
	}
	return path, nil
}

func runBinary(path, input string) (string, error) {
    cmd := exec.Command(path)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    err := cmd.Run()
    return out.String(), err
}

func genCase(rng *rand.Rand) string {
	x := rng.Int63n(1000) + 1
	d := rng.Int63n(20) + 1
	return fmt.Sprintf("%d %d\n", x, d)
}

// parse candidate output and validate correctness for input X,d
func checkOutput(out string, X, d int64) error {
    tokens := strings.Fields(out)
    if len(tokens) == 0 {
        return fmt.Errorf("empty output")
    }
    n, err := strconv.ParseInt(tokens[0], 10, 64)
    if err != nil {
        return fmt.Errorf("first token not an integer: %v", err)
    }
    if n == -1 {
        return fmt.Errorf("reported -1 but a solution always exists")
    }
    if n < 1 || n > 10000 {
        return fmt.Errorf("n out of range: %d", n)
    }
    if int64(len(tokens)) < 1+n {
        return fmt.Errorf("expected %d numbers got %d", n, int64(len(tokens))-1)
    }
    a := make([]int64, n)
    for i := int64(0); i < n; i++ {
        v, err := strconv.ParseInt(tokens[1+i], 10, 64)
        if err != nil {
            return fmt.Errorf("invalid number: %v", err)
        }
        if v < 1 || v >= 1_000_000_000_000_000_000 {
            return fmt.Errorf("value out of bounds: %d", v)
        }
        a[i] = v
    }
    // Count subsequences with max-min < d using two pointers
    sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
    // precompute powers of 2 up to n, capped at X+1 to avoid overflow
    capVal := X + 1
    pow2 := make([]int64, n+1)
    pow2[0] = 1
    for i := int64(1); i <= n; i++ {
        v := pow2[i-1] * 2
        if v > capVal {
            v = capVal
        }
        pow2[i] = v
    }
    var total int64 = 0
    r := int64(0)
    for i := int64(0); i < n; i++ {
        if r < i {
            r = i
        }
        for r+1 < n && a[r+1]-a[i] < d {
            r++
        }
        // number of subsequences where a[i] is the minimal element
        add := pow2[r-i]
        if total > X-add {
            // early exit if exceeding X
            return fmt.Errorf("count exceeds X")
        }
        total += add
    }
    if total != X {
        return fmt.Errorf("count mismatch: got %d want %d", total, X)
    }
    return nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
        os.Exit(1)
    }
    candPath, err := prepareBinary(os.Args[1], "candC")
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))
    for i := 0; i < 100; i++ {
        input := genCase(rng)
        got, err := runBinary(candPath, input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "candidate runtime error on case %d: %v\ninput:\n%s", i+1, err, input)
            os.Exit(1)
        }
        // parse X, d from input
        var X, d int64
        fmt.Sscan(strings.TrimSpace(input), &X, &d)
        if err := checkOutput(strings.TrimSpace(got), X, d); err != nil {
            fmt.Printf("case %d failed\ninput:\n%sexpected: valid construction with X=%d d=%d\ngot:\n%s\nerror: %v\n", i+1, input, X, d, got, err)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}
