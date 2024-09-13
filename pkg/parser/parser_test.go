package parser_test

import (
	"os"
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
