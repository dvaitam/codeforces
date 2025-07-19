package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// S represents an item with power p and identifier id
type S struct {
   p, id int
}

func main() {
   sc := bufio.NewReader(os.Stdin)
   var N, V int
   if _, err := fmt.Fscan(sc, &N, &V); err != nil {
       return
   }
   ks := make([]S, 0, N)
   cs := make([]S, 0, N)
   for i := 0; i < N; i++ {
       var t, p int
       fmt.Fscan(sc, &t, &p)
       if t == 1 {
           ks = append(ks, S{p: p, id: i + 1})
       } else {
           cs = append(cs, S{p: p, id: i + 1})
       }
   }
   // sort descending by power
   sort.Slice(ks, func(i, j int) bool { return ks[i].p > ks[j].p })
   sort.Slice(cs, func(i, j int) bool { return cs[i].p > cs[j].p })

   sumK := make([]int, len(ks)+1)
   sumC := make([]int, len(cs)+1)
   for i := 0; i < len(ks); i++ {
       sumK[i+1] = sumK[i] + ks[i].p
   }
   for i := 0; i < len(cs); i++ {
       sumC[i+1] = sumC[i] + cs[i].p
   }

   ans := 0
   usedK, usedC := 0, 0
   maxK := V
   if maxK > len(ks) {
       maxK = len(ks)
   }
   for useK := 0; useK <= maxK; useK++ {
       rem := V - useK
       useC := rem / 2
       if useC > len(cs) {
           useC = len(cs)
       }
       score := sumK[useK] + sumC[useC]
       if score > ans {
           ans = score
           usedK = useK
           usedC = useC
       }
   }

   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   fmt.Fprintln(w, ans)

   ids := make([]int, 0, usedK+usedC)
   for i := 0; i < usedC; i++ {
       ids = append(ids, cs[i].id)
   }
   for i := 0; i < usedK; i++ {
       ids = append(ids, ks[i].id)
   }
   if len(ids) > 0 {
       sort.Ints(ids)
       for i := 0; i < len(ids)-1; i++ {
           fmt.Fprintf(w, "%d ", ids[i])
       }
       fmt.Fprintln(w, ids[len(ids)-1])
   } else {
       // print empty line if no ids
       fmt.Fprintln(w)
   }
}
