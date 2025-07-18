#include<bits/stdc++.h>

using namespace std;

#define vec vector<int>

#define vecp vector<pair<int,int>>

#define vecc vector<vector<int>>

#define pb push_back

#define fr first

#define sc second

#define fr1(i,a,b) for(int i = a; i < b; i++)

#define fr2(i,a,b) for(int i = a; i >= b; i--)

#define fr3(i,a,b) for(int i = a; i <= b; i++)

#define pr pair<int,int>

#define mod 1000000007

#define mod1 998244353

#define all(v) (v).begin(),(v).end()

#define sz(x) (int)(x).size()

#define ppb pop_back

#define ins insert

//#define int long long



void write(vec &v){

    for(auto i:v)

        cout << i << " ";

    cout << "\n";

}



void read(vec &v){

    for(auto &i:v)

        cin >> i;

}



const int INF = 1e9;

const int64_t INFF = 1e18;

const int N = 1e6+69;

void perform(vec &op,vec &a){

    int j = 0;

    vecc order;

    for(auto i : op){

        vec temp;

        fr1(k,j,j+i){

            temp.pb(a[k]);

        }

        order.pb(temp);

        j += i;

    }

    reverse(all(order));

    a.clear();

    for(auto i : order){

        for(auto k : i){

            a.pb(k);

        }

    }

}

void solve(){

    int n;

    cin >> n;

    vec a(n);

    read(a);

    vecc op;

    for(int l = 0; l < n/2; l++){

        int r = n - l - 1;

        int ind1 = max_element(a.begin() + l,a.begin() + r + 1) - a.begin();

        int ind2 = min_element(a.begin() + l,a.begin() + r + 1) - a.begin();

        vec curr;

        int mn = min(ind1,ind2);

        int mx = max(ind1,ind2);

        fr1(i,0,l){

            curr.pb(1);

        }

        curr.pb(mn - l + 1);

        if(mx - mn > 1){

            curr.pb(mx - mn - 1);

        }

        curr.pb(r - mx + 1);

        fr1(i,r + 1,n){

            curr.pb(1);

        }

        op.pb(curr);

        perform(curr,a);

    }

    for(int i = n/2 - 1; i >= 0; i--){

        if(a[i] != i + 1){

            vec curr;

            fr1(j,0,i + 1){

                curr.pb(1);

            }

            if(n - 2*(i + 1) > 0)

                curr.pb(n - 2 * (i + 1));

            fr1(j,0,i+1){

                curr.pb(1);

            }

            op.pb(curr);

            perform(curr,a);

        }

    }

    cout << sz(op) << "\n";

    for(auto i : op){

        cout << sz(i) << " ";

        write(i);

    }

}

signed main(){

    ios_base::sync_with_stdio(false);

    cin.tie(NULL);

    int t=1;

    //cin>>t;

    fr3(i,1,t){

        // cout<<"Case #"<<i<<": ";

        solve();

    }

    return 0;

}