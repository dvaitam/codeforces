#include <bits/stdc++.h>
using namespace std;
#define forr(i, a, b) for(int i = (int) a; i < (int) b; i++)
#define forn(i, n) forr(i, 0, n)
#define dforr(i, a, b) for(int i = (int) b-1; i >= (int) a; i--)
#define dforn(i, n) dforr(i, 0, n)
#define sz(x) ((int) x.size())
#define pb push_back
#define fst first
#define snd second
typedef long long ll;
typedef long double ld;
typedef pair<int, int> ii;
const int MAXN = -1;

bool isRegular(string &s) {
	int n = sz(s);
	
	int acc = 0;
	forn(i, n) {
		if (s[i] == '(') acc++;
		else             acc--;
		
		if (acc < 0) return false;
	}
	if (acc < 0) return false;

	return acc == 0;
}

bool isBalanced(string &s) {
	int n = sz(s);
	int dif = 0;
	forn(i, n) {
		if (s[i] == '(') dif++;
		else             dif--;
	}
	return dif == 0;
}

int n;
string s, r;

int main() {
	//~ freopen("input.txt", "r", stdin);
	//~ freopen("output.txt", "w", stdout);
	ios::sync_with_stdio(false);
	int tc; cin >> tc;
	forn(cs, tc) {
		cin >> n >> s;
		r = s; reverse(r.begin(), r.end());
		
		vector<int> ans(n, 1);
		if (isRegular(s) || isRegular(r)) {
			cout << 1 << endl;
			forn(i, n) {
				if (i) cout << " ";
				cout << ans[i];
			} cout << endl;
		} else if (isBalanced(s)) {
			forn(i, n/2) if (s[i] == ')') ans[i] = 2;
			forr(i, n/2, n) if (s[i] == '(') ans[i] = 2;
			
			cout << 2 << endl;
			forn(i, n) {
				if (i) cout << " ";
				cout << ans[i];
			} cout << endl;
		} else {
			cout << -1 << endl;
		}
	}
}