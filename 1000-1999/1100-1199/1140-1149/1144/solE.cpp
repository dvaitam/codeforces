#include <bits/stdc++.h>

using namespace std;

typedef long long ll;
#define sz(a) (int)((a).size())
#define all(x) (x).begin(), (x).end()
#define X first
#define Y second
#define mp make_pair
// defines end

const int MOD = 1e9 + 7;

// Functions starts from here

inline int add(int a, int b){
    a += b;
    if(a >= MOD)a -= MOD;
    return a;
}

inline int sub(int a, int b){
    a -= b;
    if(a < 0)a += MOD;
    return a;
}

inline int mul(int a, int b){
    return (int)((long long) a * b %MOD);
}

inline int binpow(int a, int b){
    int res = 1;
    while(b > 0){
        if(b & 1)res = mul(res, a);
        a = mul(a, a);
        b /= 2;
    }
    return res;
}

inline int inv(int a){
    return binpow(a, MOD - 2);
}

// Template ends here 
// Code starts from here

const int N = 2e5 + 5;
string s, t;
int n, a[N], b[N], ans[N];

int main(){
    ios_base::sync_with_stdio(0);
    cin.tie(0); cout.tie(0);

    cin>>n;
    cin>>s; cin>>t;
    for(int i=0;i<n;++i){
        a[i] = s[i] - 'a';
        b[i] = t[i] - 'a';
    }

    for(int i=n-1;i>=0;--i){
        int val = b[i] + a[i];
        if(val%2 == 0)ans[i] = val/2;
        else{
            ans[i] = (val - 1)/2;
            ans[i + 1] += 13;
        }
    }

    int carry = 0;
    for(int i=n-1;i>=0;--i){
        // cout<<"a[i] b[i] "<<a[i]<<" "<<b[i]<<"\n";
        ans[i] += carry;
        carry = ans[i]/26;
        ans[i] %= 26;
        // cout<<ans[i]<<"\n";
    }

    string res(n, 'a'); 
    for(int i=n-1;i>=0;--i){
        res[i] = 'a' + ans[i];
    }
    cout<<res<<"\n";
    return 0;
}