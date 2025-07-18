#include <bits/stdc++.h>

using namespace std;

#define looklols    ios_base :: sync_with_stdio(0); cin.tie(0); cout.tie(0);

#define T           ll T;  cin>>T; while(T--)

#define ll          long long

#define lld         long double

#define ld          long double

#define F           first

#define S           second

#define pb          push_back

#define pf          push_front

#define all(x)      x.begin(),x.end()

#define allr(x)     x.rbegin(),x.rend()

#define ones(x) __builtin_popcountll(x)

#define sin(a) sin((a)*PI/180)

#define cos(a) cos((a)*PI/180)

#define endl        "\n"

const lld pi = 3.14159265358979323846;

const ll N=2e5+2;

const ll MOD = 998244353 , LOG = 25;

/*

  ℒ◎øкℓ☺łṧ

 */



bool solve(deque<int>& arr, int tar){

    int l=0, r=arr.size()-1;

    while(l < r){

        int curSum=arr[l]+arr[r];

        if(curSum==tar) return true;

        else if(curSum > tar) r--;

        else l++;

    }

    return false;

}









void fact(int n){

    for (int i = 2; i*i <= n; ++i) {

        int cnt=0;

        while (n%i ==0){

            n/=i;

            cnt++;

        }

        if(cnt){

            cout<<i<< " ^ "<<cnt<<endl;

        }

    }

    if(n>1)

        cout<<"n^1"<<endl;

}

vector<ll> check(ll n){

    vector<ll>b;

    for (ll i = 1; i*i<=n  ; ++i) {

        if(n%i==0){

            b.push_back(i);

            if(i!=n/i)

                b.push_back(n/i);

        }

    }

    return b;

}



set<ll> check_s(ll n){

    set<ll>b;

    for (ll i = 1; i<=n  ; ++i) {

        if(n%i==0){

            b.insert(i);

            if(i!=n/i)

                b.insert(n/i);

        }

    }

    return b;

}

bool Is_Prime(ll n){

    if(n==1) return 0;

    for (ll i = 2; i < n ; ++i) {

        if(n%i)

            return 0;

    }

    return 1;

}

vector<ll> Seive(ll n){

    vector<ll>prime(n+1,0);

    prime[0]=prime[1]=0;

    for (ll i = 2; i <=n ; ++i) {

        if (prime[i]==0){

            for (ll j = i+i; j <=n ; j+=i) {

                prime[j]++;

            }

        }

    }

    return prime;

}

ll gcd(ll a,ll b){

    if(a<b)swap(a,b);

    while (b){

        a%=b;

        swap(a,b);

    }

    return a;

}

ll lcm(ll a, ll b){

    return a/__gcd(a,b)*b;

}

ll lcm_t(ll a, ll b){

    return (a*b)/__gcd(a,b);

}



ll power(ll a,ll b){

    if(!b)return 1;// the base case of the recursive function;

    ll res =power(a,b>>1);

    return  res * res * ((b&1) ? a:1);

}

ll power_it(ll a, ll b){

    ll res=1;

    while (b){

        if(b&1){

            res*=a;

        }

        a*=a;

        b>>=1;

    }

    return res;

}

ll power_M(ll a,ll b,ll m){

    if(!b)return 1;// the base case of the recursive function;

    ll res =power_M(a%m,b>>1,m);

    return  (res * res)%m * ((b&1) ? a%m:1);

}



ll nPr(ll n, ll r){

    ll res = 1;

    for (int i = n-r+1; i <=n ; ++i) {

        res*=i;

    }

    return res;

}

ll fact(ll n){

    ll ans=1;

    for (int i = 1; i <=n ; ++i) {

        ans*=i;

    }

    return ans;

}

bool zero_or_one(int n,int index){

    (n>>=index);

    return n&1;

}





ll n,ans;

vector<ll>v;

vector<ll>taken;



void backtracking(ll ind){



    // base case ??

    if(ind==n){

        ll sum=0;

        for (ll i = 0; i < taken.size(); ++i) {

            sum|=taken[i];

        }

        ans+=sum;

        return ;

    }

    // what are the choices ??



    backtracking(ind+1);

    taken.pb(v[ind]);

    backtracking(ind+1);

    taken.pop_back();





}

void Looklols(){

    string s;

    cin >> s;

    ll l = 1, r = s.size() - 1;

    string a;

    a=s[0];

    string b ;

    for (int i = 1; i < s.size()-1; ++i) {

        b+=s[i];

    }

    string c;

    c=s[s.size()-1];

    while (l <= r && r < s.size()){

        if( ( b >= a && b >= c ) || (b <= a && b <= c)){

            cout << a << " " << b << " " << c << '\n';

            return;

        }

        else if(b >= a && b <= c){

            if(b.size() <= 1){

                cout << ":(\n";

                return;

            }

            l++;

            b = s.substr(l, r-l);

            a.push_back(s[l - 1]);

        }

        else{

            r--;

            if(b.size() <= 1){

                cout << ":(\n";

                return;

            }

            b.pop_back();

            c = s.substr(r);

        }

    }

    cout << ":(\n";

}

int main () {

    looklols

    T{

       Looklols();

    }



}