#include <bits/stdc++.h>
#define ll long long 


using namespace std;

int main()
{
    ios::sync_with_stdio(false);
    ios_base::sync_with_stdio(false);
    
    
    string str;
    cin>>str;

    vector<ll> v(str.length(),0);
    v[0] = 0;
    
    for(ll i=1;i<str.length();i++)
    {
        
        if(i==str.length()-1 && str[i]=='a')
        {
            v[i] = 1;
            for(ll j=0;j<(i/2);j++)
            {
                swap(str[j],str[i-j]);
            }    
        }
        
        else if(str[i]!=str[i+1] && i!=str.length()-1)
        {
            v[i] = 1;
            for(ll j=0;j<(i/2);j++)
            {
                swap(str[j],str[i-j]);
            }
        }
        else
        {
            v[i]  = 0;
        }
       // cout<<str<<endl;
        
    }
    
    
    for(ll i=0;i<str.length();i++)
    {
        cout<<v[i]<<" ";
    }
    
    return 0;

    
}