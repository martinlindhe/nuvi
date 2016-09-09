package news

import (
	"testing"
)

func TestImportResource(t *testing.T) {
	err := importResource(rootURL+"1473152352858.zip", true)
	if err != nil {
		t.Error("err not nil")
	}
}
