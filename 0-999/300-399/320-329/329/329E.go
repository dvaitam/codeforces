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
   // Prepare direction vectors: basic axes and extremal pairs
   type vec struct{ dx, dy int64 }
   dirs := []vec{{1, 0}, {0, 1}, {1, 1}, {1, -1}}
   // find extremal in x+y
   minSum, maxSum := xs[0]+ys[0], xs[0]+ys[0]
   minSumI, maxSumI := 0, 0
   // find extremal in x-y
   minDiff, maxDiff := xs[0]-ys[0], xs[0]-ys[0]
   minDiffI, maxDiffI := 0, 0
   for i := 1; i < n; i++ {
       if s := xs[i] + ys[i]; s < minSum { minSum, minSumI = s, i }
       if s := xs[i] + ys[i]; s > maxSum { maxSum, maxSumI = s, i }
       if d := xs[i] - ys[i]; d < minDiff { minDiff, minDiffI = d, i }
       if d := xs[i] - ys[i]; d > maxDiff { maxDiff, maxDiffI = d, i }
   }
   // direction from min to max for sums and diffs
   dirs = append(dirs, vec{xs[maxSumI] - xs[minSumI], ys[maxSumI] - ys[minSumI]})
   dirs = append(dirs, vec{xs[maxDiffI] - xs[minDiffI], ys[maxDiffI] - ys[minDiffI]})
   // track seen directions to avoid duplicates
   seen := make(map[vec]bool)
   idx := make([]int, n)
   for _, v := range dirs {
       // normalize zero vector
       if v.dx == 0 && v.dy == 0 {
           continue
       }
       if seen[v] {
           continue
       }
       seen[v] = true
       // sort by projection
       for i := range idx { idx[i] = i }
       sort.Slice(idx, func(i, j int) bool {
           return v.dx*xs[idx[i]]+v.dy*ys[idx[i]] < v.dx*xs[idx[j]]+v.dy*ys[idx[j]]
       })
       if val := calc(idx); val > ans {
           ans = val
       }
   }

   fmt.Fprintln(out, ans)
}
