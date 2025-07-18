#include<bits/stdc++.h>

#define lazy ios::sync_with_stdio(0);cin.tie(0); cout.tie(0);

#define inf 0x3f3f3f3f3f3f3f

#define endl '\n'

#define pii pair<int, int>

#define pll pair<ll ,ll>

using namespace std;



typedef long long ll;

const double eps = 1e-10;

const ll N = 4e5 + 7;

const ll M = 1e9 + 7;

const ll mod = 998244353;

ll t, n, ans[105];



void out() {

	for (int i = 1; i <= n; i++) {

		cout << ans[i] << " ";

	}

	cout << endl;

}



int main() {

	lazy;

	cin >> t;

	while (t--) {

		cin >> n;

		cout << n << endl;

		for (int i = 1; i <= n; i++)ans[i] = i;

		out();

		swap(ans[n], ans[1]);

		out();

		for (int i = 3; i <= n; i++) {

			swap(ans[i - 2], ans[i - 1]);

			out();

		}

	}

}

/*

7 5 2

2 7

*/