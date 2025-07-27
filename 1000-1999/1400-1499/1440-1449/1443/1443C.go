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

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for tc := 0; tc < t; tc++ {
       var n int
       fmt.Fscan(reader, &n)
       a := make([]int, n)
       b := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &b[i])
       }
       // pair up and sort by a ascending
       pairs := make([]struct{ a, b int }, n)
       for i := 0; i < n; i++ {
           pairs[i].a = a[i]
           pairs[i].b = b[i]
       }
       sort.Slice(pairs, func(i, j int) bool {
           return pairs[i].a < pairs[j].a
       })
       // initial sum of b (all pickups)
       sumB := 0
       for i := 0; i < n; i++ {
           sumB += pairs[i].b
       }
       ans := sumB
       // try delivering prefix [0..i], pickups suffix [i+1..]
       for i := 0; i < n; i++ {
           sumB -= pairs[i].b
           // time is max(current a, sum of remaining b)
           time := sumB
           if pairs[i].a > time {
               time = pairs[i].a
           }
           if time < ans {
               ans = time
           }
       }
       fmt.Fprintln(writer, ans)
   }
}
