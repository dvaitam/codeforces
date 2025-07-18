#include<bits/stdc++.h>
using namespace std;
int main()
{
    int n;
    cin>>n;
    string s;
    cin>>s;
    int x,y;
    if(n%2)
    {
        cout<<1<<endl;
        cout<<s<<endl;
        return 0;
    }
    int a=0,b=0;
    for(int i=0; i<n; i++)
    {
        if(s[i]=='0') a++;
        else b++;
    }
    if(a!=b)
    {
        cout<<1<<endl;
        cout<<s<<endl;
        return 0;
    }
    else
    {

        //x=n/2;
        cout<<2<<endl;
        cout<<s[0]<<" ";
        for(int i=1; i<n; i++) cout<<s[i];
        cout<<endl;



    }

    /*cout<<2<<endl;
    x=n/2+1;
    y=n-x;
    while(x) cout<<s[i],x--,i++;
    cout<<" ";
    while(y) cout<<s[i],y--,i++;
    cout<<endl;*/

  return 0;
}