package main

import (
   "bufio"
   "fmt"
   "os"
)

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   a := make([]int, n)
   var minV, maxV int
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
       if i == 0 {
           minV, maxV = a[i], a[i]
       } else {
           if a[i] < minV {
               minV = a[i]
           }
           if a[i] > maxV {
               maxV = a[i]
           }
       }
   }
   if maxV-minV <= 1 {
       fmt.Fprintln(writer, n)
       for i := 0; i < n; i++ {
           if i > 0 {
               fmt.Fprint(writer, " ")
           }
           fmt.Fprint(writer, a[i])
       }
       fmt.Fprint(writer, "\n")
       return
   }

   cnt0, cnt1, cnt2 := 0, 0, 0
   for _, v := range a {
       switch v - minV {
       case 0:
           cnt0++
       case 1:
           cnt1++
       case 2:
           cnt2++
       }
   }
   best := -1
   best0, best1, best2 := 0, 0, 0
   for i := -n; i <= n; i++ {
       if cnt1-2*i < 0 || cnt0+i < 0 || cnt2+i < 0 {
           continue
       }
       now := 0
       now += max(cnt1, cnt1-2*i) - cnt1
       now += max(cnt0, cnt0+i) - cnt0
       now += max(cnt2, cnt2+i) - cnt2
       if now > best {
           best = now
           best0 = cnt0 + i
           best1 = cnt1 - 2*i
           best2 = cnt2 + i
       }
   }
   // number of equal measurements = n - (new distinct count)
   fmt.Fprintln(writer, n-best)
   first := true
   for i := 0; i < best0; i++ {
       if !first {
           fmt.Fprint(writer, " ")
       }
       first = false
       fmt.Fprint(writer, minV)
   }
   for i := 0; i < best1; i++ {
       if !first {
           fmt.Fprint(writer, " ")
       }
       first = false
       fmt.Fprint(writer, minV+1)
   }
   for i := 0; i < best2; i++ {
       if !first {
           fmt.Fprint(writer, " ")
       }
       first = false
       fmt.Fprint(writer, minV+2)
   }
   fmt.Fprint(writer, "\n")
}
