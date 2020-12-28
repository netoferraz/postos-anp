# Crawler para coleta dos dados cadastrais de Postos de Combustíveis

**site**: https://postos.anp.gov.br/

### 1. Parametrização da execução do crawler.
A coleta pode ser parametrizada por unidade da federação (`UF`) e por tipo de posto (`tipoPosto`).

**UF:** [ **`AC`**, **`AL`**, **`AM`**, **`AP`**, **`BA`**, **`CE`**, **`DF`**, **`ES`**, **`GO`**, **`MA`**, **`MG`**, **`MS`**, **`MT`**, **`PA`**, **`PB`**, **`PE`**, **`PI`**, **`PR`**, **`RJ`**, **`RN`**, **`RO`**, **`RR`**, **`RS`**,
**`SC`**, **`SE`**, **`SP`**, **`TO`**]

**Tipo de Posto:** [**`Revendedor`**, **`Abastecimento`**, **`Escola`**, **`GNV`**, **`Flutuante`**, **`Aviação`**, **`Marítimo`**]

### 2. Exemplos de execução

#### 2.1 Coletar os dados dos postos do tipo Revendedor do estado de São Paulo

`go run ./main.go -UF SP -tipoPosto Revendedor`