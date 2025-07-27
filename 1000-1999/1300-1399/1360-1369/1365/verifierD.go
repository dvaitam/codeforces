package main

import (
    "bufio"
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "strings"
)

type test struct{ input, expected string }

func runBin(bin, input string) (string, error) {
    var cmd *exec.Cmd
    if strings.HasSuffix(bin, ".go") {
        cmd = exec.Command("go", "run", bin)
    } else { cmd = exec.Command(bin) }
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    var errb bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &errb
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
    }
    return strings.TrimSpace(out.String()), nil
}

func solveD(input string) string {
    in := bufio.NewReader(strings.NewReader(input))
    var T int
    fmt.Fscan(in, &T)
    var out strings.Builder
    dirs := [][2]int{{1,0},{-1,0},{0,1},{0,-1}}
    for ; T>0; T-- {
        var n,m int
        fmt.Fscan(in, &n, &m)
        grid := make([][]byte, n)
        for i:=0;i<n;i++ {
            var s string
            fmt.Fscan(in,&s)
            grid[i] = []byte(s)
        }
        possible := true
        for i:=0;i<n;i++ {
            for j:=0;j<m;j++ {
                if grid[i][j]=='B' {
                    for _,d := range dirs {
                        nx,ny := i+d[0], j+d[1]
                        if nx>=0 && nx<n && ny>=0 && ny<m {
                            switch grid[nx][ny] {
                            case 'G': possible=false
                            case '.': grid[nx][ny]='#'
                            }
                        }
                    }
                }
            }
        }
        if !possible {
            out.WriteString("No\n")
            continue
        }
        visited := make([][]bool, n)
        for i:=range visited { visited[i]=make([]bool,m) }
        q := make([][2]int,0)
        if grid[n-1][m-1] != '#' {
            q = append(q,[2]int{n-1,m-1})
            visited[n-1][m-1]=true
        }
        for head:=0; head<len(q); head++ {
            x,y := q[head][0], q[head][1]
            for _,d := range dirs {
                nx,ny := x+d[0], y+d[1]
                if nx>=0 && nx<n && ny>=0 && ny<m && !visited[nx][ny] && grid[nx][ny] != '#' && grid[nx][ny] != 'B' {
                    visited[nx][ny]=true
                    q = append(q,[2]int{nx,ny})
                }
            }
        }
        ok := true
        for i:=0;i<n;i++ {
            for j:=0;j<m;j++ {
                if grid[i][j]=='G' && !visited[i][j] { ok=false }
                if grid[i][j]=='B' && visited[i][j] { ok=false }
            }
        }
        if ok { out.WriteString("Yes\n") } else { out.WriteString("No\n") }
    }
    return strings.TrimSpace(out.String())
}

func genTests() []test {
    r := rand.New(rand.NewSource(1365))
    tests := make([]test,0,100)
    chars := []byte{'.', '#', 'G', 'B'}
    for len(tests) < 100 {
        n := r.Intn(4) + 1
        m := r.Intn(4) + 1
        var sb strings.Builder
        sb.WriteString("1\n")
        sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
        for i := 0; i < n; i++ {
            for j := 0; j < m; j++ {
                sb.WriteByte(chars[r.Intn(len(chars))])
            }
            sb.WriteByte('\n')
        }
        input := sb.String()
        expected := solveD(input)
        tests = append(tests, test{input, expected})
    }
    return tests
}

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    tests := genTests()
    for i, t := range tests {
        out, err := runBin(bin, t.input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, t.input)
            os.Exit(1)
        }
        if strings.TrimSpace(out) != strings.TrimSpace(t.expected) {
            fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i+1, t.input, t.expected, out)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed.\n", len(tests))
}

