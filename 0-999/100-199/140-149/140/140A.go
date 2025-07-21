package main

import (
    "fmt"
    "math"
)

func main() {
    var n int
    var Rint, rint int
    if _, err := fmt.Scan(&n, &Rint, &rint); err != nil {
        return
    }
    R := float64(Rint)
    r := float64(rint)
    // Single plate: fits if table radius >= plate radius
    if n == 1 {
        if R+1e-12 >= r {
            fmt.Println("YES")
        } else {
            fmt.Println("NO")
        }
        return
    }
    // For two or more plates, need table diameter >= 2*plate diameter
    if R < 2*r {
        fmt.Println("NO")
        return
    }
    // Plates centers lie on circle of radius (R - r)
    // Check minimal distance between centers: 2*(R - r)*sin(pi/n) >= 2*r
    angle := math.Pi / float64(n)
    if math.Sin(angle)*(R-r)+1e-12 >= r {
        fmt.Println("YES")
    } else {
        fmt.Println("NO")
    }
}
