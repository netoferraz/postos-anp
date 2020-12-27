package funcs

import (
	"entities"
	"strings"

	"github.com/gocolly/colly/v2"
)

func unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

func SetHeaders(request *colly.Request) *colly.Request {
	request.Headers.Set("Connection", "keep-alive")
	request.Headers.Set("Cache-Control", "max-age=0")
	request.Headers.Set("Upgrade-Insecure-Requests", "1")
	request.Headers.Set("Origin", "https://postos.anp.gov.br")
	request.Headers.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")
	request.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	request.Headers.Set("Sec-Fetch-Site", "same-origin")
	request.Headers.Set("Sec-Fetch-Mode", "navigate")
	request.Headers.Set("Sec-Fetch-User", "?1")
	request.Headers.Set("Sec-Fetch-Dest", "document")
	request.Headers.Set("Referer", "https://postos.anp.gov.br/consulta.asp")
	request.Headers.Set("Accept-Language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7")
	return request
}

func CollectDetails(propriedades []string, valores []string) entities.DetailsPosto {

	var container entities.DetailsPosto
	if len(propriedades) != len(valores) {
		for i := range propriedades {
			switch propriedades[i] {
			case "Autorização:":
				container.Autorizacao = valores[:len(propriedades)][i]
			case "CNPJ/CPF:":
				container.CNPJ = valores[:len(propriedades)][i]
			case "Razão Social:":
				container.RazaoSocial = valores[:len(propriedades)][i]
			case "Nome Fantasia:":
				container.NomeFantasia = valores[:len(propriedades)][i]
			case "Endereço:":
				container.Endereço = valores[:len(propriedades)][i]
			case "Complemento:":
				container.Complemento = valores[:len(propriedades)][i]
			case "Bairro:":
				container.Bairro = valores[:len(propriedades)][i]
			case "Município/UF:":
				container.Municipio_Uf = valores[:len(propriedades)][i]
			case "CEP:":
				container.CEP = valores[:len(propriedades)][i]
			case "Número Despacho:":
				container.NumeroDespacho = valores[:len(propriedades)][i]
			case "Data Publicação:":
				container.DataPublicaçao = valores[:len(propriedades)][i]
			case "Bandeira/Início:":
				container.Bandeira_Inicio = valores[:len(propriedades)][i]
			case "Tipo do Posto:":
				container.TipodoPosto = valores[:len(propriedades)][i]
			case "Sócios:":
				var uniqueSocios []string
				cleanedSocios := strings.TrimSpace(valores[:len(propriedades)][i])
				socios := strings.Split(cleanedSocios, "            ")
				uniqueSocios = unique(socios)
				for _, socio := range uniqueSocios {
					isExists := contains(container.Socios, socio)
					if !isExists {
						container.Socios = append(container.Socios, socio)
					}
				}
			}
		}
		var uniqueSocios []string
		for _, value := range valores[len(propriedades):] {
			cleanedSocios := strings.TrimSpace(value)
			socios := strings.Split(cleanedSocios, "            ")
			uniqueSocios = unique(socios)
			for _, socio := range uniqueSocios {
				isExists := contains(container.Socios, socio)
				if !isExists {
					container.Socios = append(container.Socios, socio)
				}
			}
		}
	} else {
		for i := range propriedades {
			switch propriedades[i] {
			case "Autorização:":
				container.Autorizacao = valores[i]
			case "CNPJ/CPF:":
				container.CNPJ = valores[i]
			case "Razão Social:":
				container.RazaoSocial = valores[i]
			case "Nome Fantasia:":
				container.NomeFantasia = valores[i]
			case "Endereço:":
				container.Endereço = valores[i]
			case "Complemento:":
				container.Complemento = valores[i]
			case "Bairro:":
				container.Bairro = valores[i]
			case "Município/UF:":
				container.Municipio_Uf = valores[i]
			case "CEP:":
				container.CEP = valores[i]
			case "Número Despacho:":
				container.NumeroDespacho = valores[i]
			case "Data Publicação:":
				container.DataPublicaçao = valores[i]
			case "Bandeira/Início:":
				container.Bandeira_Inicio = valores[i]
			case "Tipo do Posto:":
				container.TipodoPosto = valores[i]
			case "Socios":
				container.Socios = append(container.Socios, valores[i])
			}
		}
	}
	return container
}
