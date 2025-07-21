package main

import (
   "bufio"
   "os"
   "strings"
   "fmt"
)

func main() {
   scanner := bufio.NewScanner(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   group := make([]string, 0)
   flush := func() {
       if len(group) == 0 {
           return
       }
       var sb strings.Builder
       for _, line := range group {
           for i := 0; i < len(line); i++ {
               if line[i] != ' ' {
                   sb.WriteByte(line[i])
               }
           }
       }
       writer.WriteString(sb.String())
       writer.WriteByte('\n')
       group = group[:0]
   }
   for scanner.Scan() {
       line := scanner.Text()
       // Remove trailing carriage return if present
       if len(line) > 0 && line[len(line)-1] == '\r' {
           line = line[:len(line)-1]
       }
       // Determine if amplifying line (first non-space char is '#')
       isAmp := false
       for i := 0; i < len(line); i++ {
           if line[i] == ' ' {
               continue
           }
           if line[i] == '#' {
               isAmp = true
           }
           break
       }
       if isAmp {
           flush()
           // Write amplifying line unchanged
           writer.WriteString(line)
           writer.WriteByte('\n')
       } else {
           // Collect common lines
           group = append(group, line)
       }
   }
   // Flush any remaining common lines
   flush()
   // Check for scan error
   if err := scanner.Err(); err != nil {
       fmt.Fprintln(os.Stderr, "Error reading input:", err)
   }
}
