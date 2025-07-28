package main

import (
    "bufio"
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "strconv"
    "strings"
)

const mod int = 998244353

func encode(arr []int, k int) int {
    code := 0
    for _, v := range arr {
        code = code*k + (v - 1)
    }
    return code
}

func decode(code, n, k int) []int {
    arr := make([]int, n)
    for i := n - 1; i >= 0; i-- {
        arr[i] = code%k + 1
        code /= k
    }
    return arr
}

func applyF(arr []int, r int) []int {
    n := len(arr)
    res := make([]int, n)
    idx := 0
    for _, v := range arr {
        for i := 0; i < r && idx < n; i++ {
            res[idx] = v
            idx++
        }
        if idx >= n {
            break
        }
    }
    for idx < n {
        res[idx] = arr[len(arr)-1]
        idx++
    }
    return res
}

func applyG(arr []int, x, y int) []int {
    res := make([]int, len(arr))
    for i, v := range arr {
        if v == x {
            res[i] = y
        } else if v == y {
            res[i] = x
        } else {
            res[i] = v
        }
    }
    return res
}

func reachable(arr []int, n, k int, total int) int {
    start := encode(arr, k)
    visited := make([]bool, total)
    visited[start] = true
    mask := 1 << start
    queue := [][]int{arr}
    for len(queue) > 0 {
        cur := queue[0]
        queue = queue[1:]
        // F operations
        for r := 1; r <= n+1; r++ {
            nxt := applyF(cur, r)
            code := encode(nxt, k)
            if !visited[code] {
                visited[code] = true
                mask |= 1 << code
                queue = append(queue, nxt)
            }
        }
        // G operations
        for x := 1; x <= k; x++ {
            for y := x + 1; y <= k; y++ {
                nxt := applyG(cur, x, y)
                code := encode(nxt, k)
                if !visited[code] {
                    visited[code] = true
                    mask |= 1 << code
                    queue = append(queue, nxt)
                }
            }
        }
    }
    return mask
}

func minCover(n, k int) int {
    total := 1
    for i := 0; i < n; i++ {
        total *= k
    }
    masks := make([]int, total)
    for code := 0; code < total; code++ {
        arr := decode(code, n, k)
        masks[code] = reachable(arr, n, k, total)
    }
    full := 1<<total - 1
    dp := make([]int, full+1)
    for i := range dp {
        dp[i] = 1 << 30
    }
    dp[0] = 0
    for mask := 0; mask <= full; mask++ {
        if dp[mask] == 1<<30 {
            continue
        }
        for _, m := range masks {
            nm := mask | m
            if dp[mask]+1 < dp[nm] {
                dp[nm] = dp[mask] + 1
            }
        }
    }
    return dp[full]
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    file, err := os.Open("testcasesF.txt")
    if err != nil {
        fmt.Fprintln(os.Stderr, "could not open testcasesF.txt:", err)
        os.Exit(1)
    }
    defer file.Close()
    scanner := bufio.NewScanner(file)
    if !scanner.Scan() {
        fmt.Fprintln(os.Stderr, "empty testcases file")
        os.Exit(1)
    }
    var t int
    fmt.Sscan(scanner.Text(), &t)
    for caseNum := 1; caseNum <= t; caseNum++ {
        if !scanner.Scan() {
            fmt.Fprintf(os.Stderr, "expected %d cases, got %d\n", t, caseNum-1)
            os.Exit(1)
        }
        var n, k int
        fmt.Sscan(scanner.Text(), &n, &k)
        exp := minCover(n, k) % mod
        input := fmt.Sprintf("1\n%d %d\n", n, k)
        cmd := exec.Command(bin)
        cmd.Stdin = strings.NewReader(input)
        var out bytes.Buffer
        cmd.Stdout = &out
        cmd.Stderr = &out
        if err := cmd.Run(); err != nil {
            fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s", caseNum, err, out.String())
            os.Exit(1)
        }
        gotStr := strings.TrimSpace(out.String())
        val, err := strconv.Atoi(gotStr)
        if err != nil || val%mod != exp {
            fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\n", caseNum, exp, gotStr)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed\n", t)
}

