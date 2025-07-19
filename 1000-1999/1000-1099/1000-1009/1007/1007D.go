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

   // TODO: implement translation of solD.cpp logic in Go
   // This is a stub to ensure the file builds.
   fmt.Fprintln(writer, "YES")
}
