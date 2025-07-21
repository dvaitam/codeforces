package main

import (
   "bufio"
   "fmt"
   "math/big"
   "os"
   "sort"
   "strconv"
   "strings"
)

type interval struct {
   start, end int
   x           int
   weight      *big.Int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   line, _ := reader.ReadString('\n')
   n, _ := strconv.Atoi(strings.TrimSpace(line))
   winsList := make(map[int][]int)
   intervals := make([]interval, 0)
   for i := 0; i < n; i++ {
       line, _ = reader.ReadString('\n')
       fields := strings.Fields(line)
       if fields[0] == "win" {
           x, _ := strconv.Atoi(fields[1])
           winsList[x] = append(winsList[x], i)
       } else if fields[0] == "sell" {
           x, _ := strconv.Atoi(fields[1])
           if list, ok := winsList[x]; ok && len(list) > 0 {
               start := list[len(list)-1]
               intervals = append(intervals, interval{start: start, end: i, x: x})
           }
       }
   }
   if len(intervals) == 0 {
       fmt.Println(0)
       return
   }
   // sort by end
   sort.Slice(intervals, func(i, j int) bool {
       return intervals[i].end < intervals[j].end
   })
   m := len(intervals)
   ends := make([]int, m)
   maxX := 0
   for i, iv := range intervals {
       ends[i] = iv.end
       if iv.x > maxX {
           maxX = iv.x
       }
   }
   // precompute powers of two up to maxX
   pow2 := make([]*big.Int, maxX+1)
   pow2[0] = big.NewInt(1)
   for i := 1; i <= maxX; i++ {
       pow2[i] = new(big.Int).Lsh(pow2[i-1], 1)
   }
   // assign weights
   for i := range intervals {
       intervals[i].weight = pow2[intervals[i].x]
   }
   // dp[j]: best for first j intervals
   dp := make([]*big.Int, m+1)
   dp[0] = big.NewInt(0)
   for j := 1; j <= m; j++ {
       // option1: skip interval j-1
       best := new(big.Int).Set(dp[j-1])
       iv := intervals[j-1]
       // find last k where ends[k] < iv.start
       k := sort.Search(m, func(i int) bool { return ends[i] >= iv.start })
       // dp index for include: k is first >= start, so k-1 is last < start
       var with big.Int
       with.Set(iv.weight)
       if k > 0 {
           with.Add(&with, dp[k])
       }
       if with.Cmp(best) > 0 {
           dp[j] = &with
       } else {
           dp[j] = best
       }
   }
   fmt.Println(dp[m].String())
}
