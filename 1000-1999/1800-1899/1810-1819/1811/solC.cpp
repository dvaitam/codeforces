// author : tuanotuan //
// MSSV 23120389 //
// Nguyen Ngoc Tuan HCMUS aim ICPC quoc gia 2024 //
// Ho Chi Minh City University of Science //
#include <bits/stdc++.h>
using namespace std;
#define REP(i, a, b) for(int i=a; i<=b; i++)
#define REPD(i, b, a) for(int i=b; i>=a; i--)
#define RESET(a, b) memset(a, b, sizeof(a))
#define DEBUG(x) { cout << #x << " = "; cout << (x) << endl; }
#define db double
#define ll long long
#define ld long double
#define pii pair<int,int>
#define pll pair<long long, long long>
#define pss pair<string, string>
#define vi vector<int>
#define vll vector<long long>
#define vpii vector< pair<int,int> >
#define vpll vector< pair<long long, long long> >
#define vvi vector< vector<int> >
#define vb vector<bool>
#define all(a) (a).begin(),(a).end()
#define sz(x) (int)(x).size()
#define fi first
#define se second
#define lb lower_bound
#define ub upper_bound
#define pb push_back
#define endl "\n"
#define YES cout<<"YES";
#define NO cout<<"NO";
#define NAME "template"
long long Pow(long long a, long long n){
    if(n==0) return 1;
    ll res=1;
    while(n){
        if(n&1)
            res*=a;
        a*=a;
        n/=2;
    }
    return res;
}
long long gcd(long long x, long long y){
    while(y){
        long long r=x%y;
        x=y;
        y=r;
    }
    return x;
}
long long lcm(long long x, long long y){
    return (x/gcd(x,y))*y;
}
long long Min(long long x, long long y){
    if(x<y) return x;
    return y;
}
long long Max(long long x, long long y){
    if(x>y) return x;
    return y;
}
int dc4[]={-1,0,1,0};
int dr4[]={0,-1,0,1};
int dc8[]={-1,-1,0,1,1,1,0,-1};
int dr8[]={0,-1,-1,-1,0,1,1,1};
const int N=100005;
const long long MOD=1000000007;
const int INF=1000000005;
const long long INFF=1000000000000000005LL;
const long double PI=acos(-1.0);
int n,m,k;
//int a[N];
void solve(){
    cin>>n;
    vi b(n-1);
    REP(i, 0, n-2)
    cin>>b[i];
    vi a(n,0);
    a[0]=b[0];
   REP(i, 1, n-2)
   if(b[i]>b[i-1]){
    a[i+1]=b[i];
    i++;
    if(i<=n-2)
        a[i+1]=b[i];
   }
   else a[i]=b[i];
    for(auto e: a)
        cout<<e<<" ";
}
main(){
    // freopen(NAME".inp", "r", stdin);
    // freopen(NAME".out", "w", stdout);
    ios::sync_with_stdio(0);
    cin.tie(0); cout.tie(0);
    int tests=1;
    cin>>tests;
    while(tests--){
        solve();
        cout<<endl;
    }
}