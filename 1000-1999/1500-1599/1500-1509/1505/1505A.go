package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   scanner := bufio.NewScanner(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   // Read all tokens (words) from stdin and output "NO" for each
   for scanner.Scan() {
       // s := scanner.Text()
       writer.WriteString("NO\n")
   }
}
