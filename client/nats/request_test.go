package nats

import "testing"

var Tests = []struct {
	URL      string
	Method   string
	Items    []interface{}
	Expected string
}{
	{"/applications", "get", []interface{}{}, "applications.get"},
	{"/applications/%s", "get", []interface{}{"app1"}, "applications.get.app1"},

	{"/endpoints/%s/%s/sendMessage%s", "post",
		[]interface{}{"tech1", "rsc1", "?from=f&body=d"}, "endpoints.sendMessage.post.tech1.rsc1.from.f.body.d"},
}

func TestConvertURL(t *testing.T) {

	for _, tx := range Tests {
		out := convertURL(tx.URL, tx.Method, tx.Items...)
		if out != tx.Expected {
			t.Errorf("convertURL(%s,%s,%v) => %s, expected %s", tx.URL, tx.Method, tx.Items, out, tx.Expected)
		}
	}

}
