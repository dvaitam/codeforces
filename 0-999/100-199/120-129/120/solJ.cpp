#include <cstring>
#include <map>
#include <deque>
#include <queue>
#include <stack>
#include <sstream>
#include <iostream>
#include <iomanip>
#include <cstdio>
#include <cmath>
#include <cstdlib>
#include <ctime>
#include <algorithm>
#include <vector>
#include <set>
#include <complex>
#include <list>
#include <climits>
#include <cctype>
#include <bitset>

using namespace std;

#include <ext/hash_set>
#include <ext/hash_map>

using namespace __gnu_cxx;

#define pb push_back
#define all(v) v.begin(),v.end()
#define rall(v) v.rbegin(),v.rend()
#define sz(v) ((int)v.size())
#define rep(i,m) for(int i=0;i<(int)(m);i++)
#define rep2(i,n,m) for(int i=n;i<(int)(m);i++)
#define For(it,c) for(__typeof(c.begin()) it=c.begin();it!=c.end();++it)
#define mem(a,b) memset(a,b,sizeof(a))
#define mp make_pair
#define dot(a,b) ((conj(a)*(b)).X)
#define X real()
#define Y imag()
#define length(V) (hypot((V).X,(V).Y))
#define vect(a,b) ((b)-(a))
#define cross(a,b) ((conj(a)*(b)).imag())
#define normalize(v) ((v)/length(v))
#define rotate(p,about,theta) ((p-about)*exp(point(0,theta))+about)
#define pointEqu(a,b) (comp(a.X,b.X)==0 && comp(a.Y,b.Y)==0)

typedef stringstream ss;
typedef long long ll;
typedef pair<int, int> pii;
typedef vector<pii> vpii;
typedef vector<string> vs;
typedef vector<int> vi;
typedef vector<double> vd;
typedef vector<vector<int> > vii;
typedef long double ld;

const int oo = (int) 1e9;
const double PI = 2 * acos(0.0);
const double eps = 1e-9;

inline int comp(const double &a, const double &b) {
    if (fabs(a - b) < eps)
        return 0;
    return a > b ? 1 : -1;
}

struct point {
    int x, y, ind;
    point() {
    }
    point(int xx, int yy, int ii) {
        x = xx;
        y = yy;
        ind = ii;
    }
    bool operator <(const point &p) const {
        if (x != p.x)
            return x < p.x;
        return y < p.y;
    }
};

int n, xs[100009], ys[100009];

void print(int a, int b) {
    bool xx = 0;
    if ((xs[a] < 0 && xs[b] < 0) || (xs[a] > 0 && xs[b] > 0))
        xx = 1;
    bool yy = 0;
    if ((ys[a] < 0 && ys[b] < 0) || (ys[a] > 0 && ys[b] > 0))
        yy = 1;
    cout << a + 1 << " ";
    if (!xx && !yy)
        cout << 1;
    else if (xx && !yy)
        cout << 2;
    else if (!xx && yy)
        cout << 3;
    else
        cout << 4;
    cout << " " << b + 1 << " 1" << endl;

}

int main() {
    freopen("input.txt", "rt", stdin);
    freopen("output.txt", "wt", stdout);
    scanf("%d", &n);
    vector<point> v;
    map<point, int> m;
    rep(i,n) {
        scanf("%d%d", &xs[i], &ys[i]);
        point p(abs(xs[i]), abs(ys[i]), i);
        if (m.count(p)) {
            print(m[p], i);
            return 0;
        }
        m[p] = i;
        v.pb(p);
    }
    sort(all(v));
    double mn = oo;
    int a, b;
    map<double, set<point> > live;
    int last = 0, cur = 0;
    while (cur < n) {
        while (last < cur && v[cur].x - v[last].x > mn) {
            live[v[last].y].erase(v[last]);
            last++;
        }
        map<double, set<point> >::iterator it1 =
                live.lower_bound(v[cur].y - mn);
        map<double, set<point> >::iterator it2 =
                live.upper_bound(v[cur].y + mn);
        while (it1 != it2) {
            For(it,it1->second) {
                double d = hypot(v[cur].x - it->x, v[cur].y - it->y);
                if (comp(d, mn) < 0) {
                    mn = d;
                    a = v[cur].ind;
                    b = it->ind;
                }
            }
            it1++;
        }
        live[v[cur].y].insert(v[cur]);
        cur++;
    }
    print(a, b);
    return 0;
}