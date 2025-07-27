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

type testCase struct{nums []int}

func (tc testCase) input() string {
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("%d\n", len(tc.nums)))
    for i, v := range tc.nums {
        if i>0 { sb.WriteByte(' ') }
        sb.WriteString(fmt.Sprintf("%d", v))
    }
    sb.WriteByte('\n')
    return sb.String()
}

func solve(nums []int) string {
    if len(nums)==0 { return "" }
    max:=nums[0]
    for _,v := range nums { if v>max { max=v } }
    return fmt.Sprintf("%d", max)
}

func randomCase(rng *rand.Rand) testCase {
    n := rng.Intn(10)+1
    nums:=make([]int,n)
    for i:=0;i<n;i++{ nums[i]=rng.Intn(2000001)-1000000 }
    return testCase{nums:nums}
}

func deterministicCases() []testCase {
    return []testCase{
        {nums:[]int{1}},
        {nums:[]int{-1,0,2}},
        {nums:[]int{5,5,5,5}},
    }
}

func runCase(bin string, tc testCase) error {
    input:=tc.input()
    cmd:=exec.Command(bin)
    cmd.Stdin = strings.NewReader(input)
    var out bytes.Buffer
    cmd.Stdout=&out
    cmd.Stderr=&out
    if err:=cmd.Run(); err!=nil {
        return fmt.Errorf("runtime error: %v\n%s", err, out.String())
    }
    result:=strings.TrimSpace(out.String())
    expect:=solve(tc.nums)
    if result!=expect {
        return fmt.Errorf("expected %s got %s", expect, result)
    }
    return nil
}

func main(){
    if len(os.Args)!=2 {
        fmt.Fprintln(os.Stderr,"usage: go run verifierB.go /path/to/binary")
        os.Exit(1)
    }
    bin:=os.Args[1]
    rng:=rand.New(rand.NewSource(time.Now().UnixNano()))
    tests:=deterministicCases()
    for len(tests)<100 { tests=append(tests, randomCase(rng)) }
    for i,tc:= range tests {
        if err:=runCase(bin,tc); err!=nil {
            fmt.Fprintf(os.Stderr,"case %d failed: %v\ninput:\n%s", i+1, err, tc.input())
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}

