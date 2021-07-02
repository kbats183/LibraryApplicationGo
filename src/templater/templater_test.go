package templater

import (
	"bytes"
	"os"
	"testing"
)

func TestExecuteTemplate(t *testing.T) {
	var err error

	err = os.Mkdir("template", 0700)
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		err = os.RemoveAll("template")
		if err != nil {
			t.Error(err)
			return
		}
	}()

	testBaseTemplate, err := os.Create("template/base.html")
	if err != nil {
		t.Error(err)
		return
	}

	_, err = testBaseTemplate.WriteString("{{define \"template\"}}It's text form base.html{{end}}")
	if err != nil {
		t.Error(err)
		return
	}

	err = testBaseTemplate.Close()
	if err != nil {
		t.Error(err)
		return
	}

	testTemplate, err := os.Create("template/__test_template__.html")
	if err != nil {
		t.Error(err)
		return
	}

	_, err = testTemplate.WriteString("{{.PageName}}: {{.Data}}\n{{template \"template\" .}}")
	if err != nil {
		t.Error(err)
		return
	}

	err = testTemplate.Close()
	if err != nil {
		t.Error(err)
		return
	}

	var buf []byte
	buffer := bytes.NewBuffer(buf)

	ExecuteTemplate("__test_template__", "TestPage", buffer, "Hello world!", nil)

	t.Log(buffer.String())
	if buffer.String() != "TestPage: Hello world!\nIt's text form base.html" {
		t.Fail()
	}
}
