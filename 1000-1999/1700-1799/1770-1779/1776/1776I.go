package main

import (
    "bufio"
    "fmt"
    "os"
)

type Point struct{ x, y float64 }

func areaTri(a, b, c Point) float64 {
    v := a.x*(b.y-c.y) + b.x*(c.y-a.y) + c.x*(a.y-b.y)
    if v < 0 { v = -v }
    return v/2
}

func polygonArea(poly []Point) float64 {
    s := 0.0
    n := len(poly)
    for i := 0; i < n; i++ {
        j := (i+1)%n
        s += poly[i].x*poly[j].y - poly[i].y*poly[j].x
    }
    if s < 0 { s = -s }
    return s/2
}

func main() {
    in := bufio.NewReader(os.Stdin)
    var n int
    fmt.Fscan(in, &n)
    pts := make([]Point, n)
    for i := 0; i < n; i++ {
        fmt.Fscan(in, &pts[i].x, &pts[i].y)
    }
    total := polygonArea(pts)

    // store indices of remaining vertices
    poly := make([]int, n)
    for i := 0; i < n; i++ { poly[i] = i }

    areas := [2]float64{}
    cur := 0 // 0 Alberto, 1 Beatrice

    for len(poly) > 2 {
        // find ear with minimal area
        minIdx := 0
        minArea := 1e100
        m := len(poly)
        for i := 0; i < m; i++ {
            a := pts[poly[(i+m-1)%m]]
            b := pts[poly[i]]
            c := pts[poly[(i+1)%m]]
            ar := areaTri(a,b,c)
            if ar < minArea {
                minArea = ar
                minIdx = i
            }
        }
        areas[cur] += minArea
        // remove vertex
        poly = append(poly[:minIdx], poly[minIdx+1:]...)
        cur ^= 1
    }

    if areas[0] <= total/2 && areas[1] <= total/2 {
        fmt.Println("Either")
    } else if areas[0] <= total/2 {
        fmt.Println("Alberto")
    } else {
        fmt.Println("Beatrice")
    }
}

