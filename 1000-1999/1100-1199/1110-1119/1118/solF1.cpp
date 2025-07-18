/**************************
  * Writer : Leo101
  * Problem :
  * Tags :
**************************/
#include <iostream>
#include <cstdio>
#include <algorithm>
#include <cstring>
#define File(s) freopen(#s".in", "r", stdin); freopen(#s".out", "w", stdout)
#define gi get_int()
#define for_edge(i, x) for (int i = Head[x]; i != -1; i = Edges[i].Next)
const int Max_N = 4e5;
int get_int()
{
	int x = 0, y = 1; char ch = getchar();
	while ((ch < '0' || ch > '9') && ch != '-')
		ch = getchar();
	if (ch == '-') y = -1,ch = getchar();
	while (ch <= '9' && ch >= '0') {
		x = x * 10 + ch - '0';
		ch = getchar();
	}
	return x * y;
}

class Edge
{
public:
	int Next, To; 
} Edges[Max_N << 1];
int Head[Max_N], E_num;
void Add_edge(int From, int To)
{
	Edges[E_num] = (Edge) {Head[From], To};
	Head[From] = E_num++;
}

int Blue[Max_N], Red[Max_N], B, R, Ans, Color[Max_N];
void Dfs(int Now, int Pre = -1)
{
	for_edge(i, Now) {
		int To = Edges[i].To;
		if (To == Pre) continue;
		Dfs(To, Now);
		Blue[Now] += Blue[To];
		Red[Now] += Red[To];
	}

	if (Color[Now] == 1) Blue[Now]++;
	if (Color[Now] == 2) Red[Now]++;

	if ((Blue[Now] == 0 && Red[Now] == R) || (Red[Now] == 0 && Blue[Now] == B)) {
		Ans++;
	}

}

int main()
{
	memset(Head, -1, sizeof(Head));
	
	int n = gi;
	for (int i = 0; i < n; i++) Color[i] = gi;
	for (int i = 0; i < n; i++)
		if (Color[i] == 1) B++;
		else if (Color[i] == 2) R++;
	for (int i = 1; i < n; i++) {
		int From = gi - 1, To = gi - 1;
		Add_edge(From, To); Add_edge(To, From);
	}

	Dfs(0);

	printf("%d", Ans);

	return 0;
}