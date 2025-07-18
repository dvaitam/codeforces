#include <bits/stdc++.h>
#define _ ios_base::sync_with_stdio(0);cin.tie(0);
#define clr(a, b) memset(a, b, sizeof(a))
#define Error(x) cout << #x << " = " << x << endl
#define iDEBUG freopen("in.txt", "r", stdin)
#define oDEBUG freopen("output.txt","w",stdout)
#define PB push_back
#define MP make_pair
#define X first
#define Y second
#define SZ(x) ((int)(x).size())
#define ALL(x) x.begin(),x.end()
#define hi cout<<"hi\n"
//#define pow(x) (x)*(x)
#define pi 3.14159265359
#define INF 0x3f3f3f3f
typedef long long ll;

using namespace std;
//inline int gcd(int x, int y){ while(x^=y^=x^=y%=x); return y; }
//inline int lcm(int x, int y) { if(!x || !y) return 0; return x*y / gcd(x, y); }

int main()
{_
    ll n, k;
    scanf("%lld %lld", &n, &k);
    printf("%lld", (ll)ceil(k/(long double)n));


    return 0;
}