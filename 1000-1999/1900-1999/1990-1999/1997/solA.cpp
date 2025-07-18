#include<bits/stdc++.h>
#define int long long
using namespace std;

int helper(string s)
{
    int ans=2;
    for(int i=1;i<s.size();i++)
    {
        if(s[i]==s[i-1])
        {
            ans++;
        }
        else
        {
            ans+=2;
        }
    }
    return ans;
}

string solve()
{
    string s;
    cin>>s;
    string ans="";
    int curr=0;
    for(int i=0;i<=s.size();i++)
    {
        for(char ch='a';ch<='z';ch++)
        {
            string modified=s.substr(0,i)+ch+s.substr(i,s.size()-i);
            int time=helper(modified);
            if(time>curr)
            {
                ans=modified;
                curr=time;
            }
        }
    }
    return ans;
}

int32_t main()
{
    ios::sync_with_stdio(0);
    cin.tie(0);
    int t;
    cin>>t;
    while(t--)
    {
        cout<<solve()<<"\n";
    }
    return 0;
}