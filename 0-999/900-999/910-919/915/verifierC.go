package main

import (
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "sort"
    "strings"
    "time"
)

func solveC(aStr, bStr string) string {
    if len(aStr) < len(bStr) {
        digits := []byte(aStr)
        sort.Slice(digits, func(i, j int) bool { return digits[i] > digits[j] })
        return string(digits)
    }
    counts := make([]int, 10)
    for _, ch := range aStr {
        counts[ch-'0']++
    }
    res := make([]byte, len(bStr))
    var dfs func(pos int, limit bool) bool
    dfs = func(pos int, limit bool) bool {
        if pos == len(bStr) {
            return true
        }
        maxDigit := 9
        if limit {
            maxDigit = int(bStr[pos] - '0')
        }
        for d := maxDigit; d >= 0; d-- {
            if counts[d] == 0 {
                continue
            }
            counts[d]--
            res[pos] = byte('0' + d)
            if dfs(pos+1, limit && d == maxDigit) {
                return true
            }
            counts[d]++
        }
        return false
    }
    if dfs(0, true) {
        return string(res)
    }
    return ""
}

func generateC(rng *rand.Rand) (string, string, string) {
    l := rng.Intn(18) + 1
    digits := make([]byte, l)
    digits[0] = byte('1' + rng.Intn(9))
    for i := 1; i < l; i++ {
        digits[i] = byte('0' + rng.Intn(10))
    }
    a := string(digits)

    if rng.Intn(2) == 0 {
        // len(b) > len(a)
        bl := l + rng.Intn(3) + 1
        bDigits := make([]byte, bl)
        bDigits[0] = byte('1' + rng.Intn(9))
        for i := 1; i < bl; i++ {
            bDigits[i] = byte('0' + rng.Intn(10))
        }
        b := string(bDigits)
        return a, b, solveC(a, b)
    }
    // len(b) == len(a), choose b >= sorted digits ascending to ensure solution
    digitsB := make([]byte, l)
    for i := 0; i < l; i++ {
        digitsB[i] = byte('0' + rng.Intn(10))
    }
    b := string(digitsB)
    // to guarantee existence, set b to a large value by at least digits sorted ascending
    digitsAsc := []byte(a)
    sort.Slice(digitsAsc, func(i, j int) bool { return digitsAsc[i] < digitsAsc[j] })
    if strings.Compare(b, string(digitsAsc)) < 0 {
        b = string(digitsAsc)
    }
    return a, b, solveC(a, b)
}

func runCase(bin, input, exp string) error {
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("runtime error: %v\n%s", err, out.String())
    }
    got := strings.TrimSpace(out.String())
    if got != exp {
        return fmt.Errorf("expected %s got %s", exp, got)
    }
    return nil
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
        os.Exit(1)
    }
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))
    bin := os.Args[1]
    for i := 0; i < 100; i++ {
        a, b, exp := generateC(rng)
        input := fmt.Sprintf("%s\n%s\n", a, b)
        if err := runCase(bin, input, exp); err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}

