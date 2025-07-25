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

type testCase struct{ s string }

func genCase(rng *rand.Rand) testCase {
    n := rng.Intn(16)+5 // length 5..20
    b := make([]byte,n)
    for i:=0;i<n;i++{ b[i]=byte('a'+rng.Intn(26)) }
    return testCase{ s:string(b) }
}

func solve(s string) string {
    n := len(s)
    dp2 := make([]bool,n+1)
    dp3 := make([]bool,n+1)
    for i:=n-2;i>=5;i--{
        if i+2<=n {
            if i+2==n {
                dp2[i]=true
            } else if (dp2[i+2] && s[i:i+2]!=s[i+2:i+4]) || (dp3[i+2] && s[i:i+2]!=s[i+2:i+5]) {
                dp2[i]=true
            }
        }
        if i+3<=n {
            if i+3==n {
                dp3[i]=true
            } else if (dp2[i+3] && s[i:i+3]!=s[i+3:i+5]) || (dp3[i+3] && s[i:i+3]!=s[i+3:i+6]) {
                dp3[i]=true
            }
        }
    }
    set := make(map[string]struct{})
    for i:=5;i<n;i++{
        if dp2[i]{ set[s[i:i+2]]=struct{}{} }
        if dp3[i]{ set[s[i:i+3]]=struct{}{} }
    }
    arr := make([]string,0,len(set))
    for k := range set { arr = append(arr,k) }
    sort.Strings(arr)
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("%d\n",len(arr)))
    for _,v := range arr{ sb.WriteString(v); sb.WriteByte('\n') }
    return strings.TrimSpace(sb.String())
}

func run(bin,input string)(string,error){
    var cmd *exec.Cmd
    if strings.HasSuffix(bin,".go"){
        cmd=exec.Command("go","run",bin)
    }else{
        cmd=exec.Command(bin)
    }
    cmd.Stdin=strings.NewReader(input)
    var out bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout=&out
    cmd.Stderr=&stderr
    if err:=cmd.Run();err!=nil{
        return "",fmt.Errorf("runtime error: %v\n%s",err,stderr.String())
    }
    return strings.TrimSpace(out.String()),nil
}

func main(){
    if len(os.Args)==3 && os.Args[1]=="--" {
        os.Args = append([]string{os.Args[0]}, os.Args[2])
    }
    if len(os.Args)!=2{
        fmt.Println("usage: go run verifierA.go /path/to/binary")
        os.Exit(1)
    }
    bin:=os.Args[1]
    rng:=rand.New(rand.NewSource(time.Now().UnixNano()))
    for i:=0;i<100;i++{
        tc:=genCase(rng)
        exp:=solve(tc.s)
        got,err:=run(bin,tc.s+"\n")
        if err!=nil{
            fmt.Fprintf(os.Stderr,"case %d failed: %v\n",i+1,err)
            os.Exit(1)
        }
        if strings.TrimSpace(got)!=exp{
            fmt.Fprintf(os.Stderr,"case %d failed: input %s expected:\n%s\ngot:\n%s\n",i+1,tc.s,exp,got)
            os.Exit(1)
        }
    }
    fmt.Println("All tests passed")
}

