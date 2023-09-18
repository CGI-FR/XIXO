package xixo

import "encoding/json"

type Callback func(*XMLElement) (*XMLElement, error)

type CallbackMap func(map[string]string) (map[string]string, error)

type CallbackJSON func(string) (string, error)

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

		for _, child := range children {
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
