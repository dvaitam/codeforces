#include <bits/stdc++.h>
using namespace std;
#define Go  ios_base::sync_with_stdio(false);	cin.tie(NULL)
typedef long long ll;
#define F first
#define S second
#define PB push_back
#define all(x)   x.begin(), x.end()
#define rall(v)  v.rbegin(), v.rend()
#define getunique(v) {sort(v.begin(), v.end()); v.erase(unique(v.begin(), v.end()), v.end());}


void solve() {
	int n; cin >> n; 
	vector<int> cnt(101); 
	vector<int> v(n); 
	for (auto& I : v)	 cin >> I, cnt[I]++; 
	int c = 0; 
	vector<int> x;
	for (int i = 1; i <= 100; i++) {
		if (cnt[i] > 1) {
			c++;
			if (c <= 2) x.push_back(i);
		}
	}
	if (c < 2) {
		cout << -1 << endl; 
		return;
	}
	vector<int> b(n + 1, 1); 
	vector<int> cur(101); 
	for (int i = 0; i < n; i++) {
		cur[v[i]]++;
		if (cur[v[i]] == 1 || v[i] != x[0]) continue;
		b[i] = 2;
	}
	fill(cur.begin(), cur.end(), 0);
	for (int i = 0; i < n; i++) {
		cur[v[i]]++;
		if (cur[v[i]] == 1 || v[i] != x[1]) continue;
		b[i] = 3;
	}
	for (int i = 0; i < n; i++) cout << b[i] << ' ';
	cout << endl;
}


int main() {
	Go;

	int TES = 1;		   cin >> TES;

	while (TES--) {
		solve();
	}

	return 0;
}