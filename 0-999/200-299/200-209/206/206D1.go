package main

import (
   "bufio"
   "fmt"
   "os"
)

// A simple stub classifier for ABBYY Cup problem D1.
// Reads document id, name, and content from stdin, and outputs a default subject 1.
func main() {
   reader := bufio.NewReader(os.Stdin)
   // Read and ignore document identifier
   var id int
   if _, err := fmt.Fscan(reader, &id); err != nil {
       return
   }
   // Read and ignore document name
   var name string
   if _, err := fmt.Fscan(reader, &name); err != nil {
       return
   }
   // Consume the rest of the document text
   scanner := bufio.NewScanner(reader)
   for scanner.Scan() {
   }
   // Default classification: subject 1
   fmt.Println(1)
}
