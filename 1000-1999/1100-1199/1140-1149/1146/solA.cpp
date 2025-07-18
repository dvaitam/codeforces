#include<bits/stdc++.h>
using namespace std;

int main()
{
    ios::sync_with_stdio(false);
    cin.tie(0);
    string s;
    cin>>s;
    int length=s.size();
    int cnt=0,eraseCnt=0;
    for(int i=0; i<length; i++)
    {
        if(s[i]=='a')
            cnt++;
    }
    int Max=2*cnt-1;
    if(Max>=length) cout<<length<<endl;
    else cout<<Max<<endl;
    return 0;
}