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
        var x int
        fmt.Fscan(in,&x)
        if x%3==0 || x%5==0 { fmt.Fprintln(out,"YES") } else { fmt.Fprintln(out,"NO") }
    }
}
