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
   // TODO: Implement conversion of solF.cpp (Codeforces 1835F) to Go
   fmt.Fprintln(writer, "TODO: solution not yet implemented")
}
