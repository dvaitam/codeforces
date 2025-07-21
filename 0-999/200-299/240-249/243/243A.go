package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   arr := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &arr[i])
   }
   result := make(map[int]struct{})
   current := make(map[int]struct{})
   for _, v := range arr {
       next := make(map[int]struct{})
       next[v] = struct{}{}
       for x := range current {
           next[x|v] = struct{}{}
       }
       for y := range next {
           result[y] = struct{}{}
       }
       current = next
   }
   fmt.Println(len(result))
}
