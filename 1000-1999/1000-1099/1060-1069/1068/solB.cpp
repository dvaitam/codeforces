#include<bits/stdc++.h>
#include <ext/pb_ds/assoc_container.hpp>
#include <ext/pb_ds/tree_policy.hpp>

typedef long long int ll;
typedef long double lf;

#define pll pair<ll,ll>
#define vll vector<ll>
#define vpl vector<pair<ll,ll>>
#define mp  make_pair
#define ppll pair<ll,pair<ll,ll>>
#define pb  push_back

#define ordered_set tree< ll , null_type , less<ll> , rb_tree_tag , tree_order_statistics_node_update > 


using namespace __gnu_pbds;
using namespace std;

ll mod = 1e9+7;

ll power(ll x, ll y) 
{ 
    ll res = 1;      // Initialize result 
  
    x = x % mod;  // Update x if it is more than or  
                // equal to p 
  
    while (y > 0) 
    { 
        // If y is odd, multiply x with result 
        if (y & 1) 
            res = (res*x) % mod; 
  
        // y must be even now 
        y = y>>1; // y = y/2 
        x = (x*x) % mod;   
    } 
    return res; 
} 

// function to count the divisors 
ll countDivisors(ll n) 
{ 
    ll cnt = 0; 
    for (ll i = 1; i <= sqrt(n); i++) { 
        if (n % i == 0) { 
            // If divisors are equal, 
            // count only one 
            if (n / i == i) 
                cnt++; 
  
            else // Otherwise count both 
                cnt = cnt + 2; 
        } 
    } 
    return cnt; 
} 
  

int main(){
    
    ll b;
    cin>>b;
    
    cout<<countDivisors(b);
    
}