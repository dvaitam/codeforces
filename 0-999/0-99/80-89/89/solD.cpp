#include <cstdio>
#include <cstring>
#include <cmath>
#include <algorithm>
#include <string>
#include <vector>
#include <map>
#include <set>
#include <cctype>
#include <numeric>
#include <queue>
#define FOR(i,s,e) for(int i=(s);i<(int)(e);i++)
#define FOE(i,s,e) for(int i=(s);i<=(int)(e);i++)
#define CLR(s) memset(s,0,sizeof(s))
#define PB push_back
using namespace std;
typedef long long LL;
typedef pair<int,int> pii;
typedef vector<int> vi;
#define x first
#define y second

const double inf = 1e20, eps = 1e-9, PI = acos(-1);
#define flt(x,y) ((x)<(y)-eps)
#define fgt(x,y) ((x)>(y)+eps)
#define fle(x,y) ((x)<(y)+eps)
#define fge(x,y) ((x)>(y)-eps)
#define feq(x,y) (fabs((x)-(y))<eps)

double sq(double x){ return x*x; }

struct P {
        double a[3];
        P(){ }
        P(double x,double y,double z){ a[0]=x, a[1]=y, a[2]=z; }
        bool operator==(const P &p)const {
                FOR(i,0,3) if(!feq(a[i],p.a[i])) return 0;
                return 1;
        }
        bool operator<(const P &p)const {
                FOR(i,0,3) if(!feq(a[i],p.a[i])) return flt(a[i],p.a[i]);
                return 0;
        }
        P operator-(P p) {
                P res;
                FOR(i,0,3) res.a[i]=a[i]-p.a[i];
                return res;
        }
        P operator+(P p) {
                P res;
                FOR(i,0,3) res.a[i]=a[i]+p.a[i];
                return res;
        }
        P operator*(double k) {
                P res;
                FOR(i,0,3) res.a[i]=a[i]*k;
                return res;
        }
        P operator^(P p) {
                P res;
                FOR(i,0,3) {
                        int j = (i+1)%3;
                        int k = (i+2)%3;
                        res.a[i] = (a[j]*p.a[k] - a[k]*p.a[j]);
                }
                return res;
        }
        double operator*(P p) {
                double res = 0;
                FOR(i,0,3) res+=a[i]*p.a[i];
                return res;
        }
        void out() { FOR(i,0,3) printf(" %f ", a[i]); puts(""); }
        double mag2(){ return sq(a[0])+sq(a[1])+sq(a[2]); }
        double mag(){ return sqrt(mag2()); }
        P nor(){ double l = mag(); if (feq(l,0)) return *this; return *this*(1/l); }
        P rot(P z, double t){   // v = vcos(t) + (z^v)sin(t) + (z*v)z(1-cos(t))
                return *this*cos(t) + (z^*this)*sin(t) + z*(z*(*this))*(1-cos(t));
        }
        P rot(P z) {
                return z^(*this) + z*(z*(*this));
        }
        double dist(P p){ return (*this-p).mag(); }
        void eat(){ FOR(i,0,3) scanf("%lf",&a[i]); }
};

double ans=1e20;

// q = a + t v
bool f(P a, P o, double r1, double r2, P v, double &res) {
        double A = v.mag2();
        double B = 2*((a-o)*v);
        double C = (a-o).mag2() - (r1+r2)*(r1+r2);
        double D = B*B-4*A*C;
        if (flt(D,0)) return 0;
        if (fgt(D,0)) D=sqrt(D);
        res = (-B-D)/(2*A);
        return res > 0;
}

void upd(double x) { ans=min(ans,x); }

int main() {
        P A,V;
        double R;
        A.eat();
        V.eat();
        scanf("%lf",&R);
        int n;
        scanf("%d", &n);
        FOR(i,0,n) {
                P cen, p;
                double r;
                int m;
                cen.eat();
                scanf("%lf%d",&r,&m);
                double res;
                if (f(A,cen,R,r,V,res)) {
                        upd(res);
                }
                while (m--) {
                        p.eat();
                        if(f(A,cen,R,0,V,res))
                                upd(res);
                        if(f(A,cen+p,R,0,V,res))
                                upd(res);
                        P c=cen;
                        P d=cen+p;
                        P cd=d-c;
                        P ca=c-A;
                        double A1 = (V^cd).mag2();
                        if (feq(A1,0)) continue;

                        double B = (V^cd).mag() * (ca^cd).mag();
                        double C = (ca^cd).mag2();
                        double D = B*B-4*A1*C;
                        if (fge(D,0)) {
                                if (fgt(D,0)) D=sqrt(D);
                                double t1 = (-B+D)/(2*A1);
                                if (t1>0) {
                                        P Q=A+V*t1;
                                        double dot1=(Q-c)*(d-c);
                                        double dot2=(Q-d)*(c-d);
                                        if(fge(dot1,0) && fge(dot2,0))
                                                upd(res);
                                }
                        }
                }
        }
        if (feq(ans,1e20)) puts("-1");
        else
                printf("%.20f\n", ans);
        return 0;
}