#include<bits/stdc++.h>
using namespace std;
using ll = long long;

//Joshua
int main()
{
	ios_base::sync_with_stdio(false);
	cin.tie(NULL);
    ll t,n,m,flag=0;
    string s;
    cin>>n>>s;
    char ch=s[0];
    char a;
    for(int i=1; i<n; i++)
    {
        if(ch==s[i])
        ch=s[i];
        else
        {
            a=s[i];
            flag=1;
            break;
        }
    }
    if(flag)
    {
        cout<<"YES\n";
        cout<<ch<<a;
    }
    else
    cout<<"NO\n";
    return 0;
}