#include <bits/stdc++.h>

using namespace std;

int main()
{
    int i,j,t,d,n;
    string s;
    cin>>t;
    for(i=1;i<=t;i++)
    {
        cin>>n>>s;
        int c=0;
        int a=floor(n/2);
        for(j=0;j<n/2;j++)
        {
            if(abs(s[j]-s[n-j-1])==2 || abs(s[j]-s[n-j-1])==0)         
            c++;
        }
        if(c==a)
        cout<<"YES"<<endl;
        else
        cout<<"NO"<<endl;
        
    }
}