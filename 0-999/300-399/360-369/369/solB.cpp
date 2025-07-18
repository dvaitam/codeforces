#include <bits/stdc++.h>
#define pb(x) push_back(x)
#define bug cout<<"HERE"<<endl;
#define SSTR( x ) static_cast< std::ostringstream & >( \
        ( std::ostringstream() << std::dec << x ) ).str()
#define clr(x,y) memset(x,y,sizeof(x))
#define all(v) v.begin(),v.end()
#define FOR(i,l) for(int i=0;i<l;++i)
#define FOR1(i,s,l) for(int i=s;i<l;++i)
#define FOR2(i,s) for(int i=s;i>=0;--i)
#define fast 	ios_base::sync_with_stdio(0); cin.tie(0);
#define inp freopen("input.txt", "r", stdin);
#define out freopen("output.txt", "w", stdout);
using namespace std;

typedef long long ll;
typedef vector<int> vi;
inline int toInt(string s){int v; istringstream sin(s);sin>>v;return v;}
inline ll toll(string s){ll v; istringstream sin(s);sin>>v;return v;}

int n,k,l,r,sAll,sK;
int main() {
    fast
	cin>>n>>k>>l>>r>>sAll>>sK;
    int arr[n];
    int rem=sK%k;
    FOR(i,k)arr[i]=sK/k+(rem>0),rem=max(0,rem-1);
    if(n-k)rem=(sAll-sK)%(n-k);
    FOR1(i,k,n)arr[i]=(sAll-sK)/(n-k)+(rem>0),rem=max(0,rem-1);
    FOR(i,n)if(!i)cout<<arr[i];else cout<<" "<<arr[i];
    cout<<"\n";
	return 0;
}