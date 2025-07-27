package main

import (
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "strings"
    "time"
)

const mod = 1000000007

// run executes the binary with the given input and returns trimmed stdout.
func run(bin, input string) (string, error) {
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

// genCase generates a random small test case for problem C.
func genCase(r *rand.Rand) (string, int, int, [][2]int) {
    n := r.Intn(2) + 2        // 2..3
    d := r.Intn(2) + 1        // 1..2
    edges := make([][2]int, n-1)
    for i := 2; i <= n; i++ {
        p := r.Intn(i-1) + 1
        edges[i-2] = [2]int{i, p}
    }
    var sb strings.Builder
    fmt.Fprintf(&sb, "%d %d\n", n, d)
    for _, e := range edges {
        fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
    }
    return sb.String(), n, d, edges
}

// firstPlayerWins determines if the starting player wins with given portal placement.
func firstPlayerWins(n, d int, edges [][2]int, portals [][2]int) bool {
    total := (d+1) * n
    adj := make([][]int, total)
    for i := 0; i <= d; i++ {
        for _, e := range edges {
            u := i*n + (e[0]-1)
            v := i*n + (e[1]-1)
            adj[u] = append(adj[u], v)
            adj[v] = append(adj[v], u)
        }
    }
    for i := 0; i < d; i++ {
        a := portals[i][0] - 1
        b := portals[i][1] - 1
        u := i*n + a
        v := (i+1)*n + b
        adj[u] = append(adj[u], v)
    }
    type state struct {
        cur  int
        mask uint64
    }
    memo := make(map[state]bool)
    var dfs func(int, uint64) bool
    dfs = func(cur int, mask uint64) bool {
        st := state{cur, mask}
        if v, ok := memo[st]; ok {
            return v
        }
        for _, nx := range adj[cur] {
            if mask&(1<<uint(nx)) == 0 {
                if !dfs(nx, mask|1<<uint(nx)) {
                    memo[st] = true
                    return true
                }
            }
        }
        memo[st] = false
        return false
    }
    return dfs(0, 1)
}

// bruteForce counts winning portal placements via enumeration.
func bruteForce(n, d int, edges [][2]int) int {
    portals := make([][2]int, d)
    var ans int
    var rec func(int)
    rec = func(i int) {
        if i == d {
            if firstPlayerWins(n, d, edges, portals) {
                ans++
            }
            return
        }
        for a := 1; a <= n; a++ {
            for b := 1; b <= n; b++ {
                portals[i] = [2]int{a, b}
                rec(i + 1)
            }
        }
    }
    rec(0)
    return ans % mod
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("usage: go run verifierC.go /path/to/binary")
        os.Exit(1)
    }
    userBin := os.Args[1]
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    const tests = 100
    for i := 0; i < tests; i++ {
        input, n, d, edges := genCase(r)
        expect := bruteForce(n, d, edges)
        gotStr, err := run(userBin, input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "test %d: %v\ninput:\n%s", i+1, err, input)
            os.Exit(1)
        }
        var got int
        fmt.Sscanf(gotStr, "%d", &got)
        if got != expect {
            fmt.Printf("test %d failed\ninput:\n%sexpected: %d\ngot: %d\n", i+1, input, expect, got)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed\n", tests)
}

