package main

import (
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "strconv"
    "strings"
    "time"
)

type Test struct{
    q int
    x int
    ys []int
}

func (t Test) Input() string {
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("%d %d\n", t.q, t.x))
    for _, v := range t.ys {
        sb.WriteString(fmt.Sprintf("%d\n", v))
    }
    return sb.String()
}

func expected(t Test) string {
    cnt := make([]int, t.x)
    mex := 0
    var sb strings.Builder
    for _, y := range t.ys {
        cnt[y%t.x]++
        for cnt[mex%t.x] > 0 {
            cnt[mex%t.x]--
            mex++
        }
        sb.WriteString(fmt.Sprintf("%d\n", mex))
    }
    return strings.TrimSpace(sb.String())
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
    x := rng.Intn(5)+1
    q := rng.Intn(20)+1
    ys := make([]int,q)
    for i:=0;i<q;i++ { ys[i] = rng.Intn(30) }
    return Test{q,x,ys}
}

func main(){
    if len(os.Args)!=2 {
        fmt.Println("Usage: go run verifierD.go /path/to/binary")
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
        // candidate output may contain spaces or newline separated
        tokens := strings.Fields(got)
        got = strings.Join(tokens, "\n")
        if got != exp {
            fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\nGot:\n%s\n", i+1, tc.Input(), exp, strings.Join(tokens, " "))
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed\n", cases)
}

