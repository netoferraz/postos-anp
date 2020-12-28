package main

import (
	"entities"
	"flag"
	"fmt"
	"funcs"
	"log"
	"mongo"
	"params"
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
	var Uf string
	var categoriaPosto string
	flag.StringVar(&Uf, "UF", "", "Sigla da Unidade da Federação a ser coletada")
	flag.StringVar(&categoriaPosto, "tipoPosto", "", "Tipo do posto.")
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
	PostoDetails.Limit(&colly.LimitRule{
		Parallelism: 2,
		Delay:       5 * time.Second,
	})
	initialRequest.OnResponse(func(r *colly.Response) {
		if r.StatusCode == 200 {
			err := crawlPostos.Post(
				"https://postos.anp.gov.br/consulta.asp",
				map[string]string{
					"sCnpj":        "",
					"sRazaoSocial": "",
					"sEstado":      params.Getufs(Uf),
					"sMunicipio":   "0",
					"sBandeira":    "0",
					"sProduto":     "0",
					"sTipodePosto": params.GetTipoPosto(categoriaPosto),
					"p":            "",
					"hPesquisar":   "PESQUISAR",
				})
			if err != nil {
				fmt.Println("Deu algum erro no POST.")
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
			PostoDetails.Post("https://postos.anp.gov.br/resultado.asp", map[string]string{"Cod_inst": CodInstalacao, "estado": params.Getufs(Uf), "municipio": "0"})
		}
	})
	PostoDetails.OnRequest(func(r *colly.Request) {
		funcs.SetHeaders(r)
	})
	PostoDetails.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, " body:", r.Request.Body, " failed with response:", r.StatusCode, "\nError:", err)
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
			mongo.ReplaceDocument(client, "cnpj", container.CNPJ, container)
		}

	})
	initialRequest.Visit("https://postos.anp.gov.br/")
	initialRequest.Wait()
	crawlPostos.Wait()
}
