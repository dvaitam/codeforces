package main
import (
    "fmt"
    "os"
)
var (
    b = [4]int{1, 2, 4, 8}
    a = [100]struct{ a, b, c, d int }{}
    n   int
    an  int
)
func printAns() {
    fmt.Printf("%d\n", an-1)
    for i := 2; i <= an; i++ {
        // lea e?x, [ ... ]
        reg := byte('a' + i - 1)
        fmt.Printf("lea e%cx, [", reg)
        if a[i].b != 0 {
            base := byte('a' + a[i].b - 1)
            fmt.Printf("e%cx + ", base)
        }
        if a[i].d != 0 {
            fmt.Printf("%d*", b[a[i].d])
        }
        src := byte('a' + a[i].c - 1)
        fmt.Printf("e%cx]\n", src)
    }
    os.Exit(0)
}
func dfs(x int) {
    if x == an {
        if a[x].a == n {
            printAns()
        }
        return
    }
    for i := 0; i <= x; i++ {
        for j := 1; j <= x; j++ {
            for k := 0; k < 4; k++ {
                A := a[i].a + b[k]*a[j].a
                if A <= a[x].a || A > n {
                    continue
                }
                a[x+1] = struct{ a, b, c, d int }{A, i, j, k}
                dfs(x + 1)
            }
        }
    }
}
func main() {
    if _, err := fmt.Scan(&n); err != nil {
        return
    }
    for an = 1; ; an++ {
        a[1] = struct{ a, b, c, d int }{1, 0, 0, 0}
        dfs(1)
    }
}
