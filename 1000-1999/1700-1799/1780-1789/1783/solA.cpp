#include <bits/stdc++.h>

#define ll long long

#define ld long double

#define vi vector<int>

#define OO INT_MAX

#define endl "\n"

#define fast ios_base::sync_with_stdio(0);cin.tie(0);cout.tie(0);

using namespace std;

//    cout<<fixed<<setprecision (9);



void usaco(string s="")

{



    if(s.size()!=0)

    {

        freopen((s+".in").c_str(),"r",stdin);

        freopen((s+".out").c_str(),"w",stdout);

    }



}

void solve()

{

  int n;

  cin>>n;

  vector<int>v;

  int b=0;

  for(int i=0;i<n;i++)

  {

      int x;

      cin>>x;

      v.push_back(x);



  }

  for(int i=1;i<n;i++)

  {

      if(v[i]!=v[i-1])

      {

          b=1;

          break;

        }



  }



  if(b==0)

  {

      cout<<"NO"<<endl;

      return;

  }

    cout<<"YES"<<endl;

    sort(v.begin(),v.end());

    cout<<v[0]<<" ";

    for(int i=n-1;i>0;i--)

        cout<<v[i]<<" ";

    cout<<endl;





}



int main()

{

    fast

    //usaco("moocast");

    int t=1;

     cin>>t;

    while(t--)

        solve();



    return 0;

}