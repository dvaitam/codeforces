#include <iostream>
#include <string>
using namespace std;
bool pre(string s, string a) 
{ 
    int d=s.size()-a.size(); 
    return d<0?0:s.substr(d)==a; //string method
}
int o, O, f, z, i;
string t, s[6] = {"lios","liala","etr","etra","initis","inites"};
int V['zz'];
int main() 
{
    while (cin >> t) 
    {
        O=0;
        for(i=0;i<6;i++) 
            if (pre(t, s[i])) 
            { 
                if (z&&V[z-1]>i) 
                    break; 
                O = 1;
                V[z++]=i; //should follows grammer order
                f+=i/2==1; //noun count
                o+=i%2; //
            }
        if (!O) 
        { 
            cout << "NO\n"; return 0; 
        }
    }//o==z f==1    
    cout << (z!=1&&(o%z||f!=1)?"NO\n":"YES\n");
}