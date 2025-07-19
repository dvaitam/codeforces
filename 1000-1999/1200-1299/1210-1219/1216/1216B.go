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
   tasks := make([]struct{ id, val int }, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &tasks[i].val)
       tasks[i].id = i + 1
   }
   sort.Slice(tasks, func(i, j int) bool {
       return tasks[i].val < tasks[j].val
   })
   res := int64(0)
   ans := make([]int, n)
   t := 1
   for i := n - 1; i >= 0; i-- {
       ans[t-1] = tasks[i].id
       res += int64((t-1)*tasks[i].val + 1)
       t++
   }
   fmt.Fprintln(writer, res)
   for i := 0; i < n; i++ {
       fmt.Fprint(writer, ans[i], " ")
   }
}
