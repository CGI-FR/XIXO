package xixo

import "encoding/json"

type Callback func(*XMLElement) *XMLElement

type CallbackMap func(map[string]string) map[string]string

type CallbackJSON func(string) string

func XMLElementToMapCallback(callback CallbackMap) Callback {
	result := func(xmlElement *XMLElement) *XMLElement {
		dict := map[string]string{}

		for name, child := range xmlElement.Childs {
			dict[name] = child[0].InnerText
		}

		for name, value := range callback(dict) {
			xmlElement.Childs[name][0].InnerText = value
		}

		return xmlElement
	}

	return result
}

func XMLElementToJSONCallback(callback CallbackJSON) Callback {
	resultCallback := func(dict map[string]string) map[string]string {
		source, _ := json.Marshal(dict)

		dest := callback(string(source))

		result := map[string]string{}

		json.Unmarshal([]byte(dest), &result)

		return result
	}
	return XMLElementToMapCallback(resultCallback)
}
