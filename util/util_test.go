package util

import (
	"testing"

	"github.com/franela/goblin"
)

func Test_utils(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Tests for util/util.go", func() {

		g.Describe("The SliceContains function", func() {
			g.It("should return true if string is in string slice", func() {
				s := "test"
				list := []string{"this", "is", "a", "test"}
				g.Assert(SliceContains(s, list)).IsTrue()
			})

			g.It("should return false if string is not in string slice", func() {
				s := "bad string"
				list := []string{"this", "is", "a", "test"}
				g.Assert(SliceContains(s, list)).IsFalse()
			})
		})

		g.Describe("The TrimSuffix function", func() {
			g.It("should return a string with a matching suffix trimmed", func() {
				s := "file/path/"
				suffix := "/"
				exp := "file/path"
				actual := TrimSuffix(s, suffix)
				g.Assert(exp == actual).IsTrue()
			})

			g.It("should return a string without a matching suffix untrimmed", func() {
				s := "file/path/"
				suffix := "&"
				exp := "file/path/"
				actual := TrimSuffix(s, suffix)
				g.Assert(exp == actual).IsTrue()
			})
		})
	})
}
