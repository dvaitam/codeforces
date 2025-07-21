package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
   "strings"
)

func main() {
   scanner := bufio.NewScanner(os.Stdin)
   // map from superstition to count of countries
   counts := make(map[string]int)
   // current country's superstitions set
   currentSet := make(map[string]bool)
   firstCountry := true
   for scanner.Scan() {
       line := scanner.Text()
       if strings.HasPrefix(line, "* ") {
           // superstition line
           name := line[2:]
           // normalize: lowercase, single spaces
           fields := strings.Fields(name)
           for i, w := range fields {
               fields[i] = strings.ToLower(w)
           }
           norm := strings.Join(fields, " ")
           currentSet[norm] = true
       } else {
           // country line; flush previous country's data if not first
           if !firstCountry {
               for s := range currentSet {
                   counts[s]++
               }
           }
           firstCountry = false
           // reset for new country
           currentSet = make(map[string]bool)
       }
   }
   // flush last country
   for s := range currentSet {
       counts[s]++
   }
   // find max count
   maxCount := 0
   for _, c := range counts {
       if c > maxCount {
           maxCount = c
       }
   }
   // collect superstitions with max count
   var best []string
   for s, c := range counts {
       if c == maxCount {
           best = append(best, s)
       }
   }
   sort.Strings(best)
   // output
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   for _, s := range best {
       fmt.Fprintln(w, s)
   }
   if err := scanner.Err(); err != nil {
       // ignore
   }
}
