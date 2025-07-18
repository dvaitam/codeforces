#include <bits/stdc++.h>

using namespace std;

int main ()

{

string s1,s2;

cin>>s1;

cin>>s2;

int cnt1=0,cnt2=0;

bool bol=false;

for (int i=0;i<s1.size();i++)

{

if (s1[i]!='|'&&bol==false){cnt1++;}

else {cnt2++;bol=true;}

}

cnt2-=1;

string wor1,wor2;

int i=0;

for(i;i<cnt1;i++)

wor1+=s1[i];

i++;

for(i;cnt2--;i++)

{wor2+=s1[i];}

for (int i=0;i<s2.size();i++)

{

if (wor1.size()>wor2.size())

wor2+=s2[i];

else

wor1+=s2[i];

}

if (wor1.size()==wor2.size())

cout<<wor1<<"|"<<wor2<<endl;

else cout<<"Impossible\n";

return 0;

}