#include <cstdio>
#include <cstdlib>
#include <cstring>
#include <cmath>
#include <climits>
#include <cfloat>
#include <ctime>
#include <cassert>
#include <map>
#include <utility>
#include <set>
#include <iostream>
#include <memory>
#include <string>
#include <vector>
#include <algorithm>
#include <functional>
#include <sstream>
#include <complex>
#include <stack>
#include <queue>
#include <numeric>
#include <list>
#include <iomanip>
#include <fstream>
#include <bitset>

#ifdef LOCAL
#include "local.h"
#endif

using namespace std;

#define rep(i, n) for (int i = 0; i < (int)(n); ++i)
#define all(a) (a).begin(), (a).end()
#define clr(a, x) memset(a, x, sizeof(a))
#define sz(a) ((int)(a).size())
#define mp(a, b) make_pair(a, b)
#define ten(n) ((long long)(1e##n))

template <typename T, typename U> void upmin(T& a, const U& b) { a = min<T>(a, b); }
template <typename T, typename U> void upmax(T& a, const U& b) { a = max<T>(a, b); }
template <typename T> void uniq(T& a) { sort(a.begin(), a.end()); a.erase(unique(a.begin(), a.end()), a.end()); }
template <class T> string to_s(const T& a) { ostringstream os; os << a; return os.str(); }
template <class T> T to_T(const string& s) { istringstream is(s); T res; is >> res; return res; }
void fast_io() { cin.tie(0); ios::sync_with_stdio(false); }
bool in_rect(int x, int y, int w, int h) { return 0 <= x && x < w && 0 <= y && y < h; }

typedef long long ll;
typedef pair<int, int> pint;

const int dx[] = { 0, 1, 0, -1 };
const int dy[] = { 1, 0, -1, 0 };


void fix_pre(int n) { cout.setf(ios::fixed, ios::floatfield); cout.precision(10); }
const double PI = acos(-1.0);
typedef double gtype;
typedef complex<gtype> Point;
gtype to_rad(gtype deg)
{
    return deg * PI / 180;
}
gtype cross(const Point& a, const Point& b)
{
	return a.real() * b.imag() - a.imag() * b.real();
}
struct Line
{
	Point first, second;
	Line() {}
    Line(const Point& first, const Point& second)
        : first(first), second(second)
    {
        if (first == second)
            this->first.real(first.real() + 1e-10);
    }
};
Point ip_LL(const Line& line1, const Line& line2)
{
	Point a = line1.second - line1.first, b = line2.second - line2.first;
	gtype p = cross(b, line2.first - line1.first);
	gtype q = cross(b, a);
	return line1.first + p / q * a;
}
int main()
{
    int n;
    double l, r;
    cin >> n >> l >> r;
    Point p[22];
    Point rot[22];
    rep(i, n)
    {
        double x, y, a;
        cin >> x >> y >> a;
        p[i] = Point(x, y);

        a = to_rad(a);
        rot[i] = Point(cos(a), sin(a));
    }

    static double dp[1 << 20];
    rep(i, 1 << n)
        dp[i] = l;
    rep(mask, 1 << n)
    {
        rep(i, n)
        {
            if (!(mask >> i & 1))
            {
                Point dir = Point(dp[mask], 0) - p[i];
                dir *= rot[i];

                double np;
                static const Line axis(Line(Point(0, 0), Point(1e8, 0)));
                if (cross(dir, Point(1e8, 0)) > 1e-6)
                {
                    Point ip = ip_LL(axis, Line(p[i], p[i] + dir));
                    np = ip.real();
                }
                else
                    np = r;
                upmin(np, r);

                upmax(dp[mask | (1 << i)], np);
            }
        }
    }

    double res = dp[(1 << n) - 1] - l;
    fix_pre(8);
    cout << res << endl;
}