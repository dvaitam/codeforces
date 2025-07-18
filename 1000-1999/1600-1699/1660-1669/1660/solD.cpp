#include<bits/stdc++.h>

#define x first

#define y second

using namespace std;

using i64=long long;

void solve()

{ 

  int n;

  cin>>n;

 vector<int>a(n);

 for(int i=0;i<n;i++)

 cin>>a[i];

 int p=0,l=0,r=0;

 for(int i=0,j=-1;i<=n;i++)

 {  

    if(i==n||a[i]==0)

    {

      pair<int,int>mn[2]={pair(0,j+1),pair(n,-1)};

      int pw=0,sign=0;

      for(int k=j+1;k<i;k++)

      {

        if(a[k]<0)sign^=1;

        if(abs(a[k])==2)pw++;

        if(pw-mn[sign].x>p)

        {

            p=pw-mn[sign].x;

            l=mn[sign].y;

            r=k+1;

        }

        mn[sign]=min(mn[sign],pair(pw,k+1));

      }

      j=i;

    }

 } 

 cout<<l<<" "<<n-r<<"\n";

}

int main()

{

  std::ios::sync_with_stdio(false);

  std::cin.tie(nullptr);

  int t;

  cin>>t;

  while(t--)

  solve();

  return 0;

}