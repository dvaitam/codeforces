#include <bits/stdc++.h>

#define ll long long int

using namespace std;

int main()

{

    int t;

    cin>>t;

    while(t--)

    {

        ll n,k,b,s;

        cin>>n>>k>>b>>s;

        ll d=s-(b*k);

        ll a[n];

        a[n-1]=b*k;

        for(int i=0;i<n-1;i++)

        {

        

            if(d>=k-1)

            {a[i]=k-1;d-=k-1;}

            else

            a[i]=0;

        }

      a[n-1]+=d;d=0;

      if(n==1&&k*(b+1)<=s)

      cout<<"-1";

       else if(b*k>s||a[n-1]>=k*(b+1))

       cout<<"-1";

     

       else{

        for(int i=0;i<n;i++)

        {

            cout<<a[i]<<" ";

        }}

        cout<<endl;

      

        

    }

    return 0;

}