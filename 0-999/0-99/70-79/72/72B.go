package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
   "strconv"
   "strings"
)

func main() {
   scanner := bufio.NewScanner(os.Stdin)
   if !scanner.Scan() {
       return
   }
   n, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
   if err != nil {
       return
   }
   // Global key-value mappings
   global := make(map[string]string)
   // Section to key-value mappings
   sections := make(map[string]map[string]string)
   var sectionNames []string
   currentSection := ""
   for i := 0; i < n; i++ {
       if !scanner.Scan() {
           break
       }
       line := scanner.Text()
       // Ignore comment lines
       trimmedLeft := strings.TrimLeft(line, " \t")
       if len(trimmedLeft) > 0 && trimmedLeft[0] == ';' {
           continue
       }
       // Section header
       trimmed := strings.TrimSpace(line)
       if len(trimmed) >= 2 && trimmed[0] == '[' && trimmed[len(trimmed)-1] == ']' {
           name := strings.TrimSpace(trimmed[1 : len(trimmed)-1])
           currentSection = name
           if _, exists := sections[currentSection]; !exists {
               sections[currentSection] = make(map[string]string)
               sectionNames = append(sectionNames, currentSection)
           }
       } else {
           // Key-value line
           if idx := strings.Index(line, "="); idx != -1 {
               key := strings.TrimSpace(line[:idx])
               value := strings.TrimSpace(line[idx+1:])
               if currentSection == "" {
                   global[key] = value
               } else {
                   sections[currentSection][key] = value
               }
           }
       }
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   // Output global key-values sorted by key
   var globalKeys []string
   for k := range global {
       globalKeys = append(globalKeys, k)
   }
   sort.Strings(globalKeys)
   for _, k := range globalKeys {
       fmt.Fprintf(writer, "%s=%s\n", k, global[k])
   }
   // Output sections in sorted order
   sort.Strings(sectionNames)
   for _, sec := range sectionNames {
       fmt.Fprintf(writer, "[%s]\n", sec)
       var keys []string
       for k := range sections[sec] {
           keys = append(keys, k)
       }
       sort.Strings(keys)
       for _, k := range keys {
           fmt.Fprintf(writer, "%s=%s\n", k, sections[sec][k])
       }
   }
}
