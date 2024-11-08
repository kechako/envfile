package envfile

import (
	"errors"
	"maps"
	"os"
	"slices"
	"testing"
)

var envsTest = struct {
	envs  Envs
	kvMap map[string]string
}{
	envs: Envs{
		"AAAA=aaaa",
		"BBBB=bbbb",
		"CCCC=cccc",
		"DDDD=dddd",
		"EEEE=eeee",
	},
	kvMap: map[string]string{
		"AAAA": "aaaa",
		"BBBB": "bbbb",
		"CCCC": "cccc",
		"DDDD": "dddd",
		"EEEE": "eeee",
	},
}

func TestEnvsEnvs(t *testing.T) {
	i := 0
	for key, value := range envsTest.envs.Envs() {
		got := key + "=" + value
		want := envsTest.envs[i]
		if got != want {
			t.Errorf("Envs.Envs(): #%d: got %s, want %s", i, got, want)
		}
		i++
	}
}

func TestEnvsMap(t *testing.T) {
	got := envsTest.envs.Map()
	want := envsTest.kvMap

	if !maps.Equal(got, want) {
		t.Errorf("Envs.Map(): got %#v, want %#v", got, want)
	}
}

var parseTests = map[string]struct {
	envs Envs
	err  error
}{
	"test01.env": {
		envs: Envs{
			"FOO=foo",
			"BAR= bar bar",
			"BAZ=baz # baz",
		},
		err: nil,
	},
	"test02.env": {
		envs: Envs{
			"FOO=foo",
			"BAR=bar",
			"BAZ=baz",
		},
		err: nil,
	},
	"test03.env": {
		envs: nil,
		err:  errors.New("invalid UTF-8 bytes at line 2"),
	},
	"test04.env": {
		envs: nil,
		err:  errors.New("no variable key on line 3"),
	},
	"test05.env": {
		envs: nil,
		err:  errors.New("the variable key contains whitespace: 'FOO FOO'"),
	},
	"test06.env": {
		envs: nil,
		err:  errors.New("the variable key contains whitespace: 'FOO '"),
	},
}

func TestParseFile(t *testing.T) {
	for name, tt := range parseTests {
		t.Run(name, func(t *testing.T) {
			envs, err := ParseFile("testdata/" + name)
			if !slices.Equal(envs, tt.envs) || !errorEqual(err, tt.err) {
				t.Errorf("ParseFile(): got (%#v, %v), want (%#v, %v)", envs, err, tt.envs, tt.err)
			}
		})
	}
}

func TestLoadEnvs(t *testing.T) {
	err := LoadEnvs(envsTest.envs)
	if err != nil {
		t.Fatal(err)
	}

	for key, value := range envsTest.envs.Envs() {
		got := os.Getenv(key)
		want := value
		if got != want {
			t.Errorf("os.Getenv(\"%s\"): got %s, want %s", key, got, want)
		}
	}
}

func errorEqual(a, b error) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	return a.Error() == b.Error()
}
