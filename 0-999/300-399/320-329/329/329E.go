package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n int
   fmt.Fscan(in, &n)
   xs := make([]int64, n)
   ys := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &xs[i], &ys[i])
   }
   // Helper: compute max cycle length for given sorted indices
   calc := func(order []int) int64 {
       var best int64
       // two patterns: start low or start high
       for t := 0; t < 2; t++ {
           seq := make([]int, 0, n)
           l, r := 0, n-1
           for i := 0; i < n; i++ {
               if (i%2 == 0) == (t == 0) {
                   seq = append(seq, order[l])
                   l++
               } else {
                   seq = append(seq, order[r])
                   r--
               }
           }
           // compute cycle length
           var sum int64
           for i := 0; i < n; i++ {
               j := seq[i]
               k := seq[(i+1)%n]
               dx := xs[j] - xs[k]
               if dx < 0 {
                   dx = -dx
               }
               dy := ys[j] - ys[k]
               if dy < 0 {
                   dy = -dy
               }
               sum += dx + dy
           }
           if sum > best {
               best = sum
           }
       }
       return best
   }

   var ans int64
   // Try a set of fixed projection directions
   type vec struct{ dx, dy int64 }
   dirs := []vec{
       {1, 0}, {0, 1},
       {1, 1}, {1, -1},
       {-1, 1}, {-1, -1},
   }
   idx := make([]int, n)
   for _, v := range dirs {
       // sort by projection onto v
       for i := range idx {
           idx[i] = i
       }
       sort.Slice(idx, func(i, j int) bool {
           return v.dx*xs[idx[i]]+v.dy*ys[idx[i]] < v.dx*xs[idx[j]]+v.dy*ys[idx[j]]
       })
       if val := calc(idx); val > ans {
           ans = val
       }
   }

   fmt.Fprintln(out, ans)
}
