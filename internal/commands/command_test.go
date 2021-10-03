package commands

import (
	"fmt"
	"testing"
)

type testCommand struct {
	prefix     string
	command    string
	parameters []string
	trailing   string
}

func TestCommandStruct(t *testing.T) {
	t.Run("command parsing", func(t *testing.T) {
		tt := []struct {
			name  string
			input string
			want  *testCommand
			err   error
		}{
			{
				"PING command 2 params",
				"PING User1 example.com",
				&testCommand{"", "PING", []string{"User1", "example.com"}, ""},
				nil,
			},
			{
				"PING command 1 param",
				"PING example.com",
				&testCommand{"", "PING", []string{"example.com"}, ""},
				nil,
			},
			{
				"PONG command",
				":user@user.com PONG",
				&testCommand{"user@user.com", "PONG", []string{}, ""},
				nil,
			},
			{
				"PONG command with error",
				"PONG a:a",
				nil,
				ErrMalformedCommandString,
			},
		}

		for _, test := range tt {
			t.Run(test.name, func(t *testing.T) {
				got, err := NewCommand(test.input)

				if test.err != nil {
					if err == nil || err != test.err {
						t.Fatal("expected to get an error")
					}
					return
				}

				if got == nil {
					t.Fatal("NewCommand returned nil")
				}

				res := ""
				if got.Prefix != test.want.prefix {
					res += fmt.Sprintf("Prefix: got %s want %s\n", got.Prefix, test.want.prefix)
				}
				if got.Command != test.want.command {
					res += fmt.Sprintf("Command: got %s want %s\n", got.Command, test.want.command)
				}
				if !compareStringArrays(t, got.Parameters, test.want.parameters) {
					res += fmt.Sprintf("Parameters: got %s want %s\n", got.Parameters, test.want.parameters)
				}
				if got.Trailing != test.want.trailing {
					res += fmt.Sprintf("Trailing: got %s want %s\n", got.Trailing, test.want.trailing)
				}
				if res != "" {
					t.Fatal(res)
				}
			})
		}
	})
}

func TestFields(t *testing.T) {
	t.Run("test custom fields function", func(t *testing.T) {
		tt := []struct {
			name  string
			input string
			want  []string
			err   error
		}{
			{
				"empty prefix",
				": PING example.com",
				[]string{":", "PING", "example.com"},
				nil,
			},
			{
				"empty postfix",
				"PING example.com :",
				nil,
				ErrMalformedCommandString,
			},
			{
				"unexpected colon",
				"PING e:xample.com",
				nil,
				ErrMalformedCommandString,
			},
			{
				"well formed command",
				":example.com PING Wiz :example.com",
				[]string{":example.com", "PING", "Wiz", ":example.com"},
				nil,
			},
			{
				"PONG command",
				":user@user.com PONG",
				[]string{":user@user.com", "PONG"},
				nil,
			},
		}

		for _, test := range tt {
			t.Run(test.name, func(t *testing.T) {
				got, err := fields(test.input)

				if test.err != nil {
					if err == nil || err != test.err {
						t.Fatal("expected to get an error")
					}
					return
				}

				if got == nil {
					t.Fatal("unexpected nil")
				}

				if !compareStringArrays(t, got, test.want) {
					t.Errorf("got %s, want %s", got, test.want)
				}
			})
		}
	})
}

func ExampleFields() {
	res, _ := fields(":Angel!wings@irc.org PRIVMSG Wiz :Are you receiving this message ?")
	fmt.Printf("%q", res)
	//Output: [":Angel!wings@irc.org" "PRIVMSG" "Wiz" ":Are you receiving this message ?"]
}

func compareStringArrays(t testing.TB, arr1, arr2 []string) bool {
	t.Helper()

	if len(arr1) != len(arr2) {
		return false
	}
	for i := range arr1 {
		if arr1[i] != arr2[i] {
			return false
		}
	}
	return true
}
