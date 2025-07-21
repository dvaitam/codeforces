package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   // Generate coordinates for hexagon of radius 2
   type coord struct{ q, r int }
   var coords []coord
   for r := -2; r <= 2; r++ {
       minq := -2
       if -r-2 > minq {
           minq = -r - 2
       }
       maxq := 2
       if -r+2 < maxq {
           maxq = -r + 2
       }
       for q := minq; q <= maxq; q++ {
           coords = append(coords, coord{q, r})
       }
   }
   // Map coordinates to indices
   idx := make(map[coord]int, len(coords))
   for i, c := range coords {
       idx[c] = i
   }
   // Build all lines (q-, r-, and s-axes)
   var lines [][]int
   // q-lines
   for qv := -2; qv <= 2; qv++ {
       var line []int
       for i, c := range coords {
           if c.q == qv {
               line = append(line, i)
           }
       }
       sort.Slice(line, func(i, j int) bool {
           return coords[line[i]].r < coords[line[j]].r
       })
       lines = append(lines, line)
   }
   // r-lines
   for rv := -2; rv <= 2; rv++ {
       var line []int
       for i, c := range coords {
           if c.r == rv {
               line = append(line, i)
           }
       }
       sort.Slice(line, func(i, j int) bool {
           return coords[line[i]].q < coords[line[j]].q
       })
       lines = append(lines, line)
   }
   // s-lines (s = -q-r)
   for sv := -2; sv <= 2; sv++ {
       var line []int
       for i, c := range coords {
           if -c.q-c.r == sv {
               line = append(line, i)
           }
       }
       sort.Slice(line, func(i, j int) bool {
           return coords[line[i]].q < coords[line[j]].q
       })
       lines = append(lines, line)
   }

   // Read input and build initial mask
   mask := 0
   for _, c := range coords {
       var s string
       fmt.Fscan(reader, &s)
       if s == "O" {
           mask |= 1 << idx[c]
       }
   }
   maxMask := mask
   // dp[mask] = Grundy number
   dp := make([]int, maxMask+1)
   // seen values for mex computation
   const maxMoves = 226
   seen := make([]bool, maxMoves)
   gList := make([]int, 0, 32)

   for m := 1; m <= maxMask; m++ {
       gList = gList[:0]
       // for each line, find contiguous runs of bits
       for _, line := range lines {
           for i := 0; i < len(line); {
               if (m>>line[i])&1 == 0 {
                   i++
                   continue
               }
               j := i
               for j < len(line) && (m>>line[j])&1 == 1 {
                   j++
               }
               // run from i to j-1
               for s := i; s < j; s++ {
                   seg := 0
                   for e := s; e < j; e++ {
                       seg |= 1 << line[e]
                       child := m &^ seg
                       g := dp[child]
                       if !seen[g] {
                           seen[g] = true
                           gList = append(gList, g)
                       }
                   }
               }
               i = j
           }
       }
       // mex
       mex := 0
       for mex < maxMoves && seen[mex] {
           mex++
       }
       dp[m] = mex
       for _, g := range gList {
           seen[g] = false
       }
   }
   // Result
   if dp[mask] != 0 {
       fmt.Fprintln(writer, "Karlsson")
   } else {
       fmt.Fprintln(writer, "Lillebror")
   }
}
