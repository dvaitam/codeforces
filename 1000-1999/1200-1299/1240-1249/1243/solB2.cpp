#include<bits/stdc++.h>
using namespace std;
const int N = 1e4 + 10;
char s[N], t[N];
int tot[30];
struct node{
	int a, b;
}p[200];
int main()
{
	int T;cin>>T;
	while(T--){
		memset(tot, 0, sizeof tot); 
		int n, cnt = 0, flag = 0;
		scanf("%d", &n);
		scanf("%s", s + 1);
		scanf("%s", t + 1);
		for(int i = 1; i <= n; ++i){
			tot[s[i] - 'a']++;
			tot[t[i] - 'a']++;
		}
		for(int i = 1; i <= 26; ++i) 
			if(tot[i] % 2) flag = 1;
		if(flag) puts("No");
		else {
			puts("Yes");
			for(int i = 1; i <= n; ++i){
				if(s[i] != t[i]){
					for(int j = i + 1; j <= n; ++j){
						if(s[j] == s[i]){
							p[++cnt] = {j, i};
							swap(s[j], t[i]);
							break;
						}
						if(t[j] == s[i]){
							p[++cnt] = {j, j};
							swap(s[j], t[j]);
							p[++cnt] = {j, i};
							swap(s[j], t[i]);
							break;
						}
					}
				}
			}
			printf("%d\n", cnt);
			for(int i = 1; i <= cnt; ++i)
			printf("%d %d\n", p[i].a, p[i].b);
		}
	}
}