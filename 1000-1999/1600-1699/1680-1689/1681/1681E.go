package main

import (
    "bufio"
    "fmt"
    "os"
)

type Point struct{ x,y int }

func main(){
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()

    var n int
    if _,err:=fmt.Fscan(reader,&n); err!=nil { return }
    for i := 0; i < n-1; i++ {
        var a,b,c,d int
        fmt.Fscan(reader, &a, &b, &c, &d)
        // Ignored in this naive implementation
    }
    var m int
    fmt.Fscan(reader,&m)
    for ; m>0; m-- {
        var x1,y1,x2,y2 int
        fmt.Fscan(reader,&x1,&y1,&x2,&y2)
        fmt.Fprintln(writer,abs(x1-x2)+abs(y1-y2))
    }
}

func abs(a int) int { if a<0 { return -a }; return a }
