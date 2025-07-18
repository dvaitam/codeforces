#include<bits/stdc++.h>

#define x first

#define y second

using namespace std;

using i64=long long;

const int N=1<<17;

vector<int>sum(N);

void solve()

{ 

    int l,r;

    cin>>l>>r;

    int n=r-l+1;

    vector<int>a(n);

    int v=0;

    for(int i=0;i<n;i++)

      cin>>a[i],v^=a[i];

    auto check=[&](int x)

    { 

      int tot1=0,tot0=0;

      for(int i=0;i<n;i++)

        if(1<<x&i)

          tot1++;

        else tot0++;

      for(int i=0;i<n;i++)

        if(1<<x&a[i])

          tot0--;

        else tot1--;

        return tot1==0&&tot0==0;

    };

    if(n&1)

    {

      cout<<((sum[r+1]^sum[l])^v)<<"\n";

    }

    else

    { 

      int ans=0;

      for(int i=0;i<17;i++)

        if(check(i))ans|=(1<<i);

         cout<<ans<<"\n";

    } 

}

int main()

{

  std::ios::sync_with_stdio(false);

  std::cin.tie(nullptr);

  int t;

  for(int i=1;i+1<N;i++)

    sum[i+1]=sum[i]^i;

  cin>>t;

  while(t--)

  solve();

  return 0;

}