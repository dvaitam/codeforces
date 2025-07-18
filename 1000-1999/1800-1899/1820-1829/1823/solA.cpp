#include <bits/stdc++.h>
#include <random>
using namespace std;
using ll = long long;
using ld = long double;
#define el '\n'
#define pb(x) push_back(x)
#define mp(x,y) make_pair(x,y)
#define vi vector<int>
#define vl vector<ll>
#define popcount __builtin_popcountll
#define ios ios_base::sync_with_stdio(0);cin.tie(0);cout.tie(0)
const int mod = 1e9 + 7;
const ll INF = 1e18;
mt19937 rng(chrono::steady_clock::now().time_since_epoch().count());
ll rand(ll l , ll r) {
    return uniform_int_distribution<ll>(l, r)(rng);
}


int main() {
    ios;
    int t;cin >> t;
    while(t--){
        int n,k;cin >> n >> k;
        int a = 0,b = n;
        while(a <= n){
            if(a * (a - 1) / 2 + b * (b - 1) / 2 == k)
                break;
            a++;b--;
        }
        if(a * (a - 1) / 2 + b * (b - 1) / 2 == k){
            cout << "YES" << el;
            while(a--)cout << 1 << " ";
            while(b--)cout << -1 << " ";
            cout << el;
        }else{
            cout << "NO" << el;
        }
    }
}