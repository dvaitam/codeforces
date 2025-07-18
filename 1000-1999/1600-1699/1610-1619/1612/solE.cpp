#include <bits/stdc++.h>

#define endl "\n"

#define debug(x) cout << #x << ": -----> " << x << endl;

#define inf 0x3f3f3f3f

#define pii pair<int,int>

#define all(x) x.begin()+1,x.end()

#define _all(x) x.begin(),x.end()

#define mod 1000000007

#define ll long long

// #define int long long



using namespace std;



int n;

const int maxn=2e5+10;

int m[maxn],k[maxn],cnt[maxn];

unordered_map<int,int> mp;

struct node{

    int x,y;

    bool operator<(const node &t) const {

        return x>t.x;

    }

} a[maxn];



void sol(){

    cin>>n;

   for(int i=1;i<=n;i++){

        cin>>m[i]>>k[i];

        mp[m[i]]++;

    }



    vector<int> ans;

    double res=0.0;

    for(int i=1;i<=20;++i) {

        for(int j=1;j<=200000;++j) a[j]={0,j};

        for(int j=1;j<=n;++j) {

            if(k[j]>=i) a[m[j]].x+=i;

            else a[m[j]].x+=k[j];

        }

        sort(a+1,a+200001);

        double tmp=0.0;

        vector<int>v;

        for(int j=1;j<=i;++j) tmp+=a[j].x,v.emplace_back(a[j].y);

        tmp/=i;

        if(tmp>res) {

            res=tmp;

            ans=v;

        }

    }

    cout<<ans.size()<<'\n';

    for(auto i:ans) cout<<i<<' ';

    return;





}



signed main(){

    ios::sync_with_stdio(false); cin.tie(nullptr); cout.tie(nullptr);



    int _=1;

    // cin>>_;

    while(_--){

        sol();

    }



    return 0;

}