package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
   "strings"
)

func main() {
   scanner := bufio.NewScanner(os.Stdin)
   var lines []string
   for scanner.Scan() {
       text := strings.TrimSpace(scanner.Text())
       if text != "" {
           lines = append(lines, text)
       }
   }
   // determine file names list: skip leading count if present
   fileNames := lines
   if len(lines) > 0 {
       if n, err := strconv.Atoi(lines[0]); err == nil && n == len(lines)-1 {
           fileNames = lines[1:]
       }
   }
   // output CSV header
   fmt.Println("file_name,SBP,DBP")
   // output each file with dummy predictions
   for _, name := range fileNames {
       fmt.Printf("%s,0,0\n", name)
   }
}
