#include<bits/stdc++.h>

#define ll long long

#define se second

#define ff first

#define pb(x) push_back(x)

#define pf(x) push_front(x)

#define pp pop_back()

#define ppf pop_front()

#define T int test;cin>>test;while(test--)

#define all(x) x.begin(),x.end()

int const N=2e5+10;

int const mod=1e9+7;

ll INF=2e18;

ll POW(ll x,ll  y);

using namespace std;



int main ()

{

    ios_base::sync_with_stdio(0);cin.tie(0);cout.tie(0);



    T

    {

        int x,y;

        cin>>x>>y;

        if((x%2&&y%2==0)||(x%2==0&&y%2)) cout<<"-1 -1\n";

        else

        {

            if(x%2)

            {

                cout<<x/2<<" "<<(y+1)/2<<"\n";

            }

            else

            {

                cout<<x/2<<" "<<y/2<<"\n";

            }

        }





    }



    return 0;



}

ll POW(ll x,ll  y)

{

    if(y==0) return 1;

    ll res=POW(x,y/2)%mod;

    res=(res*res)%mod;

    if(y%2) res=(res*x)%mod;

    return res%mod;

}