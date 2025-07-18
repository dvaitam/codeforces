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



int main () {

    looklols

    T{

        ll n;

        cin>>n;

        vector<pair<ll,ll>>v;

        vector<pair<ll,ll>>ans;

        for (int i = 0; i < n; ++i) {

            int x;

            cin>>x;

            v.pb({x,i});

        }

        sort(all(v));

        for (int i = 1; i < n; ++i) {

            ans.pb({v[i].S,v[i-1].F-v[i].F%v[i-1].F});

            v[i].F+=(v[i-1].F-v[i].F%v[i-1].F);

        }

        cout<<ans.size()<<endl;

        for(auto i:ans){

            cout<<i.F+1<<" "<<i.S<<endl;

        }

    }

}