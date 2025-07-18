#include <bits/stdc++.h>



using namespace std;



void solve( int t)

{

    for(int i = 0 ; i < t ; i++)

    {

        int n , m ; cin>>n>>m;

        if(n ==1 || m == 1)

        cout<< 1<<' '<<1<<'\n';

        else 

        cout<<2<<' '<<2 <<'\n';

    }

}





int main()

{

    int t; cin>>t;

    solve(t);

}