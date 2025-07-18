// JAI SHREE RAM, HAR HAR MAHADEV, HARE KRISHNA



#include <bits/stdc++.h>



#include <ext/pb_ds/assoc_container.hpp>

#include <ext/pb_ds/tree_policy.hpp>



using namespace std;

using namespace __gnu_pbds;



typedef tree<int, null_type, less<int>, rb_tree_tag, tree_order_statistics_node_update> pbds;



#define fbo find_by_order // kth smallest element (returns iterator)

#define oof order_of_key // no. of element strictly less than element x (returns number)

using ll = long long int;

using ld = long double;

#define pll pair<ll,ll>

#define vt vector

#define um unordered_map

#define pb push_back

#define F first

#define S second

#define ub upper_bound

#define lb lower_bound

const int maxx = INT_MAX;

const int minn = INT_MIN;



#define all(a) a.begin(),a.end()

#define allr(a) a.rbegin(),a.rend()

#define rep(i,a,b) for(ll i=a;i<b;++i)

#define per(i,a,b) for(ll i=a;i>=b;--i)

#define in(a, n) rep(i,0,n) cin>>a[i];

const ll mod = 1e9 + 7;

//Living on the edge

#define getSum(a) accumulate(all(a),(ll)0);

ll mod_add(ll x, ll y) {return (x%mod + y%mod)%mod;}

ll mod_sub(ll x, ll y) {return (x%mod - y%mod)%mod;}

ll mod_mul(ll x, ll y) {return (x%mod * y%mod)%mod;}

ll mod_div(ll x, ll y) {return (x%mod / y%mod)%mod;}

ll power(ll x, ll y)

{

    if(y==0) return 1;

    ll temp = power(x,y/2);

    if(y%2==0)

        return temp*temp;

    else

        return x*temp*temp;

}

#define nline '\n'

void solve()

{

    ll n;

    cin >> n;



    vt<ll> a(n);

    in(a, n);

    

    if (n == 1) {

      for (int z = 0; z < 3; ++z) {

         cout << "1 1\n";

         cout << -a[0] << "\n";

         a[0] = 0;

      }

      return ;

   }



    cout << "1 1" << "\n"; cout << -a[0] << "\n";



    a[0] = 0;



    cout << "1 " << n << "\n";

    rep(i, 0, n) cout << -a[i] * n << " "; cout << nline;



    cout << "2 " << n << nline;

    rep(i, 1, n) cout << a[i] * (n-1) << " "; cout << nline;

}



int main()

{

    ios_base::sync_with_stdio(0), cin.tie(0), cout.tie(0);

    int T = 1;

    while(T--)

    {

        solve();

    }

    return 0;

}