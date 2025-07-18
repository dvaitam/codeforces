// Copyright © 2022 Youngmin Park. All rights reserved.

//#pragma GCC optimize("O3")

//#pragma GCC target("avx2")

#include <bits/stdc++.h>

using namespace std;



using ll = long long;

using vi = vector<int>;

using pii = pair<int, int>;

using vpi = vector<pii>;

using pll = pair<ll, ll>;

using vl = vector<ll>;

using vpl = vector<pll>;

using ld = long double;

template <typename T, size_t SZ>

using ar = array<T, SZ>;

template <typename T>

using pqg = priority_queue<T, vector<T>, greater<T>>;



#define all(v) (v).begin(), (v).end()

#define pb push_back

#define sz(x) (int)(x).size()

#define fi first

#define se second

#define lb lower_bound

#define ub upper_bound



constexpr int INF = 1e9;

constexpr ll LINF = 1e18;

const ld PI = acos((ld)-1.0);



mt19937_64 rng(chrono::steady_clock::now().time_since_epoch().count());

template <typename T>

constexpr bool ckmin(T &a, const T &b) { return b < a ? a = b, 1 : 0; }

template <typename T>

constexpr bool ckmax(T &a, const T &b) { return b > a ? a = b, 1 : 0; }

ll cdiv(ll a, ll b) { return a / b + ((a ^ b) > 0 && a % b); } // divide a by b rounded up

ll fdiv(ll a, ll b) { return a / b - ((a ^ b) < 0 && a % b); } // divide a by b rounded down



#ifdef LOCAL

#include "miscellaneous/debug.h"

#else

#define dbg(...) 42

#endif



inline namespace RecursiveLambda {

	template <typename Fun>

	struct y_combinator_result {

		Fun fun_;

		template <typename T>

		explicit constexpr y_combinator_result(T &&fun) : fun_(forward<T>(fun)) {}

		template <typename... Args>

		constexpr decltype(auto) operator()(Args &&...args) const {

			return fun_(ref(*this), forward<Args>(args)...);

		}

	};

	template <typename Fun>

	decltype(auto) y_combinator(Fun &&fun) {

		return y_combinator_result<decay_t<Fun>>(forward<Fun>(fun));

	}

};



constexpr int dx[] = {1, 0, -1, 0, 1, 1, -1, -1}, dy[] = {0, 1, 0, -1, 1, -1, 1, -1};



void solve() {

	int n, m; cin >> n >> m;

	vector<vi> a(n, vi(m)), vis(n, vi(m));

	queue<pii> q;

	vector<ar<int, 3>> ans;

	for (auto &v : a) for (auto &e : v) cin >> e;

	auto count_sq = [&](int r, int c) -> int {

		if (r == n - 1 || c == m - 1 || r < 0 || c < 0) return -1;

		vi t{a[r][c], a[r][c + 1], a[r + 1][c], a[r + 1][c + 1]};

		int cnt = 0, zr = 0;

		for (auto &e : t) {

			if (e) {

				int x{};

				for (auto &f : t) {

					if (f == e) {

						x++;

					}

				}

				ckmax(cnt, x);

			}else{

				zr++;

			}

		}

		return cnt + zr;

	};

	for (int i = 0; i < n - 1; i++) {

		for (int j = 0; j < m - 1; j++) {

			if (count_sq(i, j) == 4) {

				q.push({i, j});

				vis[i][j] = 1;

			}

		}

	}

	while (sz(q)) {

		auto [r, c] = q.front(); q.pop();

		int col{};

		if (a[r][c]) col = a[r][c];

		if (a[r+1][c]) col = a[r+1][c];

		if (a[r][c+1]) col = a[r][c+1];

		if (a[r+1][c+1]) col = a[r+1][c+1];

		if (col) ans.pb({r, c, col});

		else continue;

		a[r][c] = a[r+1][c] = a[r][c+1] = a[r+1][c+1] = 0;

		for (int i = 0; i < 8; i++) {

			int nr = r + dx[i], nc = c + dy[i];

			if (count_sq(nr, nc) == 4 && !vis[nr][nc]) {

				q.push({nr, nc});

				vis[nr][nc] = 1;

			}

		}

	}

	for (auto &v : a) {

		if (count(all(v), 0) != m) {

			cout << -1 << '\n';

			return;

		} 

	}

	reverse(all(ans));

	cout << sz(ans) << '\n';

	for (int i = 0; i < sz(ans); i++) {

		auto &[r, c, col] = ans[i];

		cout << r + 1 << ' ' << c + 1 << ' ' << col << '\n';

	}

}



int main() {

	cin.tie(0)->sync_with_stdio(0);

	cin.exceptions(cin.failbit);

	int testcase = 1;

	// cin >> testcase;

	while (testcase--) {

		solve();

	}

#ifdef LOCAL

	cerr << "Time elapsed: " << 1.0 * (double)clock() / CLOCKS_PER_SEC << " s.\n";

#endif

}