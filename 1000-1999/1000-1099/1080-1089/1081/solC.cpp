#include<bits/stdc++.h>
using namespace std;
typedef long long int ll;
typedef long double ld;
const ll mod = 998244353;

ll power(ll x, ll y, ll p)
{
    ll res = 1;      // Initialize result

    x = x % p;  // Update x if it is more than or
                // equal to p

    while (y > 0)
    {
        // If y is odd, multiply x with result
        if (y & 1)
            res = (res*x) % p;

        // y must be even now
        y = y>>1; // y = y/2
        x = (x*x) % p;
    }
    return res;
}

// Returns n^(-1) mod p
ll modInverse(ll n, ll p)
{
    return power(n, p-2, p);
}

ll ncr(ll n, ll r)
{
   // Base case
   ll p = mod;
   if (r==0)
      return 1;

    // Fill factorial array so that we
    // can find all factorial of r, n
    // and n-r
    ll fac[n+1];
    fac[0] = 1;
    for (ll i=1 ; i<=n; i++)
        fac[i] = fac[i-1]*i%p;

    return (fac[n]* modInverse(fac[r], p) % p *
            modInverse(fac[n-r], p) % p) % p;
}

int main()
{
    ios_base::sync_with_stdio(false);
    cin.tie(0);cout.tie(0);
    ll n,m,k;
    cin>>n>>m>>k;
    if(m==1 && k){
        cout<<0<<endl;
        return 0;
    }
    if(n<=k){
        cout<<0<<endl;
        return 0;
    }
    if(k==0){
        cout<<m<<endl;
        return 0;
    }
    ll ans = ncr(n-1,k);
    ll c = 0;
    ans = (ans*m)%mod;
    while(k){
        //cout<<k<<" "<<ans<<endl;
        ans = (ans*(m-1))%mod;
        k--;
    }
    cout<<ans<<endl;
}