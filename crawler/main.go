package main

import (
	"entities"
	"flag"
	"fmt"
	"funcs"
	"log"
	"mongo"
	"os"
	"params"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

func main() {
	//init database
	client, err := mongo.GetMongoClient()
	if err != nil {
		log.Fatal(err)
	}
	var SelectedUf string
	var SelectedCategoriaPosto string
	var containerCnpj []string
	var collectionNameSucess string = os.Getenv("MONGO_COLLECTION")
	if collectionNameSucess == "" {
		log.Fatal("É necessário configurar a variável de ambiente MONGO_COLLECTION.")
	}
	var collectionNameFail string = os.Getenv("MONGO_COLLECTION_ERROR")
	if collectionNameFail == "" {
		log.Fatal("É necessário configurar a variável de ambiente MONGO_COLLECTION_ERROR.")
	}
	flag.StringVar(&SelectedUf, "UF", "", "Sigla da Unidade da Federação a ser coletada")
	flag.StringVar(&SelectedCategoriaPosto, "tipoPosto", "", "Tipo do posto.")
	flag.Parse()
	initialRequest := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.80 Safari/537.36"),
		colly.Async(true),
		colly.DetectCharset(),
	)
	initialRequest.Limit(&colly.LimitRule{
		Parallelism: 2,
		Delay:       5 * time.Second,
	})
	crawlPostos := initialRequest.Clone()
	PostoDetails := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.80 Safari/537.36"),
		colly.DetectCharset(),
	)
	PostoDetails.SetRequestTimeout(20 * time.Second)
	PostoDetails.Limit(&colly.LimitRule{
		Parallelism: 2,
		Delay:       5 * time.Second,
	})
	initialRequest.OnResponse(func(r *colly.Response) {
		if r.StatusCode == 200 {
			Uf, ok := params.Getufs(SelectedUf)
			if !ok {
				log.Fatal("É necessário atribuir um valor para o parâmetro -UF.")
			}
			categoriaPosto, ok := params.GetTipoPosto(SelectedCategoriaPosto)
			if !ok {
				log.Fatal("É necessário atribuir um valor para o parâmetro -tipoPosto.")
			}
			if !strings.EqualFold(Uf, "All") && !strings.EqualFold(categoriaPosto, "All") {
				err := crawlPostos.Post(
					"https://postos.anp.gov.br/consulta.asp",
					map[string]string{
						"sCnpj":        "",
						"sRazaoSocial": "",
						"sEstado":      Uf,
						"sMunicipio":   "0",
						"sBandeira":    "0",
						"sProduto":     "0",
						"sTipodePosto": categoriaPosto,
						"p":            "",
						"hPesquisar":   "PESQUISAR",
					})
				if err != nil {
					log.Println(err)
				}

			} else if strings.EqualFold(Uf, "All") && !strings.EqualFold(categoriaPosto, "All") {
				ufs := params.BuildAllUfs()
				for _, uf := range ufs {
					err := crawlPostos.Post(
						"https://postos.anp.gov.br/consulta.asp",
						map[string]string{
							"sCnpj":        "",
							"sRazaoSocial": "",
							"sEstado":      uf,
							"sMunicipio":   "0",
							"sBandeira":    "0",
							"sProduto":     "0",
							"sTipodePosto": categoriaPosto,
							"p":            "",
							"hPesquisar":   "PESQUISAR",
						})
					if err != nil {
						log.Println(err)
					}
				}
			} else if !strings.EqualFold(Uf, "All") && strings.EqualFold(categoriaPosto, "All") {
				allCatPosto := params.BuildAllTipoPosto()
				for _, catPosto := range allCatPosto {
					err := crawlPostos.Post(
						"https://postos.anp.gov.br/consulta.asp",
						map[string]string{
							"sCnpj":        "",
							"sRazaoSocial": "",
							"sEstado":      Uf,
							"sMunicipio":   "0",
							"sBandeira":    "0",
							"sProduto":     "0",
							"sTipodePosto": catPosto,
							"p":            "",
							"hPesquisar":   "PESQUISAR",
						})
					if err != nil {
						log.Println(err)
					}

				}

			} else {
				ufs := params.BuildAllUfs()
				allCatPosto := params.BuildAllTipoPosto()
				for _, uf := range ufs {
					for _, catPosto := range allCatPosto {
						fmt.Println("Iniciando consulta para [UF]:", uf, "& [tipoPosto]", catPosto)
						err := crawlPostos.Post(
							"https://postos.anp.gov.br/consulta.asp",
							map[string]string{
								"sCnpj":        "",
								"sRazaoSocial": "",
								"sEstado":      uf,
								"sMunicipio":   "0",
								"sBandeira":    "0",
								"sProduto":     "0",
								"sTipodePosto": catPosto,
								"p":            "",
								"hPesquisar":   "PESQUISAR",
							})
						if err != nil {
							fmt.Println("Deu algum erro no POST.")
						}
					}
				}
			}

		}
	})
	crawlPostos.OnResponse(func(r *colly.Response) {
		if r.StatusCode != 200 {
			fmt.Println("Não foi bem sucedido o POST para", r.Request.URL, "com status code: ", r.StatusCode)
		}
	})
	crawlPostos.OnRequest(func(r *colly.Request) {
		funcs.SetHeaders(r)
	})
	crawlPostos.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
		r.Request.Retry()
	})
	crawlPostos.OnHTML("table tr", func(e *colly.HTMLElement) {
		var CodInstalacao string
		e.ForEach("td", func(index int, el *colly.HTMLElement) {
			if index == 0 {
				el.ForEach("input", func(index int, el *colly.HTMLElement) {
					if index == 0 {
						text := el.Attr("value")
						isProximo := strings.Contains(text, "Próximo")
						if !isProximo {
							if _, err := strconv.ParseInt(text, 10, 64); err == nil {
								CodInstalacao = text
							}
						}
					}
				})
			}
		})
		if CodInstalacao != "" {
			bodyString := fmt.Sprintf("%v", e.Request.Body)
			uf := params.GetParamsValue(bodyString, "UF")
			PostoDetails.Post("https://postos.anp.gov.br/resultado.asp", map[string]string{"Cod_inst": CodInstalacao, "estado": uf, "municipio": "0"})
		}
	})
	crawlPostos.OnHTML("form", func(elem *colly.HTMLElement) {
		formName := elem.Attr("name")
		re := regexp.MustCompile("[0-9]+")
		if strings.EqualFold(formName, "FormNext") {
			elem.ForEach("input", func(_ int, webelement *colly.HTMLElement) {
				inputName := webelement.Attr("value")
				if strings.Contains(inputName, "Próximo") {
					rawValue := webelement.Attr("onclick")
					findpag := re.FindAllString(rawValue, -1)
					if len(findpag) != 0 {
						bodyString := fmt.Sprintf("%v", elem.Request.Body)
						uf := params.GetParamsValue(bodyString, "UF")
						tipoPosto := params.GetParamsValue(bodyString, "tipoPosto")
						pag := findpag[0]
						err := crawlPostos.Post(
							"https://postos.anp.gov.br/consulta.asp",
							map[string]string{
								"sCnpj":        "",
								"sRazaoSocial": "",
								"sEstado":      uf,
								"sMunicipio":   "0",
								"sBandeira":    "0",
								"sProduto":     "0",
								"sTipodePosto": tipoPosto,
								"p":            pag,
								"hPesquisar":   "PESQUISAR",
							})
						if err != nil {
							fmt.Println("Deu algum erro no POST da página", pag)
						}

					}
				}
			})
		}
	})
	PostoDetails.OnRequest(func(r *colly.Request) {
		funcs.SetHeaders(r)
	})
	PostoDetails.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, " body:", r.Request.Body, " failed with response:", r.StatusCode, "\nError:", err)
		bodyString := fmt.Sprintf("%v", r.Request.Body)
		var payLoadError = make(map[string]string)
		reParams := regexp.MustCompile(`\w+=\w+`)
		getParams := reParams.FindAllString(bodyString, -1)
		for _, params := range getParams {
			splitParams := strings.Split(params, "=")
			payLoadError[splitParams[0]] = splitParams[1]
		}
		errString := fmt.Sprintf("%v", err)
		payLoadError["message"] = string(errString)
		stringParams := fmt.Sprintf("{\"Cod_inst\": \"%v\", \"estado\": \"%v\", \"municipio\": \"%v\"}", payLoadError["Cod_inst"], payLoadError["estado"], payLoadError["municipio"])
		generateHash := funcs.GetMD5Hash(stringParams)
		payLoadError["hash_id"] = generateHash
		mongo.ReplaceDocument(client, collectionNameFail, "hash_id", payLoadError["hash_id"], payLoadError)
		r.Request.Retry()
	})
	PostoDetails.OnHTML("table", func(el *colly.HTMLElement) {
		var ContainerEquipamentos []entities.EquipamentosPosto
		height := el.Attr("height")
		width := el.Attr("width")
		var features []string
		var values []string
		if height == "530" && width == "634" {
			el.ForEach("table", func(_ int, elem *colly.HTMLElement) {
				widthTable := elem.Attr("width")
				switch widthTable {
				case "634":
					elem.ForEach("tr td", func(_ int, line *colly.HTMLElement) {
						line.ForEach("font", func(_ int, fontline *colly.HTMLElement) {
							statusPosto := fontline.Attr("size")
							if statusPosto == "3" {
								features = append(features, "StatusPosto")
								values = append(values, fontline.Text)
							}
						})
						identifyFeature := line.Attr("align")
						switch identifyFeature {
						case "right":
							line.ForEach("b", func(_ int, _ *colly.HTMLElement) {
								if line.Text != "Nova Consulta" {
									features = append(features, line.Text)
								}

							})
						case "left":
							line.ForEach("font", func(_ int, _ *colly.HTMLElement) {
								values = append(values, line.Text)

							})
						}
					})
				case "644":
					var equipamentos entities.EquipamentosPosto
					elem.ForEach("tr td", func(_ int, line *colly.HTMLElement) {
						identifyFeature := line.Attr("width")
						switch identifyFeature {
						case "450":
							line.ForEach("font", func(_ int, font *colly.HTMLElement) {
								if font.Text != "Produtos:" && font.Text != "Tancagem (m³):" && font.Text != "Bicos:" {
									equipamentos.Produto = font.Text
								}
							})
						case "103":
							line.ForEach("font", func(_ int, font *colly.HTMLElement) {
								if font.Text != "Produtos:" && font.Text != "Tancagem (m³):" && font.Text != "Bicos:" {
									equipamentos.Tancagem = font.Text
								}
							})
						case "92":
							line.ForEach("font", func(_ int, font *colly.HTMLElement) {
								if font.Text != "Produtos:" && font.Text != "Tancagem (m³):" && font.Text != "Bicos:" {
									equipamentos.Bicos = font.Text
								}
							})
						}
						if equipamentos.Produto != "" && equipamentos.Tancagem != "" && equipamentos.Bicos != "" {
							ContainerEquipamentos = append(ContainerEquipamentos, equipamentos)
							equipamentos = entities.EquipamentosPosto{}
						}
					})
				}

			})
		} else if width == "760" {
			el.ForEach("table", func(_ int, elem *colly.HTMLElement) {
				widthTable := elem.Attr("width")
				heightTable := elem.Attr("height")
				if heightTable != "530" {
					switch widthTable {
					case "634":
						elem.ForEach("tr td", func(_ int, line *colly.HTMLElement) {
							line.ForEach("font", func(_ int, fontline *colly.HTMLElement) {
								statusPosto := fontline.Attr("size")
								if statusPosto == "3" {
									features = append(features, "StatusPosto")
									values = append(values, fontline.Text)
								}
							})
							identifyFeature := line.Attr("align")
							switch identifyFeature {
							case "right":
								line.ForEach("b", func(_ int, newline *colly.HTMLElement) {
									if line.Text != "Nova Consulta" {
										features = append(features, newline.Text)
									}

								})
							case "left":
								line.ForEach("font", func(_ int, newline *colly.HTMLElement) {
									values = append(values, newline.Text)

								})
							}
						})
					}
				}
			})
		}
		if len(features) != 0 && len(values) != 0 {
			container := funcs.CollectDetails(features, values)
			container.Equipamentos = ContainerEquipamentos
			isCnpj := funcs.Contains(containerCnpj, container.CNPJ)
			if !isCnpj {
				containerCnpj = append(containerCnpj, container.CNPJ)
				fmt.Println("Coletando os dados do CNPJ", container.CNPJ)
				mongo.ReplaceDocument(client, collectionNameSucess, "cnpj", container.CNPJ, container)
			}
		}

	})
	initialRequest.Visit("https://postos.anp.gov.br/")
	initialRequest.Wait()
	crawlPostos.Wait()
}
