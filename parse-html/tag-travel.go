package parsehtml

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func (t *Tag) IsQualifiedWithOptions(options ...string) bool {
	for _, option := range options {
		equalSignIndex := strings.Index(option, "=")
		if equalSignIndex == -1 {
			return false
		}
		optionName := option[:equalSignIndex]
		optionValue := option[equalSignIndex+1:]
		if t.Attributes[optionName] == "class" {
			classes := strings.Split(t.Attributes[optionName], " ")
			if !stringInSlice(optionValue, classes) {
				return false
			}
			continue
		}
		if t.Attributes[optionName] != optionValue {
			return false
		}
	}
	return true
}

func stringInSlice(str string, slice []string) bool {
	for _, elem := range slice {
		if elem == str {
			return true
		}
	}
	return false
}

func (h *HTMLObjectWrapper) FirstTag(name string, options ...string) *Tag {
	// implement this by using DFS
	var rootTag *Tag = h.RootTag
	return rootTag.FirstTag(name, options...)
}

func (h *HTMLObjectWrapper) AllTag(name string, options ...string) []*Tag {
	// implement this by using DFS
	var rootTag *Tag = h.RootTag
	return rootTag.AllTag(name, options...)
}

func (t *Tag) FirstTag(name string, options ...string) *Tag {
	// implement this by using DFS
	var stack []*Tag = []*Tag{t}
	for len(stack) > 0 {
		var currentTag *Tag = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if currentTag.Name == name && currentTag.IsQualifiedWithOptions(options...) {
			return currentTag
		}
		stack = append(stack, currentTag.Childrens...)
	}
	return nil
}

func (t *Tag) AllTag(name string, options ...string) []*Tag {
	// implement this by using DFS
	var result []*Tag
	var stack []*Tag = []*Tag{t}
	for len(stack) > 0 {
		var currentTag *Tag = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if currentTag.Name == name && currentTag.IsQualifiedWithOptions(options...) {
			result = append(result, currentTag)
		}
		stack = append(stack, currentTag.Childrens...)
	}
	return result
}

// query: div->div#2->div#2-4{class=hello,name=world}
// #2 means the second div tag
// #2-4{class=hello} means finding from the second div tag to the fourth div tag
// which has class attribute equals to hello
func (t *Tag) QuerySelector(query string) []*Tag {
	path := strings.Split(query, "->")
	// implement this by using DFS
	var currentResultTags []*Tag = []*Tag{t}
	for _, query := range path {
		var nextResultTags []*Tag
		for _, tag := range currentResultTags {
			queryResult, err := tag.singleQuerySelector(query)
			if err != nil {
				fmt.Println(err)
				continue
			}
			nextResultTags = append(nextResultTags, queryResult...)
		}
		currentResultTags = nextResultTags
		if len(currentResultTags) == 0 {
			fmt.Println("[query]: cannot find any tag at tag:", currentResultTags[len(currentResultTags)-1].Name)
		}
	}

	return currentResultTags
}

type qQuery struct {
	tagName string
	index   string
	attrs   []string
	query   string
}

func (q *qQuery) parse() {

	// extract next tag name, index and attributes
	var sharpIndex int = strings.Index(q.query, "#")
	if sharpIndex == -1 {
		q.tagName = q.query
		q.index = "*"
		q.attrs = []string{}
		return
	} else {
		q.tagName = q.query[:sharpIndex]
	}
	var curlyBracketIndex int = strings.Index(q.query, "{")
	if curlyBracketIndex == -1 {
		q.index = q.query[sharpIndex+1:]
		q.attrs = []string{}
	} else {
		q.index = q.query[sharpIndex+1 : curlyBracketIndex]
		q.attrs = strings.Split(q.query[curlyBracketIndex+1:len(q.query)-1], ",")
	}
}

func (t *Tag) singleQuerySelector(query string) ([]*Tag, error) {
	var currentResultTags []*Tag = []*Tag{}
	var q qQuery = qQuery{
		query: query,
	}
	q.parse()

	// find next tag
	for _, tag := range t.Childrens {
		if tag.Name == q.tagName && tag.IsQualifiedWithOptions(q.attrs...) {
			currentResultTags = append(currentResultTags, tag)
		}
	}

	// find next tag by index
	if q.tagName == "" || q.index == "*" || len(currentResultTags) == 0 {
		return currentResultTags, nil
	}

	var slashSignIndex int = strings.Index(q.index, "-")
	if slashSignIndex == -1 {
		index, err := strconv.Atoi(q.index)
		if err != nil {
			return nil, errors.New("[query]: invalid index, cannot convert to int")
		}
		if index > len(currentResultTags) {
			return nil, errors.New("[query]: invalid index, index out of range")
		}
		return []*Tag{currentResultTags[index-1]}, nil
	}
	from, err := strconv.Atoi(q.index[:slashSignIndex])
	if err != nil {
		return nil, errors.New("[query]: invalid index, cannot convert to int")
	}
	to, err := strconv.Atoi(q.index[slashSignIndex+1:])
	if err != nil {
		return nil, errors.New("[query]: invalid index, cannot convert to int")
	}
	if from > len(currentResultTags) || to > len(currentResultTags) {
		return nil, errors.New("[query]: invalid index, index out of range")
	}
	if from > to {
		return append(currentResultTags[from:], currentResultTags[:to-1]...), nil
	}
	return currentResultTags[from-1 : to], nil
}
