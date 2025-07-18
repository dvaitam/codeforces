#include <algorithm>

#include <iostream>

#include <cstring>

#include <vector>

#include <cstdio>

#include <string>

#include <cmath>

#include <queue>

#include <set>

#include <map>

#include<complex>

using namespace std;

typedef long long ll;

typedef double db;

typedef pair<int,int> pii;

typedef vector<int> vi;

#define de(x) cout << #x << "=" << x << endl

#define rep(i,a,b) for(int i=a;i<(b);++i)

#define all(x) (x).begin(),(x).end()

#define sz(x) (int)(x).size()

#define mp make_pair

#define pb push_back

#define fi first

#define se second

#define setIO(x) freopen(x".in","r",stdin);freopen(x".out","w",stdout)

typedef complex<db> P;

#define X real()

#define Y imag()

const int N = 1e5 + 1;const db eps = 1e-8;

int sgn(db x){return (x>eps)-(x<-eps);}

bool Cmp(const P&a,const P&b){

    if(sgn(a.X-b.X) != 0) return a.X < b.X;

    return a.Y < b.Y;

}

db cross(const P&a,const P&b){return (conj(a)*b).Y;}

db dot(const P&a,const P&b){return (conj(a)*b).X;}

int n;db R;

P p[N] , q[N];

int l[N],r[N],del[N],_,version[N];

priority_queue<pair<db,pii> > Q;

void Del(int x){

    del[x] = true;r[l[x]] = r[x] , l[r[x]] = l[x];

}

void Update(int x){

    if(dot(q[l[x]]-q[x],q[r[x]]-q[x])>0){

        version[x] = -1;return;

    }

    P&a=q[x],&b=q[l[x]],&c=q[r[x]];

    db r=abs(a-b)*abs(a-c)*abs(b-c)/(2*fabs(cross(a-b,c-b)));

    Q.push(mp(r,mp(x,version[x]=++_)));

}



int main(){

    scanf("%d%lf",&n,&R);db x,y;

    rep(i,0,n) scanf("%lf%lf",&x,&y) , p[i] = P(x,y);

    if(n == 1) return puts("0.000000") , 0;

    sort(p,p+n,Cmp);

    q[0] = p[0];int tp=1;rep(i,1,n){

        while(tp>=2&&cross(q[tp-1]-q[tp-2],p[i]-q[tp-2])<=0)--tp;

        q[tp++] = p[i];

    }

    int tt=tp;for(int i=n-1;i>=0;--i){

        while(tp>tt&&cross(q[tp-1]-q[tp-2],p[i]-q[tp-2])<=0)--tp;

        q[tp++] = p[i];

    }n = tp - 1;

    rep(i,0,n) l[i] = (i-1+n)%n , r[i] = (i+1)%n;

    rep(i,0,n) if(dot(q[l[i]]-q[i],q[r[i]]-q[i]) < 0) Update(i);

    while(sz(Q) && Q.top().fi > R){

        pii tp = Q.top().se;Q.pop();

        int c = tp.fi , ver = tp.se;

        if(ver != version[c]) continue;

        Del(c);Update(l[c]);Update(r[c]);

    }

    int s=0;while(del[s])++s;

    db ans = 0.;

    int i=s;

    do{

        ans += cross(q[i],q[r[i]]) / 2;

        db angle = 2 * asin(abs(q[i]-q[r[i]])/(2*R));

        ans += R*R*(angle-sin(angle)) / 2;

        i = r[i];

    }while(i!=s);

    printf("%.10f\n",ans);

    return 0;

}