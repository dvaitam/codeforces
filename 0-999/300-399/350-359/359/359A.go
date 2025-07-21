package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   grid := make([][]int, n)
   good := make([][2]int, 0)
   for i := 0; i < n; i++ {
       grid[i] = make([]int, m)
       for j := 0; j < m; j++ {
           var v int
           fmt.Fscan(reader, &v)
           grid[i][j] = v
           if v == 1 {
               good = append(good, [2]int{i, j})
           }
       }
   }
   nm := n * m
   words := (nm + 63) >> 6
   full := make([]uint64, words)
   for idx := 0; idx < nm; idx++ {
       full[idx>>6] |= 1 << (uint(idx) & 63)
   }
   type rect struct {
       m1, m2    int
       r1, r2    int
       c1, c2    int
       bits      []uint64
   }
   rects := make([][]rect, 4)
   // build rectangles for each corner
   for _, g := range good {
       x, y := g[0], g[1]
       // corner (1,1)
       {
           r := rect{m1: x + 1, m2: y + 1, r1: 0, r2: x, c1: 0, c2: y}
           r.bits = make([]uint64, words)
           for i := r.r1; i <= r.r2; i++ {
               for j := r.c1; j <= r.c2; j++ {
                   idx := i*m + j
                   r.bits[idx>>6] |= 1 << (uint(idx) & 63)
               }
           }
           rects[0] = append(rects[0], r)
       }
       // corner (1,m)
       {
           r := rect{m1: x + 1, m2: m - y, r1: 0, r2: x, c1: y, c2: m - 1}
           r.bits = make([]uint64, words)
           for i := r.r1; i <= r.r2; i++ {
               for j := r.c1; j <= r.c2; j++ {
                   idx := i*m + j
                   r.bits[idx>>6] |= 1 << (uint(idx) & 63)
               }
           }
           rects[1] = append(rects[1], r)
       }
       // corner (n,1)
       {
           r := rect{m1: n - x, m2: y + 1, r1: x, r2: n - 1, c1: 0, c2: y}
           r.bits = make([]uint64, words)
           for i := r.r1; i <= r.r2; i++ {
               for j := r.c1; j <= r.c2; j++ {
                   idx := i*m + j
                   r.bits[idx>>6] |= 1 << (uint(idx) & 63)
               }
           }
           rects[2] = append(rects[2], r)
       }
       // corner (n,m)
       {
           r := rect{m1: n - x, m2: m - y, r1: x, r2: n - 1, c1: y, c2: m - 1}
           r.bits = make([]uint64, words)
           for i := r.r1; i <= r.r2; i++ {
               for j := r.c1; j <= r.c2; j++ {
                   idx := i*m + j
                   r.bits[idx>>6] |= 1 << (uint(idx) & 63)
               }
           }
           rects[3] = append(rects[3], r)
       }
   }
   // pareto prune
   for c := 0; c < 4; c++ {
       arr := rects[c]
       sort.Slice(arr, func(i, j int) bool {
           if arr[i].m1 != arr[j].m1 {
               return arr[i].m1 > arr[j].m1
           }
           return arr[i].m2 > arr[j].m2
       })
       filtered := make([]rect, 0, len(arr))
       best2 := -1
       for _, r := range arr {
           if r.m2 > best2 {
               filtered = append(filtered, r)
               best2 = r.m2
           }
       }
       rects[c] = filtered
   }
   // check coverage
   // k = 1
   for c := 0; c < 4; c++ {
       for _, r := range rects[c] {
           ok := true
           for w := 0; w < words; w++ {
               if r.bits[w] != full[w] {
                   ok = false
                   break
               }
           }
           if ok {
               fmt.Println(1)
               return
           }
       }
   }
   // k = 2
   for c1 := 0; c1 < 4; c1++ {
       for c2 := c1 + 1; c2 < 4; c2++ {
           for _, r1 := range rects[c1] {
               for _, r2 := range rects[c2] {
                   ok := true
                   for w := 0; w < words; w++ {
                       if (r1.bits[w] | r2.bits[w]) != full[w] {
                           ok = false
                           break
                       }
                   }
                   if ok {
                       fmt.Println(2)
                       return
                   }
               }
           }
       }
   }
   // k = 3
   for c1 := 0; c1 < 4; c1++ {
       for c2 := c1 + 1; c2 < 4; c2++ {
           for c3 := c2 + 1; c3 < 4; c3++ {
               for _, r1 := range rects[c1] {
                   for _, r2 := range rects[c2] {
                       for _, r3 := range rects[c3] {
                           ok := true
                           for w := 0; w < words; w++ {
                               if (r1.bits[w] | r2.bits[w] | r3.bits[w]) != full[w] {
                                   ok = false
                                   break
                               }
                           }
                           if ok {
                               fmt.Println(3)
                               return
                           }
                       }
                   }
               }
           }
       }
   }
   // otherwise 4
   fmt.Println(4)
}
