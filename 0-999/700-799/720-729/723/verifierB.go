package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var rawTestcases = []string{
	"100",
	"217",
	"____GFzYt___(Gisi)ZqITZMjtgUEJgwBu(JECH)_(f)(T_YQOaN)___TemKopCffuGFgtJ_(_InZMJ)(CfMZy)_____lmlcN(EefRW)cfS(THrHZpn)(_ALrCFQPSY)______LOvmpbUrxYkv_Soc(MRebhOmMK)fxhcMbmEn(R)BNgqeoePt____G_(Sz)wUEKk(XdYR)vHqCQlaERAKGty",
	"215",
	"(jJSaD)(vV)ripWEwNsR_____(iTtyVAPfaM)kpoOCy(AczSKAXQTc)eqSkCHFJ(c)_______mJOfUia_____TNgmhMPmtr(gEz)(br)___(_PPwhjrbc)RqJu_____PFTPDOBxIlLsaijrv___(vX)(c)__sxzJishEUplHUet_______JE_____(BctvVRjk)(yZOfeZfmVo)agzJHsC_",
	"109",
	"oqLXkBheSbHC(hFzq)ZjgmDyxIjgM__(ORBHFRu)____(INoa)____(qMY)LsTTZEeceoictaWC__xGyHGcKf(WMeVBWnsI)(_EYyMLo)VGK_",
	"17",
	"___________PiUvdA",
	"37",
	"_____(RS)AcNDyDEXc(iOugSJPwmy)_SDNOvR",
	"183",
	"(YiyZs)(HY)Y_________(cc)___QKKnofXOX(AGthjBKBf)egAXjUb____(bFuUqfw)wSwlaoxe(nan)(Va)___(M)lDh___(ibnxv)________(lLfgIL)____uGppxAQci(z)iA_______(tOwfpCO)_____AaAUuCn____rh(YSC)______",
	"45",
	"vHjwDOOf___(NDanthXO)_____(BTW)_(WoIWzrOb)Qcz",
	"135",
	"(CgVqwsW)(Mfce)___(hHpWke)________HOnIgAOIzV______(khS)zLDiJQt____(nEFSGuF)(C)__(dNnbwEza)(Rf)(axchNar)(UojWK)__DTvykvA(jCTjHui)Cw_____",
	"187",
	"CnLTcoOflxd(oNt)(TG)_______QPJVBLDFq(nvrccdkw)PaieY__(JoDmvMg)(Yu)_________(m)_wmmqRUVtt(qEwTpct)(aD)____FDChffpAnC(_B)(clpFoir)____(NIYmT)____(DIOqrobhNY)(lU)__VQaVI(dhyPrhV)(oRTTIQ)__eR",
	"84",
	"UE_(poQLvQvy)(cRdF)M(Q)(e__kCyue)ncWEzeULFKGFRjk(_FztgrXxwO)_GPHpmwQe__ZG(gmxiNI)P____K_DXIGvM(t)TUcU__QHWVVb_cD__s(AiCSuxR)tAG(yurVriw)C_T__TNV___",
	"180",
	"(dM(dAgC)wCe(iNA_HQs)K)ztt(z)j(T_M)SRwX___Bqxo_EQh(V_(Ef_ZYVMSZD)Z)UCuvU_NKkwGXrNhv)dJ_Uk(w)zADysgnKi__w(bryvOhs(mQ(nh_QeO)EXu_pTcEUbZtJiOOjWrYQqXWtuU)pMZH)Uej_ixKtsdr__U(nD)(gOtY)h_cM(H_XWAgz)hNSVVbNXm",
	"158",
	"jxnHML(FKHZYg)eAJKEx__DvbKX__kFRQXpyQtrUxXNjx__cnLM(dwtJ)_oPDC_____NIRFrzMza(_g)SBVniLnm_Wnv(RF__rT_L_QvwYH(IFmYtn_)(kH)WGvvpB(cZzg)DQU_fgk__bSphNNSBTtrzJfcT__(aHdJYp)DGA_sMrVl(_fk)NQx_rzuwiIF",
	"164",
	"zJJYFX_qYhKgYRGKAtEepwJCqJTp_rrxoxlcKvc__n__ueKMuv__RD_Io((Yg)zLfH)(iq)mxwJBFZTwEU_PrQmCuoFOF(cK)NxuOhDQrhXSXTHxe(L)(rI)RsflZJTBOjbihcRQmKInIK(nSa)ckekNwxsW__oyItOz_owmxzQSyGMD_ncrbJqok",
	"34",
	"zyz(P)KkSOSiNdE(DK)stIWEoZnlmZtNdyj",
	"222",
	"azy___sZdFCaNX(If)i__YxqknGdKnvgXXZQSrZYkyjJpbgNviMIIQWjxsBVOvFDZeuNv(fkt)BdRCGXMNGOsuA___uzf___wDBQcPNOFhn___S(yACkGW)HKjFLcMpz__pwuPWqExCGReQMD___ECtpXFNLX_____FvDdxEr_-rZKYDyGJHvubTwZqurvmkhzOs_dUIZCzhk(J)bRqjk_mDc____xnI_(_TQW)B",
	"64",
	"q_jowXPLnKDnjZU(rnmP)Zacazgzfu(YDPx)kS___HZTBghfVPbwMlGqre",
	"37",
	"___EyiYKCb(fC__BAOM)MveFuxYz_pXKSM",
	"165",
	"LbdzH__lPLWctFz(XHiEhoV)zAOD(zE)amgXeVXpyDG(Whau)RcVzpLE__VUXjbUG_______rkKcFjkXzz_EgMzjIBfzZH(A)ZVMniXMDdMnuvQKlhPIlAmHnOnbZ(Vx)JcKx_EFRijBTAOWbRTnAkJE_rqc_sYGUItIPJQ",
	"162",
	"OVXVULjEIQ_etuJgHW__e_kuHrdWYBBPCz(LoCkA(MNU)X)HTSkUHYxSa_eWZAgNWWCzqpfz__PS_cipxOoPuFvGzGWIDzsZ__LPgJw____QI_O_ipKQe_gS________vK__PnkTvlfLtFl_Ab(yWHW)jGYhxjIa____tVz__Iw_",
	"168",
	"CrEvpaD___(tlL)JPz__M__sePcwTmSrE___V_PrFV_oJGZC(H)UXho_pBS(YFo)kJFGC__NaXiuUZi_HQXg__(Cn)pn____XF(jhvL(fyFAu)JYQQW(_gZlFWZQtNL)HXsmJxe)YJycdz_(pH)aXkJn(iPcywAq)Bb__HrUUXtLEb_SlsXfsCN",
	"69",
	"xddYwLLMo_gwb_RmKI(VAPk)fnpzD__IYg_wcuVn_____",
	"26",
	"nm_qjo(ujzfD)sPDrxi_jxJm_Vrnm",
	"59",
	"Myl_yX_ydlhD___EN__(wYysVJp)HZ___aMsghpRtVoeIz",
	"78",
	"iLmYirccDgZGzKJjCGsDwGQf_mfjZkYaaYVKjRyCIZ(ol)G_oCaAWfZlgmnV",
	"57",
	"FoKfcV_pbwJM_RrZgsXvaavp(qWt)YVXxvYFlxgXMOgkDXmOeRQWvlwG",
	"180",
	"Mzrv(xY)EmwOtkMlqzZwzFSd_SXFCH_BwRZeRXxsG_wsIRg(eLiGceylAY)jSP(QoIn)RYiMQji(LK)YnpP(hzYLw_T)eQSYuXJ_KCkfTLAWMAmzunWd__J(ntG)kbshEdype(NfGMjpWWAT)(KMwrn)yFyFsjOxUlkXdiFiShWv_E_RIGUIKPZY",
	"64",
	"TqFLEWkg_RsmZvgKkl(YTTkyTR)HaGUb(QTQ)KpfYCGnSfCzWagJ",
	"201",
	"YeuylZnmtrK_mV(OEyfclBfvD)aagqiwyksmTHy_ryYjVOhNfRzhtYvVEUYCFpFhWUAvnG(etxdSbKrAA)RwGGNbMiy___LHqCrORY_fpkphcrXbU_ThpYSzZy(JfdYo)hHQNhZxfxcT_mxpjFfiGNmJtWcIUSZyaRaJDrsYO_cfLxnbRLtLOJWhzfYRAkYhK",
	"215",
	"yrS(eMZ(omRHfOAW)gAEaaFjigJ)RBmwGyyzRFsTvgUhUw__Knq(yvUedcFPGn)ZgYEPdxhnxHJLBrQxnIpOgxQnLpzUQgjjOeghsuwwsFMsxCLtnZRSvk(q)cyySKdPp_tKF(tEyaAU)UyyLjVMBvXcRTIQIF_vTRSZRpwCdbsVVyRJI__Yy_urL",
	"113",
	"tRJiSXYUzUmdgViFFvaIPGXlJP_jNUpKWlUz_XgZNCpdx__aAgTrBMMjcP(Dp)qpNyVldzhJOKItSHuuHcPzWUYpLdufDJi(JGS)(",
	"31",
	"SV_zlKRWGmOIbuwuyYzGMpbqEXFkbiN",
	"51",
	"tDxpJwIXdWt____(MEwU)_L_He_(ZHCbe)(KY)gQOS",
	"127",
	"RfaW_rhfnqf_ee_hwPY_bnzZNRJNJOeYFJkhUTrjJIHupMkPGyuYvwIeyxGPCxVeTDQcKWrFjkbdRJUymkPuZxfxxJoFRIoOSrgaTSL_lHIkuvpwkrbhveGKWELCrwbdTH_BF_zAeMWfHj__XLr(",
	"37",
	"QZv(Cywtx)kjEsMd_EbNg_AuEwsb",
	"166",
	"FVzPfQlhCOmsIyNjOWPGBFMJGyKskMExymZnEUQlXXoSGYDDcUNSaBrYOtCWoIgvCu_nYMFsKRuTpAnGoLh__RUJcnVabQqHYpNKmrufzRA(nDmRyAUK)((RJf)NndcUQJ_m)H(o)____lnVC_mpoB_O(kKCYe)eHfcs_yXC",
	"14",
	"k_Zb(WN)zNzLjM",
	"80",
	"_NxOcJP_____CkTGWgbLMw__ldLmuICD__ItLoPvBzzgK",
	"160",
	"TNfXsnMcdezzvFmYDsCuHvmLuNl__evQszbnb(QFmrVqzN)cTTlYhmuFwuNxppmrYjgNPZS__kVqvydXV(ZgVbw)GJfDSJhDkrhjk_hLyqbv_dzgD__kX(scJw)flqkbYopcY(bp)b___TUuYtoN_otWjFcZZUWV(TJoxYFmCE)KQNjzUcD_MjIpPUmSaGbDMziukKCuJ",
	"151",
	"yPKxMdCZnaevGIYlHCx_lHU(HTdrioV)UnVoRebAkAzdPPYIqDKyk___tS_dwwMzyQMPCqgYPmgeOmQnwa_lClHHkIUVEh_CHLPXK_eGjNBLKgzwP(hZSWCxoR)ktCcFS",
	"60",
	"BOBLSFAk_BcYE__xM___rUElnTCh___vJubJrqVRKWcXm",
	"213",
	"rRQtN(reYBi)pjl__FvMrs_UcAJl_oRnn_yjpb___XTOSsxohlbzaOg____LKkya__Aj(AagkOQnXkX)dKLVulHlrDHLvyFHmSOQBKnYdfUdGdricLZlBV_WqgQqgTliWoOYZh_agXdMoaEWBhq(Zw)cMDVwekfjl(SEpr)ZLgdbuMfFl_zHVndsPGsiKl(CctvxLZM)lbKn",
	"57",
	"W(TjB)_i_jfSUFmceNfvnhqhzQwfVUiLZ_MJHnx_vnW",
	"29",
	"NlEtYIqjCahqTfgUAvikLvJlvqIJo",
	"129",
	"tWXYQCveOfaDTOWszGihMUtkCbDqmTZlPboiWzzP(ZbaNyEAtx)AOPbVJFUd__CBcZjuVDPhZnym___XHkpC()YpmZRZMXQNXZXqcJhKKV__fOnBAWKrOfytvXcykTnXMWqKZcqFlvp(W)RUYQ_CnpcdKJtj_f",
	"56",
	"ZjbQ_TVhONudYVe(Nwii)MNINQdDSaOUYHuvM_",
	"48",
	"NeEsGCdyqNDYpDMERNFyYYARWsvPjjbnnxdlEdrsnCe",
	"192",
	"xUdWdw_Hl___VvqPyjjAzd_mrVkqYNj_(ijRzt)fDdYkDchIFmxagKhWkioThFnWBBrSoRxJFBC_ayTzOKZPNaXOwJ(nxKrT__hpHqpGE)cBHtQuBwCKN(iM)tJWyKDX___nVrhJtGjF(_C)jVgGjpWoTdrSGwmPzplePuNXdZjRXzOOLwjNpYbzZIjtZ_",
	"90",
	"hD(viKrz)yQOOsklSlCsaILu_P__zms_CeTwvaNdnD___wklta",
	"23",
	"(w)seJT_H_PtwZQz(fEN)ZkM",
	"29",
	"zccQytTYGZDlu_RuUPQ_rzWPEcINz",
	"95",
	"wBzXkKioPLWAlMqzSiiUFQEXFzEQI_rTVMrzqTyJetmgGZJboznRAKbL",
	"135",
	"xEUdUhF_Tm_WcIgXxLVTxzoQUM_mcWzWyEsEJkN__bPzGXwWSMpAly_owLQMR_ZHgabCoXjPjwFxysmbSQPPPjEFNTpDyjQtNPImN",
	"204",
	"rbLvS_lgbgQmKTLthtZKzMjpaNPZpDXOOulOIYjvKrHqjkcygAqAfGrFMgmOJHPhU_nY(ncnm)NRTzSfFpoN__mC_(SQa)kQPEPyr_grJKWABtktO(tJteHWQ)EMfTDTshux___ZRLo((poy)oozThKDEUQv_aTyc)ZGVPudANgdV(IHBdwbz)KcdqW__ZqdEMZI_l_SNuToa",
	"211",
	"HmLISvtQz____GbBgVTxz_TVCV____RYEVxXgKbKdLmCGdHmYFlVKNkMuRDCU_MDQq_zgG__gkwtydvScp_nvWnk_lSllaeBEidBGFJiIDlnHAWzxOvRgxAGFlXFXdYJgiWLwEy_ERCYGxNRbjURke_bFFpNHisCxntJ_yqmYaJVv",
	"55",
	"XXMrG_xGxz(EmAIVuItn)US_Ex(EbtkLtqO)D_",
	"52",
	"(mfkN)cavH_vIG_XXXXUD(K)jyjxe_X",
	"121",
	"(qtHdXkLfzRh_)uRqIUwrFWpGKpbQ_rBhgtlcOhZIoRYcxykJsDQigbiPnUmHKbYmqaPmeVUkXh________L(_tOR(p)faipf)GkFgtLehAs",
	"46",
	"yTDAY_____niX(J)kRQOVPQrFR(U)qkHjSgvQX",
	"28",
	"J(omkXT)gxruMlSnXaxT(k)aiIbObGb",
	"63",
	"PcK____NZaNgmtYFbJG(TCBNUD)__ul__WFihavAkyvkx",
	"202",
	"UHRzogxKdtSZqpTjUF__jEXlQPpcAlGyhFdbKohpNBdBRRSYsbxKjJJsnqu_TvWPxAJN__aIqxvqysQCy_lxKobFjplyABm__zriQkt_nTFDye_gUCwMyFZXJpS_lWfPqcLuNATixbQvaSXOkrlFezCEepljq_uAKKm__rfwwpWiYZRWw",
	"124",
	"JsqtzDxxzCzcgNxrpeNwJhSZMUSioVVajiGbJkdLzssZbcvNzZjkoRBkqTLTmnFubsbNWsMHpLUNGPcGyC_kQWdCarVG_FKxeGJEr__",
	"98",
	"UbcCoXKMXNhMLfrAaIogBcixtD_rbGeVzRLnqfaUhQwTQfPQzzbzuuEEE",
	"175",
	"VquP_CO(jbcjbo)pTTrrytiKlOFiS(bQ)jgWYBIai(w)pFwwB(nYyRPa_E)Zo_hkSuVtVGRx_rFHQhgyRdqEYpu(dTtSj)CdwGuKlvWHRgPxGMN_LiwQzdZm_KA_eHoCkhA_vOThzcUCpIMuHbLlTjaI",
	"16",
	"kQc__OJIHvD__lCBn",
	"66",
	"CBDwgIoPDdIgScG__kBLdDumurSfKNhWTlgBtyeQnDYlIFKmPmxt_qy",
	"29",
	"_bnCeGQ_vTBbPQVWPGNhuFGaivQfI",
	"190",
	"xidypbkKBXtmzsWqpRItZKXvToxPpeEvIHCMmE_HioTmfwFSzDYemIIEaVzAqryKTg(Pqk)qqpShk(zRa)iCK_hsIDwTtKiiiAGmwDDb_TpjfWngd_BYDhBVevnAmcOAkFiyNrDywknviXRUYTntCaRPxnAkQfIHMgvrSpmVEpcAJMlQOKWKiM",
	"35",
	"WIcWTBHmWRvAyifnjzKrdHxZRlFRFTGG",
	"80",
	"Btz(Edmhmq)L(FLiZZ)nghvBPlpDVkUzB____fhYWRiPKd",
	"76",
	"jxnWUfeOryfc_rJuwhus_rJWkc(u_op_tY__)iR_pJiWeOrqgQzpxhXXpZo_I",
	"136",
	"pNPl__MxukvkrUDBXcFfvokauGWYjvNXU(vSkM(vCYkfX)QdbvbKcmo)mAO(zqVg_JkE)LZcHAreFx(gCDSTy)gfJufPZoyQibpxzGMHty(ThUys)IPVPbVH(kCp)AxNQRUvmap",
	"208",
	"ZAiYcMMrOetjaPy(hRUs_iMfVV)KZVQtOhtBaPLFUKNvhvRKmCoOpYQsnifSpsJyH_aQsgShm_UqWMnTuXwdAvAkTArGpZpwJ__)wW_SvZ_iNZtpKUIgBDjDXaeqVqyIeN_kICLN(hBQzlucPX)JUZHO___PIWIBl(DrpVde)RzN",
	"36",
	"hXyoBauUC_uaiU_GMpJYAjmCDkqdVEH__OEn",
	"82",
	"aaJHAoIucRAOZpWiiVoMJKWvGwXo__rRPqwqKGCvPOoOtQiEJyvRCgooGq__m",
	"60",
	"UKQZkDvKeczFcGFXQuLBlTlwUpOt_BZoGttNEsFaKXJQfHDLCyNkvRmzEESj",
	"58",
	"Q_uVzWXlFfRFZMUYRJMu_nyyWRnvuvJqeCw(gK)bkUFShtXjpXSdqzYbQaW",
	"129",
	"KVhOSKXKpDehAHFbRRvyfNjaXPSySwzERWELfemkHXSqNGGjKBxtpw_(nkWGOOgGeItv)e_tgMRDzWRoXKgZxPAB(_OcFAmubGa)KBW(LrRZ)I_tl_mfJGnW",
	"87",
	"B(DQLKfEUI)lOMFD_vmeCHOC____vYA_____OK__QfCTEZdlvypfzLZqpJDNiKoJOHsT",
	"26",
	"fr(lWq)bZJhtPBOwwANwoBMuxOk",
	"145",
	"rOj_jdznnxGWswBiMiuzCVlYgmndsmKhtmtoSuFX_xlDbGwvUMSZjSyuheKzHA_imbvGJdhobFfmxZmxZMPZUWjMjzuhCa_tgsczqxdGMhcrxZGXqZNiKgJ__ruwFGgusBFk",
	"12",
	"(Zb)_Br_KOgV",
	"64",
	"(IJllkqsQy)(uVD)S_V(ecv_RiiBBXsDJ)z__FOgRHnWF(WyHH)QkQtwBiL",
	"27",
	"fYmlcbpJfDCxnVJLkEG_Exzwjzr",
	"196",
	"vq_e__gNubJeOQPaUDstTHxLQd_ouKXz__yPgxCJVOmtXDn(KBIu)yFjJmXbSxDVzDMunfxqtjDprA_YJaAGYfI_C__JVvChXJR_uIyd(oL_G_e)miA(wiWTBYF(waH)gJGSxxfYnwDyR_AGSrpBZyo_lCIyoeKYTuoyN_X_JgloXMHDcm_",
	"104",
	"xDMZCgQHkMoYrMnfJhmlolKGXPRqzLEzzjbvnUYunCxqmFCnOYVLJXmLvLaLAH__oXolPVgdNmRShHElNEWNvXgUo",
	"181",
	"PzKCaijzsOwByWBcHySLXiEPdksXsFg_ygIJXyADmBNkwnfVLFuGKnPZkV(SPJ)RdggCZRJoYsbJj__ukNrSvsBQpxiUeYsGoJ_YoABmt",
	"26",
	"JPlJAgTTREHnNYnNOzChd_tBpY",
	"75",
	"zNfBejDObakbXnxVYCWvngbmJIijBjJrfHwlrQoQPAOhkjCGzPzeJkBYlNle",
	"116",
	"VeBMmGcy_ejqtWovodWdrYjPjz_fZdBnYBbGqsTDeXzSDTvi(_HgLIDLOJP)zZjDTUTUiBubkDJeDLVwfwBVJK__ihiYdrglzUbxuAl",
	"139",
	"BD_RlTFVPh_VHmiRgQXSaFmToV(Eq)JtZhhELpXLscH(DI)nxoapvVCVaMiwsYsYrrqLLETNmxVdIGavfoBLbwPtJuCOcORExK",
	"118",
	"gEeXWPvdURvgWqfGgP_krWfdoHHElTvyZVgJFJTdKmxZCibcVLrTLCCvnHYtQttnlItWpd_oScE_KnS_UDxbPV_vYQuuX",
}

