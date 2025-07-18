#include <iostream>

#include <fstream>

#include <stack>

#include <set>

#include <map>

#include <stack>

#include <vector>

#include <queue>

#include <string>

#include <algorithm>

#include <numeric>

#include <cmath>

#include <array>

#include <bitset>

#include <queue>

#include <cstring>

#include <iomanip>

#define int long long

#define all(v) begin(v), end(v)

#define ve vector

#define vi vector<int>

#define vd vector<double>

#define pb push_back

#define pii pair<int,int>

#define rep(i, n) for(int i = 0; i < (n); i++)

using namespace std;

using ll = long long;

using ull = unsigned long long;



const double pi = atan(1) * 4;



void fast() {

	ios::sync_with_stdio(0);

	cin.tie(0);

	cout.tie(0);

	cout << fixed; cout.precision(10);

}



string s[9];



void change(int x, int y) {

	if (s[x][y] == '9') {

		s[x][y] = '1';

	}

	else {

		s[x][y]++;

	}

}



void solve() {

	for (auto& i : s)

		cin >> i;

	change(0, 0);

	change(1, 3);

	change(2, 6);

	change(3, 1);

	change(4, 4);

	change(5, 7);

	change(6, 2);

	change(7, 5);

	change(8, 8);

	for (auto& i : s)

		cout << i << "\n";

}



signed main() {

#ifdef LOCAL

	freopen("local.in", "r", stdin);

	freopen("local.out", "w", stdout);

#endif

	fast();

	int T = 1;

	cin >> T;

	while (T--)

		solve();

	return 0;

}