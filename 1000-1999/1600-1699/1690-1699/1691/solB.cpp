#include<bits/stdc++.h>

using namespace std;



/// ----------------------------- (Debug) ------------------------------

#define sim template < class c

#define ris return * this

#define dor > debug & operator <<

#define eni(x) sim > typename enable_if<sizeof dud<c>(0) x 1, debug&>::type operator<<(c i) {

sim > struct rge { c b, e; }; sim > rge<c> range(c i, c j) { return rge<c>{i, j}; }

sim > auto dud(c* x) -> decltype(cerr << *x, 0); sim > char dud(...);

struct debug {

#ifndef ONLINE_JUDGE

eni(!=) cerr << boolalpha << i; ris; }

eni(==) ris << range(begin(i), end(i));}

sim, class b dor(pair < b, c > d) {ris << "(" << d.first << ", " << d.second << ")";}

sim dor(rge<c> d) {*this << "["; for (auto it = d.b; it != d.e; ++it) *this << ", " + 2 * (it == d.b) << *it; ris << "]";}

#else

sim dor(const c&) { ris; }

#endif

};

vector<char*> tokenizer(const char* args) {

    char *token = new char[111]; strcpy(token, args); token = strtok(token, ", ");

    vector<char*> v({token});

    while(token = strtok(NULL,", ")) v.push_back(token);

    return reverse(v.begin(), v.end()), v;

}

void debugg(vector<char*> args) { cerr << "\b\b "; }

template <typename Head, typename... Tail>

void debugg(vector<char*> args, Head H, Tail... T) {

    debug() << " [" << args.back() << ": " << H << "] ";

    args.pop_back(); debugg(args, T...);

}

#define harg(...) #__VA_ARGS__

#ifndef ONLINE_JUDGE

#define dbg(...) { debugg(tokenizer(harg(__VA_ARGS__, \b\t-->Line)), __VA_ARGS__, __LINE__); cerr << endl;}

#else

#define dbg(...) { }

#endif

/// -----------------------------------------------------------------------



typedef long long ll;

typedef vector<int> vi;

typedef vector<ll> vl;

typedef pair<int,int> pi;

typedef vector<pi> vpi; 

#define IOS ios::sync_with_stdio(0); cin.tie(0)  

#define F first

#define S second

#define PB push_back

#define EB emplace_back

#define MP make_pair

#define REP(i,a,b) for(i=a;i<=b;i++)

#define RAP(i,a,b) for(i=a;i>=b;i--)

#define spa <<" "<<

#define all(x) (x).begin(), (x).end()

#define sz(x) (int)x.size()

const ll mod=998244353;

const int MX=0x3f3f3f3f;

const int maxn=10;



int main(){

    IOS;

    int t;

    cin>>t;

    while(t--){

        int n,i;

        cin>>n;

        vi s(n);



        map<int,int> mp;

        REP(i,0,n-1) {

            cin>>s[i];

            mp[s[i]]++;

        }



        int flag = 1;

        for(auto u:mp){

            if(u.S == 1){

                flag = 0;

                break;

            }

        }



        if(!flag) cout<<-1<<"\n";

        else{

            int prev = 0;

            REP(i,0,n-2){

                if(s[i] == s[i+1]) cout<<i+1+1<<" ";

                else {

                    cout<<prev+1<<" ";

                    prev = i+1;

                }

            }

            cout<<prev+1<<"\n";

        }

    }

    return 0;

}