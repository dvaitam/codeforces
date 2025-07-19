package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const inf = int64(1e18)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   N := n + 2
   x := make([]int64, N)
   for i := 0; i < N; i++ {
       fmt.Fscan(reader, &x[i])
   }
   d := make([]int64, N)
   for i := 1; i <= n; i++ {
       var v int64
       fmt.Fscan(reader, &v)
       d[i] = 2 * v
   }
   d[0] = x[n+1] - x[0]
   d[n+1] = d[0]

   r := make([]int64, N)
   st := make([]int, N)
   // forward sweep
   r[0] = d[0]
   e := 1
   st[0] = 0
   for i := 1; i < N; i++ {
       for e >= 2 && d[st[e-1]] < x[i]-x[st[e-2]] {
           e--
       }
       r[i] = d[i] - (x[i] - x[st[e-1]])
       st[e] = i
       e++
   }
   if e > 2 {
       fmt.Fprintln(writer, "0.0")
       return
   }
   // backward sweep
   l := make([]int64, N)
   l[n+1] = d[n+1]
   e = 1
   st[0] = n + 1
   for i := n; i >= 0; i-- {
       for e >= 2 && d[st[e-1]] < x[st[e-2]]-x[i] {
           e--
       }
       l[i] = d[i] - (x[st[e-1]] - x[i])
       st[e] = i
       e++
   }
   r[n+1] = -inf
   l[0] = -inf

   // events
   type event struct{ key int64; idx int }
   ev := make([]event, 0, N)
   for i := 0; i < N; i++ {
       if l[i] <= 0 {
           continue
       }
       ev = append(ev, event{key: x[i] - l[i], idx: i})
   }
   sort.Slice(ev, func(i, j int) bool {
       if ev[i].key != ev[j].key {
           return ev[i].key < ev[j].key
       }
       return ev[i].idx < ev[j].idx
   })
   esz := len(ev)
   pos := make([]int, N)
   for i, evi := range ev {
       pos[evi.idx] = i + 1
   }
   // fenwick tree: store min
   fenw := make([]int64, esz+1)
   for i := 1; i <= esz; i++ {
       fenw[i] = inf
   }
   // functions
   modify := func(i int, v int64) {
       for i <= esz {
           if v < fenw[i] {
               fenw[i] = v
           }
           i = (i | (i - 1)) + 1
       }
   }
   findMin := func(i int) int64 {
       res := inf
       for i > 0 {
           if fenw[i] < res {
               res = fenw[i]
           }
           i &= i - 1
       }
       return res
   }

   ans := inf
   // process from end to start
   for i := n + 1; i >= 0; i-- {
       if r[i] > 0 {
           to := x[i] + r[i]
           // find first ev.key >= to+1
           toPos := sort.Search(esz, func(j int) bool {
               return ev[j].key >= to+1
           })
           u := findMin(toPos)
           if u != inf {
               if u-x[i] < ans {
                   ans = u - x[i]
               }
           }
       }
       if l[i] > 0 {
           modify(pos[i], x[i])
       }
   }
   // output ans/2 with .0 or .5
   whole := ans / 2
   frac := (ans % 2) * 5
   fmt.Fprintf(writer, "%d.%d\n", whole, frac)
}
