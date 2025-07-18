// .... .... .....!
// ...... ......!
// .... ....... ..... ..!
// ...... ... ... .... ... .... .....!
// ... .. ... .... ...?

#include<bits/stdc++.h>
using namespace std;

#define rep(i, n) for (int i = 0, _n = (int)(n); i < _n; i++)
#define fer(i, x, n) for (int i = (int)(x), _n = (int)(n); i < _n; i++)
#define rof(i, n, x) for (int i = (int)(n), _x = (int)(x); i-- > _x; )
#define sz(x) (int((x).size()))
#define pb push_back
#define all(X) (X).begin(),(X).end()
#define X first
#define Y second
//#define endl '\n'

template<class P, class Q> inline void smin(P &a, Q b) { if (b < a) a = b; }
template<class P, class Q> inline void smax(P &a, Q b) { if (a < b) a = b; }

typedef long long ll;
typedef pair<int, int> pii;

////////////////////////////////////////////////////////////////////////////////

const int maxn = 1000 + 100;

int n, k;
string s;
int t;

bool can(int pos, int lo, int hi) {
	if(pos < n) smax(lo, -k + 1), smin(hi, k - 1);
	if(lo > hi) return false;
	if(pos == 0) {
		if(lo <= 0 && 0 <= hi) { t = 0; return true; }
		return false;
	}
	pos--;
	if(s[pos] == 'L') if(can(pos, lo + 1, hi + 1)) { t--; return true; }
	if(s[pos] == 'W') if(can(pos, lo - 1, hi - 1)) { t++; return true; }
	if(s[pos] == 'D') if(can(pos, lo + 0, hi + 0)) { return true; }
	if(s[pos] == '?') if(can(pos, lo - 1, hi + 1)) { 
		if(lo <= t && t <= hi) s[pos] = 'D';
		else if(t <= lo) s[pos] = 'W', t++;
		else if(t >= hi) s[pos] = 'L', t--;
		return true;
	}

	return false;
}

int main() {
	ios_base::sync_with_stdio(false); cin.tie(0);

	cin >> n >> k;
	cin >> s;

	if(can(n, k, k)) { cout << s << endl; return 0; }
	if(can(n, -k, -k)) { cout << s << endl; return 0; }
	cout << "NO" << endl;

	return 0;
}