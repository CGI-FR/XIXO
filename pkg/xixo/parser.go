package xixo

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/rs/zerolog/log"
)

type XMLParser struct {
	reader            *bufio.Reader
	writer            *bufio.Writer
	loopElements      map[string]Callback
	resultChannel     chan *XMLElement
	skipElements      map[string]bool
	attrOnlyElements  map[string]bool
	skipOuterElements bool
	xpathEnabled      bool
	scratch           *scratch
	scratch2          *scratch
	scratchWriter     *scratch
	deffer            bool
	TotalReadSize     uint64
	nextWrite         *byte
}

func NewXMLParser(reader io.Reader, writer io.Writer) *XMLParser {
	return &XMLParser{
		reader: bufio.NewReader(reader), writer: bufio.NewWriter(writer),
		loopElements:     map[string]Callback{},
		attrOnlyElements: map[string]bool{},
		resultChannel:    make(chan *XMLElement, 256),
		skipElements:     map[string]bool{},
		scratch:          &scratch{data: make([]byte, 1024)},
		scratch2:         &scratch{data: make([]byte, 1024)},
		scratchWriter:    &scratch{data: make([]byte, 1024)},
	}
}

func (x *XMLParser) Stream() error {
	defer func() {
		// write pending byte
		if x.nextWrite != nil {
			err := x.writer.WriteByte(*x.nextWrite)
			if err != nil {
				log.Error().Err(err).Msg("error when closing file")
			}
		}

		x.writer.Flush()
	}()

	err := x.parse()
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	return nil
}

func (x *XMLParser) RegisterCallback(match string, callback Callback) {
	x.loopElements[match] = callback
}

func (x *XMLParser) RegisterJSONCallback(match string, callback CallbackJSON) {
	x.loopElements[match] = XMLElementToJSONCallback(callback)
}

func (x *XMLParser) SkipElements(skipElements []string) *XMLParser {
	if len(skipElements) > 0 {
		for _, s := range skipElements {
			x.skipElements[s] = true
		}
	}

	return x
}

func (x *XMLParser) ParseAttributesOnly(loopElements ...string) *XMLParser {
	for _, e := range loopElements {
		x.attrOnlyElements[e] = true
	}

	return x
}

// by default skip elements works for stream elements childs
// if this method called parser skip also outer elements.
func (x *XMLParser) SkipOuterElements() *XMLParser {
	x.skipOuterElements = true

	return x
}

func (x *XMLParser) EnableXpath() *XMLParser {
	x.xpathEnabled = true

	return x
}

func (x *XMLParser) parse() error {
	defer close(x.resultChannel)

	var element *XMLElement

	var tagClosed bool

	var err error

	var b byte

	var iscomment bool

	err = x.skipDeclerations()

	if err != nil {
		return err
	}

	for {
		b, err = x.readByte()

		if err != nil {
			return err
		}

		if x.isWS(b) {
			continue
		}

		if b == '<' {
			iscdata, _, err := x.isCDATA()
			if err != nil {
				return err
			}

			if iscdata {
				continue
			}

			iscomment, err = x.isComment()

			if err != nil {
				return err
			}

			if iscomment {
				continue
			}

			x.defferWrite()

			element, tagClosed, err = x.startElement()

			if err != nil {
				return err
			}

			if _, found := x.loopElements[element.Name]; found {
				if tagClosed {
					continue
				}

				if _, ok := x.attrOnlyElements[element.Name]; !ok {
					element = x.getElementTree(element)
				}
				x.resultChannel <- element

				if callback, ok := x.loopElements[element.Name]; ok {
					mutatedElement, err := callback(element)
					if err != nil {
						return err
					}

					_, err = x.writer.WriteString(mutatedElement.String()[1:])
					if err != nil {
						return err
					}

					x.cancelDefferWrite()
				}

				if element.Err != nil {
					return element.Err
				}
			} else if x.skipOuterElements {
				if _, ok := x.skipElements[element.Name]; ok && !tagClosed {
					err = x.skipElement(element.Name)
					if err != nil {
						return err
					}

					continue
				}
			} else {
				err = x.commitDefferWrite()
				if err != nil {
					return err
				}
			}
		}
	}
}

