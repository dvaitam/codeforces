package main

import (
    "bytes"
    "context"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "strings"
    "time"
)

const mod int64 = 998244353

func dfs(idx int, letters []byte, N int) ([]byte, int64) {
    if idx*2 > N {
        return []byte{letters[idx]}, 1
    }
    ls, lc := dfs(idx*2, letters, N)
    rs, rc := dfs(idx*2+1, letters, N)
    if bytes.Compare(ls, rs) > 0 {
        ls, rs = rs, ls
        lc, rc = rc, lc
    }
    count := (lc * rc) % mod
    if !bytes.Equal(ls, rs) {
        count = (count * 2) % mod
    }
    res := make([]byte, 1+len(ls)+len(rs))
    res[0] = letters[idx]
    copy(res[1:], ls)
    copy(res[1+len(ls):], rs)
    return res, count
}

func solve(n int, s string) int64 {
    N := (1 << n) - 1
    letters := make([]byte, N+1)
    for i := 1; i <= N; i++ {
        letters[i] = s[i-1]
    }
    _, ans := dfs(1, letters, N)
    return ans % mod
}

func runBinary(bin, input string) (string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    cmd := exec.CommandContext(ctx, bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = os.Stderr
    err := cmd.Run()
    return strings.TrimSpace(out.String()), err
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("usage: go run verifierE.go /path/to/binary")
        return
    }
    bin := os.Args[1]
    r := rand.New(rand.NewSource(5))
    tests := 100
    for i := 0; i < tests; i++ {
        n := r.Intn(3) + 2 // 2..4
        N := (1 << n) - 1
        sb := make([]byte, N)
        for j := 0; j < N; j++ {
            if r.Intn(2) == 0 {
                sb[j] = 'A'
            } else {
                sb[j] = 'B'
            }
        }
        s := string(sb)
        input := fmt.Sprintf("%d\n%s\n", n, s)
        expected := fmt.Sprintf("%d", solve(n, s))
        out, err := runBinary(bin, input)
        if err != nil {
            fmt.Printf("test %d runtime error: %v\n", i+1, err)
            return
        }
        out = strings.TrimSpace(out)
        if out != expected {
            fmt.Printf("test %d failed: n=%d s=%s expected=%s got=%s\n", i+1, n, s, expected, out)
            return
        }
    }
    fmt.Println("All tests passed")
}
