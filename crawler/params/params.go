package params

import (
	"log"
	"regexp"
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
	tipoPosto["All"] = "All"
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
	ufs["ALL"] = "ALL"
	return ufs
}

//BuildAllUfs get All instances of Uf
func BuildAllUfs() []string {
	return []string{"AC", "AL", "AM", "AP", "BA", "CE", "DF",
		"ES", "GO", "MA", "MG", "MS", "MT", "PA", "PB",
		"PE", "PI", "PR", "RJ", "RN", "RO", "RR", "RS", "SC",
		"SE", "SP", "TO"}
}

//BuildAllTipoPosto get an instance of tipoPosto
func BuildAllTipoPosto() []string {
	return []string{"1", "2", "3", "4", "5", "6", "7"}
}

//Getufs returns one instance of Uf
func Getufs(uf string) (string, bool) {
	ufs := buildUfParams()
	uf = strings.ToUpper(uf)
	getuf, ok := ufs[uf]
	if !ok {
		log.Fatal("Não existe na base a UF ", uf)
	}
	return getuf, ok
}

//GetTipoPosto returns an instance of Uf
func GetTipoPosto(categoria string) (string, bool) {
	codPosto := buildCombustivelParams()
	getCodigo, ok := codPosto[categoria]
	if !ok {
		log.Fatal("Os parâmetros aceitos são: Revendedor, Abastecimento, Escola, GNV, Flutuante, Aviação, Marítimo.")
	}
	return getCodigo, ok
}

//GetParamsValue get a value from body parameter
func GetParamsValue(bodyString string, param string) string {
	if !strings.EqualFold(param, "UF") && !strings.EqualFold(param, "tipoPosto") {
		log.Fatal("O parâmetro param somente aceitam os valores UF e tipoPosto")
	}
	var payLoad = make(map[string]string)
	reParams := regexp.MustCompile(`\w+=\w+`)
	getParams := reParams.FindAllString(bodyString, -1)
	for _, params := range getParams {
		splitParams := strings.Split(params, "=")
		payLoad[splitParams[0]] = splitParams[1]
	}
	if strings.EqualFold(param, "UF") {
		keyParam := "sEstado"
		value, ok := payLoad[keyParam]
		if !ok {
			log.Fatal("Não foi possível encontrar o parâmetro", keyParam, "no body da requisição.")
		}
		return value
	}
	keyParam := "sTipodePosto"
	value, ok := payLoad[keyParam]
	if !ok {
		log.Fatal("Não foi possível encontrar o parâmetro", keyParam, "no body da requisição.")
	}
	return value
}
