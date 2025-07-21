package main

import "fmt"

func main() {
    var grid [3]string
    for i := 0; i < 3; i++ {
        fmt.Scan(&grid[i])
    }
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            if grid[i][j] != grid[2-i][2-j] {
                fmt.Println("NO")
                return
            }
        }
    }
    fmt.Println("YES")
}
