package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const INF = 1070000000

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m, k int
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   sort.Ints(a)
   type task struct{ d, id int }
   tasks := make([]task, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &tasks[i].d)
       tasks[i].id = i + 1
   }
   sort.Slice(tasks, func(i, j int) bool { return tasks[i].d < tasks[j].d })

   headA, headT := 0, 0
   var ans []int
   // iterate days
   for d := 0; ; d++ {
       // skip tasks with deadline < d
       for headT < m && tasks[headT].d < d {
           headT++
       }
       // if next a slot expired
       if headA < n && a[headA] < d {
           fmt.Fprintln(writer, -1)
           return
       }
       // schedule up to k items
       cnt := 0
       for cnt < k && (headT < m || headA < n) {
           var td int
           if headT < m {
               td = tasks[headT].d
           } else {
               td = INF
           }
           if headA >= n || td < a[headA] {
               // take task
               ans = append(ans, tasks[headT].id)
               headT++
           } else {
               // use slot
               headA++
           }
           cnt++
       }
       if headT >= m && headA >= n {
           break
       }
   }
   // output result
   sort.Ints(ans)
   fmt.Fprintln(writer, len(ans))
   for i, v := range ans {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   writer.WriteByte('\n')
}
