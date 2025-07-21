package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
   "strings"
)

type rowData struct {
   cells []string
   idx   int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   // Read column names
   headerLine, err := reader.ReadString('\n')
   if err != nil {
       return
   }
   headerLine = strings.TrimSpace(headerLine)
   cols := strings.Fields(headerLine)
   // Read sort rules
   rulesLine, err := reader.ReadString('\n')
   if err != nil {
       return
   }
   rulesLine = strings.TrimSpace(rulesLine)
   ruleParts := strings.Split(rulesLine, ", ")
   // Map column names to indices
   colIndex := make(map[string]int, len(cols))
   for i, name := range cols {
       colIndex[name] = i
   }
   // Parse rules
   type rule struct{
       idx int
       asc bool
   }
   var rules []rule
   for _, part := range ruleParts {
       parts := strings.Fields(part)
       if len(parts) != 2 {
           continue
       }
       name, ord := parts[0], parts[1]
       ci, ok := colIndex[name]
       if !ok {
           continue
       }
       asc := true
       if ord == "DESC" {
           asc = false
       }
       rules = append(rules, rule{idx: ci, asc: asc})
   }
   // Read table rows
   var data []rowData
   scanner := bufio.NewScanner(reader)
   rowIdx := 0
   for scanner.Scan() {
       line := scanner.Text()
       if strings.TrimSpace(line) == "" {
           continue
       }
       cells := strings.Fields(line)
       data = append(data, rowData{cells: cells, idx: rowIdx})
       rowIdx++
   }
   // Stable sort according to rules
   sort.SliceStable(data, func(i, j int) bool {
       a, b := data[i], data[j]
       for _, r := range rules {
           ai := a.cells[r.idx]
           bi := b.cells[r.idx]
           if ai == bi {
               continue
           }
           if r.asc {
               return ai < bi
           }
           return ai > bi
       }
       // Preserve input order if all keys equal
       return a.idx < b.idx
   })
   // Print sorted table (including header)
   fmt.Println(strings.Join(cols, " "))
   for _, d := range data {
       fmt.Println(strings.Join(d.cells, " "))
   }
}
