package typer

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type AEvent struct {}

func TestIdentify(t *testing.T) {
	assertThat := require.New(t)

	t.Run("Given a AEvent variable", func(t *testing.T) {
		event := AEvent{}
		t.Run("When its type is identified", func(t *testing.T) {
			result := Identify(event)
			t.Run("Then a 'AEvent' name is returned", func(t *testing.T) {
				assertThat.Equal("AEvent", result)
			})
		})
	})
	t.Run("Given a pointer to an AEvent variable", func(t *testing.T) {
		event := new(AEvent)
		t.Run("When its type is identified", func(t *testing.T) {
			result := Identify(event)
			t.Run("Then a 'AEvent' name is returned", func(t *testing.T) {
				assertThat.Equal("AEvent", result)
			})
		})
	})
}
