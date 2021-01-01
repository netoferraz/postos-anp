package entities

import "time"

type Posto struct {
	CNPJ          string
	RazaoSocial   string
	NomeFantasia  string
	Uf            string
	Municipio     string
	Bandeira      string
	DataInicio    string
	CodInstalacao string
}

type EquipamentosPosto struct {
	Produto  string
	Tancagem string
	Bicos    string
}

type DetailsPosto struct {
	Autorizacao       string
	CNPJ              string
	RazaoSocial       string
	NomeFantasia      string
	Endereco          string
	Complemento       string
	Bairro            string
	MunicipioUf       string
	CEP               string
	NumeroDespacho    string
	DataPublicacao    string
	BandeiraInicio    string
	TipodoPosto       string
	StatusPosto       string
	Socios            []string
	Equipamentos      []EquipamentosPosto
	DatetimeCollected time.Time
}

type Pessoajuridica interface {
	getcnpj() string
}

func (p Posto) getcnpj() string {
	return p.CNPJ
}

func (p DetailsPosto) getcnpj() string {
	return p.CNPJ
}
