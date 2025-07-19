package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Node stores the original index and value
type Node struct {
   x int
   v int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   nodes := make([]Node, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &nodes[i].v)
       nodes[i].x = i + 1
   }
   // sort by value
   sort.Slice(nodes, func(i, j int) bool {
       return nodes[i].v < nodes[j].v
   })
   // track top index for each value
   top := make([]int, n+5)
   for i := range top {
       top[i] = -1
   }
   // check contiguous values and record positions
   prev := -1
   for i := 0; i < n; i++ {
       cur := nodes[i].v
       if cur-prev > 1 {
           fmt.Fprintln(writer, "Impossible")
           return
       }
       top[cur] = i
       prev = cur
   }
   ans := make([]int, n)
   m, now := 0, 0
   // first pass
   for i := 0; i < n; i++ {
       if now < 0 || now >= len(top) || top[now] == -1 {
           break
       }
       idx := top[now]
       ans[m] = nodes[idx].x
       m++
       t := nodes[idx].v
       top[now]--
       if top[now] < 0 || nodes[top[now]].v != t {
           top[now] = -1
       }
       now++
       if now >= 3 && top[now-3] != -1 && top[now-2] != -1 && top[now-1] != -1 {
           now -= 3
       }
   }
   // second pass
   for now >= 3 {
       now -= 3
       if top[now] == -1 {
           continue
       }
       idx := top[now]
       ans[m] = nodes[idx].x
       m++
       top[now] = -1
       now++
       for now < len(top) && top[now] != -1 {
           idx2 := top[now]
           ans[m] = nodes[idx2].x
           m++
           top[now] = -1
           now++
       }
   }
   if m != n {
       fmt.Fprintln(writer, "Impossible")
       return
   }
   fmt.Fprintln(writer, "Possible")
   for i := 0; i < n; i++ {
       if i == n-1 {
           fmt.Fprintln(writer, ans[i])
       } else {
           fmt.Fprint(writer, ans[i], " ")
       }
   }
}
