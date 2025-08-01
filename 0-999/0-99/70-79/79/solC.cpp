#include <cstdio>
#include <cstring>
#include <algorithm>
using namespace std;

const int MAXN = 100000 + 86;

char s[MAXN], b[12][12];
int l[12];

int ok(const int k, const int n) {
	int mx = -1;
	for (int t=0,i; t<n; ++t) {
		for (i=l[t]-1; i>=0&&k-l[t]+i+1>=0; --i) if (b[t][i] != s[k-l[t]+i+1])	break;
		if (i < 0) mx = max(mx, k-l[t]+1);
	}
	return mx;
}

int main() {
	int n;

	scanf("%s%d", s, &n);
	for (int i=0; i<n; ++i) scanf("%s", b[i]), l[i] = strlen(b[i]);
	int p = 0, mx = -1, st = 0, i;	
	for (i=0; s[i]; ++i) {
		int k = ok(i, n);	
		if (k == -1) continue;
		if (i - p > mx) mx = i - p, st = p;
		p = max(p, k + 1);
	}	
	if (i - p > mx) mx = i - p, st = p;
	printf("%d %d\n", mx, st);

	return 0;
}