func (x *XMLParser) getElementTree(result *XMLElement) *XMLElement {
	if result.Err != nil {
		return result
	}

	var (
		cur       byte
		next      byte
		err       error
		element   *XMLElement
		tagClosed bool
		iscomment bool
	)

	x.scratch2.reset() // this hold the inner text

	for {
		cur, err = x.readByte()

		if err != nil {
			result.Err = err

			return result
		}

		if cur == '<' {
			iscdata, cddata, err := x.isCDATA()
			if err != nil {
				result.Err = err

				return result
			}

			if iscdata {
				for _, cd := range cddata {
					x.scratch2.add(cd)
				}

				continue
			}

			iscomment, err = x.isComment()

			if err != nil {
				result.Err = err

				return result
			}

			if iscomment {
				continue
			}

			next, err = x.readByte()

			if err != nil {
				result.Err = err

				return result
			}

			if next == '/' { // close tag
				tag, err := x.closeTagName()
				if err != nil {
					result.Err = err

					return result
				}

				if tag == result.Name {
					if len(result.Childs) == 0 {
						result.InnerText = string(x.scratch2.bytes())
					}

					return result
				}
			} else {
				err = x.unreadByte()
				if err != nil {
					return nil
				}
			}

			x.defferWrite()
			element, tagClosed, err = x.startElement()

			if err != nil {
				result.Err = err

				return result
			}

			if _, ok := x.skipElements[element.Name]; ok && !tagClosed {
				err = x.skipElement(element.Name)
				if err != nil {
					result.Err = err

					return result
				}

				continue
			}

			if !tagClosed {
				element = x.getElementTree(element)
			}

			if x.xpathEnabled {
				element.parent = result
			}

			if _, ok := result.Childs[element.Name]; ok {
				result.Childs[element.Name] = append(result.Childs[element.Name], *element)
				if x.xpathEnabled {
					result.childs = append(result.childs, element)
				}
			} else {
				var childs []XMLElement
				childs = append(childs, *element)
				if result.Childs == nil {
					result.Childs = map[string][]XMLElement{}
				}
				result.Childs[element.Name] = childs

				if x.xpathEnabled {
					result.childs = append(result.childs, element)
				}
			}
		} else {
			x.scratch2.add(cur)
		}
	}
}

func (x *XMLParser) skipElement(elname string) error {
	var (
		c       byte
		next    byte
		err     error
		curname string
	)

	for {
		c, err = x.readByte()

		if err != nil {
			return err
		}

		if c == '<' {
			next, err = x.readByte()

			if err != nil {
				return err
			}

			if next == '/' {
				curname, err = x.closeTagName()
				if err != nil {
					return err
				}

				if curname == elname {
					return nil
				}
			}
		}
	}
}

func (x *XMLParser) startElement() (*XMLElement, bool, error) {
	x.scratch.reset()

	var (
		cur  byte
		prev byte
		err  error

		// a tag have 3 forms * <abc > ** <abc type="foo" val="bar"/> *** <abc />
		attr    string
		attrVal string
	)

	result := &XMLElement{}

	for {
		cur, err = x.readByte()

		if err != nil {
			return nil, false, x.defaultError()
		}

		if x.isWS(cur) {
			result.Name = string(x.scratch.bytes())

			if x.xpathEnabled {
				names := strings.Split(result.Name, ":")
				if len(names) > 1 {
					result.prefix = names[0]
					result.localName = names[1]
				} else {
					result.localName = names[0]
				}
			}

			x.scratch.reset()

			goto search_close_tag
		}

		//nolint: nestif
		if cur == '>' {
			if prev == '/' {
				result.Name = string(x.scratch.bytes()[:len(x.scratch.bytes())-1])

				if x.xpathEnabled {
					names := strings.Split(result.Name, ":")
					if len(names) > 1 {
						result.prefix = names[0]
						result.localName = names[1]
					} else {
						result.localName = names[0]
					}
				}

				return result, true, nil
			}

			result.Name = string(x.scratch.bytes())

			if x.xpathEnabled {
				names := strings.Split(result.Name, ":")
				if len(names) > 1 {
					result.prefix = names[0]
					result.localName = names[1]
				} else {
					result.localName = names[0]
				}
			}

			return result, false, nil
		}

		x.scratch.add(cur)
		prev = cur
	}

search_close_tag:
	for {
		cur, err = x.readByte()

		if err != nil {
			return nil, false, x.defaultError()
		}

		if x.isWS(cur) {
			continue
		}

		if cur == '=' {
			if result.Attrs == nil {
				result.Attrs = map[string]string{}
			}

			cur, err = x.readByte()

			if err != nil {
				return nil, false, x.defaultError()
			}

			if !(cur == '"' || cur == '\'') {
				return nil, false, x.defaultError()
			}

			attr = string(x.scratch.bytes())
			attrVal, err = x.string(cur)
			if err != nil {
				return nil, false, x.defaultError()
			}
			result.Attrs[attr] = attrVal
			if x.xpathEnabled {
				result.attrs = append(result.attrs, &xmlAttr{name: attr, value: attrVal})
			}
			x.scratch.reset()

			continue
		}

		if cur == '>' { // if tag name not found
			if prev == '/' { // tag special close
				return result, true, nil
			}

			return result, false, nil
		}

		x.scratch.add(cur)
		prev = cur
	}
}

