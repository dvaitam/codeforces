package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for i := 0; i < t; i++ {
       var a, b, c int
       fmt.Fscan(reader, &a, &b, &c)
       first := "1" + strings.Repeat("0", a-1)
       second := strings.Repeat("9", b-c+1) + strings.Repeat("0", c-1)
       fmt.Fprint(writer, first, " ", second, "\n")
   }
