#define _CRT_SECURE_NO_WARNINGS
#include<iostream>
#include<string>
#include<algorithm>
#include<memory.h>
#include<vector>
#include<stdio.h>
#include<cstdlib>
#include<map>
#include <set>
#include<stack>
#include<cmath>
#include<queue>

using namespace std;


int main() {
	int h1, a1, c1;
	scanf("%d%d%d", &h1, &a1, &c1);
	int h2, a2;
	scanf("%d%d", &h2, &a2);
	int Ineed = h2 / a1;
	if (h2 %a1 != 0)
		++Ineed;
	int HeNeed = h1 / a2;
	if (h1% a2 != 0)
		++HeNeed;
	int cntH = 0;
	while (Ineed > HeNeed){
		++cntH;
		h1 = (h1 + c1) - a2;
		HeNeed = h1 / a2;
		if (h1%a2 != 0)
			++HeNeed;	
	}
	printf("%d\n", cntH + Ineed);
	for (int i = 0; i < cntH; i++)
		puts("HEAL");
	for (int i = 0; i < Ineed; i++)
		puts("STRIKE");


	return 0;
}