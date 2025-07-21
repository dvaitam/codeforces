#include<bits/stdc++.h>
using namespace std;

typedef long long ll;
void sol() {
	int n, m;
	cin >> n >> m;
	unordered_set<int>se;		//	存储x的公因数.
	vector<int>cf;
	for (int i = 2;i * i <= m;i++) {
		if (m % i == 0) {
			cf.push_back(i);
			if (m / i != i) {
				cf.push_back(m / i);
			}
		}
	}

	sort(cf.begin(), cf.end());
	auto is = [&](int x)->bool {
		auto it = lower_bound(cf.begin(), cf.end(), x);
		if (it != cf.end()) {
			return (*it) == x;
		}
		return false;
	};

	vector<int>a(n + 1);
	for (int i = 1;i <= n;i++) {
		cin >> a[i];
	}

	int ans = 0;
	for (int i = 1;i <= n;i++) {
		if (a[i] == 1 || a[i] > m || !is(a[i])) continue;
		vector<int>tmp;
		for (int p : se) {
			int res = p * a[i];
			if (res == m) {
				ans++;
				se.clear();
				tmp.clear();
				break;
			}
			if (is(res)) {
				tmp.push_back(res);
			}
		}
		for (auto p : tmp) se.insert(p);
		se.insert(a[i]);
	}
	ans++;
	cout << ans << '\n';
}

int main() {
	ios::sync_with_stdio(0), cin.tie(0);
	int t;
	cin >> t;
	while (t--) {
		sol();
	}
	return 0;
}