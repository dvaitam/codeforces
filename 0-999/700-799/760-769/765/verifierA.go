package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func solveCase(n int, home string, flights []string) string {
	diff := 0
	for _, f := range flights {
		if strings.HasPrefix(f, home) {
			diff++
		} else {
			diff--
		}
	}
	if diff == 0 {
		return "home"
	}
	return "contest"
}

// Embedded testcases from testcasesA.txt.
const testcasesRaw = `15
4
YNB
YNB->MZJ
LSG->JEY
YNB->RWZ
EJD->VKP
5
DLN
DLN->GRP
DLN->BZR
CXM->ZVU
DLN->DLN
HXK->CGS
2
HZE
OCC->QPD
HZE->HZE
1
RKR
RKR->RSJ
4
CTZ
CTZ->JFG
BTV->CTZ
CTZ->EEB
WRV->CTZ
5
IQZ
IQZ->IQZ
NSI->IQZ
WZL->IQZ
PSU->IQZ
IQZ->IQZ
3
DWH
DWH->DWH
DWH->DWH
WHB->URT
1
ADU
ADU->DMC
3
DBT
DBT->FWD
DBT->BVA
TDI->DBT
1
UJL
UJL->UJL
4
BTD
MGI->BTD
SFW->YBZ
FKQ->BTD
OVF->BTD
4
SQJ
MVI->SQJ
OXC->SQJ
SQJ->SQJ
SQJ->LTJ
3
SUT
SUT->SUT
UCA->WKF
SUT->MWV
5
NBM
SNY->BFO
NBM->OQP
TYA->PKJ
BZZ->NBM
UCX->NBM
1
MVN
MVN->MVN`

type testcase struct {
	n       int
	home    string
	flights []string
}

func parseTestcases() ([]testcase, error) {
	scan := bufio.NewScanner(strings.NewReader(testcasesRaw))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return nil, fmt.Errorf("missing test count")
	}
	t, err := strconv.Atoi(scan.Text())
	if err != nil {
		return nil, fmt.Errorf("bad test count: %v", err)
	}
	var cases []testcase
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			return nil, fmt.Errorf("case %d missing n", i+1)
		}
		n, err := strconv.Atoi(scan.Text())
		if err != nil {
			return nil, fmt.Errorf("case %d bad n: %v", i+1, err)
		}
		if !scan.Scan() {
			return nil, fmt.Errorf("case %d missing home", i+1)
		}
		home := scan.Text()
		flights := make([]string, n)
		for j := 0; j < n; j++ {
			if !scan.Scan() {
				return nil, fmt.Errorf("case %d incomplete flights", i+1)
			}
			flights[j] = scan.Text()
		}
		cases = append(cases, testcase{n: n, home: home, flights: flights})
	}
	return cases, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]

	cases, err := parseTestcases()
	if err != nil {
		fmt.Println("failed to parse embedded testcases:", err)
		os.Exit(1)
	}

	for caseNum, tc := range cases {
		n := tc.n
		home := tc.home
		flights := tc.flights
		expected := solveCase(n, home, flights)

		var input bytes.Buffer
		fmt.Fprintf(&input, "%d\n%s\n", n, home)
		for _, fl := range flights {
			input.WriteString(fl)
			input.WriteByte('\n')
		}

		cmd := exec.Command(exe)
		cmd.Stdin = bytes.NewReader(input.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", caseNum+1, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != expected {
			fmt.Printf("case %d failed: expected %s got %s\n", caseNum+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
