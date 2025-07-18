#include<bits/stdc++.h>

#include<string>

using namespace std;



typedef long long int ll;

typedef vector<ll> vll;

typedef map<ll,ll> mll;

typedef pair<ll,ll> pll;

#define pi acos(-1)

#define MAX 1e6+6

#define pb push_back

#define eb emplace_back

#define print(x) cout << x << endl

#define isset(x,i) (x&(1<<i))

#define set(x,i) (x|=(1<<i))

#define unset(x,i) {if(isset(x,i)) x^=(1<<i);}

#define unsetbefore(x,k) (x&=((1<<k)-1))

#define loop(start,i,end) for(ll i=start;i<end;i++)

#define rloop(start,i,end) for(ll i=start;i>=end;i--)

#define all(a) (a).begin(),(a).end()

#define rall(a) (a).rbegin(),(a).rend()

#define allset(a,value) memset(a,value,sizeof(a)) //Value= 0 or -1



void solve()

{

    ll n, x, y, d, t;cin>>n>>x>>y;

    rloop(n-1,i,1)

    {

        if(!((y-x)%i)) 

        {

            d=i;

            break;

        }

    }

    t=(y-x)/d;

    for(ll i=x;i<=y;i+=t)

    {

        cout<<i<<' ';

    }

    d=n-1-d;

    for(ll i=x-t;i>0&&d>0;i-=t,d--)

    {

        cout<<i<<' ';

    }

    for(ll i=y+t;d>0;i+=t,d--)

    {

        cout<<i<<' ';

    }

    cout<<endl;

    return;

}



int main()

{

    ios::sync_with_stdio(false);

    cin.tie(nullptr);

    //Close before submission...

    //freopen("input.txt","r",stdin);

    //freopen("output.txt","w",stdout);

    ll t;

    cin >> t;

    while(t--) solve();

    return 0;

}

/*** Output **



*/