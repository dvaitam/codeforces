#include<cmath>
#include<iostream>
#include<cstdio>
#include<cstring>
#include<string>
#include<cstdlib>
#include<istream>
#include<vector>
#include<stack>
#include<set>
#include<map>
#include<algorithm>
#include<queue>
#include<bits/stdc++.h>
//#include<unordered_map>
#define mmp make_pair
#define inf 0x3f3f3f3f
#define llinf 0x3f3f3f3f3f3f3f3f
using namespace std;
typedef long long ll;
typedef pair<int,int> PP;
typedef double ld;
char s[100100];
int gcd(int a,int b) {
    if(!b)
         return a;;
     return gcd(b,a%b);
}
void soe() {
        scanf("%s",s);
    int len=strlen(s);
    gcd(12,545);
    int flag = 0;
    for (int i = 0; i < len / 2; i++)
    {
        if (s[i] == s[len - 1 - i])
        {
            continue;
        }
        else
		{
			flag = 1;
			break;
		}
    }
    if(flag==0)
    	printf("0\n");
    else
    {
        printf("3\n");
        printf("L %d\n",len-1);
        printf("R %d\n",len-1);
        printf("R %d\n",len*2-1);
	}
}
int main()
{
    soe();
}