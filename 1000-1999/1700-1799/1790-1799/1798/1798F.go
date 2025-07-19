package main

import (
   "bufio"
   "fmt"
   "math/rand"
   "os"
   "sort"
   "time"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   fmt.Fscan(reader, &n, &k)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   type pair struct{ ff, ss int }
   s := make([]pair, k)
   for i := 0; i < k; i++ {
       fmt.Fscan(reader, &s[i].ff)
       s[i].ss = i
   }
   sort.Slice(s, func(i, j int) bool {
       return s[i].ff < s[j].ff
   })
   vis := make([]bool, n)
   g := make([][]int, k)
   rand.Seed(time.Now().UnixNano())
   // process first k-1 groups
   for i := 0; i < k-1; i++ {
       id := s[i].ss
       O := s[i].ff
       // build list of unvisited elements
       type elem struct{ val, idx int }
       elems := make([]elem, 0, n)
       for j := 0; j < n; j++ {
           if !vis[j] {
               elems = append(elems, elem{a[j], j})
           }
       }
       m := len(elems)
       ok := false
       c := 0
       for !ok {
           c++
           if c == 1 {
               sort.Slice(elems, func(x, y int) bool {
                   return elems[x].val%O < elems[y].val%O
               })
           } else {
               for t := 0; t < 20; t++ {
                   x := rand.Intn(m)
                   y := rand.Intn(m)
                   elems[x], elems[y] = elems[y], elems[x]
               }
           }
           // prefix sums
           pre := make([]int, m+1)
           for j := 1; j <= m; j++ {
               pre[j] = pre[j-1] + elems[j-1].val
           }
           // find window
           for j := O; j <= m; j++ {
               if (pre[j]-pre[j-O])%O == 0 {
                   ok = true
                   for t := j - O; t < j; t++ {
                       idx := elems[t].idx
                       g[id] = append(g[id], idx)
                       vis[idx] = true
                   }
                   break
               }
           }
       }
   }
   // last group: unvisited plus filler
   sum := 0
   lastID := s[k-1].ss
   for i := 0; i < n; i++ {
       if !vis[i] {
           sum += a[i]
           g[lastID] = append(g[lastID], i)
       }
   }
   ff := s[k-1].ff
   rem := sum % ff
   ljd := ff - rem
   if rem == 0 {
       ljd = ff
   }
   // append filler
   a = append(a, ljd)
   g[lastID] = append(g[lastID], n)
   // output
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, ljd)
   for i := 0; i < k; i++ {
       for _, idx := range g[i] {
           fmt.Fprint(writer, a[idx], " ")
       }
       fmt.Fprintln(writer)
   }
}