func (x *XMLParser) isComment() (bool, error) {
	var (
		c   byte
		err error
	)

	c, err = x.readByte()

	if err != nil {
		return false, err
	}

	if c != '!' {
		err := x.unreadByte()
		if err != nil {
			return false, err
		}

		return false, nil
	}

	var d, e byte

	d, err = x.readByte()

	if err != nil {
		return false, err
	}

	e, err = x.readByte()

	if err != nil {
		return false, err
	}

	if d != '-' || e != '-' {
		err = x.defaultError()

		return false, err
	}

	// skip part
	x.scratch.reset()

	for {
		c, err = x.readByte()

		if err != nil {
			return false, err
		}

		if c == '>' &&
			len(x.scratch.bytes()) > 1 &&
			x.scratch.bytes()[len(x.scratch.bytes())-1] == '-' &&
			x.scratch.bytes()[len(x.scratch.bytes())-2] == '-' {
			return true, nil
		}

		x.scratch.add(c)
	}
}

func (x *XMLParser) isCDATA() (bool, []byte, error) {
	var (
		c   byte
		b   []byte
		err error
	)

	b, err = x.reader.Peek(2)

	if err != nil {
		return false, nil, err
	}

	if b[0] != '!' {
		return false, nil, nil
	}

	if err != nil {
		return false, nil, err
	}

	if b[1] != '[' {
		return false, nil, nil
	}

	// read peaked byte
	_, err = x.readByte()

	if err != nil {
		return false, nil, err
	}

	_, err = x.readByte()

	if err != nil {
		return false, nil, err
	}

	c, err = x.readByte()

	if err != nil {
		return false, nil, err
	}

	if c != 'C' {
		err = x.defaultError()

		return false, nil, err
	}

	c, err = x.readByte()

	if err != nil {
		return false, nil, err
	}

	if c != 'D' {
		err = x.defaultError()

		return false, nil, err
	}

	c, err = x.readByte()

	if err != nil {
		return false, nil, err
	}

	if c != 'A' {
		err = x.defaultError()

		return false, nil, err
	}

	c, err = x.readByte()

	if err != nil {
		return false, nil, err
	}

	if c != 'T' {
		err = x.defaultError()

		return false, nil, err
	}

	c, err = x.readByte()

	if err != nil {
		return false, nil, err
	}

	if c != 'A' {
		err = x.defaultError()

		return false, nil, err
	}

	c, err = x.readByte()

	if err != nil {
		return false, nil, err
	}

	if c != '[' {
		err = x.defaultError()

		return false, nil, err
	}

	// this is possibly cdata // ]]>
	x.scratch.reset()

	for {
		c, err = x.readByte()

		if err != nil {
			return false, nil, err
		}

		if c == '>' &&
			len(x.scratch.bytes()) > 1 &&
			x.scratch.bytes()[len(x.scratch.bytes())-1] == ']' &&
			x.scratch.bytes()[len(x.scratch.bytes())-2] == ']' {
			return true, x.scratch.bytes()[:len(x.scratch.bytes())-2], nil
		}

		x.scratch.add(c)
	}
}

