#include <bits/stdc++.h>

#define IOS ios::sync_with_stdio(0);cin.tie(0);cout.tie(0);

using namespace std;

 

using _T= long double; //long long

typedef double dd;

constexpr _T eps=1e-8;

constexpr long double PI=3.1415926535897932384l;

 

template<typename T> struct point

{

    T x,y;

 

    bool operator==(const point &a) const {return (abs(x-a.x)<=eps && abs(y-a.y)<=eps);}

    bool operator<(const point &a) const {if (abs(x-a.x)<=eps) return y<a.y-eps; return x<a.x-eps;}

    bool operator>(const point &a) const {return !(*this<a || *this==a);}

    point operator+(const point &a) const {return {x+a.x,y+a.y};}

    point operator-(const point &a) const {return {x-a.x,y-a.y};}

    point operator-() const {return {-x,-y};}

    point operator*(const T k) const {return {k*x,k*y};}

    point operator/(const T k) const {return {x/k,y/k};}

    T operator*(const point &a) const {return x*a.x+y*a.y;}

    T operator^(const point &a) const {return x*a.y-y*a.x;}

    int toleft(const point &a) const {const auto t=(*this)^a; return (t>eps)-(t<-eps);}

    T len2() const {return (*this)*(*this);}

    T dis2(const point &a) const {return (a-(*this)).len2();}

    double len() const {return sqrt(len2());}

    double dis(const point &a) const {return sqrt(dis2(a));}

    double ang(const point &a) const {return acos(max(-1.0,min(1.0,((*this)*a)/(len()*a.len()))));}

    //point rot(const double rad) const {return {x*cos(rad)-y*sin(rad),x*sin(rad)+y*cos(rad)};}

    //point rot(const double cosr,const double sinr) const {return {x*cosr-y*sinr,x*sinr+y*cosr};}

};

 

using Point=point<_T>;

 

struct argcmp

{

    bool operator()(const Point &a,const Point &b) const

    {

        const auto quad=[](const Point &a)

        {

            if (a.y<-eps) return 1;

            if (a.y>eps) return 4;

            if (a.x<-eps) return 5;

            if (a.x>eps) return 3;

            return 2;

        };

        const int qa=quad(a),qb=quad(b);

        if (qa!=qb) return qa<qb;

        const auto t=a^b;

        //if (abs(t)<=eps) return a*a<b*b-eps;

        return t>eps;

    }

};

 

template<typename T> struct line

{

    point<T> p,v;

 

    bool operator==(const line &a) const {return v.toleft(a.v)==0 && v.toleft(p-a.p)==0;}

    int toleft(const point<T> &a) const {return v.toleft(a-p);}

    point<T> inter(const line &a) const {return p+v*((a.v^(p-a.p))/(v^a.v));}//直线与直线的交点

    double dis(const point<T> &a) const {return abs(v^(a-p))/v.len();}

    point<T> proj(const point<T> &a) const {return p+v*((v*(a-p))/(v*v));}//点在直线上的投影点

    bool operator<(const line &a) const

    {

        if (abs(v^a.v)<=eps && v*a.v>=-eps) return toleft(a.p)==-1;

        return argcmp()(v,a.v);

    }

};

 

using Line=line<_T>;

 

template<typename T> struct segment

{

    point<T> a,b;

    bool operator< (const segment &s) const {return make_pair(a,b)<make_pair(s.a,s.b);}

    int is_on(const point<T> &p) const//点在线段的端点返回-1, 在线段上返回1，否则返回0

    {

        if (p==a || p==b) return -1;

        return (p-a).toleft(p-b)==0 && (p-a)*(p-b)<-eps;

    }

 

    int is_inter(const line<T> &l) const

    {

        if (l.toleft(a)==0 || l.toleft(b)==0) return -1;

        return l.toleft(a)!=l.toleft(b);

    }

     

    int is_inter(const segment<T> &s) const

    {

        if (is_on(s.a) || is_on(s.b) || s.is_on(a) || s.is_on(b)) return -1;

        const line<T> l{a,b-a},ls{s.a,s.b-s.a};

        return l.toleft(s.a)*l.toleft(s.b)==-1 && ls.toleft(a)*ls.toleft(b)==-1;

    }

 

    double dis(const point<T> &p) const

    {

        if ((p-a)*(b-a)<-eps || (p-b)*(a-b)<-eps) return min(p.dis(a),p.dis(b));

        const line<T> l{a,b-a};

        return l.dis(p);

    }

 

    double dis(const segment<T> &s) const

    {

        if (is_inter(s)) return 0;

        return min({dis(s.a),dis(s.b),s.dis(a),s.dis(b)});

    }

};

 

using Segment=segment<_T>;



bool segs_inter(const vector<Segment> &segs)

{

    if (segs.empty()) return false;

    using seq_t=tuple<_T,int,Segment>;

    const auto seqcmp=[](const seq_t &u, const seq_t &v)

    {

        const auto [u0,u1,u2]=u;

        const auto [v0,v1,v2]=v;

        if (abs(u0-v0)<=eps) return make_pair(u1, u2) < make_pair(v1, v2);

        return u0<v0-eps;

    };

    vector<seq_t> seq;

    for (auto seg:segs)

    {

        if (seg.a.x>seg.b.x+eps) swap(seg.a,seg.b);

        seq.push_back({seg.a.x,0,seg});

        seq.push_back({seg.b.x,1,seg});

    }

    sort(seq.begin(),seq.end(),seqcmp);

    const auto cal=[](const Segment &u, const _T x)

    {

        if (abs(u.a.x-u.b.x)<=eps) return u.a.y;

        return ((x-u.a.x)*(u.b.y-u.a.y)+u.a.y*(u.b.x-u.a.x))/(u.b.x-u.a.x);

    };

    _T x_now;

    auto cmp=[&](const Segment &u, const Segment &v){return cal(u,x_now)<cal(v,x_now)-eps;};

    multiset<Segment,decltype(cmp)> s{cmp};

    for (const auto [x,o,seg]:seq)

    {

        x_now=x;

        const auto it=s.lower_bound(seg);

        if (o==0)

        {

            if (it!=s.end() && seg.is_inter(*it)) return true;

            if (it!=s.begin() && seg.is_inter(*prev(it))) return true;

            s.insert(seg);

        }

        else s.erase(s.find(seg));

    }

    return false;

}

const int N = 250010;

struct node {

    dd x, y, dx, dy, s;

}nd[N];

int n; 

bool check(dd mid) {

    vector<Segment>seg; 

    for (int i = 0; i < n; i++) {

        Point p = {nd[i].x, nd[i].y};

        Point vec = {nd[i].dx, nd[i].dy};

        Point pp = p + vec / vec.len() * mid * nd[i].s;

        seg.push_back({p, pp});

    }

    return segs_inter(seg);

}

int main() {

    IOS;

    cin >> n;

    for (int i = 0; i < n; i++) {

        cin >> nd[i].x >> nd[i].y >> nd[i].dx >> nd[i].dy >> nd[i].s;

    }

    dd l = 0, r = 1e10, ans = 0;

    while(r - l > eps && (r - l) / r > eps) {

        dd mid = (l + r) / 2;

        if(check(mid)) {

            r = mid;

            ans = mid;

        }

        else {

            l = mid;

        }

    }

    cout.setf(ios::fixed);

    if(ans < eps) {

        cout << "No show :(\n";

    }

    else {

        cout << setprecision(8) << ans << "\n";

    }

}