#include <cstdio>
#include <algorithm>
#include <queue>
#define fir first
#define sec second
#define mp make_pair
#define pii pair<int, int>
using namespace std;
typedef long long ll;

inline void R(int &x) {
	char ch = getchar(); x = 0;
	for (; ch<'0'; ch=getchar());
	for (; ch>='0'; ch=getchar()) x = x*10+ch-'0';
}
const int N = 100005;
int n, m, lim, a[N], b[N], c[N], pa[N], pb[N], bel[N];
inline bool cmpb(int x, int y) {return b[x] < b[y];}
inline bool cmpa(int x, int y) {return a[x] < a[y];}
bool chk(int k) {
	priority_queue< int, vector<int>, greater<int> > H;
	int pp = n;
	int cst = 0;
	for (int i=m; i>0; i-=k) {
		for (; pp && b[pb[pp]]>=a[pa[i]]; --pp) H.push(c[pb[pp]]);
		if (H.empty()) return 0;
		cst += H.top();
		if (cst > lim) return 0;
		H.pop();
	}
	return 1;
}
bool gen(int k) {
	priority_queue< pii, vector<pii>, greater<pii> > H;
	int pp = n;
	int cst = 0;
	for (int i=m; i>0; i-=k) {
		for (; pp && b[pb[pp]]>=a[pa[i]]; --pp) H.push(mp(c[pb[pp]], pb[pp]));
		if (H.empty()) return 0;
		cst += H.top().first;
		for (int j=0; j<k && i-j>0; ++j)
			bel[pa[i-j]] = H.top().second;
		if (cst > lim) return 0;
		H.pop();
	}
	return 1;
}
int main() {
	R(n); R(m); R(lim);
	for (int i=1; i<=m; ++i) R(a[i]), pa[i] = i;
	for (int i=1; i<=n; ++i) R(b[i]), pb[i] = i;
	for (int i=1; i<=n; ++i) R(c[i]);
	sort(pa+1, pa+m+1, cmpa);
	sort(pb+1, pb+n+1, cmpb);
	int l = 1, r = m+1, mid;
	while (l < r) {
		mid = (l + r) >> 1;
		if (chk(mid))
			r = mid; else
			l = mid + 1;
	}
	if (r > m)
		puts("NO");
	else {
		puts("YES");
		gen(l);
		for (int i=1; i<=m; ++i) printf("%d ", bel[i]);
	}
	return 0;
}