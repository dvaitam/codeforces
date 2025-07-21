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

   var door string
   var a int
   if _, err := fmt.Fscan(reader, &door); err != nil {
       return
   }
   fmt.Fscan(reader, &a)

   var res string
   if (door == "front" && a == 1) || (door == "back" && a == 2) {
       res = "L"
   } else {
       res = "R"
   }
   fmt.Fprintln(writer, res)
}
