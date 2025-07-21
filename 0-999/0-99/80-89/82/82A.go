package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int64
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   names := []string{"Sheldon", "Leonard", "Penny", "Rajesh", "Howard"}
   group := int64(1)
   count := int64(len(names)) * group
   for n > count {
       n -= count
       group *= 2
       count = int64(len(names)) * group
   }
   idx := (n - 1) / group
   fmt.Println(names[idx])
}
