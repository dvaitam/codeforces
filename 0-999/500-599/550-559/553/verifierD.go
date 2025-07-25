package main

import (
    "bufio"
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "strings"
)

type testCaseD struct{input string}

func parseTestcases(path string) ([]testCaseD,error){
    f,err:=os.Open(path)
    if err!=nil{return nil,err}
    defer f.Close()
    in:=bufio.NewReader(f)
    var T int
    if _,err:=fmt.Fscan(in,&T); err!=nil { return nil,err }
    cases:=make([]testCaseD,T)
    for i:=0;i<T;i++{
        var N,M,K int
        fmt.Fscan(in,&N,&M,&K)
        var sb strings.Builder
        fmt.Fprintf(&sb,"%d %d %d\n",N,M,K)
        if K>0{
            for j:=0;j<K;j++{
                var x int
                fmt.Fscan(in,&x)
                if j>0{sb.WriteByte(' ')}
                sb.WriteString(fmt.Sprint(x))
            }
            sb.WriteByte('\n')
        } else { sb.WriteByte('\n') }
        for j:=0;j<M;j++{
            var u,v int
            fmt.Fscan(in,&u,&v)
            fmt.Fprintf(&sb,"%d %d\n",u,v)
        }
        cases[i]=testCaseD{input:sb.String()}
    }
    return cases,nil
}

func run(bin,input string)(string,error){
    var cmd *exec.Cmd
    if strings.HasSuffix(bin, ".go"){
        cmd=exec.Command("go","run",bin)
    }else{ cmd=exec.Command(bin) }
    cmd.Stdin=strings.NewReader(input)
    var out,errb bytes.Buffer
    cmd.Stdout=&out
    cmd.Stderr=&errb
    if err:=cmd.Run(); err!=nil{
        return "",fmt.Errorf("runtime error: %v\n%s",err,errb.String())
    }
    return strings.TrimSpace(out.String()),nil
}

func main(){
    if len(os.Args)!=2{
        fmt.Println("usage: go run verifierD.go /path/to/binary")
        os.Exit(1)
    }
    bin:=os.Args[1]
    cases,err:=parseTestcases("testcasesD.txt")
    if err!=nil{fmt.Fprintf(os.Stderr,"failed to parse testcases: %v\n",err);os.Exit(1)}
    for idx,tc:=range cases{
        exp,err:=run("553D.go",tc.input)
        if err!=nil{
            fmt.Fprintf(os.Stderr,"oracle failure on case %d: %v\n",idx+1,err)
            os.Exit(1)
        }
        got,err:=run(bin,tc.input)
        if err!=nil{
            fmt.Fprintf(os.Stderr,"case %d failed: %v\n",idx+1,err)
            os.Exit(1)
        }
        if strings.TrimSpace(got)!=strings.TrimSpace(exp){
            fmt.Fprintf(os.Stderr,"case %d failed: expected %s got %s\n",idx+1,exp,got)
            os.Exit(1)
        }
    }
    fmt.Printf("All %d tests passed\n", len(cases))
}

