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

   var n, k int
   fmt.Fscan(reader, &n)
   stations := make([]struct{ x int; idx int }, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &stations[i].x)
       stations[i].idx = i + 1
   }
   fmt.Fscan(reader, &k)
   // sort stations by coordinate
   sort.Slice(stations, func(i, j int) bool {
       return stations[i].x < stations[j].x
   })
   // initial window [0..k-1]
   var sumX int64
   var S int64
   for t := 0; t < k; t++ {
       xi := int64(stations[t].x)
       sumX += xi
       // weight w = 2*t - k + 1
       S += xi * int64(2*t - k + 1)
   }
   best := S
   bestL := 0
   // slide window from l=0 to n-k
   for l := 0; l + k < n; l++ {
       leftX := int64(stations[l].x)
       rightX := int64(stations[l+k].x)
       // update S: remove left, include right with weights shift
       S = S + int64(k+1)*leftX - 2*sumX + int64(k-1)*rightX
       // update sum of X
       sumX = sumX - leftX + rightX
       if S < best {
           best = S
           bestL = l + 1
       }
   }
   // output the original indices of chosen stations
   for t := 0; t < k; t++ {
       if t > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, stations[bestL+t].idx)
   }
   writer.WriteByte('\n')
}
