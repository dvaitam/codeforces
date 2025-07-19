package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
       a[i]--
   }
   indeg := make([]int, n)
   for _, v := range a {
       indeg[v]++
   }
   typ := make([]int, n)
   type pair struct{ i, t int }
   queue := make([]pair, 0, n)
   // push function
   var push = func(i, t int) {
       if typ[i] != 0 {
           if typ[i] != t {
               fmt.Fprint(writer, -1)
               writer.Flush()
               os.Exit(0)
           }
           return
       }
       typ[i] = t
       queue = append(queue, pair{i, t})
   }
   // initial roots
   for i := 0; i < n; i++ {
       if indeg[i] == 0 {
           push(i, 1)
       }
   }
   // BFS
   for head := 0; head < len(queue); head++ {
       i, t := queue[head].i, queue[head].t
       if t == 1 {
           push(a[i], 2)
       } else if t == 2 {
           indeg[a[i]]--
           if indeg[a[i]] == 0 {
               push(a[i], 1)
           }
       }
   }
   // handle cycles
   for i := 0; i < n; i++ {
       if typ[i] != 0 {
           continue
       }
       j, t := i, 1
       for {
           if typ[j] != 0 {
               if typ[j] != t {
                   fmt.Fprint(writer, -1)
                   writer.Flush()
                   os.Exit(0)
               }
               break
           }
           typ[j] = t
           j = a[j]
           t = 3 - t
       }
   }
   // collect answer
   ans := make([]int, 0, n)
   for i := 0; i < n; i++ {
       if typ[i] == 1 {
           ans = append(ans, a[i]+1)
       }
   }
   fmt.Fprintln(writer, len(ans))
   for idx, v := range ans {
       if idx > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   fmt.Fprintln(writer)
}
