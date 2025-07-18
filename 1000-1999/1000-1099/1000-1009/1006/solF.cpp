#include <bits/stdc++.h>

using namespace std;

#include <ext/pb_ds/assoc_container.hpp>
#include <ext/pb_ds/assoc_container.hpp>
#include <ext/pb_ds/tree_policy.hpp>

using namespace __gnu_pbds;

#define rand() ((rand()<<15)+rand())

template<typename Key, typename Value>
using hash_map = cc_hash_table<Key, Value>;
template<typename T, typename Cmp = less<T>>
using oset = tree<T, null_type, Cmp, rb_tree_tag, tree_order_statistics_node_update>;

typedef long long lint;
typedef complex<double> point;

const int INF = 0x3f3f3f3f;
const lint INFL = 0x3f3f3f3f3f3f3f3fLL;
const double E = 1e-9, pi = 2 * acos(0);

template<typename T, typename U> void amin(T &x, U y) { if (y < x) x = y; }
template<typename T, typename U> void amax(T &x, U y) { if (x < y) x = y; }

vector<vector<lint>> A;
vector<vector<unordered_map<lint, int>>> C;
lint ans, k;

void dfs(int r, int c, int d, int inc, int t, lint sum) {
	if (r < 0 or r >= (int)A.size() or c < 0 or c >= (int)A[r].size()) {
		return;
	}
	sum ^= A[r][c];
	if (d == t) {
		if (inc == 1) {
			C[r][c][sum]++;
		} else {
			if (r + inc >= 0) {
				ans += C[r + inc][c][sum ^ k];
			}
			if (c + inc >= 0) {
				ans += C[r][c + inc][sum ^ k];
			}
		}
	} else {
		dfs(r + inc, c, d + 1, inc, t, sum);
		dfs(r, c + inc, d + 1, inc, t, sum);
	}
}

bool solve(int testNumber) {
	int n, m;
	if (!(cin >> n >> m >> k)) {
		return false;
	}
	A.assign(n, vector<lint>(m));
	C.assign(n, vector<unordered_map<lint, int>>(m));
	for (auto &vec: A) {
		for (lint &x: vec) {
			cin >> x;
		}
	}
	int d = min(n, m);
	if (n * m == 1) {
		ans = A[0][0] == k;
	} else {
		ans = 0;
		dfs(0, 0, 0, 1, d - 1, 0);
		dfs(n - 1, m - 1, 0, -1, n + m - 2 - d, 0);
	}
	printf("%lld\n", ans);
	return true;
}

void init(const char *in = nullptr, const char *out = nullptr) {
	if (in) freopen(in, "r", stdin);
	if (out) freopen(out, "w", stdout);
	srand(unsigned((long long)new char));
}

int main() {
	ios_base :: sync_with_stdio(false);
	cin.tie(nullptr);
	init();
	for (int i=1; solve(i); i++);
	return 0;
}