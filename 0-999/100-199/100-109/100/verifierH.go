package main

import (
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "strings"
)

type TestCase struct {
    input string
    ans   string
}

func isValid(grid [10]string) bool {
    var g [10][10]byte
    for i := 0; i < 10; i++ {
        for j := 0; j < 10; j++ {
            if j < len(grid[i]) && grid[i][j] == '*' {
                g[i][j] = '*'
            } else {
                g[i][j] = '0'
            }
        }
    }
    diag := [4][2]int{{-1, -1}, {-1, 1}, {1, -1}, {1, 1}}
    for i := 0; i < 10; i++ {
        for j := 0; j < 10; j++ {
            if g[i][j] != '*' {
                continue
            }
            for _, d := range diag {
                ni, nj := i+d[0], j+d[1]
                if ni >= 0 && ni < 10 && nj >= 0 && nj < 10 && g[ni][nj] == '*' {
                    return false
                }
            }
        }
    }
    var vis [10][10]bool
    counts := map[int]int{1: 0, 2: 0, 3: 0, 4: 0}
    dirs := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
    for i := 0; i < 10; i++ {
        for j := 0; j < 10; j++ {
            if g[i][j] == '*' && !vis[i][j] {
                stack := [][2]int{{i, j}}
                vis[i][j] = true
                var cells [][2]int
                cells = append(cells, [2]int{i, j})
                for len(stack) > 0 {
                    cur := stack[len(stack)-1]
                    stack = stack[:len(stack)-1]
                    ci, cj := cur[0], cur[1]
                    for _, d := range dirs {
                        ni, nj := ci+d[0], cj+d[1]
                        if ni >= 0 && ni < 10 && nj >= 0 && nj < 10 && g[ni][nj] == '*' && !vis[ni][nj] {
                            vis[ni][nj] = true
                            stack = append(stack, [2]int{ni, nj})
                            cells = append(cells, [2]int{ni, nj})
                        }
                    }
                }
                size := len(cells)
                if size < 1 || size > 4 {
                    return false
                }
                minI, maxI, minJ, maxJ := 10, -1, 10, -1
                setC := make(map[[2]int]bool)
                for _, c := range cells {
                    setC[c] = true
                    if c[0] < minI {
                        minI = c[0]
                    }
                    if c[0] > maxI {
                        maxI = c[0]
                    }
                    if c[1] < minJ {
                        minJ = c[1]
                    }
                    if c[1] > maxJ {
                        maxJ = c[1]
                    }
                }
                if minI == maxI {
                    if maxJ-minJ+1 != size {
                        return false
                    }
                    for jj := minJ; jj <= maxJ; jj++ {
                        if !setC[[2]int{minI, jj}] {
                            return false
                        }
                    }
                } else if minJ == maxJ {
                    if maxI-minI+1 != size {
                        return false
                    }
                    for ii := minI; ii <= maxI; ii++ {
                        if !setC[[2]int{ii, minJ}] {
                            return false
                        }
                    }
                } else {
                    return false
                }
                counts[size]++
            }
        }
    }
    required := map[int]int{1: 4, 2: 3, 3: 2, 4: 1}
    for k, v := range required {
        if counts[k] != v {
            return false
        }
    }
    return true
}

func genBoard(r *rand.Rand) [10]string {
    var g [10]string
    for i := 0; i < 10; i++ {
        var row strings.Builder
        for j := 0; j < 10; j++ {
            if r.Intn(5) == 0 {
                row.WriteByte('*')
            } else {
                row.WriteByte('.')
            }
        }
        g[i] = row.String()
    }
    return g
}

func genTests() []TestCase {
    r := rand.New(rand.NewSource(8))
    cases := make([]TestCase, 100)
    for i := 0; i < 100; i++ {
        b := genBoard(r)
        valid := isValid(b)
        var sb strings.Builder
        sb.WriteString("1\n")
        for j := 0; j < 10; j++ {
            sb.WriteString(b[j])
            sb.WriteByte('\n')
        }
        ans := "NO"
        if valid {
            ans = "YES"
        }
        cases[i] = TestCase{input: sb.String(), ans: ans}
    }
    return cases
}

func run(bin, in string) (string, error) {
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(in)
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    return strings.TrimSpace(out.String()), err
}

func main() {
    if len(os.Args) < 2 {
        fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    tests := genTests()
    for i, tc := range tests {
        got, err := run(bin, tc.input)
        if err != nil {
            fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
            os.Exit(1)
        }
        if strings.TrimSpace(got) != tc.ans {
            fmt.Fprintf(os.Stderr, "test %d failed expected %q got %q\n", i+1, tc.ans, got)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}

