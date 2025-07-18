#include <cstdio>
#include <cstring>
#include <cstdlib>
#include <cmath>
#include <cassert>
#include <iostream>
#include <sstream>
#include <algorithm>
#include <vector>
#include <set>
#include <map>
#include <queue>

using namespace std;

#define sz(v) ((int) (v).size())
#define all(v) (v).begin(), (v).end()
#define mp make_pair
#define pb push_back
#define forn(i,n) for (int i=0; i<n; i++) 

typedef long long ll;
typedef long long int64;
typedef pair<int,int> ii;
typedef vector<int> vi;
typedef vector<string> vs;

template<typename T> T abs(T x) { return x>0 ? x : -x; }
template<typename T> T sqr(T x) { return x*x;          }

typedef double D;
const D eps=1e-8;

struct P {
    D x,y;
    P() {}
    P(D x, D y): x(x), y(y) {}
    P operator -(P a) {
        return P(x-a.x,y-a.y);
    }
    P operator +(P a) {
        return P(x+a.x,y+a.y);
    }
    P operator *(D k) {
        return P(x*k,y*k);
    }
    D len() {
        return sqrt(sqr(x)+sqr(y));
    }
    void save() {
        printf("%.5lf %.5lf\n",x,y);
    }
};

typedef vector<P> vp;

struct L {
    D a,b,c;
    L() {}
    L(D a, D b, D c): a(a), b(b), c(c) {}
};

struct C {
    D x,y,r;
    C() {}
    C(D x, D y, D r): x(x), y(y), r(r) {}
    P o() {
        return P(x,y);
    }
    void load() {
        cin>>x>>y>>r;
    }
};

struct CL {
    int type; // 0 - C, 1 - L
    C c;
    L l;
    CL() {
        type=-1;
    }
};

vp cross(L l1, L l2) {
    D det=l1.a*l2.b-l1.b*l2.a;
    vp res;
    if (abs(det)<eps) return res;
    D det1=-(l1.c*l2.b-l1.b*l2.c);
    D det2=-(l1.a*l2.c-l1.c*l2.a);
    res.pb(P(det1/det,det2/det));
    return res;
}

vp cross(C c, L l) {
    vp res;
    D al,be;
    al=l.b;
    be=-l.a;
    D x0,y0;
    if (abs(l.a)<abs(l.b)) {
        x0=0;
        y0=-l.c/l.b;
    } else {
        y0=0;
        x0=-l.c/l.a;
    }
    D A,B,C,t;
    A=sqr(al)+sqr(be);
    B=2*al*(x0-c.x)+2*be*(y0-c.y);
    C=sqr(x0-c.x)+sqr(y0-c.y)-sqr(c.r);
    D d=B*B-4*A*C;
    if (d<-eps) return res;
    if (d<0) d=0;
    t=(-B+sqrt(d))/(2*A);
    res.push_back(P(x0+al*t,y0+be*t));
    t=(-B-sqrt(d))/(2*A);
    res.push_back(P(x0+al*t,y0+be*t));
    return res;
}

vp cross(L l, C c) {
    return cross(c,l);
}

vp cross(C c1, C c2) {
    D a,b,c;
    a=2*(c2.x-c1.x);
    b=2*(c2.y-c1.y);
    c=sqr(c2.r)-sqr(c1.r)+sqr(c1.x)-sqr(c2.x)+sqr(c1.y)-sqr(c2.y);
    return cross(c1,L(a,b,c));
}

vp cross(CL cl1, CL cl2) {
    if (cl1.type==0 && cl2.type==0) return cross(cl1.c,cl2.c);
    if (cl1.type==0 && cl2.type==1) return cross(cl1.c,cl2.l);
    if (cl1.type==1 && cl2.type==0) return cross(cl1.l,cl2.c);
    if (cl1.type==1 && cl2.type==1) return cross(cl1.l,cl2.l);
    assert(false);
}

C c[3];

CL getL(C c1, C c2) {
    D a,b,c;
    a=2*c2.x-2*c1.x;
    b=2*c2.y-2*c1.y;
    c=sqr(c1.x)-sqr(c2.x)+sqr(c1.y)-sqr(c2.y);
    CL res;
    res.l=L(a,b,c);
    res.type=1;
    return res;
}

CL getC(C c1, C c2) {
    if (c1.r>c2.r) return getC(c2,c1);
    D c=c1.r/c2.r;
    P o1=c1.o();
    P o2=c2.o();
    P v=o2-o1;
    P p1,p2;
    p1=o1+v*(c/(1+c));
    p2=o1+v*(c/(c-1));
    P o=(p1+p2)*0.5;
    D r=(p1-o).len();
    CL res;
    res.c=C(o.x,o.y,r);
    res.type=0;
    return res;
}

CL getCL(C c1, C c2) {
    if (abs(c1.r-c2.r)<eps) return getL(c1,c2);
    return getC(c1,c2);
}

int main()
{
    forn(i,3) c[i].load();
    CL cl1=getCL(c[0],c[1]);
    CL cl2=getCL(c[1],c[2]);
    vp cr=cross(cl1,cl2);
    double mi=1e100;
    P ans;
    forn(i,sz(cr)) {
        vector<double> q;
        forn(j,3) q.pb((cr[i]-c[j].o()).len()/c[j].r);
        bool ok=true;
        forn(j,3) if (abs(q[j]-q[0])>eps) ok=false;
        if (q[0]<1-eps) ok=false;
        if (!ok) continue;
        if (q[0]<mi) {
            mi=q[0];
            ans=cr[i];
        }
    }
    if (mi<1e50) ans.save();

    return 0;
}