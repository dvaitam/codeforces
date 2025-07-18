#include <algorithm>
#include <iostream>
#include <valarray>
#include <iomanip>
#include <fstream>
#include <sstream>
#include <cstdlib>
#include <cstring>
#include <cassert>
#include <numeric>
#include <complex>
#include <cstdio>
#include <string>
#include <vector>
#include <bitset>
#include <ctime>
#include <cmath>
#include <queue>
#include <stack>
#include <deque>
#include <map>
#include <set>

using namespace std;

#define FOREACH(i, c) for(__typeof((c).begin()) i = (c).begin(); i != (c).end(); ++i)
#define FOR(i, a, n) for (int i = (a); i < int(n); ++i)
#define error(x) cout << #x << " = " << (x) << endl;
#define all(n) (n).begin(), (n).end()
#define Size(n) ((int)(n).size())
#define mk make_pair
#define pb push_back
#define F first
#define S second
#define X real()
#define Y imag()

typedef long long ll;
typedef pair<int, int> pii;
typedef pair<ll, ll> pll;
typedef complex<double> point;

template <class P, class Q> void smin(P &a, Q b) { if (b < a) a = b; }
template <class P, class Q> void smax(P &a, Q b) { if (b > a) a = b; }
template <class P, class Q> bool in(const P &a, const Q &b) { return a.find(b) != a.end(); }

int c(int h, int w) {
	int tot = h-1;
	FOR(i, 0, h) {
		if (i%2 == 0) tot += w-1;
		else {
			if (i+1 < h) tot += 2*(w-1);
			else tot += w-1;
		}
	}
	return tot;
}

const int MAXN = 1024;
int h, w, k, board[MAXN][MAXN];
string row[MAXN], col[MAXN];

void fill1(int i, int j) {
	if (j && board[i][j-1]) {
		if (row[i][j-1] == 'E') board[i][j] = board[i][j-1]; else board[i][j] = 3-board[i][j-1];
	} else if (i && board[i-1][j]) {
		if (col[i-1][j] == 'E') board[i][j] = board[i-1][j]; else board[i][j] = 3-board[i-1][j];
	}
}

void fill2(int r, int c) {
	int poll[] = {0, 0, 0};
	if (r > 0 && board[r-1][c]) poll[col[r-1][c]=='E'?board[r-1][c]:3-board[r-1][c]]++;
	if (c > 0 && board[r][c-1]) poll[row[r][c-1]=='E'?board[r][c-1]:3-board[r][c-1]]++;
	if (r+1 < h && board[r+1][c]) poll[col[r][c]=='E'?board[r+1][c]:3-board[r+1][c]]++;
	if (c+1 < w && board[r][c+1]) poll[row[r][c]=='E'?board[r][c+1]:3-board[r][c+1]]++;
	board[r][c] = poll[1] >= poll[2]?1:2;
}

int main() {
	ios::sync_with_stdio(false);
	cin >> h >> w >> k;
	if (k == 1) {
		int tot = 0;
		FOR(i, 0, 2*h-1) {
			string s;
			cin >> s;
			tot += count(all(s), 'E');
		}
		if (tot*4 >= 3*(w*(h-1)+h*(w-1))) {
			cout << "YES" << endl;
			FOR(i, 0, h) {
				FOR(j, 0, w) cout << 1 << " ";
				cout << "\n";
			}
		} else cout << "NO" << endl;
		return 0;
	}
	FOR(i, 0, 2*h-1) {
		if (i%2) cin >> col[i/2];
		else cin >> row[i/2];
	}
	board[0][0] = 1;
	if (c(h, w) >= c(w, h)) {
		FOR(i, 1, h) fill1(i, 0);
		FOR(i, 0, h) if (i%2 == 0) FOR(j, 1, w) fill1(i, j);
		FOR(i, 0, h) if (i%2 == 1) FOR(j, 1, w) fill2(i, j);
	} else {
		FOR(i, 1, w) fill1(0, i);
		FOR(j, 0, w) if (j%2 == 0) FOR(i, 1, h) fill1(i, j);
		FOR(j, 0, w) if (j%2 == 1) FOR(i, 1, h) fill2(i, j);
	}
	cout << "YES" << endl;
	FOR(i, 0, h) {
		FOR(j, 0, w) cout << board[i][j] << " ";
		cout << "\n";
	}
	return 0;
}