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
        var a,b,c,d int
        fmt.Fscan(in,&a,&b,&c,&d)
        // matrix [[a b],[c d]] rotate clockwise -> [[c a],[d b]]
        fmt.Fprintf(out,"%d %d %d %d\n", c,a,d,b)
    }
}
