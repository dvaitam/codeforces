#include<iostream>
#include<cstdio>
#include<algorithm>
#include<cstring>
#include<cmath>
using namespace std;
char s[105], ans[105],co[255];

int main() {
	int n;
	cin >> n;
	while (n--) {
		memset(co, 0, sizeof(co));
		cin >> s;
		int l = strlen(s);
		ans[l] = 0;
		sort(s, s + l);
		for (int i = 0; i < l; i++)	co[s[i]]++;
		int sum = 0;
		for (int i = l - 1; i > 0; i--)	sum += s[i] - s[i - 1];
		//if (sum > 0 && sum <= 2) { cout << "No answer" << endl; continue; }
		int now = 0;
		for (char i = 'a'; i <= 'z'; i += 2)
			while (co[i]--)ans[now++] = i;
		for (char i = 'b'; i <= 'z'; i += 2)
			while (co[i]--)ans[now++] = i;
		bool flag = 1;
		for(int i=0;i<l;i++)
			if (abs(ans[i + 1] - ans[i]) == 1) { flag = 0; break; }
		if (!flag) {
			memset(co, 0, sizeof(co));
			now = 0;
			for (int i = 0; i < l; i++)	co[s[i]]++;
			for (char i = 'b'; i <= 'z'; i += 2)
				while (co[i]--)ans[now++] = i;
			for (char i = 'a'; i <= 'z'; i += 2)
				while (co[i]--)ans[now++] = i;
			flag = 1;
			for (int i = 0; i < l; i++)
				if (abs(ans[i + 1] - ans[i]) == 1) { flag = 0; break; }
		}
		if (!flag)cout << "No answer" << endl;
		else cout << ans << endl;

	}

#ifdef _DEBUG
	system("pause");
#endif
	return 0;
}