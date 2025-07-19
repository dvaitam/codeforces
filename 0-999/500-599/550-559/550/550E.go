package main

import "fmt"

func main() {
    var n int
    if _, err := fmt.Scan(&n); err != nil {
        return
    }
    a := make([]int, n)
    for i := 0; i < n; i++ {
        fmt.Scan(&a[i])
    }
    // If last element is not zero, no solution
    if a[n-1] != 0 {
        fmt.Print("NO")
        return
    }
    // Find first zero in the prefix
    flag := -1
    for i := 0; i < n-1; i++ {
        if a[i] == 0 {
            flag = i
            break
        }
    }
    // If the only zero in prefix is at the last position, no solution
    if flag == n-2 {
        fmt.Print("NO")
        return
    }
    // Construct and print the expression
    fmt.Println("YES")
    if flag >= 0 {
        if flag > 0 {
            fmt.Print("(")
            for i := 0; i < flag-1; i++ {
                fmt.Printf("%d->", a[i])
            }
            fmt.Printf("%d)->", a[flag-1])
        }
        fmt.Print("(0)->(")
        for i := flag + 1; i < n-1; i++ {
            fmt.Printf("%d->", a[i])
        }
        fmt.Printf("%d)", a[n-2])
        fmt.Print("->0")
    } else {
        for i := 0; i < n-1; i++ {
            fmt.Printf("%d->", a[i])
        }
        fmt.Print("0")
    }
}
