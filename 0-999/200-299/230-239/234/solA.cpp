#include <bits/stdc++.h>



using namespace std;



int main()

{

    freopen ("input.txt","r",stdin);

    freopen ("output.txt","w",stdout);

    int n;

    cin>>n;

    string s;

    cin>>s;

    for(int i=0;i<n/2;i++){

        if(s[i]==s[i+n/2] || (s[i]!=s[i+n/2] && s[i]!='R')){

            cout<<i+1<<" "<<i+1+n/2<<"\n";

        }

        else{

            cout<<i+1+n/2<<" "<<i+1<<"\n";

        }

    }

    return 0;

}