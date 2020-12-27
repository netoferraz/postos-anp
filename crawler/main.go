package main

import (
	"entities"
	"fmt"
	"funcs"
	"log"
	"mongo"
	"params"
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
	PostoDetails := initialRequest.Clone()
	initialRequest.OnResponse(func(r *colly.Response) {
		if r.StatusCode == 200 {
			err := crawlPostos.Post(
				"https://postos.anp.gov.br/consulta.asp",
				map[string]string{
					"sCnpj":        "",
					"sRazaoSocial": "",
					"sEstado":      params.Getufs("AC"),
					"sMunicipio":   "0",
					"sBandeira":    "0",
					"sProduto":     "0",
					"sTipodePosto": params.GetTipoPosto("Revendedor"),
					"p":            "",
					"hPesquisar":   "PESQUISAR",
				})
			if err != nil {
				fmt.Println("Deu algum erro no POST.")
			}
		}
	})
	crawlPostos.OnResponse(func(r *colly.Response) {
		if r.StatusCode == 200 {
			fmt.Println("POST bem sucedido: ", r.StatusCode)
		} else {
			fmt.Println("Não foi bem sucedido o POST: ", r.StatusCode)
		}
	})
	crawlPostos.OnRequest(func(r *colly.Request) {
		funcs.SetHeaders(r)
	})
	crawlPostos.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	crawlPostos.OnHTML("table tr", func(e *colly.HTMLElement) {
		posto := entities.Posto{}
		e.ForEach("td", func(index int, el *colly.HTMLElement) {
			if index != 0 {
				switch index {
				case 1:
					posto.RazaoSocial = el.Text
				case 2:
					posto.NomeFantasia = el.Text
				case 3:
					posto.Uf = el.Text
				case 4:
					posto.Municipio = el.Text
				case 5:
					splitInfo := strings.Split(el.Text, "-")
					if len(splitInfo) == 2 {
						bandeira := splitInfo[0]
						data_inicio := splitInfo[1]
						posto.Bandeira = bandeira
						posto.DataInicio = data_inicio

					} else {
						bandeira := splitInfo[0]
						posto.Bandeira = bandeira
					}
				}
			} else {
				el.ForEach("input", func(index int, el *colly.HTMLElement) {
					switch index {
					case 0:
						text := el.Attr("value")
						isProximo := strings.Contains(text, "Próximo")
						if !isProximo {
							posto.CodInstalacao = text
						}
					case 1:
						cnpj := el.Attr("value")
						posto.CNPJ = cnpj
					}
				})
			}
		})
		if !(entities.Posto{} == posto) {
			PostoDetails.Post("https://postos.anp.gov.br/resultado.asp", map[string]string{"Cod_inst": posto.CodInstalacao, "estado": posto.Uf, "municipio": "0"})

		}
	})
	PostoDetails.OnRequest(func(r *colly.Request) {
		funcs.SetHeaders(r)
	})
	// Set error handler
	// VOLTAR AQUI DEPOIS
	//PostoDetails.OnError(func(r *colly.Response, err error) {
	//	fmt.Println("Request URL:", r.Request.URL, " body:", r.Request.Body, "Headers: ", r.Request.Headers, " failed with response:", r.StatusCode, "\nError:", err)
	//})
	PostoDetails.OnHTML("table", func(el *colly.HTMLElement) {
		var ContainerEquipamentos []entities.EquipamentosPosto
		height := el.Attr("height")
		width := el.Attr("width")
		var features []string
		var values []string
		if height == "530" && width == "634" {
			el.ForEach("table", func(_ int, elem *colly.HTMLElement) {
				width_table := elem.Attr("width")
				switch width_table {
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
		}
		if len(features) != 0 && len(values) != 0 {
			container := funcs.CollectDetails(features, values)
			container.Equipamentos = ContainerEquipamentos
			if len(container.Equipamentos) != 0 {
				mongo.CreateDocument(client, container)
			}
		}

	})
	initialRequest.Visit("https://postos.anp.gov.br/")
	initialRequest.Wait()
	crawlPostos.Wait()
	PostoDetails.Wait()
}
