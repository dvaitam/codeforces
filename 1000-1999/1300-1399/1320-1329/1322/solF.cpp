#include <bits/stdc++.h>

using namespace std;

const int N = 5e5 + 5;

using pii = pair<int, int>;

 

void upd(pii &x, const pii &y) {

	auto &[a, b] = x;

	const auto &[c, d] = y;

	x = {max(a, c), min(b, d)};

}

 

int n, m, val[N], repr[N], dis[N];

vector<int> G[N];

int fa[N][20], ord[N], dep[N];

 

int find(int x) {

	if (repr[x] == x) {

		return x;

	} else {

		int r = find(repr[x]);

		dis[x] ^= dis[repr[x]];

		repr[x] = r;

		return r;

	}

}

 

void merge(int u, int v, int d) {

	int fu = find(u), fv = find(v);

	if (fu != fv) {

		repr[fu] = fv;

		dis[fu] = (d ^ dis[u] ^ dis[v]);

	} else if ((dis[u] ^ dis[v]) != d) {

		cout << -1 << "\n";

		exit(0);

	}

}

 

void dfs(int u, int p) {

	static int dfs_clock = 0;

	if (p) {

		G[u].erase(find(G[u].begin(), G[u].end(), p));

	}

	dep[u] = dep[p] + 1;

	ord[++dfs_clock] = u;

	fa[u][0] = p;

	for (int i = 1; i < 20; i++) {

		fa[u][i] = fa[fa[u][i - 1]][i - 1];

	}

	for (int v : G[u]) {

		dfs(v, u);

	}

}

 

int kfa(int u, int k) {

	for (int i = 20; i--; ) {

		if ((k >> i) & 1) {

			u = fa[u][i];

		}

	}

	return u;

} 

 

int lca(int u, int v) {

	if (dep[u] < dep[v]) {

		swap(u, v);

	}

	u = kfa(u, dep[u] - dep[v]);

	if (u == v) {

		return u;

	}

	for (int i = 20; i--; ) {

		if (fa[u][i] != fa[v][i]) {

			u = fa[u][i], v = fa[v][i];

		}

	}

	return fa[u][0];

}

 

int dp[N], type[N];

pii s[N];

bool solve(int u, int k) {

	pii x(1, k);

	for (int v : G[u]) {

		pii &t = (find(v) == find(u) ? x : s[find(v)]);

		if (dis[v] == dis[u]) {

			type[v] = 0;

			upd(t, { dp[v] + 1, k });

		} else {

			type[v] = 1;

			upd(t, { 1, (k + 1 - dp[v]) - 1 });

		}

	}

	

	pii y(1, k);

	for (int v : G[u]) {

		if (find(v) != find(u)) {

			auto [l, r] = s[find(v)];

			if (l > k + 1 - r) {

				type[v] ^= 1;

				tie(l, r) = pii(k + 1 - r, k + 1 - l);

			}

			upd(y, { l, r });

		}

	}

	

	auto [a, b] = x;

	auto [c, d] = y;

	if (max(a, c) <= min(b, d)) {

		dp[u] = max(a, c);

	} else {

		tie(c, d) = pii(k + 1 - d, k + 1 - c);

		if (max(a, c) <= min(b, d)) {

			dp[u] = max(a, c);

			for (int v : G[u]) {

				if (find(v) != find(u)) {

					type[v] ^= 1;

				}

			}

		} else {

			return false;

		}

	}

	return true;

}

 

bool check(int k) {

	fill(s + 1, s + n + 1, pii(1, k));

	for (int i = n; i; i--) {

		if (!solve(ord[i], k)) {

			return false;

		}

	}

	return true;

}

 

int main() {

	ios::sync_with_stdio(false);

	cin.tie(nullptr);

	

	cin >> n >> m;

	iota(repr + 1, repr + n + 1, 1);

	

	for (int i = 1, u, v; i <= n - 1; i++) {

		cin >> u >> v;

		G[u].push_back(v);

		G[v].push_back(u);

	}

	

	dfs(1, 0);

	while (m--) {

		int u, v;

		cin >> u >> v;

		int l = lca(u, v);

		if (u != l) {

			val[u]++;

			u = kfa(u, dep[u] - dep[l] - 1);

			val[u]--;

		}

		if (v != l) {

			val[v]++;

			v = kfa(v, dep[v] - dep[l] - 1);

			val[v]--;

		}

		if (u != l && v != l) {

			merge(u, v, 1);

		}

	}

	

	for (int i = n; i; i--) {

		int u = ord[i];

		val[fa[u][0]] += val[u];

		if (val[u]) {

			merge(u, fa[u][0], 0);

		}

	}

	

	int l = 1, r = n, ans = 0;

	while (l <= r) {

		int mid = (l + r) / 2;

		if (check(mid)) {

			ans = mid;

			r = mid - 1;

		} else {

			l = mid + 1;

		}

	}

	

	check(ans);

	cout << ans << "\n";

	for (int i = 1; i <= n; i++) {

		int u = ord[i];

		type[u] ^= type[fa[u][0]];

	}

	for (int i = 1; i <= n; i++) {

		cout << (type[i] == 0 ? dp[i] : ans + 1 - dp[i]) << " ";

	}

	return 0;

}