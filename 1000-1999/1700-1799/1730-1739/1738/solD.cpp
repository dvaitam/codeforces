#include <bits/stdc++.h>

using namespace std;

#define endl '\n'

#define int long long

#define lowbit(x) (x&(-x))

#define x first

#define y second

const int N = 2e5 + 10;

int gcd(int a, int b) { return b ? gcd(b, a % b) : a; }

bool flag = false;

int n, b[N], k;

vector<int> v[N];

vector<int> a;

void slove() {

	cin >> n; k = 0; a.clear();

	for (int i = 0; i <= n + 1; i++)v[i].clear();

	for (int i = 1; i <= n; i++){

		cin >> b[i];

		k += b[i] > i;

		v[b[i]].push_back(i);

	}

	int cur = 0, cnt = 0;

	if (v[n + 1].size()) cur = n + 1;

	while (cnt < n){

		cnt += v[cur].size();

		for (auto& it : v[cur])if (v[it].size())swap(it, v[cur].back());

		a.insert(a.end(), v[cur].begin(), v[cur].end());

		cur = v[cur].back();

	}

	cout << k << '\n';

	for (auto it : a)cout << it << " ";cout << '\n';

}

signed main() {

	ios::sync_with_stdio(false);

	cin.tie(0); cout.tie(0);

	int T = 1; cin >> T;

	while (T--) slove();

}