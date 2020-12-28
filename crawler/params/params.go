package params

import (
	"log"
	"strings"
)

var ufs = make(map[string]string)
var tipoPosto = make(map[string]string)

func buildCombustivelParams() map[string]string {
	tipoPosto["Revendedor"] = "1"
	tipoPosto["Abastecimento"] = "2"
	tipoPosto["Escola"] = "3"
	tipoPosto["GNV"] = "4"
	tipoPosto["Flutuante"] = "5"
	tipoPosto["Aviação"] = "6"
	tipoPosto["Marítimo"] = "7"
	return tipoPosto
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
	uf = strings.ToUpper(uf)
	getuf := ufs[uf]
	if getuf == "" {
		log.Fatal("Não existe na base a UF ", uf)
	}
	return getuf
}

func GetTipoPosto(categoria string) string {
	codPosto := buildCombustivelParams()
	getCodigo := codPosto[categoria]
	if getCodigo == "" {
		log.Fatal("Os parâmetros aceitos são: Revendedor, Abastecimento, Escola, GNV, Flutuante, Aviação, Marítimo.")
	}
	return getCodigo
}
