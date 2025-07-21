package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   intervals := make([][2]int, 0, n)
   for i := 0; i < n; i++ {
       var t int
       fmt.Fscan(reader, &t)
       if t == 1 {
           var x, y int
           fmt.Fscan(reader, &x, &y)
           intervals = append(intervals, [2]int{x, y})
       } else if t == 2 {
           var a, b int
           fmt.Fscan(reader, &a, &b)
           a--
           b--
           if isReachable(intervals, a, b) {
               fmt.Println("YES")
           } else {
               fmt.Println("NO")
           }
       }
   }
}

func isReachable(intervals [][2]int, start, target int) bool {
   n := len(intervals)
   visited := make([]bool, n)
   queue := make([]int, 0, n)
   visited[start] = true
   queue = append(queue, start)
   for len(queue) > 0 {
       u := queue[0]
       queue = queue[1:]
       a0, b0 := intervals[u][0], intervals[u][1]
       for v := 0; v < n; v++ {
           if !visited[v] {
               c, d := intervals[v][0], intervals[v][1]
               if (c < a0 && a0 < d) || (c < b0 && b0 < d) {
                   visited[v] = true
                   if v == target {
                       return true
                   }
                   queue = append(queue, v)
               }
           }
       }
   }
   return visited[target]
}
