package main

import (
   "bufio"
   "fmt"
   "os"
)

type edge struct {
   to, id int
}

var (
   pr []int
)

func init() {
   const maxP = 200005
   mark := make([]bool, maxP)
   for i := 2; i < maxP; i++ {
       if mark[i] {
           continue
       }
       pr = append(pr, i)
       for j := i * 2; j < maxP; j += i {
           mark[j] = true
       }
   }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var t, n int
   fmt.Fscan(in, &t)
   for t > 0 {
       t--
       fmt.Fscan(in, &n)
       if n == 2 {
           fmt.Fprintln(out, "2 2")
           continue
       }
       // build vec of primes
       var vec []int
       for i := 0; i < len(pr); i++ {
           vec = append(vec, pr[i])
           k := len(vec)
           if k&1 == 1 {
               if k*(k+1)/2 >= n-1 {
                   break
               }
           } else {
               if k*(k+1)/2 - k/2 + 1 >= n-1 {
                   break
               }
           }
       }
       k := len(vec)
       // build graph
       adj := make([][]edge, k)
       nume := 0
       if k&1 == 1 {
           for i := 0; i < k; i++ {
               for j := i + 1; j < k; j++ {
                   nume++
                   adj[i] = append(adj[i], edge{j, nume})
                   adj[j] = append(adj[j], edge{i, nume})
               }
           }
       } else {
           for i := 0; i < k; i++ {
               for j := i + 1; j < k; j++ {
                   if i&1 == 0 && j == i+1 {
                       continue
                   }
                   nume++
                   adj[i] = append(adj[i], edge{j, nume})
                   adj[j] = append(adj[j], edge{i, nume})
               }
           }
           // add extra edge (1-based 1-2) => (0-1)
           nume++
           adj[0] = append(adj[0], edge{1, nume})
           adj[1] = append(adj[1], edge{0, nume})
       }
       // prepare for tour
       ptr := make([]int, k)
       used := make([]bool, k)
       evis := make([]bool, nume+2)
       a := make([]int, n)
       ptr2 := 0

       var tour func(v int)
       tour = func(v int) {
           if ptr2 == n {
               return
           }
           for ptr[v] < len(adj[v]) {
               e := adj[v][ptr[v]]
               if evis[e.id] {
                   ptr[v]++
                   continue
               }
               evis[e.id] = true
               tour(e.to)
               if ptr2 == n {
                   return
               }
               ptr[v]++
           }
           if !used[v] {
               used[v] = true
               if ptr2 < n {
                   a[ptr2] = vec[v]
                   ptr2++
               }
           }
           if ptr2 == n {
               return
           }
           if ptr2 < n {
               a[ptr2] = vec[v]
               ptr2++
           }
       }
       tour(0)
       // output
       for i := 0; i < n; i++ {
           if i > 0 {
               out.WriteByte(' ')
           }
           fmt.Fprint(out, a[i])
       }
       out.WriteByte('\n')
   }
}
