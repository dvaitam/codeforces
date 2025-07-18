#include <bits/stdc++.h>
using namespace std;
#define fir first
#define sec second
#define pb push_back
#define mp make_pair
#define Size(x) ((int)((x).size()))
#define rep(i, l, r) for(int (i) = (l); (i) <= (r); ++(i))
#define per(i, r, l) for(int (i) = (r); (i) >= (l); --(i))
#define chkmax(x, y) (x) = max((x), (y))
#define chkmin(x, y) (x) = min((x), (y))
typedef long long ll;
typedef pair<int,int> pii;
typedef vector<int> vii;
inline char gc() {
	static char now[1<<16], *S, *T;
	if(S == T) {T = (S = now) + fread(now, 1, 1<<16, stdin); if(S == T) return EOF;}
	return *S++;
}
#define gc getchar
inline int read() {
	int x = 0, f = 1; char c = gc();
	while(c < '0' || c > '9') {(c ^ '-') ? 0 : (f = 0); c = gc();}
	while(c >= '0' && c <= '9') {x = (x << 3) + (x << 1) + (c ^ 48); c = gc();}
	return f ? x : ((~x) + 1);
}
int main() {
	int n = read();
	ll res = 0, ans = 0;
	for (int i = 1; i <= n; ++ i) {
		int x = read();
		ll tmp = min(res, 1ll * x / 2);
		tmp += (x - tmp * 2) / 3;
		res += x - 3 * tmp;
		ans += tmp;
	}
	printf("%I64d\n", ans);
	return 0;
}