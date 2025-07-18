#include<bits/stdc++.h>

using namespace std;

using ll = long long;



#define int			long long

#define ff			first

#define ss			second

#define pob			pop_back

#define pof			pop_front

#define pb			push_back

#define pf			push_front

#define lb			lower_bound

#define ub			upper_bound

#define pii			pair<int,int>

#define sz(x)			(int)(x.size())

#define all(x)			x.begin(),x.end()

#define PI			3.1415926535897932384626

#define precise(i)		cout<<fixed<<setprecision(i)

#define uniq(v)           	(v).erase(unique(all(v)),(v).end())



#ifndef ONLINE_JUDGE

#include "debug.h"

#else

#define deb(x...)

#endif



const int mod = 1e9 + 7;



void solve() {

	int n;

	cin >> n;

	vector<int> a(n), b(n);



	for (auto &x : a) cin >> x;

	for (auto &x : b) cin >> x;



	if (is_sorted(all(a)) and is_sorted(all(b))) {

		cout << "0\n";

		return;

	}

	vector<pii> operation;

	for (int i = 0; i <= n - 2; i++) {

		int index = -1, prev = (i ? b[i - 1] : 0);



		int Min_a = *min_element(a.begin() + i, a.end());

		int Min_b = *min_element(b.begin() + i, b.end());



		for (int j = i + 1; j <= n - 1; j++) {

			if (a[j] == Min_a and b[j] == Min_b) {

				index = j;

				break;

			}

		}

		if (index != -1) {

			swap(a[i], a[index]);

			swap(b[i], b[index]);

			operation.pb({i, index});

		}

	}

	if (is_sorted(all(a)) and is_sorted(all(b))) {

		cout << sz(operation) << "\n";

		for (auto &[x, y] : operation) cout << x + 1 << " " << y + 1 << "\n";

	}

	else cout << "-1\n";

}



signed main() {



	ios_base::sync_with_stdio(false);

	cin.tie(NULL); cout.tie(NULL);



	//freopen("input.txt","r",stdin);

	//freopen("output.txt","w",stdout);



	int TESTS = 1;

	cin >> TESTS;

	for (int T = 1; T <= TESTS; T++) {

		solve();

	}



	return 0;

}