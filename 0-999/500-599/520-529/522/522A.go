package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   depths := make(map[string]int)
   depths["polycarp"] = 1
   maxDepth := 1
   for i := 0; i < n; i++ {
       var name1, rep, name2 string
       if _, err := fmt.Fscan(reader, &name1, &rep, &name2); err != nil {
           break
       }
       l1 := strings.ToLower(name1)
       l2 := strings.ToLower(name2)
       parentDepth := depths[l2]
       depths[l1] = parentDepth + 1
       if depths[l1] > maxDepth {
           maxDepth = depths[l1]
       }
   }
   fmt.Println(maxDepth)
}
