
# XIXO - XML Input XML Output

**XIXO** is a versatile command-line tool designed for seamless XML input and XML output operations. It empowers you to manipulate XML content effortlessly, offering a range of features for both basic and advanced operations. This README provides a comprehensive guide to **XIXO**, including installation instructions, usage examples, and acknowledgments to the open-source community.

## Installation

To install xixo, follow these steps:

- Download the right release for your operating system and architecture from the xixo releases page.

- Extract the xixo binary from the downloaded archive.

- Move the xixo binary to a directory that is included in your system's PATH. This step is essential to make xixo accessible from any location in your terminal.

Note: If you are not sure which directory to use, you can typically place it in /usr/local/bin or ~/bin (for Linux/macOS) or C:\Windows\System32 (for Windows). You may need administrator privileges to move the binary to some directories.

- Verify the installation by running the following command in your terminal:

```bash
$ xixo --version
```

If installed correctly, it should display the version of xixo.

Now, you can start using xixo to edit XML files with ease.

## Example

### Input XML

```xml
<root>
    <foo>
        <bar>a</bar>
        <baz>z</baz>
    </foo>
    <foo>
        <bar>b</bar>
    </foo>
    <baz>
        <bar>c</bar>
    </baz>
</root>
```

### Command

```shell
$ xixo  --subscribers foo="tee debug.jsonl | jq --unbuffered -c '.bar |= ascii_upcase' " < test/data/foo_bar_baz.xml
<root>
    <foo>
        <bar>A</bar>
        <baz>z</baz>
    </foo>
    <foo>
        <bar>B</bar>
    </foo>
    <baz>
        <bar>c</bar>
    </baz>
</root>
```

### Process Description

1. **Initialization**: **xixo** begins by parsing the input XML file in a streaming manner. It identifies the structure of the XML and locates elements that match the subscriber criteria (`foo` elements in this case).

2. **Subscriber Script Execution**: The subscriber script (`tee debug.jsonl | jq --unbuffered -c '.bar |= ascii_upcase'`) is executed once at the beginning of the parsing process, as indicated. It's important to note that the script is not called separately for each `foo` element but rather only once as the input is piped from **xixo**. The script performs the following steps:

    - It starts by using `tee` to write matched elements (e.g., `<foo><bar>a</bar><baz>z</baz></foo>`) as JSON lines, such as `{"bar":"a","baz":"z"}`, to the `debug.jsonl` file.
    - It then uses `jq` to apply the transformation to the `bar` element within the JSON content. The `ascii_upcase` function is used to convert the text within the `bar` element to uppercase. The `--unbuffered` flag ensures that **jq** processes the input line by line.
    - The script generates modified JSON lines, such as `{"bar":"A","baz":"z"}` and `{"bar":"B"}`, and writes them to the standard output.

3. **Merging JSON Output**: As the XML parsing progresses, **xixo** reads the JSON lines from the script's standard output. Each line of JSON data corresponds to a `foo` element in the XML. **xixo** combines this JSON data with the current matching XML element to produce the updated XML structure with the transformations applied. The modified XML is then emitted as output in a streaming manner.

4. **Final Output**: The final XML output, with the transformations applied to the `bar` elements within the `foo` elements.

### Key Points

- **Subscriber Script Execution**: The subscriber script is executed only once at the beginning of the parsing process and then stopped at the end. It processes the XML elements as they match the criteria and serializes them into JSON lines, which are then merged with the corresponding XML elements to produce the modified XML output.

- **Performance Optimization**: **xixo** optimizes performance by not calling the subscriber script for each `foo` element separately but rather processing the input in a stream and merging the results efficiently.

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
