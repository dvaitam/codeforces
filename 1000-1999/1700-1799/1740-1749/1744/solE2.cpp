#include <bits/stdc++.h>

using namespace std;

#define IOS ios::sync_with_stdio(false), cin.tie(nullptr), cout.tie(nullptr)

#define xx first

#define yy second

#define lowbit(x) (x&(-x))

#define int long long

#define debug(a) cout<<"#a = "<<a<<endl

mt19937 rng((unsigned int) chrono::steady_clock::now().time_since_epoch().count());

typedef pair<int,int> pii;

const int inf=0x3f3f3f3f,N=1,M=2*N;

const int mod = 1e9+7;

int t,a,b,c,d;

map<int,int> cnt,ccc;

vector<pii> q;

vector<int> v;

void dfs(int u,int res){

    if(res>c) return ;

    if(u == q.size()){

        if(a*b/res<=d) v.push_back(res);

        return ;

    }

    for(int i = 0;i<=q[u].yy;i++){

        dfs(u+1,res);

        res *= q[u].xx;

    }

}

signed main(){

    IOS;

    cin>>t;

    while(t--){

        cin>>a>>b>>c>>d;

        int x = a;

        cnt = ccc;

        q.clear();

        v.clear();

        for(int i = 2;i<=x/i;i++){

            while(x%i==0){

                cnt[i]++;

                x/=i;

            }

        }

        if(x>1) cnt[x]++;

        x = b;

        for(int i = 2;i<=x/i;i++){

            while(x%i==0){

                cnt[i]++;

                x/=i;

            }

        }

        if(x>1) cnt[x]++;

        for(auto it : cnt) q.push_back(it);

        dfs(0,1);

        int e = a*b;

        pii ans =  {-1,-1};

        for(auto x : v){

            int y = e/x;

            x = (a/x + 1)*x;

            y = (b/y + 1)*y;

            if(x<=c&&y<=d){

                ans = {x,y};

            }

        }

        cout<<ans.xx<<" "<<ans.yy<<"\n";

    }

    

return 0;

}