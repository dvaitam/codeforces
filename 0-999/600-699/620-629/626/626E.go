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
   fmt.Fscan(reader, &n)
   // read and sort input
   arr := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &arr[i])
   }
   sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
   // 1-indexed k array
   k := make([]int64, n+2)
   for i := 1; i <= n; i++ {
       k[i] = arr[i-1]
   }

   var (
       g1, g2, g3 int
       now, tot  int
       all       int64
       ans       float64
   )
   // first: odd count
   ans = 0.0
   g1, g2, g3 = 1, 1, 0
   now, tot = 0, 1
   all = k[1]
   for i := 2; i <= n; i++ {
       all -= k[i-now-1]
       all += k[i]
       for i-now-1 > 0 && n+1-now-1 > i && (k[i-now-1]+k[n+1-now-1])*int64(tot) >= 2*all {
           now++; tot += 2
           all += k[i-now] + k[n+1-now]
       }
       for now > 0 && (k[i-now]+k[n+1-now])*int64(tot) < 2*all {
           all -= k[i-now] + k[n+1-now]
           tot -= 2
           now--
       }
       cur := float64(all)/float64(tot) - float64(k[i])
       if ans < cur {
           ans = cur
           g1, g2, g3 = 1, i, now
       }
       if n+1-now == i+1 {
           all -= k[i-now] + k[n+1-now]
           tot -= 2
           now--
       }
   }
   // second: even count
   now, tot = 0, 2
   all = k[1] + k[2]
   for i := 2; i < n; i++ {
       all -= k[i-now-1]
       all += k[i+1]
       for i-now-1 > 0 && n+1-now-1 > i+1 && (k[i-now-1]+k[n+1-now-1])*int64(tot) >= 2*all {
           now++; tot += 2
           all += k[i-now] + k[n+1-now]
       }
       for now > 0 && (k[i-now]+k[n+1-now])*int64(tot) < 2*all {
           all -= k[i-now] + k[n+1-now]
           tot -= 2
           now--
       }
       cur := float64(all)/float64(tot) - float64(k[i]+k[i+1])/2.0
       if ans < cur {
           ans = cur
           g1, g2, g3 = 2, i, now
       }
       if n+1-now == i+2 {
           all -= k[i-now] + k[n+1-now]
           tot -= 2
           now--
       }
   }
   // output
   if g1 == 1 {
       cnt := g3*2 + 1
       fmt.Fprintln(writer, cnt)
       // left side
       for i := g3; i >= 1; i-- {
           fmt.Fprint(writer, k[g2-i], " ")
       }
       // median
       fmt.Fprint(writer, k[g2])
       // right side
       for i := 1; i <= g3; i++ {
           fmt.Fprint(writer, " ", k[n+1-i])
       }
       fmt.Fprintln(writer)
   } else {
       cnt := g3*2 + 2
       fmt.Fprintln(writer, cnt)
       for i := g3; i >= 1; i-- {
           fmt.Fprint(writer, k[g2-i], " ")
       }
       fmt.Fprint(writer, k[g2], " ", k[g2+1])
       for i := 1; i <= g3; i++ {
           fmt.Fprint(writer, " ", k[n+1-i])
       }
       fmt.Fprintln(writer)
   }
}
