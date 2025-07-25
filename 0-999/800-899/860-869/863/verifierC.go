package main

import (
    "bytes"
    "fmt"
    "math/rand"
    "os"
    "os/exec"
    "strings"
)

type matrices struct{
    A [4][4]int
    B [4][4]int
}

func expectedC(k int64, a, b int, mat matrices) (int64,int64) {
    nextA := mat.A
    nextB := mat.B
    visitedStep := [4][4]int64{}
    visitedA := [4][4]int64{}
    visitedB := [4][4]int64{}
    for i:=1;i<=3;i++{
        for j:=1;j<=3;j++{
            visitedStep[i][j] = -1
        }
    }
    var scoreA, scoreB int64
    var step int64
    curA, curB := a, b
    for step < k {
        if visitedStep[curA][curB] != -1 {
            prev := visitedStep[curA][curB]
            cycleLen := step - prev
            cycleScoreA := scoreA - visitedA[curA][curB]
            cycleScoreB := scoreB - visitedB[curA][curB]
            if cycleLen > 0 {
                times := (k - step) / cycleLen
                if times > 0 {
                    scoreA += cycleScoreA * times
                    scoreB += cycleScoreB * times
                    step += cycleLen * times
                }
            }
        }
        if step >= k {
            break
        }
        visitedStep[curA][curB] = step
        visitedA[curA][curB] = scoreA
        visitedB[curA][curB] = scoreB
        if curA != curB {
            if curA == 1 && curB == 3 || curA == 2 && curB == 1 || curA == 3 && curB == 2 {
                scoreA++
            } else {
                scoreB++
            }
        }
        step++
        na := nextA[curA][curB]
        nb := nextB[curA][curB]
        curA, curB = na, nb
    }
    return scoreA, scoreB
}

func genTestsC() []string {
    rand.Seed(3)
    tests := make([]string, 0, 100)
    for len(tests) < 100 {
        k := rand.Int63n(50) + 1
        a := rand.Intn(3) + 1
        b := rand.Intn(3) + 1
        matA := [4][4]int{}
        matB := [4][4]int{}
        var sb strings.Builder
        sb.WriteString(fmt.Sprintf("%d %d %d\n", k, a, b))
        for i:=1;i<=3;i++{
            for j:=1;j<=3;j++{
                matA[i][j] = rand.Intn(3)+1
                sb.WriteString(fmt.Sprintf("%d ", matA[i][j]))
            }
            sb.WriteString("\n")
        }
        for i:=1;i<=3;i++{
            for j:=1;j<=3;j++{
                matB[i][j] = rand.Intn(3)+1
                sb.WriteString(fmt.Sprintf("%d ", matB[i][j]))
            }
            sb.WriteString("\n")
        }
        tests = append(tests, sb.String())
    }
    return tests
}

func runBinary(bin string, input string) (string, error) {
    cmd := exec.Command(bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = os.Stderr
    err := cmd.Run()
    return strings.TrimSpace(out.String()), err
}

func parseMatrices(lines []string) matrices {
    var res matrices
    idx := 0
    for i:=1;i<=3;i++{
        parts := strings.Fields(lines[idx])
        for j:=1;j<=3;j++{
            fmt.Sscanf(parts[j-1], "%d", &res.A[i][j])
        }
        idx++
    }
    for i:=1;i<=3;i++{
        parts := strings.Fields(lines[idx])
        for j:=1;j<=3;j++{
            fmt.Sscanf(parts[j-1], "%d", &res.B[i][j])
        }
        idx++
    }
    return res
}

func main(){
    if len(os.Args) != 2 {
        fmt.Fprintf(os.Stderr, "Usage: go run verifierC.go <binary>\n")
        os.Exit(1)
    }
    bin := os.Args[1]
    tests := genTestsC()
    for idx, t := range tests {
        lines := strings.Split(strings.TrimSpace(t), "\n")
        var k int64
        var a,b int
        fmt.Sscanf(lines[0], "%d %d %d", &k, &a, &b)
        mats := parseMatrices(lines[1:])
        sa,sb := expectedC(k,a,b,mats)
        want := fmt.Sprintf("%d %d", sa,sb)
        got, err := runBinary(bin, t)
        if err != nil {
            fmt.Printf("Test %d: runtime error: %v\n", idx+1, err)
            os.Exit(1)
        }
        if strings.TrimSpace(got) != want {
            fmt.Printf("Test %d failed.\nInput:\n%s\nExpected: %s\nGot: %s\n", idx+1, t, want, got)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed.\n", len(tests))
}

