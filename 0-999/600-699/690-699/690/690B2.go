package main

import (
    "bufio"
    "fmt"
    "os"
    "sort"
)

type Point struct{ x, y int }

func cross(a, b, c Point) int {
    return (b.x-a.x)*(c.y-a.y) - (b.y-a.y)*(c.x-a.x)
}

func convexHull(pts []Point) []Point {
    if len(pts) <= 1 {
        return pts
    }
    sort.Slice(pts, func(i, j int) bool {
        if pts[i].x == pts[j].x {
            return pts[i].y < pts[j].y
        }
        return pts[i].x < pts[j].x
    })
    var lower []Point
    for _, p := range pts {
        for len(lower) >= 2 && cross(lower[len(lower)-2], lower[len(lower)-1], p) <= 0 {
            lower = lower[:len(lower)-1]
        }
        lower = append(lower, p)
    }
    var upper []Point
    for i := len(pts)-1; i >= 0; i-- {
        p := pts[i]
        for len(upper) >= 2 && cross(upper[len(upper)-2], upper[len(upper)-1], p) <= 0 {
            upper = upper[:len(upper)-1]
        }
        upper = append(upper, p)
    }
    hull := append(lower, upper[1:len(upper)-1]...)
    return hull
}

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()
    for {
        var N int
        if _, err := fmt.Fscan(in, &N); err != nil {
            return
        }
        if N == 0 {
            break
        }
        lines := make([]string, N)
        for i := 0; i < N; i++ {
            fmt.Fscan(in, &lines[i])
        }
        // build candidate points
        var pts []Point
        for y := 0; y <= N; y++ {
            for x := 0; x <= N; x++ {
                sum := 0
                if x > 0 && y > 0 {
                    sum += int(lines[N-y][x-1] - '0')
                }
                if x < N && y > 0 {
                    sum += int(lines[N-y][x] - '0')
                }
                if x > 0 && y < N {
                    sum += int(lines[N-y-1][x-1] - '0')
                }
                if x < N && y < N {
                    sum += int(lines[N-y-1][x] - '0')
                }
                if sum >= 8 {
                    pts = append(pts, Point{x, y})
                }
            }
        }
        if len(pts) == 0 {
            fmt.Fprintln(out, 0)
            continue
        }
        hull := convexHull(pts)
        if len(hull) == 0 {
            fmt.Fprintln(out, 0)
            continue
        }
        // find lexicographically smallest point
        start := 0
        for i := 1; i < len(hull); i++ {
            if hull[i].x < hull[start].x || (hull[i].x == hull[start].x && hull[i].y < hull[start].y) {
                start = i
            }
        }
        fmt.Fprintln(out, len(hull))
        for k := 0; k < len(hull); k++ {
            idx := (start - k + len(hull)) % len(hull)
            p := hull[idx]
            fmt.Fprintf(out, "%d %d\n", p.x, p.y)
        }
    }
}

