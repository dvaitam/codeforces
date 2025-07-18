// LUOGU_RID: 158552455
#include <bits/stdc++.h>
using namespace std;

namespace SOL {

using i64 = long long;
void debug(const char *msg, ...) {
#ifdef CLESIP
    va_list arg;
    static char pbString[512];
    va_start(arg,msg);
    vsprintf(pbString,msg,arg);
    cerr << "[DEBUG] " << pbString << "\n";
    va_end(arg);
#endif    
}
template<typename T, typename L>
bool chkmax(T &x, L y) { if (x < y) return x = y, true; return false; }
template<typename T, typename L>
bool chkmin(T &x, L y) { if (y < x) return x = y, true; return false; }

const int N = 55;
int n, m, k, a[N][N];
vector<array<int, 4>> ans;

void solve() {
	cin >> n >> m >> k;
	for (int i = 1; i <= n; i ++)
		for (int j = 1; j <= m; j ++)
			cin >> a[i][j];

	for (int i = n; i >= 1; i --) {
		for (int j = 1; j <= m; j ++) {
			for (int l = 2; l <= n - i + 1; l ++)
				ans.push_back({ i, j, l, a[i][j] });
			for (int l = i + 1; l <= 2 * i - 1; l ++)
				ans.push_back({ l, j, n - i + 1, a[i][j] });
			for (int l = n - i + 2; l <= n + a[i][j]; l ++)
				ans.push_back({ 2 * i - 1, j, l, a[i][j] });
		}
	}
	for (int t = 1; t <= k; t ++) {
		for (int i = 1; i < n; i ++) if ((i & 1) || i + 1 == n) {
			for (int j = 1; j <= m; j ++)
				ans.push_back({ 2 * i, j, t + n - 1, t });
		}
		for (int i = 1; i < n; i ++) {
			ans.push_back({ 2 * i, m + 1, t + n - 1, t });
			if (i + 1 < n)
				ans.push_back({ 2 * i + 1, m + 1, t + n - 1, t });
		}
	}
	cout << ans.size() << "\n";
	for (auto [x, y, z, c] : ans)
		cout << x << " " << y << " " << z << " " << c << "\n";
}

}
int main() {
	ios::sync_with_stdio(false), cin.tie(0), cout.tie(0);
	SOL::solve();
	return 0;
}