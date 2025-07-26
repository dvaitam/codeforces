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
    n int64
}

func (t Test) Input() string {
    return fmt.Sprintf("1\n%d\n", t.n)
}

func expected(t Test) string {
    n0 := t.n
    var p1, p2, p3 int64
    for i := int64(2); i*i <= n0; i++ {
        if n0%i == 0 {
            p1 = i
            break
        }
    }
    if p1 == 0 {
        return "NO"
    }
    n1 := n0 / p1
    for j := p1 + 1; j*j <= n1; j++ {
        if n1%j == 0 {
            p2 = j
            break
        }
    }
    if p2 == 0 {
        return "NO"
    }
    p3 = n1 / p2
    if p3 <= 1 || p3 == p1 || p3 == p2 {
        return "NO"
    }
    return fmt.Sprintf("YES\n%d %d %d", p1, p2, p3)
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
    return Test{n: rng.Int63n(1000000000-2) + 2}
}

func main(){
    if len(os.Args)!=2 {
        fmt.Println("Usage: go run verifierC.go /path/to/binary")
        os.Exit(1)
    }
    bin := os.Args[1]
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))
    const cases = 100
    for i:=0;i<cases;i++ {
        tc := genTest(rng)
        exp := strings.TrimSpace(expected(tc))
        got, err := runProg(bin, tc.Input())
        if err != nil {
            fmt.Printf("case %d: %v\n", i+1, err)
            os.Exit(1)
        }
        if strings.TrimSpace(got) != exp {
            fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\nGot:\n%s\n", i+1, tc.Input(), exp, got)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed\n", cases)
}

