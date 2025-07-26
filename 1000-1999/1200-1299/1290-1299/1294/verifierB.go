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

type point struct{ x, y int }

type Test struct{
    pts []point
}

func (t Test) Input() string {
    var sb strings.Builder
    sb.WriteString("1\n")
    sb.WriteString(fmt.Sprintf("%d\n", len(t.pts)))
    for _, p := range t.pts {
        sb.WriteString(fmt.Sprintf("%d %d\n", p.x, p.y))
    }
    return sb.String()
}

func expected(t Test) string {
    pts := append([]point(nil), t.pts...)
    sort.Slice(pts, func(i,j int) bool {
        if pts[i].x == pts[j].x { return pts[i].y < pts[j].y }
        return pts[i].x < pts[j].x
    })
    cx, cy := 0, 0
    path := make([]byte,0)
    for _, p := range pts {
        if p.y < cy {
            return "NO"
        }
        for cx < p.x { path = append(path,'R'); cx++ }
        for cy < p.y { path = append(path,'U'); cy++ }
    }
    return "YES\n"+string(path)
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
    n := rng.Intn(5)+1
    used := make(map[[2]int]bool)
    pts := make([]point,n)
    for i:=0;i<n;i++ {
        for {
            x := rng.Intn(6)
            y := rng.Intn(6)
            if x==0 && y==0 { continue }
            key := [2]int{x,y}
            if !used[key] {
                used[key]=true
                pts[i]=point{x,y}
                break
            }
        }
    }
    return Test{pts}
}

func main(){
    if len(os.Args)!=2 {
        fmt.Println("Usage: go run verifierB.go /path/to/binary")
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

