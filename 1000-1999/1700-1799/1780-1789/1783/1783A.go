package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
   "strconv"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var n int
       fmt.Fscan(reader, &n)
       v := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &v[i])
       }
       allSame := true
       for i := 1; i < n; i++ {
           if v[i] != v[i-1] {
               allSame = false
               break
           }
       }
       if allSame {
           fmt.Println("NO")
           continue
       }
       fmt.Println("YES")
       sort.Ints(v)
       out := make([]string, 0, n)
       out = append(out, strconv.Itoa(v[0]))
       for i := n - 1; i > 0; i-- {
           out = append(out, strconv.Itoa(v[i]))
       }
       fmt.Println(strings.Join(out, " "))
   }
}