func expectedB(n int, s string) (int, int) {
	inside := false
	current := 0
	maxOut := 0
	inWords := 0
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') {
			current++
		} else {
			if current > 0 {
				if inside {
					inWords++
				} else if current > maxOut {
					maxOut = current
				}
				current = 0
			}
			if ch == '(' {
				inside = true
			} else if ch == ')' {
				inside = false
			}
		}
	}
	if current > 0 {
		if inside {
			inWords++
		} else if current > maxOut {
			maxOut = current
		}
	}
	return maxOut, inWords
}

type testCase struct {
	n     int
	s     string
	input string
	want  string
}

func parseCases() []testCase {
	var cases []testCase
	for i := 1; i < len(rawTestcases); i += 2 {
		n, _ := strconv.Atoi(strings.TrimSpace(rawTestcases[i]))
		s := rawTestcases[i+1]
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n%s\n", n, s)
		out, in := expectedB(n, s)
		want := fmt.Sprintf("%d %d\n", out, in)
		cases = append(cases, testCase{n: n, s: s, input: sb.String(), want: want})
	}
	return cases
}

func runCase(exe, input string) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	cases := parseCases()
	for i, tc := range cases {
		got, err := runCase(exe, tc.input)
		if err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		if got != strings.TrimSpace(tc.want) {
			fmt.Printf("case %d failed: expected %s got %s\ninput:\n%s", i+1, strings.TrimSpace(tc.want), got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
