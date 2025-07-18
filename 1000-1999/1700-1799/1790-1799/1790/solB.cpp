#include <bits/stdc++.h>

#define int long long

#define fs first

#define sd second

#define MAX LONG_MAX

#define pb push_back

#define MP make_pair

#define optimizar_rend ios_base::sync_with_stdio(0);

#define input cin.tie(0);

#define fixed cout.setf(ios::fixed);cout.precision(0);

using namespace std;

int Mod = 1e9 + 7, INF=1e18;



int32_t main()

{

    optimizar_rend input

        int t; cin>>t;

        while(t--)

        {

            int n,s,r; cin>>n>>s>>r;

            cout<<s-r<<" ";

            if(r%(n-1)==0) for(int i=1;i<n;i++) cout<<r/(n-1)<<" ";

            else{

                vector<int> v(n-1,0);

                int i=0;

                while(r!=0){

                    r--;

                    v[i]++; i++;

                    if(i==v.size()) i=0;

                }

                for(int i=0;i<v.size();i++) cout<<v[i]<<" ";

            }

            cout<<"\n";

        }      

                       

    return 0;

}