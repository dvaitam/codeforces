package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Point represents a point in 3D space with an identifier
type Point struct {
   x, y, z, id int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   for {
       var n int
       if _, err := fmt.Fscan(reader, &n); err != nil {
           break
       }
       ps := make([]Point, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &ps[i].x, &ps[i].y, &ps[i].z)
           ps[i].id = i + 1
       }
       // initial slice copy (if needed for debugging)
       // first pass: match by x and y
       ans := make([][2]int, 0, n/2)
       sort.Slice(ps, func(i, j int) bool {
           if ps[i].x != ps[j].x {
               return ps[i].x < ps[j].x
           }
           if ps[i].y != ps[j].y {
               return ps[i].y < ps[j].y
           }
           return ps[i].z < ps[j].z
       })
       var nps []Point
       for i := 0; i < len(ps); {
           if i+1 < len(ps) && ps[i].x == ps[i+1].x && ps[i].y == ps[i+1].y {
               ans = append(ans, [2]int{ps[i].id, ps[i+1].id})
               i += 2
           } else {
               nps = append(nps, ps[i])
               i++
           }
       }
       ps = nps

       // second pass: match by x
       sort.Slice(ps, func(i, j int) bool {
           if ps[i].x != ps[j].x {
               return ps[i].x < ps[j].x
           }
           if ps[i].y != ps[j].y {
               return ps[i].y < ps[j].y
           }
           return ps[i].z < ps[j].z
       })
       nps = nps[:0]
       for i := 0; i < len(ps); {
           if i+1 < len(ps) && ps[i].x == ps[i+1].x {
               ans = append(ans, [2]int{ps[i].id, ps[i+1].id})
               i += 2
           } else {
               nps = append(nps, ps[i])
               i++
           }
       }
       ps = nps

       // final pass: match remaining arbitrarily
       sort.Slice(ps, func(i, j int) bool {
           if ps[i].x != ps[j].x {
               return ps[i].x < ps[j].x
           }
           if ps[i].y != ps[j].y {
               return ps[i].y < ps[j].y
           }
           return ps[i].z < ps[j].z
       })
       for i := 0; i+1 < len(ps); i += 2 {
           ans = append(ans, [2]int{ps[i].id, ps[i+1].id})
       }

       // output
       for _, p := range ans {
           fmt.Fprintln(writer, p[0], p[1])
       }
   }
}
