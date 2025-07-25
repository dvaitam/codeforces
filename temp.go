package main
import (
    "fmt"
    "math"
)
func main() {
    a:=1.0
    b:=1.0
    c:=1.0
    d:=1.0
    for i:=0;i<1000;i++ {
        a = math.Sqrt(7 - math.Sqrt(6 - a))
        b = math.Sqrt(7 - math.Sqrt(6 + b))
        c = math.Sqrt(7 + math.Sqrt(6 - c))
        d = math.Sqrt(7 + math.Sqrt(6 + d))
    }
    prod := a*b*c*d
    fmt.Printf("%.12f\n", a)
    fmt.Printf("%.12f\n", b)
    fmt.Printf("%.12f\n", c)
    fmt.Printf("%.12f\n", d)
    fmt.Printf("%.12f\n", prod)
}
