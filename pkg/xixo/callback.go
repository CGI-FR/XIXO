package xixo

import (
	"encoding/json"
	"regexp"
)

type Callback func(*XMLElement) (*XMLElement, error)

type CallbackMap func(map[string]string) (map[string]string, error)

type CallbackJSON func(string) (string, error)

type Attributs struct {
	Attr    string
	AttrVal string
}

func XMLElementToMapCallback(callback CallbackMap) Callback {
	result := func(xmlElement *XMLElement) (*XMLElement, error) {
		dict := map[string]string{}
		var AttributsList map[string][]Attributs
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

		// check dict[name] include "@"
		re := regexp.MustCompile("@")
		for key, value := range dict {
			if re.MatchString(key) {
				// if include, use regexp to get element:before@ ,attr:after@
				parts := re.Split(key, 2)
				tagName := parts[0]
				newAttribut := Attributs{Attr: parts[1], AttrVal: value}
				if existingElement, ok := AttributsList[tagName]; ok {
					// if key already in attributs
					existingElement = append(existingElement, newAttribut)
					AttributsList[tagName] = existingElement
				} else {
					// if key not in attributs yet
					AttributsList[tagName] = []Attributs{newAttribut}
				}
			}
		}

		for _, child := range children {
			// if attrlist, ok := AttributsList[child.Name]; ok {
			// 	if child.Attrs == nil {

			// 	}
			// 	child.InnerText = value
			// }
			// if xmlElement.Attrs == nil

			// creat one
			// apprend {} to Attrs
			if value, ok := dict[child.Name]; ok {
				child.InnerText = value
			}
		}

		return xmlElement, nil
	}

	return result
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
