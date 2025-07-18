#pragma GCC optimize("O3,unroll-loops")

#pragma GCC target("avx2,bmi,bmi2,lzcnt,popcnt")

#pragma GCC optimize("O3,fast-math,inline")

#include<bits/stdc++.h>

using namespace std;

//#define int long long int

#define fastio ios::sync_with_stdio(0); cin.tie(nullptr);

#define MOD ((int)1e9+7)

int n,m,q,k,cnt(int,int);

int ans;

#define X return void(cout<<-1);

void GETAC()

{

    cin>>n>>k; int k2(k);

    if(!(n%k))

    {

        cout<<"yEs\n";

        while(k2--) cout<<n/k<<' ';

        return;

    }

    k2=(k-1);

    int s(0);

    if((n-k2)&1 and n-k2>0)

    {

        cout<<"yEs\n";

        while(k2--)

            cout<<1<<' ';

        cout<<n-(k-1);

    }

    else if(!((n-2*k2)&1) and n-2*k2>0)

    {

        cout<<"yEs\n";

        while(k2--)

            cout<<2<<' ';

        cout<<n-2*(k-1);

    }

    else cout<<"nO";

}



signed main()

{

    fastio

    int t(1); cin>>t;

    while(t--)GETAC(),cout<<'\n';

}