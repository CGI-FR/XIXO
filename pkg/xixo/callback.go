package xixo

import (
	"encoding/json"
	"strings"
)

type Callback func(*XMLElement) (*XMLElement, error)

type CallbackMap func(map[string]string) (map[string]string, error)

type CallbackJSON func(string) (string, error)

// XMLElementToMapCallback transforms an XML element into a map, applies a callback function,
// adds parent attributes, and updates child elements.
func XMLElementToMapCallback(callback CallbackMap) Callback {
	result := func(xmlElement *XMLElement) (*XMLElement, error) {
		dict := map[string]string{}
		for name, child := range xmlElement.Childs {
			dict[name] = child[0].InnerText
		}

		extractExistedAttributes(xmlElement, dict)

		dict, err := callback(dict)
		if err != nil {
			return nil, err
		}

		if len(xmlElement.Childs) > 0 {
			existingChilds := make(map[string]bool)
			for key := range dict {
				existingChilds[key] = true
			}
			// Check if the attribute is already present
			for name := range xmlElement.Childs {
				if !existingChilds[name] {
					xmlElement.RemoveChild(name)
				}
			}
		}

		// Extract parent attributes and add them to the XML element.
		parentAttributes := extractParentAttributes(dict)
		for _, attr := range parentAttributes {
			xmlElement.AddAttribute(attr)
		}
		// Apply remove on parentAttributes
		removeAttributes(parentAttributes, xmlElement)
		children := xmlElement.childs

		// Select child elements and update their text content and attributes.
		childAttributes := extractChildAttributes(dict)

		for _, child := range children {
			if value, ok := dict[child.Name]; ok {
				child.InnerText = value
			}

			attributes, ok := childAttributes[child.Name]
			if ok {
				// Add new attributes
				for _, attr := range attributes {
					child.AddAttribute(attr)
				}
			}

			// Apply remove on child attributes
			removeAttributes(attributes, child)
		}

		return xmlElement, nil
	}

	return result
}

func removeAttributes(attributes []Attribute, element *XMLElement) {
	// Check if attributes are available for the current child
	existingAttributes := make(map[string]bool)
	for _, existAttribute := range attributes {
		existingAttributes[existAttribute.Name] = true
	}
	// Check if the attribute is already present
	for _, xmlAttributeName := range element.AttrKeys {
		if !existingAttributes[xmlAttributeName] {
			element.RemoveAttribute(xmlAttributeName)
		}
	}
}

func extractExistedAttributes(xmlElement *XMLElement, dict map[string]string) {
	for name, child := range xmlElement.Childs {
		for attrName, attr := range child[0].Attrs {
			dict[name+"@"+attrName] = attr.Value
		}
	}

	for attrName, attr := range xmlElement.Attrs {
		dict["@"+attrName] = attr.Value
	}
}

// extractChildAttributes extracts child attributes from the dictionary.
func extractChildAttributes(dict map[string]string) map[string][]Attribute {
	childAttributes := make(map[string][]Attribute)
	// check dict[name] include "@"
	for key, value := range dict {
		parts := strings.SplitN(key, "@", 2)

		if len(parts) == 2 {
			tagName := parts[0]
			newAttribut := Attribute{Name: parts[1], Value: value}
			// if key already in attributes
			if existingElement, ok := childAttributes[tagName]; ok {
				childAttributes[tagName] = append(existingElement, newAttribut)
			} else {
				childAttributes[tagName] = []Attribute{newAttribut}
			}
		}
	}

	return childAttributes
}

// extractParentAttributes extracts parent attributes from the dictionary.
func extractParentAttributes(dict map[string]string) []Attribute {
	parentAttributes := []Attribute{}

	for key, value := range dict {
		if strings.HasPrefix(key, "@") {
			attributeKey := key[1:]
			attribute := Attribute{Name: attributeKey, Value: value}
			parentAttributes = append(parentAttributes, attribute)
		}
	}

	return parentAttributes
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
