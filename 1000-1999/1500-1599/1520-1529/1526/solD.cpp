#include <bits/stdc++.h>
using namespace std;

typedef unsigned long long ull;
typedef long long ll;
typedef long double ld;
typedef pair<int,int> pii;
typedef pair<ll,ll> pll;

#define fastio() ios_base::sync_with_stdio(0);cin.tie(0);cout.tie(0)
#define test() int t;cin>>t;for(int test=1;test<=t;test++)
#define pb push_back
#define nl cout<<"\n"
#define F first
#define S second
#define all(x) x.begin(),x.end()

template<class C> void min_self( C &a, C b ){ a = min(a,b); }
template<class C> void max_self( C &a, C b ){ a = max(a,b); }

const ll MOD = 1000000007;
ll mod( ll n, ll m=MOD ){ n%=m;if(n<0)n+=m;return n; }

const int MAXN = 2e5+5;
const int LOGN = 21;
const ll INF = 1e14;
int dx[] = {1,0,-1,0};
int dy[] = {0,1,0,-1};

template<class T1, class T2> void add( T1 &x, T2 y, ll m = MOD )
{
    x += y;
    if( x >= m )
        x -= m;
}

template<class T1, class T2> void sub( T1 &x, T2 y, ll m = MOD )
{
    x -= y;
    if( x < 0 )
        x += m;
}
map<char,int> val;
map<int,char> conv_back;

void solve()
{
    string s; cin>>s;
    int n = s.size();

    vector<int> idx[4];

    for(int i=0; i<n; i++)
        idx[val[s[i]]].pb(i);

    ll swaps[4][4] = {};

    for(int i=0; i<4; i++)
    {
        for(int j=0; j<4; j++)
        {
            if(i==j) continue;

            int ii = 0, jj = 0, ni = idx[i].size(), nj = idx[j].size();
            while(ii<ni && jj<nj)
            {
                if(idx[i][ii] < idx[j][jj])
                {
                    swaps[i][j] += jj;
                    ii++;
                }
                else
                    jj++;
            }
            swaps[i][j] += ((ll)(ni-ii)*(ll)jj);
            // cout<<i<<" "<<j<<" "<<swaps[i][j],nl;
        }
    }

    ll mx = 0;
    string final = s;

    for(int i=0; i<4; i++)
    {
        for(int j=0; j<4; j++)
        {
            if(i==j) continue;
            for(int k=0; k<4; k++)
            {
                if(i==k || j==k) continue;
                int l = 6-i-j-k;

                ll tans = swaps[i][j] + swaps[i][k] + swaps[i][l];
                tans += (swaps[j][k] + swaps[j][l]);
                tans += swaps[k][l];

                if(tans > mx)
                {
                    mx = tans;
                    string now = "";
                    for(auto cnt: idx[i])
                        now += conv_back[i];
                    for(auto cnt: idx[j])
                        now += conv_back[j];
                    for(auto cnt: idx[k])
                        now += conv_back[k];
                    for(auto cnt: idx[l])
                        now += conv_back[l];
                    final = now;
                }

            }
        }
    }

    cout<<final,nl;
}

int main() 
{
    fastio();
    #ifndef ONLINE_JUDGE
    freopen("input.txt","r",stdin);
    freopen("output.txt","w",stdout);
    #endif

    val['A'] = 0;
    val['N'] = 1;
    val['T'] = 2;
    val['O'] = 3;

    conv_back[0] = 'A';
    conv_back[1] = 'N';
    conv_back[2] = 'T';
    conv_back[3] = 'O';

    test()
    {
        solve();
    }

    

    cerr << "\nTime elapsed: " << 1000 * clock() / CLOCKS_PER_SEC << "ms\n";
    return 0;
}