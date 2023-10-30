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
		parentAttributs := extractAttributesParent(dict)
		if parentAttributs != nil {
			for _, attr := range parentAttributs {
				xmlElement.AddAttribut(attr.Name, attr.Value)
			}
		}
		children, err := xmlElement.SelectElements("child::*")
		if err != nil {
			return nil, err
		}

		AttributsList := exctratAttributsChild(dict)

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

func exctratAttributsChild(dict map[string]string) map[string][]Attributs {
	AttributsList := make(map[string][]Attributs)
	// check dict[name] include "@"
	for key, value := range dict {
		parts := strings.SplitN(key, "@", 2)
		// if include, use split to get element:before@ ,attr:after@
		if len(parts) == 2 {

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

func extractAttributesParent(dict map[string]string) []Attributs {
	AttributesMap := []Attributs{}
	for key, value := range dict {
		if strings.HasPrefix(key, "@") {
			attributeKey := key[1:]
			attribute := Attributs{Name: attributeKey, Value: value}
			AttributesMap = append(AttributesMap, attribute)
		}
	}

	return AttributesMap
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
