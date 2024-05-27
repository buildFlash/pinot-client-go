package pinot

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetQueryTemplate(t *testing.T) {
	assert.Equal(t, "http://%s/query/sql", getQueryTemplate("sql", "localhost:8000"))
	assert.Equal(t, "http://%s/query", getQueryTemplate("pql", "localhost:8000"))
	assert.Equal(t, "%s/query/sql", getQueryTemplate("sql", "http://localhost:8000"))
	assert.Equal(t, "%s/query", getQueryTemplate("pql", "http://localhost:8000"))
	assert.Equal(t, "%s/query/sql", getQueryTemplate("sql", "https://localhost:8000"))
	assert.Equal(t, "%s/query", getQueryTemplate("pql", "https://localhost:8000"))
}

func TestCreateHTTPRequest(t *testing.T) {
	r, err := createHTTPRequest(context.Background(), "localhost:8000", []byte(`{"sql": "select * from baseballStats limit 10"}`), map[string]string{"a": "b"})
	assert.Nil(t, err)
	assert.Equal(t, "POST", r.Method)
	_, err = createHTTPRequest(context.Background(), "localhos\t:8000", []byte(`{"sql": "select * from baseballStats limit 10"}`), map[string]string{"a": "b"})
	assert.NotNil(t, err)
}

func TestCreateHTTPRequestWithTrace(t *testing.T) {
	r, err := createHTTPRequest(context.Background(), "localhost:8000", []byte(`{"sql": "select * from baseballStats limit 10", "trace": "true"}`), map[string]string{"a": "b"})
	assert.Nil(t, err)
	assert.Equal(t, "POST", r.Method)
	_, err = createHTTPRequest(context.Background(), "localhos\t:8000", []byte(`{"sql": "select * from baseballStats limit 10", "trace": "true"}`), map[string]string{"a": "b"})
	assert.NotNil(t, err)
}

func TestJsonAsyncHTTPClientTransport(t *testing.T) {
	transport := &jsonAsyncHTTPClientTransport{
		client: http.DefaultClient,
		header: map[string]string{"a": "b"},
	}
	_, err := transport.execute(context.Background(), "localhos\t:8000", &Request{
		queryFormat: "sql",
		query:       "select * from baseballStats limit 10",
	})
	assert.NotNil(t, err)
	assert.True(t, strings.HasPrefix(err.Error(), "parse "))

	_, err = transport.execute(context.Background(), "randomhost", &Request{
		queryFormat: "sql",
		query:       "select * from baseballStats limit 10",
	})
	assert.NotNil(t, err)

	_, err = transport.execute(context.Background(), "localhost:18000", &Request{
		queryFormat:         "sql",
		query:               "select * from baseballStats limit 10",
		useMultistageEngine: true,
	})
	assert.NotNil(t, err)
	assert.True(t, strings.HasPrefix(err.Error(), "Post "))
}
