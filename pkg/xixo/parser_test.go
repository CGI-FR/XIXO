package xixo_test

import (
	"bytes"
	"testing"

	"github.com/CGI-FR/xixo/pkg/xixo"
	"github.com/stretchr/testify/assert"
)

// TestCopyXMLWithoutCallback vérifie que le parser copie simplement le fichier XML sans callback.
func TestCopyXMLWithoutCallback(t *testing.T) {
	t.Parallel()
	// Fichier XML en entrée
	inputXML := `
	<root>
		<element1>Contenu1</element1>
		<element2>Contenu2</element2>
	</root>`

	// Lisez les résultats du canal et construisez le XML résultant
	var resultXMLBuffer bytes.Buffer

	// Créez un bufio.Reader à partir du XML en entrée
	reader := bytes.NewBufferString(inputXML)

	// Créez une nouvelle instance du parser XML sans enregistrer de fonction de rappel
	parser := xixo.NewXMLParser(reader, &resultXMLBuffer)

	// Créez un canal pour collecter les résultats du parser
	err := parser.Stream()
	assert.Nil(t, err)

	// Vérifiez si le résultat XML est identique à l'entrée
	resultXML := resultXMLBuffer.String()

	if resultXML != inputXML {
		t.Errorf("Le résultat XML ne correspond pas à l'entrée.\nEntrée:\n%s\nSortie:\n%s", inputXML, resultXML)
	}
}

// TestModifyElement1ContentWithCallback vérifie que la fonction de rappel modifie correctement les nœuds <element1>.
func TestModifyElement1ContentWithCallback(t *testing.T) {
	t.Parallel()
	// Fichier XML en entrée
	inputXML := `
	<root>
		<element1>Hello <name>world</name> !</element1>
		<element2>Contenu2 <name> </name> ! </element2>
	</root>`

	// Lisez les résultats du canal et construisez le XML résultant
	var resultXMLBuffer bytes.Buffer

	// Créez un bufio.Reader à partir du XML en entrée
	reader := bytes.NewBufferString(inputXML)

	// Créez une nouvelle instance du parser XML avec la fonction de rappel
	parser := xixo.NewXMLParser(reader, &resultXMLBuffer)
	parser.RegisterCallback("element1", modifyElement1Content)

	// Créez un canal pour collecter les résultats du parser
	err := parser.Stream()
	assert.Nil(t, err)

	// Résultat XML attendu avec le contenu modifié
	expectedResultXML := `
	<root>
		<element1>ContenuModifie</element1>
		<element2>Contenu2 <name> </name> ! </element2>
	</root>`

	// Vérifiez si le résultat XML correspond à l'attendu
	resultXML := resultXMLBuffer.String()
	assert.Equal(t, expectedResultXML, resultXML)
}

// Callback pour modifier le contenu des nœuds <element1>.
func modifyElement1Content(elem *xixo.XMLElement) (*xixo.XMLElement, error) {
	elem.InnerText = "ContenuModifie"

	return elem, nil
}

// TestModifyElement1ContentWithCallback vérifie que la fonction de rappel modifie correctement les nœuds <element1>.
func TestModifyElementWrappedWithTextWithCallback(t *testing.T) {
	t.Parallel()
	// Fichier XML en entrée
	inputXML := `
	<root>
		<element1>Hello <name>world</name> !</element1>
		<element2>Contenu2 <name> </name> ! </element2>
	</root>`

	// Lisez les résultats du canal et construisez le XML résultant
	var resultXMLBuffer bytes.Buffer

	// Créez un bufio.Reader à partir du XML en entrée
	reader := bytes.NewBufferString(inputXML)

	// Créez une nouvelle instance du parser XML avec la fonction de rappel
	parser := xixo.NewXMLParser(reader, &resultXMLBuffer)
	parser.RegisterCallback("name", modifyElement1Content)

	// Créez un canal pour collecter les résultats du parser
	err := parser.Stream()
	assert.Nil(t, err)

	// Résultat XML attendu avec le contenu modifié
	expectedResultXML := `
	<root>
		<element1>Hello <name>ContenuModifie</name> !</element1>
		<element2>Contenu2 <name>ContenuModifie</name> ! </element2>
	</root>`

	// Vérifiez si le résultat XML correspond à l'attendu
	resultXML := resultXMLBuffer.String()

	if resultXML != expectedResultXML {
		t.Errorf("Le résultat XML ne correspond pas à l'attendu.\nAttendu:\n%s\nObtenu:\n%s", expectedResultXML, resultXML)
	}
}

func TestAttributsShouldSavedAfterParser(t *testing.T) {
	t.Parallel()
	// Fichier XML en entrée
	inputXML := `
	<root name="start">
		<name age="12" gender="male">Hello</name>
	</root>`

	// Lisez les résultats du canal et construisez le XML résultant
	var resultXMLBuffer bytes.Buffer

	// Créez un bufio.Reader à partir du XML en entrée
	reader := bytes.NewBufferString(inputXML)

	// Créez une nouvelle instance du parser XML avec la fonction de rappel et xPath
	parser := xixo.NewXMLParser(reader, &resultXMLBuffer).EnableXpath()
	parser.RegisterCallback("name", modifyElement1Content)
	// Créez un canal pour collecter les résultats du parser
	err := parser.Stream()
	assert.Nil(t, err)

	// Résultat XML attendu avec le contenu modifié et attributes restés
	expectedResultXML := `
	<root name="start">
		<name age="12" gender="male">ContenuModifie</name>
	</root>`

	// Vérifiez si le résultat XML correspond à l'attendu
	resultXML := resultXMLBuffer.String()

	if resultXML != expectedResultXML {
		t.Errorf("Le résultat XML ne correspond pas à l'attendu.\nAttendu:\n%s\nObtenu:\n%s", expectedResultXML, resultXML)
	}
}

func TestModifyAttributsWithMapCallback(t *testing.T) {
	t.Parallel()
	// Fichier XML en entrée
	inputXML := `
	<root>
		<element1 age="22" sex="male">Hello world!</element1>
		<element2>Contenu2 !</element2>
	</root>`

	// Lisez les résultats du canal et construisez le XML résultant
	var resultXMLBuffer bytes.Buffer

	// Créez un bufio.Reader à partir du XML en entrée
	reader := bytes.NewBufferString(inputXML)

	// Créez une nouvelle instance du parser XML avec la fonction de rappel
	parser := xixo.NewXMLParser(reader, &resultXMLBuffer).EnableXpath()
	parser.RegisterMapCallback("root", mapCallbackAttributsWithParent)

	// Créez un canal pour collecter les résultats du parser
	err := parser.Stream()
	assert.Nil(t, err)

	// Résultat XML attendu avec le contenu modifié
	expectedResultXML := `
	<root type="bar">
  <element1 age="50" sex="male">newChildContent</element1>
  <element2 age="25">Contenu2 !</element2>
</root>`

	// Vérifiez si le résultat XML correspond à l'attendu
	resultXML := resultXMLBuffer.String()

	if resultXML != expectedResultXML {
		t.Errorf("Le résultat XML ne correspond pas à l'attendu.\nAttendu:\n%s\nObtenu:\n%s", expectedResultXML, resultXML)
	}
}
