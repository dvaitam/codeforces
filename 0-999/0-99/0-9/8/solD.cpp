//go:build ignore
#include <cstdio>
#include <cstring>
#include <algorithm>
#include <iostream>
#include <cmath>
#include <queue>
#define ll long long

using namespace std;
double t1, t2, ab, ac, bc;
struct P
{
    double x, y;
    P(double x=0,double y=0):x(x),y(y){};
}a, b, c, p, p1;
P operator +(P a,P b){ return P(a.x+b.x,a.y+b.y); }
P operator *(P a,double b){ return P(a.x*b,a.y*b); }
inline double dis(P a,P b){ return sqrt((a.x-b.x)*(a.x-b.x)+(a.y-b.y)*(a.y-b.y)); }
double cal(double k)
{
    p=b*(1-k)+c*k;
    double ap=dis(a, p);
    if (ap+(k+1)*bc<t1&&ap+(1-k)*bc<t2) return min(t1-(k+1)*bc, t2-(1-k)*bc);
    double l=0, r=1, mid;
    while (r-l>1e-15){
        mid=(l+r)/2;
        p1=a*(1-mid)+p*mid;
        if (ap*mid+dis(p1, b)+bc<t1&&ap*mid+dis(p1, c)<t2) l=mid;
        else r=mid;
    }
    return (r+l)/2*ap;
}
int main()
{
    double l, r, mid, midmid;
    scanf("%lf%lf", &t1, &t2);
    scanf("%lf%lf%lf%lf%lf%lf", &a.x, &a.y, &c.x, &c.y, &b.x, &b.y);
    ab=dis(a, b); bc= dis(b, c); ac=dis(a, c);
    t1+=ab+bc+1e-12; t2+=ac+1e-12;
    if (ab+bc<t2) {printf("%.10lf\n", min(t1, t2)); return 0;}
    l=0; r=1;
    while(r-l>1e-15){
        mid=(l+r+l)/3;
        midmid=(r+l+r)/3;
        if (cal(mid)-cal(midmid)<1e-12) l=mid;
        else r=midmid;
    }
    printf("%.10lf\n", cal((r+l)/2));
    return 0;
}