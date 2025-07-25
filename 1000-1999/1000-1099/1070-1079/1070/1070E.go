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

   var c int
   fmt.Fscan(reader, &c)
   for ; c > 0; c-- {
       var n, m int
       var tLimit int64
       fmt.Fscan(reader, &n, &m, &tLimit)
       p := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &p[i])
       }
       // unique sorted difficulties
       uniq := make([]int, n)
       copy(uniq, p)
       sort.Ints(uniq)
       u := 0
       for i := 0; i < n; i++ {
           if i == 0 || uniq[i] != uniq[i-1] {
               uniq[u] = uniq[i]
               u++
           }
       }
       uniq = uniq[:u]
       // binary search max d index
       lo, hi := 0, len(uniq)-1
       bestIdx := -1
       var bestCount int
       for lo <= hi {
           mid := (lo + hi) / 2
           d := uniq[mid]
           cnt, used := simulate(p, m, int64(d), tLimit)
           if used <= tLimit {
               bestIdx = mid
               bestCount = cnt
               lo = mid + 1
           } else {
               hi = mid - 1
           }
       }
       if bestIdx < 0 {
           // no tasks
           fmt.Fprintln(writer, 0, 1)
       } else {
           fmt.Fprintln(writer, bestCount, uniq[bestIdx])
       }
   }
}

// simulate returns number of tasks completed and total time used (including breaks) for given d
// simulate calculates number of tasks completed and minimal total time (tasks + necessary breaks)
func simulate(p []int, m int, d int64, _tLimit int64) (count int, totalTime int64) {
   var taskTime, breakTime, curSum, lastBreakSum int64
   var fullGroups int
   for _, pi := range p {
       if int64(pi) <= d {
           count++
           taskTime += int64(pi)
           curSum += int64(pi)
           if count% m == 0 {
               // full group completed: add break
               breakTime += curSum
               lastBreakSum = curSum
               fullGroups++
               curSum = 0
           }
       }
   }
   // skip break after last group
   if fullGroups > 0 && count% m == 0 {
       breakTime -= lastBreakSum
   }
   totalTime = taskTime + breakTime
   return count, totalTime
}
