#include<bits/stdc++.h>
using namespace std;
typedef long long int ll;

vector<ll>v;


int main()
{
    string s,s1;
    int n;
    cin>>n;
    cin>>s>>s1;
    string ss=s,ss1=s1;

    sort(ss.begin(),ss.end());
    sort(ss1.begin(),ss1.end());


    if(ss==ss1)
    {
        for(int i=0; i<n; i++)
        {
            char ch=s1[i];

            if(s[i]!=s1[i])
            {
                for(int j=i+1; j<n; j++)
                {
                    if(s[j]==ch)
                    {
                        swap(s[j],s[j-1]);
                        v.push_back(j);
                        break;
                    }
                }
                if(s1[i]!=s[i])i--;
            }
        }
        cout<<v.size()<<endl;

        for(int i=0; i<v.size(); i++)
            cout<<v[i]<<" ";
    }

    else
        cout<<-1<<endl;

    return 0;

}