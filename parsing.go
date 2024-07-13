package main

import (
	"encoding/json"
	"encoding/xml"
	"os"
)

type Assets struct {
	XMLName    xml.Name `xml:"biblioteche"`
	Text       string   `xml:",chardata"`
	DataExport string   `xml:"data-export,attr"`
	Biblioteca []struct {
		Text          string `xml:",chardata"`
		CodiceIsil    string `xml:"codice-isil,attr"`
		Denominazione string `xml:"denominazione,attr"`
		Materiale     []struct {
			Text      string `xml:",chardata" json:"asset"`
			Categoria string `xml:"categoria,attr" json:"type"`
			Posseduto string `xml:"posseduto,attr" json:"items,omitempty"`
		} `xml:"materiale"`
	} `xml:"biblioteca"`
}

func ParseAssets(source string) (map[string]string, error) {
	xmlfile, err := os.ReadFile(source)
	if err != nil {
		return nil, err
	}
	patrimonio := Assets{}
	err = xml.Unmarshal(xmlfile, &patrimonio)
	if err != nil {
		return nil, err
	}

	assets := make(map[string]string)
	for _, biblioteca := range patrimonio.Biblioteca {
		j, _ := json.Marshal(biblioteca.Materiale)
		assets[biblioteca.CodiceIsil] = string(j)
	}
	return assets, nil
}

type Fonds struct {
	XMLName    xml.Name `xml:"biblioteche"`
	Text       string   `xml:",chardata"`
	DataExport string   `xml:"data-export,attr"`
	Biblioteca []struct {
		Text          string `xml:",chardata"`
		CodiceIsil    string `xml:"codice-isil,attr"`
		Denominazione string `xml:"denominazione,attr"`
		FondoSpeciale []struct {
			Text          string `xml:",chardata" json:"-"`
			Denominazione string `xml:"denominazione" json:"name"`
			Dewey         struct {
				Text   string `xml:",chardata" json:"label,omitempty"`
				Codice string `xml:"codice,attr" json:"code,omitempty"`
			} `xml:"dewey" json:"dewey"`
			Descrizione string `xml:"descrizione" json:"description"`
		} `xml:"fondo-speciale"`
	} `xml:"biblioteca"`
}

func ParseFonds(source string) (map[string]string, error) {
	xmlfile, err := os.ReadFile(source)
	if err != nil {
		return nil, err
	}
	fondi := Fonds{}
	err = xml.Unmarshal(xmlfile, &fondi)
	if err != nil {
		return nil, err
	}

	fonds := make(map[string]string)
	for _, biblioteca := range fondi.Biblioteca {
		j, _ := json.Marshal(biblioteca.FondoSpeciale)
		fonds[biblioteca.CodiceIsil] = string(j)
	}
	return fonds, nil
}
