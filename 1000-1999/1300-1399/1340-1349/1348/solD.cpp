#include <algorithm>

#include <array>

#include <bitset>

#include <cmath>

#include <complex>

#include <cstdlib>

#include <iomanip>

#include <iostream>

#include <map>

#include <numeric>

#include <queue>

#include <limits.h> 

#include <random>

#include <set>

#include <string>

#include <unordered_map>

#include <unordered_set>

#include <vector>

#define ll long long

#define int long long 

using namespace std;

template <typename T>

inline T gcd(T a, T b) { while (b != 0) swap(b, a %= b); return a; } 





void solve() {

	int n; cin >> n;

	int days = 0;

	while((1 << days) - 1 < n) days++;

	days--; 

	cout << days << "\n"; 

	int last = 1; 

	n--; 

	for(int i = days; i; i--) {

		// cout << "n: " << n << "\n"; 

		int lb = last, ub = last * 2;

		int ans = ub;

		while(lb <= ub) {

			int m = (lb + ub) / 2;

			int mn = i * m, mx = ((1 << i) - 1) * m; 			

			// cout << m << " " << mn << " " << mx << "\n"; 

			if(mn > n) 	ub = m - 1; 

			else if(mx < n) lb = m + 1;

			else {

				ans = m;

				break;

			}

		}

		// cout << "ans: "  << ans << "\n"; 

		n -= ans; 

		cout << ans - last << " ";

		last = ans; 

	}

	cout << "\n"; 

}







signed main() {

	ios::sync_with_stdio(false);

	cin.tie(0);

	int tt = 1; 

	cin >> tt;



	for(int cas = 1; cas <= tt; cas++) {

		solve();

	}

}