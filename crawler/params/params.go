package params

import "log"

var ufs = make(map[string]string)
var combustivel = make(map[string]string)

func buildCombustivelParams() map[string]string {
	combustivel["Revendedor"] = "1"
	combustivel["Abastecimento"] = "2"
	combustivel["Escola"] = "3"
	combustivel["GNV"] = "4"
	combustivel["Flutuante"] = "5"
	combustivel["Aviação"] = "6"
	combustivel["Marítimo"] = "7"
	return combustivel
}

func buildUfParams() map[string]string {
	ufs["AC"] = "AC"
	ufs["AL"] = "AL"
	ufs["AM"] = "AM"
	ufs["AP"] = "AP"
	ufs["BA"] = "BA"
	ufs["CE"] = "CE"
	ufs["DF"] = "DF"
	ufs["ES"] = "ES"
	ufs["GO"] = "GO"
	ufs["MA"] = "MA"
	ufs["MG"] = "MG"
	ufs["MS"] = "MS"
	ufs["MT"] = "MT"
	ufs["PA"] = "PA"
	ufs["PB"] = "PB"
	ufs["PE"] = "PE"
	ufs["PI"] = "PI"
	ufs["PR"] = "PR"
	ufs["RJ"] = "RJ"
	ufs["RN"] = "RN"
	ufs["RO"] = "RO"
	ufs["RR"] = "RR"
	ufs["RS"] = "RS"
	ufs["SC"] = "SC"
	ufs["SE"] = "SE"
	ufs["SP"] = "SP"
	ufs["TO"] = "TO"
	return ufs
}

func Getufs(uf string) string {
	ufs := buildUfParams()
	getuf := ufs[uf]
	if getuf == "" {
		log.Fatal("Não existe na base a UF ", uf)
	}
	return getuf
}

func GetTipoPosto(combustivel string) string {
	codPosto := buildCombustivelParams()
	getCodigo := codPosto[combustivel]
	if getCodigo == "" {
		log.Fatal("Não existe na base o combustível ", combustivel)
	}
	return getCodigo
}
