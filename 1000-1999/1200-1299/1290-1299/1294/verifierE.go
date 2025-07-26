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

type Test struct{
    n, m int
    mat [][]int
}

func (t Test) Input() string {
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("%d %d\n", t.n, t.m))
    for i := 0; i < t.n; i++ {
        for j := 0; j < t.m; j++ {
            if j > 0 { sb.WriteByte(' ') }
            sb.WriteString(fmt.Sprint(t.mat[i][j]))
        }
        sb.WriteByte('\n')
    }
    return sb.String()
}

func expected(t Test) string {
    n, m := t.n, t.m
    nm := n*m
    total := 0
    for col := 0; col < m; col++ {
        cnt := make([]int, n)
        for row := 0; row < n; row++ {
            val := t.mat[row][col]
            if val >= 1 && val <= nm && (val-1)%m == col {
                targetRow := (val - 1) / m
                shift := (row - targetRow + n) % n
                cnt[shift]++
            }
        }
        best := n + 1
        for shift := 0; shift < n; shift++ {
            moves := shift + (n - cnt[shift])
            if moves < best { best = moves }
        }
        total += best
    }
    return fmt.Sprint(total)
}

func runProg(bin,input string) (string,error){
    var cmd *exec.Cmd
    if strings.HasSuffix(bin,".go") { cmd = exec.Command("go","run",bin) } else { cmd = exec.Command(bin) }
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    var errBuf bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &errBuf
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
    }
    return strings.TrimSpace(out.String()), nil
}

func genTest(rng *rand.Rand) Test {
    n := rng.Intn(3)+1
    m := rng.Intn(3)+1
    nm := n*m
    mat := make([][]int, n)
    for i:=0;i<n;i++ {
        row := make([]int,m)
        for j:=0;j<m;j++ {
            row[j] = rng.Intn(nm*2)+1
        }
        mat[i] = row
    }
    return Test{n,m,mat}
}

func main(){
    if len(os.Args)!=2 {
        fmt.Println("Usage: go run verifierE.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))
    const cases = 100
    for i:=0;i<cases;i++ {
        tc := genTest(rng)
        exp := expected(tc)
        got, err := runProg(bin, tc.Input())
        if err != nil {
            fmt.Printf("case %d: %v\n", i+1, err)
            os.Exit(1)
        }
        if strings.TrimSpace(got) != exp {
            fmt.Printf("case %d failed\ninput:\n%sexpected: %s\nGot: %s\n", i+1, tc.Input(), exp, got)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed\n", cases)
}

