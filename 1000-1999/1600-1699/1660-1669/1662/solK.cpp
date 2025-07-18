#pragma GCC optimize "O3"

#include <bits/stdc++.h>

using namespace std;



#define FOR(i, b, e) for(int i = (b); i < (e); i++)

#define SZ(x) ((int)x.size())

#define PB push_back

#define X first

#define Y second

#define f X

#define s Y

using vi = vector<int>;

using ll = long long;

using k = double;

using kk = pair<k, k>;



kk PT[3];

constexpr k PI = acos(-1.0);



const double EPS = 1E-9;



struct pt {

    double x, y;



    pt(){}

    pt(kk pkt) : x(pkt.f), y(pkt.s) {}



    bool operator<(const pt& p) const

    {

        return x < p.x - EPS || (abs(x - p.x) < EPS && y < p.y - EPS);

    }

};



struct line {

    double a, b, c;



    line() {}

    line(pt p, pt q){

        a = p.y - q.y;

        b = q.x - p.x;

        c = -a * p.x - b * p.y;

        norm();

    }



    void norm()

    {

        double z = sqrt(a * a + b * b);

        if (abs(z) > EPS)

            a /= z, b /= z, c /= z;

    }



    double dist(pt p) const { return a * p.x + b * p.y + c; }

};



double det(double a, double b, double c, double d)

{

    return a * d - b * c;

}



double det(kk a, kk b)

{

    return det(a.f, a.s, b.f, b.s);

}



inline bool betw(double l, double r, double x)

{

    return min(l, r) <= x + EPS && x <= max(l, r) + EPS;

}



inline bool intersect_1d(double a, double b, double c, double d)

{

    if (a > b)

        swap(a, b);

    if (c > d)

        swap(c, d);

    return max(a, c) <= min(b, d) + EPS;

}



bool intersect(pt a, pt b, pt c, pt d, pt& left, pt& right)

{

    if (!intersect_1d(a.x, b.x, c.x, d.x) || !intersect_1d(a.y, b.y, c.y, d.y))

        return false;

    line m(a, b);

    line n(c, d);

    double zn = det(m.a, m.b, n.a, n.b);

    if (abs(zn) < EPS) {

        if (abs(m.dist(c)) > EPS || abs(n.dist(a)) > EPS)

            return false;

        if (b < a)

            swap(a, b);

        if (d < c)

            swap(c, d);

        left = max(a, c);

        right = min(b, d);

        return true;

    } else {

        left.x = right.x = -det(m.c, m.b, n.c, n.b) / zn;

        left.y = right.y = -det(m.a, m.c, n.a, n.c) / zn;

        return betw(a.x, b.x, left.x) && betw(a.y, b.y, left.y) &&

               betw(c.x, d.x, left.x) && betw(c.y, d.y, left.y);

    }

}



kk rotate(kk pkt, k angle){

    return {pkt.f * cos(angle) - pkt.s * sin(angle), pkt.f * sin(angle) + pkt.s * cos(angle)};

}



kk operator - (kk a, kk b){

    return {a.f - b.f, a.s - b.s};

}



kk operator + (kk a, kk b){

    return {a.f + b.f, a.s + b.s};

}



ostream & operator << (ostream & os, kk & pkt){

    os << pkt.f << ' ' << pkt.s;

    return os;

}



k dist(kk &a, kk &b) {

    return sqrt((a.X - b.X) * (a.X - b.X) + (a.Y - b.Y) * (a.Y - b.Y));

}



kk fermat_point(kk a, kk b, kk c) {

    // cout << a << "     " << b << "     " << c << endl;



    if(det({b - a}, {c - a}) > 0){

        // cout << "HYC SWAP " << endl;

        swap(b, c);

    }



    auto f = [&](kk x, kk y, kk z){

        kk bok = z - y;

        bok = rotate(bok, PI / 3.0);

        kk PKT = y + bok;



        return make_pair(PKT, x);

    };



    pair<kk, kk> odc1 = f(a, b, c);

    pair<kk, kk> odc2 = f(b, c, a);



    // cout << "ODCINEK (" << odc1.f << ';' << odc1.s << ")         (" << odc2.f << ';' << odc2.s << ")" << endl; 



    pt A, B;

    if(not intersect(pt(odc1.f), pt(odc1.s), pt(odc2.f), pt(odc2.s), A, B)){

        vector<kk> punkty = {a, b, c};

        kk ret;

        k min_len = 1e18;

        for(int i = 0; i < 3; i++){

            k now = dist(punkty[i], punkty[(i - 1 + 3) % 3]) + dist(punkty[i], punkty[(i + 1) % 3]);

            if(now < min_len){

                min_len = now;

                ret = punkty[i];

            }

        }

        return ret;

    }

    // cout << "WYNIK POLICZYŁ SIĘ  " << A.x << ' ' << A.y << endl;

    // cout << endl;

    return make_pair(A.x, A.y);

}





k ans(kk &a) {

    k ret = 0;

    FOR(i, 0, 3) {

        kk sr = fermat_point(a, PT[i], PT[(i + 1) % 3]);

        ret = max(ret, dist(sr, a) + dist(sr, PT[i]) + dist(sr, PT[(i + 1) % 3]));

    }

    return ret;

}



k los() {

    return 1ll * rand() * rand() % 1'000'000 * 1.0 / 1'000'000;

}



void solve() {

    srand(123);

    FOR(i, 0, 3) cin >> PT[i].X >> PT[i].Y;

    kk akt = fermat_point(PT[0], PT[1], PT[2]);

    k best = ans(akt);

    for(k eps = 1e5; eps >= 1e-9; eps /= 1.01) {

        k ang = los() * 2 * PI;

        kk nowy = {akt.X + eps * cos(ang), akt.Y + eps * sin(ang)};

        k wyn = ans(nowy);

        if(wyn < best) best = wyn, akt = nowy;

    }

    cout << fixed << setprecision(10);

    // cout << akt.X << ' ' << akt.Y << '\n';

    cout << best << '\n';   

}



int main() {

    ios::sync_with_stdio(0);

    cin.tie(0);

    int tt = 1; //cin >> tt;

    FOR(te, 0, tt) solve();

    return 0;

}