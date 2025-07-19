#include <bits/stdc++.h>

using namespace std;
constexpr auto MOD = 1'000'000'007LL;
long long t, b, a, i;
string s;
int main() {
	for (cin >> t; t-- and cin >> s; cout << a << '\n') {
		map M = {pair{0LL, 1LL}};
		i = a = b = 0;
		for (char c : s) {
			b += c == '1' ? +1 : -1;
			a = (a + (s.size() - i) * M[b]) % MOD;
			M[b] = (M[b] + (++i) + 1) % MOD;
		}
	}
