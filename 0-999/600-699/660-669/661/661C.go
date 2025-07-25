package main
import(
    "bufio"
    "fmt"
    "os"
)
func main(){
    in:=bufio.NewReader(os.Stdin)
    out:=bufio.NewWriter(os.Stdout)
    defer out.Flush()
    var t int
    if _,err:=fmt.Fscan(in,&t); err!=nil{ return }
    for i:=0;i<t;i++{
        var n int
        fmt.Fscan(in,&n)
        sum:=0
        for j:=0;j<n;j++{
            var x int
            fmt.Fscan(in,&x)
            sum+=x
        }
        fmt.Fprintln(out,sum)
    }
}