func (x *XMLParser) skipDeclerations() error {
	var (
		a, b []byte
		c, d byte
		err  error
	)

scan_declartions:
	for {
		// when identifying a xml declaration we need to know 2 bytes ahead.
		// Unread works 1 byte at a time so we use Peek and read together.
		a, err = x.reader.Peek(1)

		if err != nil {
			return err
		}

		if a[0] == '<' {
			b, err = x.reader.Peek(2)

			if err != nil {
				return err
			}

			if b[1] == '!' || b[1] == '?' { // either comment or declaration
				// read 2 peaked byte
				_, err = x.readByte()

				if err != nil {
					return err
				}

				_, err = x.readByte()
				if err != nil {
					return err
				}

				c, err = x.readByte()

				if err != nil {
					return err
				}

				d, err = x.readByte()

				if err != nil {
					return err
				}

				if c == '-' && d == '-' {
					goto skipComment
				}

				goto skipDecleration
			}

			return nil
		}

		// read peaked byte
		_, err = x.readByte()

		if err != nil {
			return err
		}
	}

skipComment:
	x.scratch.reset()

	for {
		c, err = x.readByte()

		if err != nil {
			return err
		}

		if c == '>' &&
			len(x.scratch.bytes()) > 1 &&
			x.scratch.bytes()[len(x.scratch.bytes())-1] == '-' &&
			x.scratch.bytes()[len(x.scratch.bytes())-2] == '-' {
			goto scan_declartions
		}

		x.scratch.add(c)
	}

skipDecleration:
	depth := 1

	for {
		c, err = x.readByte()

		if err != nil {
			return err
		}

		if c == '>' {
			depth--
			if depth == 0 {
				goto scan_declartions
			}

			continue
		}

		if c == '<' {
			depth++
		}
	}
}

func (x *XMLParser) closeTagName() (string, error) {
	x.scratch.reset()

	var (
		c   byte
		err error
	)

	for {
		c, err = x.readByte()

		if err != nil {
			return "", err
		}

		if c == '>' {
			return string(x.scratch.bytes()), nil
		}

		if !x.isWS(c) {
			x.scratch.add(c)
		}
	}
}

func (x *XMLParser) defferWrite() {
	x.scratchWriter.reset()
	x.deffer = true
}

func (x *XMLParser) cancelDefferWrite() {
	x.scratchWriter.reset()
	x.deffer = false
}

func (x *XMLParser) commitDefferWrite() error {
	_, err := x.writer.Write(x.scratchWriter.bytes())
	if err != nil {
		return err
	}

	x.scratchWriter.reset()
	x.deffer = false

	return nil
}

func (x *XMLParser) readByte() (byte, error) {
	by, err := x.reader.ReadByte()
	if err != nil {
		return 0, err
	}

	if !x.deffer {
		if x.nextWrite != nil {
			err = x.writer.WriteByte(*x.nextWrite)
			if err != nil {
				return 0, err
			}
		}

		x.nextWrite = &by
	} else {
		x.scratchWriter.add(by)
	}

	x.TotalReadSize++

	return by, nil
}

func (x *XMLParser) unreadByte() error {
	err := x.reader.UnreadByte()
	if err != nil {
		return err
	}

	if x.nextWrite != nil {
		x.nextWrite = nil
	}

	x.TotalReadSize--

	return nil
}

func (x *XMLParser) isWS(in byte) bool {
	if in == ' ' || in == '\n' || in == '\t' || in == '\r' {
		return true
	}

	return false
}

func (x *XMLParser) defaultError() error {
	err := fmt.Errorf("invalid xml")

	return err
}

func (x *XMLParser) string(start byte) (string, error) {
	x.scratch.reset()

	var (
		err error
		c   byte
	)

	for {
		c, err = x.readByte()
		if err != nil {
			if err != nil {
				return "", err
			}
		}

		if c == start {
			return string(x.scratch.bytes()), nil
		}

		x.scratch.add(c)
	}
}

// scratch taken from
// https://github.com/bcicen/jstream
type scratch struct {
	data []byte
	fill int
}

// reset scratch buffer.
func (s *scratch) reset() { s.fill = 0 }

// bytes returns the written contents of scratch buffer.
func (s *scratch) bytes() []byte {
	return s.data[0:s.fill]
}

// grow scratch buffer.
func (s *scratch) grow() {
	ndata := make([]byte, cap(s.data)*2)
	copy(ndata, s.data)
	s.data = ndata
}

// append single byte to scratch buffer.
func (s *scratch) add(c byte) {
	if s.fill+1 >= cap(s.data) {
		s.grow()
	}

	s.data[s.fill] = c
	s.fill++
}
