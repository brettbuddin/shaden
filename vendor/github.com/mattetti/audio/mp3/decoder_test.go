package mp3_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/mattetti/audio/mp3"
)

func Test_SeemsValid(t *testing.T) {
	testCases := []struct {
		input   string
		isValid bool
	}{
		{"fixtures/frame.mp3", false},
		{"fixtures/HousyStab.mp3", true},
		{"../wav/fixtures/bass.wav", false},
		{"fixtures/nullbytes.mp3", true},
		{"fixtures/idv3-24.mp3", true},
	}

	for i, tc := range testCases {
		t.Logf("test case %d - %s\n", i, tc.input)
		f, err := os.Open(tc.input)
		if err != nil {
			panic(err)
		}
		if o := mp3.SeemsValid(f); o != tc.isValid {
			t.Fatalf("expected %t\ngot\n%t\n", tc.isValid, o)
		}
		f.Close()
	}
}

func Test_Decoder_Duration(t *testing.T) {
	testCases := []struct {
		input    string
		duration string
	}{
		{"fixtures/HousyStab.mp3", "16.483264688s"},
		{"fixtures/slayer.mp3", "28.447345872s"},
		{"fixtures/nullbytes.mp3", "13.505305616s"},
		{"fixtures/idv3-24.mp3", "11.128162848s"},
		{"fixtures/weird_duration.mp3", "5.714122448s"},
	}

	for i, tc := range testCases {
		t.Logf("duration test case %d - %s\n", i, tc.input)
		f, err := os.Open(tc.input)
		if err != nil {
			panic(err)
		}
		d := mp3.New(f)
		dur, err := d.Duration()
		if err != nil {
			t.Fatal(err)
		}
		if o := fmt.Sprintf("%s", dur); o != tc.duration {
			t.Fatalf("expected %s\ngot\n%s\n", tc.duration, o)
		}
	}
}

func ExampleDecoder_Duration() {
	f, err := os.Open("fixtures/HousyStab.mp3")
	if err != nil {
		panic(err)
	}
	d := mp3.New(f)
	dur, err := d.Duration()
	if err != nil {
		panic(err)
	}
	fmt.Println(dur)
	//Output: 16.483264688s
}
