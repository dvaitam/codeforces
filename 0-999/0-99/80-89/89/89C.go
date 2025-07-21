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
   var n, m int
   fmt.Fscan(reader, &n, &m)
   grid := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &grid[i])
   }
   // collect chips
   type chip struct{ r, c int; dir byte }
   chips := make([]chip, 0)
   id := make([][]int, n)
   for i := range id {
       id[i] = make([]int, m)
       for j := range id[i] {
           id[i][j] = -1
       }
   }
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           d := grid[i][j]
           if d == 'L' || d == 'R' || d == 'U' || d == 'D' {
               id[i][j] = len(chips)
               chips = append(chips, chip{i, j, d})
           }
       }
   }
   k := len(chips)
   // neighbor pointers
   l := make([]int, k)
   r := make([]int, k)
   u := make([]int, k)
   d := make([]int, k)
   for i := 0; i < k; i++ {
       l[i], r[i], u[i], d[i] = -1, -1, -1, -1
   }
   // rows: left/right
   rowMap := make(map[int][]int)
   for i, ch := range chips {
       rowMap[ch.r] = append(rowMap[ch.r], i)
   }
   for _, idxs := range rowMap {
       sort.Slice(idxs, func(i, j int) bool {
           return chips[idxs[i]].c < chips[idxs[j]].c
       })
       for t := 0; t < len(idxs); t++ {
           if t > 0 {
               l[idxs[t]] = idxs[t-1]
           }
           if t+1 < len(idxs) {
               r[idxs[t]] = idxs[t+1]
           }
       }
   }
   // columns: up/down
   colMap := make(map[int][]int)
   for i, ch := range chips {
       colMap[ch.c] = append(colMap[ch.c], i)
   }
   for _, idxs := range colMap {
       sort.Slice(idxs, func(i, j int) bool {
           return chips[idxs[i]].r < chips[idxs[j]].r
       })
       for t := 0; t < len(idxs); t++ {
           if t > 0 {
               u[idxs[t]] = idxs[t-1]
           }
           if t+1 < len(idxs) {
               d[idxs[t]] = idxs[t+1]
           }
       }
   }
   // simulate
   maxPoints := 0
   ways := 0
   // copy neighbor arrays per simulation
   for start := 0; start < k; start++ {
       // local copies
       lc := make([]int, k); copy(lc, l)
       rc := make([]int, k); copy(rc, r)
       uc := make([]int, k); copy(uc, u)
       dc := make([]int, k); copy(dc, d)
       points := 0
       curr := start
       for {
           // find next
           var nxt int
           switch chips[curr].dir {
           case 'L': nxt = lc[curr]
           case 'R': nxt = rc[curr]
           case 'U': nxt = uc[curr]
           case 'D': nxt = dc[curr]
           }
           // remove curr
           // fix row neighbors
           if lc[curr] != -1 {
               rc[lc[curr]] = rc[curr]
           }
           if rc[curr] != -1 {
               lc[rc[curr]] = lc[curr]
           }
           // fix column neighbors
           if uc[curr] != -1 {
               dc[uc[curr]] = dc[curr]
           }
           if dc[curr] != -1 {
               uc[dc[curr]] = uc[curr]
           }
           points++
           if nxt == -1 {
               break
           }
           curr = nxt
       }
       if points > maxPoints {
           maxPoints = points
           ways = 1
       } else if points == maxPoints {
           ways++
       }
   }
   fmt.Fprint(writer, maxPoints, " ", ways)
}
