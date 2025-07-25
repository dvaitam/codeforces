package main
import(
    "bufio"
    "fmt"
    "os"
)
func main(){
    in:=bufio.NewScanner(os.Stdin)
    out:=bufio.NewWriter(os.Stdout)
    defer out.Flush()
    in.Scan()
    var t int
    fmt.Sscan(in.Text(),&t)
    for i:=0;i<t;i++{
        in.Scan()
        s:=in.Text()
        fmt.Fprintln(out,"g_"+s)
    }
}
