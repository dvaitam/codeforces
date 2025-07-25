package main

import (
    "bufio"
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "strings"
)

type testCaseE struct{input string}

func parseTestcases(path string)([]testCaseE,error){
    f,err:=os.Open(path)
    if err!=nil{return nil,err}
    defer f.Close()
    in:=bufio.NewReader(f)
    var T int
    if _,err:=fmt.Fscan(in,&T); err!=nil {return nil,err}
    cases:=make([]testCaseE,T)
    for i:=0;i<T;i++{
        var n,m,t,x int
        fmt.Fscan(in,&n,&m,&t,&x)
        var sb strings.Builder
        fmt.Fprintf(&sb,"%d %d %d %d\n",n,m,t,x)
        for j:=0;j<m;j++{
            var u,v,w int
            fmt.Fscan(in,&u,&v,&w)
            fmt.Fprintf(&sb,"%d %d %d\n",u,v,w)
            for k:=0;k<t;k++{
                var p int
                fmt.Fscan(in,&p)
                if k>0{sb.WriteByte(' ')}
                sb.WriteString(fmt.Sprint(p))
            }
            sb.WriteByte('\n')
        }
        cases[i]=testCaseE{input:sb.String()}
    }
    return cases,nil
}

func run(bin,input string)(string,error){
    var cmd *exec.Cmd
    if strings.HasSuffix(bin, ".go") {
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
        fmt.Println("usage: go run verifierE.go /path/to/binary")
        os.Exit(1)
    }
    bin:=os.Args[1]
    cases,err:=parseTestcases("testcasesE.txt")
    if err!=nil{fmt.Fprintf(os.Stderr,"failed to parse testcases: %v\n",err);os.Exit(1)}
    for idx,tc:=range cases{
        exp,err:=run("553E.go",tc.input)
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

