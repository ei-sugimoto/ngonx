package parser_test

import (
	"os"
	"sort"
	"testing"

	"github.com/ei-sugimoto/ngonx/pkg/parser"
)

func TestParse(t *testing.T) {
	// テスト用のTOMLファイルを作成
	testConfig := `
[server1]
host = "localhost"
port = 8080

[server2]
host = "example.com"
port = 9090
`
	// 一時ファイルを作成
	tmpfile, err := os.CreateTemp("", "config.toml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // テスト終了後に削除

	// 一時ファイルに書き込む
	if _, err := tmpfile.Write([]byte(testConfig)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// 元のファイルパスを一時ファイルのパスに変更
	originalPath := "/var/lib/ngonx/config.toml"
	os.Setenv("CONFIG_PATH", tmpfile.Name())
	defer os.Setenv("CONFIG_PATH", originalPath)

	// Parse関数を呼び出す
	p := parser.NewServer()
	p.Parse()
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// 結果を検証
	expected := map[string]parser.Server{
		"server1": {Host: "localhost", Port: 8080},
		"server2": {Host: "example.com", Port: 9090},
	}

	for key, expectedServer := range expected {
		if server, ok := p[key]; !ok {
			t.Errorf("expected server %s not found", key)
		} else if server != expectedServer {
			t.Errorf("server %s = %v, want %v", key, server, expectedServer)
		}
	}
}

func TestGetURLList(t *testing.T) {
	serverMap := parser.ServerMap{
		"server1": {Host: "localhost", Port: 8080, EndPoint: ""},
		"server2": {Host: "example.com", Port: 9090, EndPoint: ""},
	}

	expected := parser.URLAndEndPointList{
		{URL: "http://localhost:8080", EndPoint: ""},
		{URL: "http://example.com:9090", EndPoint: ""},
	}

	urlList := serverMap.GetURLList()

	if len(urlList) != len(expected) {
		t.Fatalf("expected %d URLs, got %d", len(expected), len(urlList))
	}

	extractedURLs := make([]string, len(urlList))
	for i, ue := range urlList {
		extractedURLs[i] = ue.URL
	}

	expectedURLs := make([]string, len(expected))
	for i, ue := range expected {
		expectedURLs[i] = ue.URL
	}

	sort.Strings(extractedURLs)
	sort.Strings(expectedURLs)

	for i, url := range extractedURLs {
		if url != expectedURLs[i] {
			t.Errorf("expected URL %s, got %s", expectedURLs[i], url)
		}
	}
}
