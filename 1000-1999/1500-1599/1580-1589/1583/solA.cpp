#include <iostream>
#include <algorithm>
#include <vector>
#include <set>
#include <map>
#include <cmath>
#define ll long long

using namespace std;

ll normal_pow(ll a, ll n){
    for(int i = 0; i < n - 1; i++)
        a *= 2;
    return a;
}

bool prime(ll n){
    for(ll i = 2; i <= sqrt(n); i++)
        if(n % i == 0)
            return false;
    return true;
}

void solve(){
    ll n, res = 0, k = 0;
    cin >> n;
    vector<ll> mas(n);
    for(int i = 0; i < n; i++){
        cin >> mas[i];
        res += mas[i];
        if (mas[i] % 2)
            k = i;
    }
    if (prime(res)) {
        cout << n -1 << "\n";
        for(int i =0; i< n; i++) {
            if (i != k)
                cout << i + 1 << " ";
        }
    }
    else{
        cout << n << "\n";
        for(int i =0; i< n; i++) {
            cout << i + 1 << " ";
        }
    }
    cout << endl;
}

int main() {
    ios_base::sync_with_stdio(false);
    cin.tie(0);
    cout.tie(0);
    ll t;
    cin >> t;
    while(t--){
        solve();
    }
}