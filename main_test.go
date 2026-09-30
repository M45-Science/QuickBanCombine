package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"
)

func runCombineTest(t *testing.T, inputs ...string) []banDataType {
	t.Helper()
	dir := t.TempDir()
	oldArgs := os.Args
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		os.Args = oldArgs
		if err := os.Chdir(oldDir); err != nil {
			t.Error(err)
		}
	})
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	os.Args = []string{"QuickBanCombine"}
	for i, body := range inputs {
		name := filepath.Join(dir, strconv.Itoa(i)+".json")
		if err := os.WriteFile(name, []byte(body), 0600); err != nil {
			t.Fatal(err)
		}
		os.Args = append(os.Args, name)
	}
	main()
	data, err := os.ReadFile("composite.json")
	if err != nil {
		t.Fatal(err)
	}
	var bans []banDataType
	if err := json.Unmarshal(data, &bans); err != nil {
		t.Fatal(err)
	}
	return bans
}

func TestCombineBansAcrossFiles(t *testing.T) {
	got := runCombineTest(t, `["ALICE","bob"]`, `[{"username":"alice","reason":"first"},{"username":"ALICE","reason":"second"},{"username":"Carol"}]`)
	want := []banDataType{{UserName: "alice", Reason: "[dup] first, second"}, {UserName: "bob"}, {UserName: "carol"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("composite = %#v, want %#v", got, want)
	}
}

func TestCombineIdenticalDuplicatesKeepsOneBan(t *testing.T) {
	got := runCombineTest(t, `["Alice","alice","Bob",""]`)
	want := []banDataType{{UserName: "alice"}, {UserName: "bob"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("composite = %#v, want %#v", got, want)
	}
}

func TestCombineInvalidInputPreservesPreviousOutput(t *testing.T) {
	runCombineTest(t, `["alice"]`)
	before, err := os.ReadFile("composite.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile("invalid.json", []byte(`{"not":"a ban list"}`), 0600); err != nil {
		t.Fatal(err)
	}
	if err := combineFiles([]string{"invalid.json"}); err == nil {
		t.Fatal("invalid input accepted")
	}
	after, err := os.ReadFile("composite.json")
	if err != nil || string(after) != string(before) {
		t.Fatalf("previous output changed: %q, error = %v", after, err)
	}
}
