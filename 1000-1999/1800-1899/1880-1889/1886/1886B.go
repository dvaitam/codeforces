package main

import (
    "bufio"
    "fmt"
    "math"
    "os"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()

    var t int
    if _, err := fmt.Fscan(reader, &t); err != nil {
        return
    }
    for i := 0; i < t; i++ {
        var px, py, ax, ay, bx, by float64
        fmt.Fscan(reader, &px, &py, &ax, &ay, &bx, &by)
        // candidate using point A
        dPA := (ax-px)*(ax-px) + (ay-py)*(ay-py)
        dOA := ax*ax + ay*ay
        // candidate using point B
        dPB := (bx-px)*(bx-px) + (by-py)*(by-py)
        dOB := bx*bx + by*by

        ans := math.Min(
            math.Max(dPA, dOA),
            math.Max(dPB, dOB),
        )
        // initial radius
        actual := math.Sqrt(ans)

        // distance AB
        dAB := math.Hypot(ax-bx, ay-by)
        // scenario 1
        r1 := math.Max(dAB/2, math.Max(math.Sqrt(dOA), math.Sqrt(dPB)))
        // scenario 2
        r2 := math.Max(dAB/2, math.Max(math.Sqrt(dOB), math.Sqrt(dPA)))
        if r1 < actual {
            actual = r1
        }
        if r2 < actual {
            actual = r2
        }

        // print with high precision
        fmt.Fprintf(writer, "%.20f\n", actual)
    }
}
