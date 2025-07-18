#define _CRT_SECURE_NO_WARNINGS

#pragma comment(linker, "/STACK:268435456")



#include <cstdio>

#include <iostream>

#include <iomanip>

#include <map>

#include <queue>

#include <set>

#include <stack>

#include <string>

#include <vector>

#include <algorithm>

#include <cassert>

#include <memory.h>

#include <ctime>

#include <cmath>

#include <complex>



//#include <unordered_set>

//#include <unordered_map>



using namespace std;



#define forn(i, n) for(int i = 0; i < int(n); ++i)

#define for1(i, n) for(int i = 1; i < int(n); ++i)

#define fork(i, k, n) for(int i = int(k); i <= int(n); ++i)

#define forft(i, from, to) for(int i = int(from); i < int(to); ++i)

#define forr(i, n) for(int i = int(n) - 1; i >= 0; --i)

#define pb push_back

#define mp make_pair

#define mnp(a, b) make_pair((a) < (b) ? (a) : (b), (a) < (b) ? (b) : (a))

#define sz(d) int(d.size())

#define all(a) a.begin(), a.end()

#define ms(a, v) memset(a, v, sizeof(a))

#define X first

#define Y second

#define correct(x, y, xmax, ymax) ((x) >= 0 && (x) < (xmax) && (y) >= 0 && (y) < (ymax))



template<typename T> T sqr(const T &x) {

	return x * x;

}



template<typename T> T my_abs(const T &x) {

	return x < 0 ? -x : x;

}



typedef long long li;

typedef long long ll;

typedef unsigned long long ull;

typedef unsigned long long uli;

typedef long double ld;

typedef pair<int,int> pt;

typedef pair<ld,ld> pd;



const int INF = (int)1e9;

const li LINF = (li)6e18;

const li INF64 = LINF;

const li INFLL = LINF;

const ld EPS = 1e-7;

const ld PI = 3.1415926535897932384626433832795;

const ld ME = 1e-5;

const li MOD = (int)1e9 + 7;

const li MOD2 = (int)1e9 + 21;

const int dx[] = {-1, 0, 1, 0};

const int dy[] = {0, 1, 0, -1};

const int SQN = 500;

const int LOGN = 17;



const int N = (int)1e3 + 10;



struct line {

	ld A, B, C;

};



inline line get_line(const pd& a, const pd& b) {

	line result;

	result.A = b.Y - a.Y;

	result.B = a.X - b.X;

	result.C = -result.A * a.X - result.B * a.Y;

	return result;

}



inline pd operator +(const pd& a, const pd& b) {

	return pd(a.X + b.X, a.Y + b.Y);

}



inline pd operator -(const pd& a, const pd& b) {

	return pd(a.X - b.X, a.Y - b.Y);

}



inline pd operator *(const pd& a, ld k) {

	return pd(a.X * k, a.Y * k);

}



inline ld cp(const pd& a, const pd& b) {

	return a.X * b.Y - a.Y * b.X;

}



inline bool intersect(const line& a, const line& b, pd& res) {

	ld d = cp(pd(a.A, a.B), pd(b.A, b.B));

	if (abs(d) < EPS) {

		return false;

	}

	res.X = -cp(pd(a.C, a.B), pd(b.C, b.B)) / d;

	res.Y = -cp(pd(a.A, a.C), pd(b.A, b.C)) / d;

	return true;

}



inline ld dist(const pd& a, const pd& b) {

	return sqrt(sqr(a.X - b.X) + sqr(a.Y - b.Y));

}



inline bool equals(const pd& a, const pd& b) {

	return dist(a, b) < EPS;

}



inline ld get_area(vector<pd>& p) {

	if (sz(p) < 3) {

		return 0;

	}

	pd c(0, 0);

	forn(i, sz(p)) {

		c = c + (p[i] * (1.0 / (ld)sz(p)));

	}

	forn(i, sz(p)) {

		p[i] = p[i] - c;

	}

	vector<pair<ld, pd> > pp;

	forn(i, sz(p)) {

		pp.pb(mp((ld)atan2(p[i].Y, p[i].X), p[i]));

	}

	sort(all(pp));

	forn(i, sz(p)) {

		p[i] = pp[i].Y;

	}

	p.erase(unique(all(p), equals), p.end());

	ld res = 0;

	forn(i, sz(p)) {

		int j = (i + 1) % sz(p);

		res += cp(p[i] - c, p[j] - c);

	}

	return abs(res) / 2.0;

}



