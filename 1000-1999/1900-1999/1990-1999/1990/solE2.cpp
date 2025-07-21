#include <bits/stdc++.h>
using namespace std; const int N = 5005;
vector<int> g[N], ord; int d[N], s = 73, c, t;
bool check(int x) { return cout << "? " << x << '\n', cin >> x, x; }
void dfs(int v, int p = -1) {
	ord.push_back(v);
	for (auto u: g[v]) if (u != p && d[u] >= s) c = u;
	for (auto u: g[v]) if (u != p && d[u] >= s) if (u == c || check(u)) return dfs(u, v);
}
void cfs(int v, int p = -1) { for (auto u: g[v]) if (u != p) cfs(u, v), d[v] = max(d[u] + 1, d[v]); }
int main() {
	for (cin >> t; t--;) {
		int n, q = 0, fl = 0;
		cin >> n;
		for (int v = 1; v <= n; ++v) d[v] = 0, g[v].clear();
		for (int i = 1, x, y; i < n && cin >> x >> y; ++i) g[x].push_back(y), g[y].push_back(x);
		for (int i = 2; i <= n; ++i) if (g[i].size() == 1) q = i;
		for (int i = 0; i < s; ++i) fl |= check(q);
		if (fl) {
			cout << "! " << q << '\n';
			continue;
		}
		ord.clear(), cfs(1), dfs(1);
		int l = 0, r = ord.size() - 1, mid;
		while (l != r) mid = (l + r + 1) >> 1, check(ord[mid]) ? l = mid : (r = mid - 2, --l, l = max(l, 0), r = max(r, 0));
		cout << "! " << ord[l] << '\n';
	}
}