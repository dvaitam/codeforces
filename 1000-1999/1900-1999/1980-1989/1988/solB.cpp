#include<iostream>
using namespace std;
int main()
{int t;
cin>>t;
while(t--)
{int n;
cin>>n;
string s;
cin>>s;
int a=0,b=0;
for(int i=0;i<n;i++)
{if(s[i]=='1')b++;
else if(s[i]=='0'&&s[i-1]!='0')a++;}
if(b>a)cout<<"YES"<<endl;
else cout<<"NO"<<endl;}}