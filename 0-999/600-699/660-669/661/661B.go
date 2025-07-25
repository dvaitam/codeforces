package main
import(
    "bufio"
    "fmt"
    "os"
)
func season(m int, w *bufio.Writer){
    var s string
    switch m {
    case 12,1,2:
        s = "Winter"
    case 3,4,5:
        s = "Spring"
    case 6,7,8:
        s = "Summer"
    default:
        s = "Autumn"
    }
    w.WriteString(s)
}
func main(){
    in:=bufio.NewReader(os.Stdin)
    out:=bufio.NewWriter(os.Stdout)
    defer out.Flush()
    var t int
    if _,err:=fmt.Fscan(in,&t); err!=nil{ return }
    for i:=0;i<t;i++{
        var m int
        fmt.Fscan(in,&m)
        season(m, out)
        out.WriteByte('\n')
    }
}
