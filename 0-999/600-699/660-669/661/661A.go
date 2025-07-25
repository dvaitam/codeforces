package main
import (
    "bufio"
    "fmt"
    "os"
)
func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()
    var t int
    if _, err := fmt.Fscan(in, &t); err != nil { return }
    for i := 0; i < t; i++ {
        var n int64
        fmt.Fscan(in, &n)
        res := n*(n+1)/2 + 1
        fmt.Fprintln(out, res)
    }
}
