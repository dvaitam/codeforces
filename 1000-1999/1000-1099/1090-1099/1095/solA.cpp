#include<bits/stdc++.h>
using namespace std;
int main()
{
    int n;
    cin >> n;
    string s,a= "";
    cin >> s;
    int cnt=1,l = s.size();
    for(int i=0; i<l; i+=cnt)
    {
        a += s[i];


        cnt++;
    }
cout << a << endl;

}