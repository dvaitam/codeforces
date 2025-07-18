#include <bits/stdc++.h> 

using namespace std;

 

#define ll long long int

#define ld long double

#define mod 1000000007

#define inf 1e18

#define pb push_back

#define vi vector<ll>

#define vs vector<string>

#define pii pair<ll,ll>

#define ump unordered_map

#define mp make_pair

#define pq_max priority_queue<ll>

#define pq_min priority_queue<ll,vi,greater<ll> >

#define all(n) n.begin(),n.end()

#define ff first

#define ss second

#define mid(l,r) (l+(r-l)/2)

#define bitc(n) builtin_popcount(n)

#define log(args...) { string _s = #args; replace(_s.begin(), _s.end(), ',', ' '); stringstream _ss(_s); istream_iterator<string> _it(_ss); err(_it, args); }

template <typename T> T gcd(T a, T b){if(a%b) return gcd(b,a%b);return b;}

template <typename T> T lcm(T a, T b){return (a*(b/gcd(a,b)));}

void file_i_o()

{

    ios_base::sync_with_stdio(0);

    cin.tie(0);

    cout.tie(0);

    #ifndef ONLINE_JUDGE

        freopen("input.txt", "r", stdin);

        freopen("output.txt", "w", stdout);

    #endif

}





int main(int argc, char const *argv[]) {

    file_i_o();

    //write your code here

    int n;

    cin>>n;

    vector<int> a(n);

    for(int i=0; i<n; i++){

        cin>>a[i];

    }

    sort(a.begin(), a.end());

    vector<int> res(n, -1);

    int i = 0, j = 1;

    while(j<n){

        res[j] = a[i++];

        j += 2;

    }

    j = 0;

    while(i<n){

        if(res[j]!=-1) j++;

        res[j++] = a[i++];

    }

    int cheap = 0;

    for(int i=1; i<n-1; i++){

        if(res[i]<res[i-1] && res[i]<res[i+1]){

            cheap++;

        }

    }

    cout<<cheap<<endl;

    for(int i=0; i<n; i++){

        cout<<res[i]<<' ';

    }

    cout<<endl;

    



    return 0;

}