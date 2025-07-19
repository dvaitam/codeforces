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

   var n, m int
   fmt.Fscan(reader, &n, &m)

   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
       a[i]--
   }

   cnt := make([]int, m)
   for i := 0; i < n; i++ {
       if a[i] >= 0 && a[i] < m {
           cnt[a[i]]++
       }
   }

   per := n / m
   more := n % m
   ans := 0

   id := make([]int, m)
   for i := range id {
       id[i] = i
   }
   sort.Slice(id, func(i, j int) bool {
       return cnt[id[i]] > cnt[id[j]]
   })

   tar := make([]int, m)
   for idx, v := range id {
       if idx < more {
           tar[v] = per + 1
       } else {
           tar[v] = per
       }
   }

   // Reduce surplus counts
   for i := 0; i < m; i++ {
       if cnt[i] > tar[i] {
           left := cnt[i] - tar[i]
           for left > 0 {
               // find an index with value i
               idx := -1
               for j := 0; j < n; j++ {
                   if a[j] == i {
                       idx = j
                       break
                   }
               }
               // reassign to a deficit value
               for j := 0; j < m; j++ {
                   if cnt[j] < tar[j] {
                       a[idx] = j
                       cnt[j]++
                       cnt[i]--
                       ans++
                       break
                   }
               }
               left--
           }
       }
   }

   // Assign values outside [0, m)
   for i := 0; i < n; i++ {
       if a[i] < 0 || a[i] >= m {
           for j := 0; j < m; j++ {
               if cnt[j] < tar[j] {
                   a[i] = j
                   cnt[j]++
                   ans++
                   break
               }
           }
       }
   }

   // Output results
   fmt.Fprintln(writer, per, ans)
   for i := 0; i < n; i++ {
       if i > 0 {
           fmt.Fprint(writer, " ")
       }
       fmt.Fprint(writer, a[i]+1)
   }
   fmt.Fprintln(writer)
}
