package parsehtml

import "fmt"

// from tag content, return the query string from current tag
func (h *HTMLObjectWrapper) ReverseFinding(tc string) string {
	var query string = ""
	var currentTag *Tag = h.RootTag
	var stack []*Tag = []*Tag{currentTag}

	for len(stack) > 0 {
		currentTag = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if h.HTMLContent[currentTag.TagPosition.Start:currentTag.TagPosition.End+1] == tc {
			break
		}
		stack = append(stack, currentTag.Childrens...)
	}

	if currentTag == nil || len(stack) == 0 {
		return ""
	}

	for currentTag.Name != "html" {
		query = currentTag.getQueryString() + "->" + query
		currentTag = currentTag.Father
	}

	return query[:len(query)-2]
}

func (t *Tag) getQueryString() string {
	var query string = t.Name
	var currentTag *Tag = t.Father
	if currentTag == nil {
		return query
	}
	var index int = 0
	for i := range currentTag.Childrens {
		if currentTag.Childrens[i].Name == t.Name {
			index++
		}
		if currentTag.Childrens[i] == t {
			query = fmt.Sprintf("%s#%d", query, index)
		}
	}
	if len(t.Attributes) == 0 {
		return query
	}
	var attrs string = ""
	for k, v := range t.Attributes {
		attrs = fmt.Sprintf("%s%s=%s,", attrs, k, v)
	}

	return fmt.Sprintf("%s{%s}", query, attrs[:len(attrs)-1])
}
