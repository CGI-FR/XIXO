
# XIXO - XML Input XML Output

**XIXO** is a Go library that allows you to parse XML files and extract custom attributes. It empowers you to manipulate XML content effortlessly, offering a range of features for both basic and advanced operations. This README provides a comprehensive guide to **XIXO**, including installation instructions, usage examples, and acknowledgments to the open-source community.

## Installation

To install xixo,  you can use go get:

```
go get github.com/CGI-FR/XIXO
```

Now, you can start using xixo to edit XML files with ease.

## Example

To use **XIXO**, you need to create a Parser object with the path of the XML file to parse and the name of element, here is a emexple of XML file: (same exemple in Unit testing **TestMapCallbackWithAttributsParentAndChilds()** in callback_test.go )

```xml
<root type="foo">
    <element1 age="22" sex="male">Hello world !</element1>
    <element2>Contenu2 </element2>
</root>
```

### Process Description

1. **Initialization**: **xixo** begins by parsing the input XML file in a streaming manner. It identifies the structure of the XML and locates elements that match the begin element name. (`root` elements in this case).

2. **Transform XML to Go map** As the XML parsing progresses, **xixo** reads the XML file and creat a element tree of root element. In exemple will be:

```go
{"@type":"foo","element1":"Hello world !","element1@age":"22","element1":"male","element2":"Contenu2 "}.
```

3. **Edite element value and attribute value with callback**: For modify the data, we need give callback function a map with the element name as key, new data as value. In exemple we give a map like this:

```go
{"@type":"bar","element1@age":"50","element1":"newChildContent","element2@age":"25"}.
```

4. **Final Output**: The final XML output will be:

```xml
<root type="bar">
    <element1 age="50" sex="male">newChildContent</element1>
    <element2 age="25">Contenu2 </element2>
</root>
```

### Key Points

- **Performance Optimization**: **xixo** optimizes performance by not calling the subscriber script for each `root` element separately but rather processing the input in a stream and merging the results efficiently.

This detailed process demonstrates how **xixo** processes XML files in a streaming and efficient manner, applying custom transformations to specific elements using subscribers.

## License

**xixo** is licensed under the MIT License. See the [LICENSE](https://github.com/youen/xixo/blob/main/LICENSE) file for details.

## Acknowledgments

The XML parsing functionality in **xixo** is based on the excellent work of the developers behind the [xml-stream-parser](https://github.com/tamerh/xml-stream-parser) project. We would like to extend our gratitude to the following contributors:

- [Tamer Gür](https://github.com/tamerh) (Tamer Gür)
- [Jiří Setnička](https://github.com/setnicka) (Jiří Setnička)
- [tsak](https://github.com/tsak) (tsak)
- [Ilia Mirkin](https://github.com/imirkin) (Ilia Mirkin)

Their work has been instrumental in enabling efficient XML parsing within **xixo**, and we appreciate their contributions to the open-source community.
