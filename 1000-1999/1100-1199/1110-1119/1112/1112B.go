package main

import (
   "bufio"
   "fmt"
   "os"
   "math"
)

type kebab struct {
   start, finish int64
   idx           int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   t := make([]kebab, k)
   // initialize servers
   for j := 0; j < k; j++ {
       t[j].idx = -1
   }
   ans := make([]bool, n)
   var m int
   var currentTime, nextTime int64

   // assign tasks
   for i := 0; i < n; i++ {
       // find server with min finish time
       v := 0
       mn := t[0].finish
       for j := 1; j < k; j++ {
           if t[j].finish < mn {
               mn = t[j].finish
               v = j
           }
       }
       if t[v].finish > 0 {
           m++
       }
       currentTime = t[v].finish
       t[v].start = currentTime
       t[v].finish = currentTime + a[i]
       t[v].idx = i

       // compute next finish
       v = 0
       mn = t[0].finish
       for j := 1; j < k; j++ {
           if t[j].finish < mn {
               mn = t[j].finish
               v = j
           }
       }
       nextTime = t[v].finish

       // mark answers
       u := int64(math.Floor(float64(100*m)/float64(n) + 0.5))
       for j := 0; j < k; j++ {
           if t[j].idx < 0 {
               continue
           }
           if u+t[j].start-1 >= currentTime && u+t[j].start <= nextTime {
               ans[t[j].idx] = true
           }
       }
   }

   // finish remaining
   for {
       // find next to finish
       v := -1
       mn := int64(1<<62 - 1)
       for j := 0; j < k; j++ {
           if t[j].idx < 0 {
               continue
           }
           if t[j].finish < mn {
               mn = t[j].finish
               v = j
           }
       }
       if v < 0 {
           break
       }
       m++
       currentTime = t[v].finish
       // reset server
       t[v].start = 0
       t[v].finish = 0
       t[v].idx = -1

       // next finish
       v = -1
       mn = int64(1<<62 - 1)
       for j := 0; j < k; j++ {
           if t[j].idx < 0 {
               continue
           }
           if t[j].finish < mn {
               mn = t[j].finish
               v = j
           }
       }
       if v < 0 {
           break
       }
       nextTime = t[v].finish

       u := int64(math.Floor(float64(100*m)/float64(n) + 0.5))
       for j := 0; j < k; j++ {
           if t[j].idx < 0 {
               continue
           }
           if u+t[j].start-1 >= currentTime && u+t[j].start <= nextTime {
               ans[t[j].idx] = true
           }
       }
   }

   // output
   count := 0
   for i := 0; i < n; i++ {
       if ans[i] {
           count++
       }
   }
   fmt.Fprint(writer, count)
}
