package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const solution112ASource = `package main

import (
	"fmt"
	"strings"
)

func main() {
	var a, b string
	// read two strings
	if _, err := fmt.Scan(&a, &b); err != nil {
		return
	}
	ua := strings.ToUpper(a)
	ub := strings.ToUpper(b)
	switch {
	case ua > ub:
		fmt.Print("1")
	case ua < ub:
		fmt.Print("-1")
	default:
		fmt.Print("0")
	}
}
`

// Keep the embedded reference solution reachable so it is preserved in the binary.
var _ = solution112ASource

type testCase struct {
	a string
	b string
}

var testcases = []testCase{
	{a: "bV", b: "rp"},
	{a: "iVgR", b: "VIfL"},
	{a: "cbfnoGM", b: "bJmTPSI"},
	{a: "oCLrZaW", b: "ZkSBvrj"},
	{a: "Wvgf", b: "ygww"},
	{a: "qZcUDIhyfJ", b: "sONxKmTecQ"},
	{a: "Xsfo", b: "gyrD"},
	{a: "kxwnQr", b: "SRPeMO"},
	{a: "IUp", b: "kDy"},
	{a: "OSJoR", b: "uXXdo"},
	{a: "Z", b: "u"},
	{a: "renKTun", b: "PFzPDjq"},
	{a: "pVJ", b: "IqV"},
	{a: "BLzxoiGFfW", b: "dhjOkYRBMe"},
	{a: "yMDHqJa", b: "RUhRIWr"},
	{a: "hsBkDa", b: "UUqGWl"},
	{a: "gOtOGMmjx", b: "WkIXHaMuF"},
	{a: "h", b: "x"},
	{a: "pdpKf", b: "fUFeW"},
	{a: "XiiQEJkqH", b: "MBnIWUSmT"},
	{a: "zQPxC", b: "HChpo"},
	{a: "vb", b: "LJ"},
	{a: "Loae", b: "TOdo"},
	{a: "cv", b: "eG"},
	{a: "rQFn", b: "IiUK"},
	{a: "EpYEZAmggQ", b: "BwBADUdRPP"},
	{a: "dz", b: "Uv"},
	{a: "pm", b: "mI"},
	{a: "iBlrDpeC", b: "ZJgdPIaf"},
	{a: "kAFE", b: "nzdk"},
	{a: "ayqYYDs", b: "BSUYJQT"},
	{a: "jmsndLVI", b: "dVuddLEG"},
	{a: "kdGfleMeR", b: "pzhKpLMcN"},
	{a: "AQ", b: "LK"},
	{a: "uqnQTupqz", b: "iQPtDuWea"},
	{a: "NKgeInGq", b: "iwepxskC"},
	{a: "TtNZPHaQJ", b: "tQgiqhgVJ"},
	{a: "rsM", b: "nTv"},
	{a: "ROqG", b: "Fqdf"},
	{a: "rcavXiO", b: "qkVCJTB"},
	{a: "aheSjIcxL", b: "JjBictxYc"},
	{a: "nRpQgw", b: "XJANVj"},
	{a: "kZZl", b: "AblV"},
	{a: "YAZQVZ", b: "prkYSg"},
	{a: "cEomDwt", b: "YoobQmz"},
	{a: "reXrwP", b: "GzRIvb"},
	{a: "ql", b: "Lq"},
	{a: "g", b: "M"},
	{a: "wUYuBMG", b: "hyKmqcT"},
	{a: "aHZIRUV", b: "VQmxBeQ"},
	{a: "NuQhUt", b: "GtQAuz"},
	{a: "JimAQ", b: "yRVlN"},
	{a: "tzJatsnBYL", b: "MPuDCCRnGE"},
	{a: "Qfs", b: "GQO"},
	{a: "vfWpRtoZmj", b: "bcpENXeDAO"},
	{a: "mTSyFzpjPS", b: "aWXgXBolZS"},
	{a: "DdJphDiZD", b: "QHJMuWCNU"},
	{a: "BJCkVECqW", b: "pOrXXHFOp"},
	{a: "CeTsp", b: "rvuIf"},
	{a: "joy", b: "SjT"},
	{a: "eAAv", b: "IDAd"},
	{a: "AyXL", b: "SbWK"},
	{a: "EawtWyA", b: "IVVIZMo"},
	{a: "orBFbyvQ", b: "RZzUkDiN"},
	{a: "bzLKQbfPB", b: "iDldqyunD"},
	{a: "vWyrWA", b: "qfEbVI"},
	{a: "w", b: "o"},
	{a: "XP", b: "cW"},
	{a: "p", b: "m"},
	{a: "N", b: "j"},
	{a: "iEQh", b: "KnDS"},
	{a: "XxkMM", b: "VThXk"},
	{a: "gLbtK", b: "RyzTm"},
	{a: "LS", b: "Op"},
	{a: "SX", b: "tR"},
	{a: "ZhYKYcwIBQ", b: "xeGPvaAFgB"},
	{a: "ODTjBl", b: "UHPrNZ"},
	{a: "XEDBULrup", b: "frCpWDKNQ"},
	{a: "vbFulFn", b: "wZqvrMS"},
	{a: "JaHmf", b: "pUAFJ"},
	{a: "SEPT", b: "FCYb"},
	{a: "so", b: "zS"},
	{a: "tQLx", b: "EJHw"},
	{a: "VJvwSDr", b: "tqohUmu"},
	{a: "VI", b: "WS"},
	{a: "mnV", b: "ErU"},
	{a: "WHMsgmsoxl", b: "taTIircdJs"},
	{a: "OWF", b: "gaK"},
	{a: "EECvl", b: "dqEhe"},
	{a: "FeKORdj", b: "jZKtfph"},
	{a: "WAMMYNoXH", b: "yCCtLBtKN"},
	{a: "N", b: "V"},
	{a: "Wn", b: "On"},
	{a: "Qfkpl", b: "JekaA"},
	{a: "SMEscosT", b: "sSDeRoqY"},
	{a: "QZmBhIoPjr", b: "jedkYtMVKs"},
	{a: "hDStSzrG", b: "IFCfMcBV"},
	{a: "MqbfoR", b: "KLbWRr"},
	{a: "cWWlEHPCrl", b: "LBOFfEwAvu"},
	{a: "kv", b: "AS"},
}

func solveCase(tc testCase) string {
	ua := strings.ToUpper(tc.a)
	ub := strings.ToUpper(tc.b)
	switch {
	case ua > ub:
		return "1"
	case ua < ub:
		return "-1"
	default:
		return "0"
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for idx, tc := range testcases {
		input := tc.a + "\n" + tc.b + "\n"
		want := solveCase(tc)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx+1, err, string(out))
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != want {
			fmt.Printf("test %d failed: expected %s got %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
