package parsehtml

import (
	"fmt"
	"strings"
)

type TagPosition struct {
	Start int
	End   int
}

type HTMLObjectWrapper struct {
	RootTag      *Tag
	DeepestLevel int
	HTMLContent  string
	tagIDs       map[string]*Tag
	tables       []*Tag
}

type Tag struct {
	Name              string
	Attributes        map[string]string
	Childrens         []*Tag
	Father            *Tag
	Text              []string
	TagPosition       *TagPosition
	ClosedTagPosition *TagPosition
}

type TaggingCharacter struct {
	TagChar  string
	Position int
	Level    int
}

func (t *Tag) ShowTagBrief() {
	fmt.Println("Tag Name:", t.Name)
	fmt.Println("\tTag Attributes:", t.Attributes)
	fmt.Println("\tTag Children:", len(t.Childrens))
	fmt.Println("\tTag Father:", t.Father)
	fmt.Println("\tTag Position:", t.TagPosition)
	fmt.Println("\tTag Closed Position:", t.ClosedTagPosition)
	fmt.Println("\tTag Text:", t.Text)
}

func getListOfTaggingCharacters(html string) []TaggingCharacter {
	var result []TaggingCharacter
	var inSingleQuoteFlag bool = false
	var inDoubleQuoteFlag bool = false

	for i, char := range html {
		if char == '\'' {
			inSingleQuoteFlag = !inSingleQuoteFlag
			continue
		}
		if char == '"' {
			inDoubleQuoteFlag = !inDoubleQuoteFlag
			continue
		}
		if char == '<' && !inSingleQuoteFlag && !inDoubleQuoteFlag {
			if html[i+1] == '/' {
				result = append(result, TaggingCharacter{TagChar: "</", Position: i})
				continue
			}
			result = append(result, TaggingCharacter{TagChar: "<", Position: i})
			continue
		}
		if char == '>' && !inSingleQuoteFlag && !inDoubleQuoteFlag {
			if html[i-1] == '/' {
				result = append(result, TaggingCharacter{TagChar: "/>", Position: i})
				continue
			}
			result = append(result, TaggingCharacter{TagChar: ">", Position: i})
			continue
		}
	}
	return result
}

func setLevelForTaggingCharacters(lt *[]TaggingCharacter) int {
	var currLevel int = -1
	var deepestLevel int = -1

	for i, tagChar := range *lt {
		if tagChar.TagChar == "<" {
			currLevel++
			(*lt)[i].Level = currLevel
			if currLevel > deepestLevel {
				deepestLevel = currLevel
			}
			continue
		}
		if tagChar.TagChar == "/>" {
			(*lt)[i].Level = currLevel
			currLevel--
			continue
		}
		if tagChar.TagChar == "</" {
			(*lt)[i].Level = currLevel
			currLevel--
			continue
		}
		if tagChar.TagChar == ">" {
			(*lt)[i].Level = (*lt)[i-1].Level
			continue
		}
	}
	return deepestLevel
}

func createTagObject(html string, lt *[]TaggingCharacter) (*Tag, map[string]*Tag, []*Tag) {
	// init the html tag (root of the tree)
	var rootTag *Tag = new(Tag)
	rootTag.Name = "html"
	rootTag.Childrens = make([]*Tag, 0)
	rootTag.Father = nil
	rootTag.Text = make([]string, 0)
	rootTag.TagPosition = &TagPosition{
		Start: (*lt)[0].Position,
		End:   (*lt)[1].Position,
	}

	var currTag *Tag = rootTag
	listOfTag := make([]*Tag, 1)
	listOfTag[0] = currTag

	var numberOfTag int = len(*lt) / 2

	// init the list of tag
	// exclude the html tag
	for i := 1; i < numberOfTag-1; i++ {
		text := getTextBetweenTag((*lt)[2*i-1].Position+1, (*lt)[2*i].Position, html)
		if text != "" {
			currTag.Text = append(currTag.Text, text)
		}
		if (*lt)[2*i].TagChar == "<" {
			newTag := new(Tag)
			newTag.Childrens = make([]*Tag, 0)
			newTag.Father = currTag
			newTag.TagPosition = &TagPosition{
				Start: (*lt)[2*i].Position,
				End:   (*lt)[2*i+1].Position,
			}
			newTag.Text = make([]string, 0)
			listOfTag = append(listOfTag, newTag)
			// fmt.Println("[debug]: ", html[currTag.TagPosition.Start:currTag.TagPosition.End], "->", html[newTag.TagPosition.Start:newTag.TagPosition.End])
			currTag = newTag
		}

		if (*lt)[2*i].TagChar == "</" || (*lt)[2*i+1].TagChar == "/>" {
			currTag.ClosedTagPosition = &TagPosition{
				Start: (*lt)[2*i].Position,
				End:   (*lt)[2*i+1].Position,
			}
			currTag.Father.Childrens = append(currTag.Father.Childrens, currTag)
			// fmt.Println("[debug]: ", html[currTag.TagPosition.Start:currTag.TagPosition.End], ">-", html[currTag.Father.TagPosition.Start:currTag.Father.TagPosition.End])
			currTag = currTag.Father
		}
	}

	// extract all tag's id and find all tables
	tagIds := make(map[string]*Tag)
	var tables []*Tag
	for i, tag := range listOfTag {
		extractTagAttributes(html[listOfTag[i].TagPosition.Start:listOfTag[i].TagPosition.End], listOfTag[i])
		if tag.Attributes["id"] != "" {
			tagIds[tag.Attributes["id"]] = tag
		}
		if tag.Name == "table" {
			tables = append(tables, tag)
		}
	}

	return rootTag, tagIds, tables
}

func (h *HTMLObjectWrapper) ParseHTML() {
	listOfTaggingCharacters := getListOfTaggingCharacters(h.HTMLContent)
	h.DeepestLevel = setLevelForTaggingCharacters(&listOfTaggingCharacters)
	h.RootTag, h.tagIDs, h.tables = createTagObject(h.HTMLContent, &listOfTaggingCharacters)
}

func extractTagAttributes(tag string, TagObj *Tag) {
	tag = strings.Trim(tag, "<>")
	attributes := strings.Split(tag, " ")
	TagObj.Name = attributes[0]
	if len(attributes) == 1 {
		return
	}
	TagObj.Attributes = make(map[string]string)
	for _, attribute := range attributes[1:] {
		keyValue := strings.Split(attribute, "=")
		for i, elem := range keyValue {
			keyValue[i] = strings.Trim(elem, "\"'")
		}
		TagObj.Attributes[keyValue[0]] = keyValue[1]
	}
}

func getTextBetweenTag(start int, end int, html string) string {
	text := html[start:end]
	text = strings.Trim(text, " \n\t")
	return text
}
