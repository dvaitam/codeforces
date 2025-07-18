//Codeforces Problem B
//Antoni D�ugosz
#include <bits/stdc++.h>
using namespace std;

string s,s2;

int main()
{
    mt19937 rng(chrono::high_resolution_clock::now().time_since_epoch().count());
    ios_base::sync_with_stdio(false);
    cin.tie(0);
    cout.tie(0);
    int i,j,p=0;
    cin>>s;
    for(i=0;i<s.size();i++)
        if(s[i]!='a')
            s2.push_back(s[i]);
    if(s2.empty())
    {
        cout<<s;
        return 0;
    }
    if(s2.size()%2==0&&s2.substr(0,s2.size()/2)==s.substr(s.size()-s2.size()/2,s2.size()/2))
        cout<<s.substr(0,s.size()-s2.size()/2);
    else
        cout<<":(";

}