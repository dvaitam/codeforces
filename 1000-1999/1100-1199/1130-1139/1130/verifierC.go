package main

import (
    "bytes"
    "container/list"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "strconv"
    "strings"
)

func bfs(start [2]int, grid [][]byte) ([][2]int, map[[2]int]bool) {
    n := len(grid)
    visited := make(map[[2]int]bool)
    q := list.New()
    q.PushBack(start)
    visited[start] = true
    var comp [][2]int
    dirs := [][2]int{{1,0},{-1,0},{0,1},{0,-1}}
    for q.Len() > 0 {
        e := q.Front()
        q.Remove(e)
        p := e.Value.([2]int)
        comp = append(comp, p)
        for _, d := range dirs {
            nx, ny := p[0]+d[0], p[1]+d[1]
            if nx<0 || nx>=n || ny<0 || ny>=n {continue}
            if visited[[2]int{nx,ny}] {continue}
            if grid[nx][ny] == '1' {continue}
            visited[[2]int{nx,ny}] = true
            q.PushBack([2]int{nx,ny})
        }
    }
    return comp, visited
}

func solve(input string) string {
    reader := strings.NewReader(strings.TrimSpace(input))
    var n int
    fmt.Fscan(reader, &n)
    var sx, sy, ex, ey int
    fmt.Fscan(reader, &sx, &sy)
    fmt.Fscan(reader, &ex, &ey)
    sx--; sy--; ex--; ey--
    grid := make([][]byte, n)
    for i := 0; i < n; i++ {
        var row string
        fmt.Fscan(reader, &row)
        grid[i] = []byte(row)
    }
    comp1, vis1 := bfs([2]int{sx, sy}, grid)
    if vis1[[2]int{ex, ey}] {
        return "0"
    }
    comp2, _ := bfs([2]int{ex, ey}, grid)
    ans := 1<<60
    for _, a := range comp1 {
        for _, b := range comp2 {
            dx := a[0]-b[0]
            dy := a[1]-b[1]
            d := dx*dx + dy*dy
            if d < ans {
                ans = d
            }
        }
    }
    return strconv.Itoa(ans)
}

type test struct{ input, expected string }

func randGrid(n int) (grid []string, sx, sy, ex, ey int) {
    for {
        grid = make([]string, n)
        for i := 0; i < n; i++ {
            var sb strings.Builder
            for j := 0; j < n; j++ {
                if rand.Intn(4) == 0 {
                    sb.WriteByte('1')
                } else {
                    sb.WriteByte('0')
                }
            }
            grid[i] = sb.String()
        }
        sx, sy = rand.Intn(n), rand.Intn(n)
        ex, ey = rand.Intn(n), rand.Intn(n)
        if grid[sx][sy]=='0' && grid[ex][ey]=='0' {
            break
        }
    }
    return
}

func generateTests() []test {
    rand.Seed(1130)
    var tests []test
    for len(tests) < 5 {
        n := len(tests) + 1
        grid, sx, sy, ex, ey := randGrid(n)
        var sb strings.Builder
        fmt.Fprintf(&sb, "%d\n%d %d\n%d %d\n", n, sx+1, sy+1, ex+1, ey+1)
        for _, row := range grid {
            fmt.Fprintln(&sb, row)
        }
        inp := sb.String()
        tests = append(tests, test{inp, solve(inp)})
    }
    for len(tests) < 100 {
        n := rand.Intn(6) + 1
        grid, sx, sy, ex, ey := randGrid(n)
        var sb strings.Builder
        fmt.Fprintf(&sb, "%d\n%d %d\n%d %d\n", n, sx+1, sy+1, ex+1, ey+1)
        for _, row := range grid {
            fmt.Fprintln(&sb, row)
        }
        inp := sb.String()
        tests = append(tests, test{inp, solve(inp)})
    }
    return tests
}

func runBinary(bin, input string) (string, error) {
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    err := cmd.Run()
    return strings.TrimSpace(out.String()), err
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    tests := generateTests()
    for i, t := range tests {
        got, err := runBinary(bin, t.input)
        if err != nil {
            fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
            os.Exit(1)
        }
        if got != strings.TrimSpace(t.expected) {
            fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected: %s\nGot: %s\n", i+1, t.input, strings.TrimSpace(t.expected), got)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed.\n", len(tests))
}

