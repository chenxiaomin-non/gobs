# gobs
A Simple Web Crawler by Golang

# Crawl a web's HTML Content
Step 1: Get the HTML Content from the Website/Web page using HTTP Connection (in get-html module)
```
import (
    gh "github.com/chenxiaomin-non/gobs/gethtml"
    "fmt"
)

func main() {
    var req gh.HttpRequest = gh.HttpRequest{
      Url: "fap.fpt.edu.vn/Student.aspx",
      Method: "GET",
      Header: nil,
      Body: nil,
    }
    res, err := gh.GetHttp(req)
    if err != nil {
        fmt.Printf("[error]: HTTP Request failed!")
        return
    }
    html := res.Body
}
```

Step 2: Create a HTML Object Wrapper and parse the HTML Content in to searchable and search on it
```
import (
    parser "github.com/chenxiaomin-non/gobs/parsehtml"
    kmp "github.com/chenxiaomin-non/gobs/kmp"
    "fmt"
)

func main() {
    var html string = "<html> <body>hello, world<div>Some text here!</div> <div id="hello"></div></body></html>"
    // init the Wrapper Object 
    wrapper := parser.HTMLObjectWrapper{
        HTMLContent: html,
    }
	  wrapper.ParseHTML()

    // some simple operations on this
    wrapper.RootTag.ShowTagBrief()   // show html tag info
    wrapper.RootTag.QuerySelector("body->div#2{id=hello}")[0].ShowTagBried()  // find div tag by the path
    wrapper.FindAllTag("div")[0].ShowTagBrief()   // find all div tag
    wrapper.RootTag.QuerySelector("body->div#1-2")[0].ShowTagBried()  // search tag by range
    fmt.Println(wrapper.ReverseFinding("<div id=\"hello\">"))  // from the tag content, find the path (for optimize finding content in large scale)

    // also provide a way to find table by KMP substring searching (which is faster, maybe)
    tableTags, _ := kmp.ProcessTableTag(html)
}
```
