#include <bits/stdc++.h>

using namespace std;

#define int long long int



void solve()

{

    //cout<<"--------------------------------"<<endl;

    int a,b; cin>>a>>b;

    if(b==1){

        cout<<"YES"<<endl;

        for(int i=0;i<a;i++){cout<<i+1<<endl;}}



    else if(a%2!=0){cout<<"NO"<<endl;}

    else{

        cout<<"YES"<<endl; int x=1,y=b;

        for(int i=0;i<a/2;i++){

            vector<int> v1,v2;

            while(y--){v1.push_back(x);x++;v2.push_back(x);x++;}

            y=b;

            for(int j=0;j<b;j++){cout<<v1[j]<<" ";} cout<<endl;

            for(int j=0;j<b;j++){cout<<v2[j]<<" ";} cout<<endl;

        }

    }

}



int32_t main()

{

   ios_base::sync_with_stdio(false);

   int t; cin>>t;

   while(t--){solve();}

   //solve();

}