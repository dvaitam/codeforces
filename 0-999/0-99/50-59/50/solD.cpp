#include<bits/stdc++.h>

#define ll long long

#define ld long double

#define fi first

#define se second

ll mpow(ll a, ll n,ll mod)

{ll ret=1;ll b=a;while(n) {if(n&1)

    ret=(ret*b)%mod;b=(b*b)%mod;n>>=1;}

return (ll)ret;

}

using namespace std;

#define sd(x) scanf("%d",&x)

#define pd(x) printf("%d",x)

#define sl(x) scanf("%lld",&x)

#define pl(x) printf("%lld",x)

#define mem(x,a) memset(x,a,sizeof(x))

#define pii pair<int,int>

#define mp make_pair

#define pb push_back

#define all(v) v.begin(),v.end()

#define N (int)(1e2+25)

pii pts[N];

int n,k,e,ix,iy;

double dp[N][N];

double prob[N];

double dis(int x1,int y1,int x2,int y2){

    return sqrt((x2-x1)*(x2-x1)+(y2-y1)*(y2-y1));

}

bool chk(double x){

    for(int i=0;i<=n;i++){

        for(int j=0;j<=n;j++){

            dp[i][j]=0.0;

        }

    }

    for(int i=1;i<=n;i++){

        double d=dis(ix,iy,pts[i].fi,pts[i].se);

        if(d<x)

            prob[i]=1;

        else if(d>x*1000.0)prob[i]=0;

        else{

            prob[i]=exp(1-(d*d)/(x*x));

        }

    }

    dp[0][0]=1.0;

    for(int i=1;i<=n;i++){

        for(int j=0;j<=i;j++){

            dp[i][j+1]+=dp[i-1][j]*prob[i];

            dp[i][j]+=dp[i-1][j]*(1-prob[i]);

        }

    }

    double sum=0.0;

    for(int i=0;i<k;i++){

        sum+=dp[n][i];

    }

    return sum<=(e/1000.0);

}

void solve(){

    sd(n);sd(k);sd(e);

    sd(ix);sd(iy);

    double lo=0,hi=10000.0;

    for(int i=1;i<=n;i++){

        sd(pts[i].fi);sd(pts[i].se);

    }

    double eps=0.0000001;

    while(hi-lo>eps){

        double mid=(lo+hi)/2.0;

        if(chk(mid)){

            hi=mid;

        }

        else{

            lo=mid;

        }

    }

    cout<<fixed<<setprecision(10)<<lo;

}

int main(){

   //freopen("C-large-practice.IN","r",stdin);

   //freopen("out.txt","w",stdout);

    int t=1;

 //  sd(t);

   for(int i=1;i<=t;i++){

       //printf("Case #%d:\n",i);

       solve();

   }

   return 0;

}