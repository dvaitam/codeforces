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

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   sectors := make([][]int, n)
   for i := 0; i < n; i++ {
       var k int
       fmt.Fscan(reader, &k)
       sectors[i] = make([]int, k)
       for j := 0; j < k; j++ {
           fmt.Fscan(reader, &sectors[i][j])
       }
       sort.Ints(sectors[i])
   }

   unstable := 0
   for i := 0; i < n; i++ {
       dist := sectors[i]
       left := (i - 1 + n) % n
       right := (i + 1) % n
       distL := sectors[left]
       distR := sectors[right]
       for j := 0; j+1 < len(dist); j++ {
           a, b := dist[j], dist[j+1]
           // count in left sector distances in (a, b)
           l1 := sort.Search(len(distL), func(x int) bool { return distL[x] > a })
           l2 := sort.Search(len(distL), func(x int) bool { return distL[x] >= b })
           countL := l2 - l1
           // count in right sector distances in (a, b)
           r1 := sort.Search(len(distR), func(x int) bool { return distR[x] > a })
           r2 := sort.Search(len(distR), func(x int) bool { return distR[x] >= b })
           countR := r2 - r1
           if countL != countR {
               unstable++
           }
       }
   }

   fmt.Fprint(writer, unstable)
}
