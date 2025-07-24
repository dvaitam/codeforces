package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   vals := make(map[int]struct{})
   for i := 0; i < n; i++ {
       var a int
       fmt.Fscan(reader, &a)
       vals[a] = struct{}{}
   }
   uniq := make([]int, 0, len(vals))
   for k := range vals {
       uniq = append(uniq, k)
   }
   switch len(uniq) {
   case 0:
       fmt.Println("YES")
   case 1, 2:
       fmt.Println("YES")
   case 3:
       sort.Ints(uniq)
       if uniq[2]-uniq[1] == uniq[1]-uniq[0] {
           fmt.Println("YES")
       } else {
           fmt.Println("NO")
       }
   default:
       fmt.Println("NO")
   }
}
