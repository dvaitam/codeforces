#include <iostream>
#include <bits/stdc++.h>
#include <string>
#include <bitset>
using namespace std;

#define pb push_back
#define ull unsigned long long
#define ll long long
#define F first
#define S second
#define PI acos(-1)
#define EPS 1e-9
#define ld double
#define MAX 1000000000
#define NIL 0
#define INF 1e15
#define infi 1000000000
#define rd(a) scanf("%d",&a)
#define rd2(a,b) scanf("%d %d",&a,&b)
#define rd3(a,b,c) scanf("%d %d %d",&a,&b,&c)
#define rdll(a) scanf("%I64d",&a)
#define sz(a) (int) a.size()
#define lp(i,a,n) for(int i=(a); i<=(n) ; ++i)
#define lpd(i,n,a) for(int i=(n); i>=(a) ; --i)
#define pi acos(-1)
#define lc (x << 1)
#define rc (x << 1 | 1)
#define cp(a,b)                 ( (conj(a)*(b)).imag() )	// a*b sin(T), if zero -> parllel
#define dp(a,b)                 ( (conj(a)*(b)).real() )	// a*b cos(T), if zero -> prep
#define angle(a)                (atan2((a).imag(), (a).real()))
#define X real()
#define Y imag()
#define vec(a,b)                ((b)-(a))
#define vvi vector<vector<int>>
#define f first
#define s second
#define clr(a,b) memset(a,b,sizeof a)

#define mod1 1000500001ll
#define mod2 1000300001ll
#define base1 137ll
#define base2 31ll

typedef complex<double>CX;
typedef pair<int,int>ii;
typedef pair<ii,ll>tri;
typedef pair<vector<int>,int>vii;
typedef pair<int, int> pii;
typedef pair<ii,ii> line;
typedef pair<double,double>point;
typedef pair<ll, ll> pll;
typedef vector<int> vi;
typedef pair<int,pair<int,int>>edge;


const int N=200005;
const int M=22;

ll mod=1e9+9;

ll gcd(ll a, ll b) { return b == 0 ? a : gcd(b, a % b); }
ll lcm(ll a, ll b) { return a * (b / gcd(a, b)); }
bool is_vowel(char c){if(c=='a'||c=='e'||c=='i'||c=='o'||c=='u')return 1;return 0;}
ll extended_euclidean(ll a,ll b,ll &x,ll &y){if(b==0){x=1;y=0;return a;}ll g = extended_euclidean(b,a%b,y,x);y -= (a/b)*x;return g;}
ll power(ll base,ll p,ll mod){if(p==1)return base;if(!p)return 1ll;ll ret=power(base,p/2,mod);ret*=ret;ret%=mod;if(p&1)ret*=base;return ret%mod;}


int cnt[1005],n;
int main(){
  //   freopen("test.in","r",stdin);
   rd(n);
   lp(i,1,n-1){
      int a,b;
      rd2(a,b);
      cnt[a]++;
      cnt[b]++;
   }
   if(cnt[n]!=n-1){
    cout<<"NO";
    return 0;
   }
   lp(i,1,n){
     if(cnt[i]>i){
        cout<<"NO";
        return 0;
     }
   }
   vector<int>all;
   lp(i,1,n){
     if(cnt[i]==0)all.pb(i);
   }
   sort(all.rbegin(),all.rend());
   vector<ii>sol;
   int prv=n,idx=0;
   lpd(i,n-1,1){
      if(cnt[i]==0)continue;
      while(cnt[i]>1){
        sol.pb(ii(prv,all[idx]));
        if(all[idx]>i){
            cout<<"NO";
            return 0;
        }
       // printf("%d %d\n",prv,all[idx++]);
        prv=all[idx];
        idx++;
        cnt[i]--;
      }
      sol.pb(ii(i,prv));
      //printf("%d %d\n",i,prv);
      prv=i;
   }
   cout<<"YES\n";
   for(auto it : sol)printf("%d %d\n",it.F,it.S);
   return 0;
}