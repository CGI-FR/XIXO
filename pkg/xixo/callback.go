package xixo

import (
	"encoding/json"
	"strings"
)

type Callback func(*XMLElement) (*XMLElement, error)

type CallbackMap func(map[string]string) (map[string]string, error)

type CallbackJSON func(string) (string, error)

type Attributs struct {
	Name  string
	Value string
}

func XMLElementToMapCallback(callback CallbackMap) Callback {
	result := func(xmlElement *XMLElement) (*XMLElement, error) {
		dict := map[string]string{}
		for name, child := range xmlElement.Childs {
			dict[name] = child[0].InnerText
		}

		dict, err := callback(dict)
		if err != nil {
			return nil, err
		}

		children, err := xmlElement.SelectElements("child::*")
		if err != nil {
			return nil, err
		}

		AttributsList := exctratAttributs(dict, xmlElement.Name)

		for _, child := range children {

			if value, ok := dict[child.Name]; ok {
				child.InnerText = value
			}

			if attributes, ok := AttributsList[child.Name]; ok {
				for _, attr := range attributes {
					child.AddAttribut(attr.Name, attr.Value)
				}
			}

		}

		return xmlElement, nil
	}

	return result
}

func exctratAttributs(dict map[string]string, parentName string) map[string][]Attributs {
	AttributsList := make(map[string][]Attributs)
	// check dict[name] include "@"
	for key, value := range dict {
		parts := strings.SplitN(key, "@", 2)
		// if include, use split to get element:before@ ,attr:after@
		if len(parts) == 1 {
			tagName := parentName
			newAttribut := Attributs{Name: parts[0], Value: value}

			if existingElement, ok := AttributsList[tagName]; ok {
				existingElement = append(existingElement, newAttribut)
				AttributsList[tagName] = existingElement
			} else {
				AttributsList[tagName] = []Attributs{newAttribut}
			}
		} else if len(parts) == 2 {

			tagName := parts[0]
			newAttribut := Attributs{Name: parts[1], Value: value}

			// if key already in attributs
			if existingElement, ok := AttributsList[tagName]; ok {

				existingElement = append(existingElement, newAttribut)
				AttributsList[tagName] = existingElement
			} else {
				AttributsList[tagName] = []Attributs{newAttribut}
			}
		}
	}
	return AttributsList
}

func XMLElementToJSONCallback(callback CallbackJSON) Callback {
	resultCallback := func(dict map[string]string) (map[string]string, error) {
		source, err := json.Marshal(dict)
		if err != nil {
			return nil, err
		}

		dest, err := callback(string(source))
		if err != nil {
			return nil, err
		}

		result := map[string]string{}

		err = json.Unmarshal([]byte(dest), &result)
		if err != nil {
			return nil, err
		}

		return result, nil
	}

	return XMLElementToMapCallback(resultCallback)
}
