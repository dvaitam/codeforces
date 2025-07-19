package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var tst int
   if _, err := fmt.Fscan(reader, &tst); err != nil {
       return
   }
   for ; tst > 0; tst-- {
       // TODO: implement conversion logic from solG.cpp
       // For now, output -1 to compile successfully
       fmt.Fprintln(writer, -1)
   }
}
