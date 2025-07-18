#include <cstdio>
#include <algorithm>
#include <vector>
#define mp make_pair
#define pb push_back
#define fir first
#define sec second
#define NE(x, y) ++ne, e[ne]=y, h[ne]=s[x], s[x]=ne
using namespace std;
typedef long long ll;

inline void R(int &x) {
	char ch = getchar(); x = 0;
	for (; ch<'0'; ch=getchar());
	for (; ch>='0'; ch=getchar()) x = x*10+ch-'0';
}
const int N = 100005;
vector<int> ans;
int n, ne = 0, s[N], e[N*2], h[N*2], a[N], b[N];
void dfs(int x, int fa, int fc, int fs) {
	int cur = a[x] ^ fc ^ b[x];
	if (cur) ans.pb(x);
	for (int w=s[x]; w; w=h[w]) if (e[w]!=fa)
		dfs(e[w], x, fs, fc ^ cur);
}
int main() {
	R(n);
	int x, y;
	for (int i=1; i<n; ++i) {
		R(x); R(y);
		NE(x, y); NE(y, x);
	}
	for (int i=1; i<=n; ++i) R(a[i]);
	for (int i=1; i<=n; ++i) R(b[i]);
	dfs(1, 0, 0, 0);
	sort(ans.begin(), ans.end());
	printf("%d\n", (int)ans.size());
	for (int i=0; i<ans.size(); ++i)
		printf("%d\n", ans[i]);
	return 0;
}