vector<pd> s;



ld ans = 0;



inline void solve() {

	int n;

	ld h, f;

	cin >> n >> h >> f;

	forn(i, n) {

		int l, r;

		cin >> l >> r;

		if (l < 0 && r > 0) {

			s.pb(pd(l, 0));

			s.pb(pd(0, r));

		} else {

			s.pb(pd(l, r));

		}

	}

	pd cu = pd(0, f);

	pd cd = pd(0, -f);

	line ul = get_line(pd(0, h), pd(1, h));

	line dl = get_line(pd(0, -h), pd(1, -h));

	forn(i, sz(s)) {

		pd lp = pd(s[i].X, -h);

		pd rp = pd(s[i].Y, -h);

		line l = get_line(cd, lp);

		line r = get_line(cd, rp);

		pd ulp;

		intersect(l, ul, ulp);

		pd urp;

		intersect(r, ul, urp);

		ans += h * (s[i].Y - s[i].X + urp.X - ulp.X);

	}

	ans *= 2.0;

	forn(i, sz(s)) {

		forn(j, i + 1) {

			vector<pd> pts;

			pd lp1 = pd(s[i].X, -h);

			pd rp1 = pd(s[i].Y, -h);

			line l1 = get_line(cd, lp1);

			line r1 = get_line(cd, rp1);

			pd ulp1;

			intersect(l1, ul, ulp1);

			pd urp1;

			intersect(r1, ul, urp1);

			pd lp2 = pd(s[j].X, h);

			pd rp2 = pd(s[j].Y, h);

			line l2 = get_line(cu, lp2);

			line r2 = get_line(cu, rp2);

			pd dlp2;

			intersect(l2, dl, dlp2);

			pd drp2;

			intersect(r2, dl, drp2);

			if (ulp1.X >= lp2.X && ulp1.X <= rp2.X) {

				pts.pb(ulp1);

			}

			if (urp1.X >= lp2.X && urp1.X <= rp2.X) {

				pts.pb(urp1);

			}

			if (dlp2.X >= lp1.X && dlp2.X <= rp1.X) {

				pts.pb(dlp2);

			}

			if (drp2.X >= lp1.X && drp2.X <= rp1.X) {

				pts.pb(drp2);

			}

			if (lp2.X >= ulp1.X && lp2.X <= urp1.X) {

				pts.pb(lp2);

			}

			if (rp2.X >= ulp1.X && rp2.X <= urp1.X) {

				pts.pb(rp2);

			}

			if (lp1.X >= dlp2.X && lp1.X <= drp2.X) {

				pts.pb(lp1);

			}

			if (rp1.X >= dlp2.X && rp1.X <= drp2.X) {

				pts.pb(rp1);

			}

			pd cll;

			if (intersect(l1, l2, cll) && cll.Y >= -h && cll.Y <= h) {

				pts.pb(cll);

			}

			pd clr;

			if (intersect(l1, r2, clr) && clr.Y >= -h && clr.Y <= h) {

				pts.pb(clr);

			}

			pd crl;

			if (intersect(r1, l2, crl) && crl.Y >= -h && crl.Y <= h) {

				pts.pb(crl);

			}

			pd crr;

			if (intersect(r1, r2, crr) && crr.Y >= -h && crr.Y <= h) {

				pts.pb(crr);

			}

			if (i == j) {

				ans -= get_area(pts);

			} else {

				ans -= 2.0 * get_area(pts);

			}

		}

	}

	cout << ans << endl;

}



int main() {

	//ios_base::sync_with_stdio(false);

	//freopen("input.txt", "r", stdin);

	//freopen("output.txt", "w", stdout);



	srand((unsigned int)time(NULL));



	cout << setprecision(15) << fixed;



	solve();



	cerr << "time: " << clock() << " ms" << endl;



	return 0;

}