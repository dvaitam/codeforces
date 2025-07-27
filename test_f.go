package main
import "fmt"

func f(n uint64) uint64 {
    if n == 0 {
        return 0
    }
    if n == 1 {
        return 1
    }
    if n == 2 {
        return 1
    }
    if n == 3 {
        return 2
    }
    if n&1 == 1 {
        if n&3 == 1 {
            return 2 + f((n-3)/2)
        } else {
            return 1 + f((n-1)/2)
        }
    } else {
        if n&3 == 2 {
            return n - f(n/2)
        } else {
            return n - 1 - f(n/2-1)
        }
    }
}

func main() {
    for i := uint64(1); i <= 30; i++ {
        fmt.Printf("%d %d\n", i, f(i))
    }
}
