package dockref

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInvalid(t *testing.T) {
	ref, e := FromOriginal("invalid:reference:format")
	assert.Nil(t, ref)
	assert.Error(t, e)
}

func TestWellknownNames(t *testing.T) {
	originals := []string{"nginx", "alpine", "httpd"}
	for _, original := range originals {
		t.Run("Parses "+original, func(t *testing.T) {
			ref, e := FromOriginal(original)
			assert.Nil(t, e)
			assert.Equal(t, "", ref.Tag())

			expected := "docker.io/library/" + original
			assert.Equal(t, expected, ref.Name())
		})
	}
}

func TestWellknownTaggedNames(t *testing.T) {
	originals := []string{"nginx:latest", "nginx:1.15.2-alpine-perl", "mongo:3.4.16-windowsservercore-ltsc2016"}
	for _, original := range originals {
		t.Run("Parses "+original, func(t *testing.T) {
			ref, e := FromOriginal(original)
			assert.Nil(t, e)

			expected := "docker.io/library/" + original
			assert.Equal(t, expected, ref.Name()+":"+ref.Tag())
		})
	}
}

func TestOriginalsAreUnchanged(t *testing.T) {
	originals := []string{"nginx", "nginx:latest", "nginx:1.15.2-alpine-perl"}
	for _, original := range originals {
		t.Run(original+" remains "+original, func(t *testing.T) {
			ref, e := FromOriginal(original)
			assert.Nil(t, e)

			expected := original
			assert.Equal(t, expected, ref.Original())
		})
	}
}

func TestDigestOnly(t *testing.T) {
	originals := []string{"d21b79794850b4b15d8d332b451d95351d14c951542942a816eea69c9e04b240"}
	for _, original := range originals {
		t.Run("Parses "+original, func(t *testing.T) {
			ref, e := FromOriginal(original)
			assert.Nil(t, e)
			assert.Empty(t, ref.Name())
			assert.Empty(t, ref.Tag())

			expected := "sha256:" + original
			assert.Equal(t, expected, ref.DigestString())
			assert.Equal(t, expected, string(ref.Digest()))
		})
	}
}

func TestNameAndDigest(t *testing.T) {
	originals := []string{"nginx@sha256:d21b79794850b4b15d8d332b451d95351d14c951542942a816eea69c9e04b240"}
	for _, original := range originals {
		t.Run("Parses "+original, func(t *testing.T) {
			ref, e := FromOriginal(original)
			assert.Nil(t, e)
			assert.Equal(t, "docker.io/library/nginx", ref.Name())
			assert.Empty(t, ref.Tag())

			assert.Equal(t, "sha256:d21b79794850b4b15d8d332b451d95351d14c951542942a816eea69c9e04b240", ref.DigestString())
			assert.Equal(t, "sha256:d21b79794850b4b15d8d332b451d95351d14c951542942a816eea69c9e04b240", string(ref.Digest()))
		})
	}
}

func TestNameAndTagAndDigest(t *testing.T) {
	originals := []string{"nginx:1.2@sha256:d21b79794850b4b15d8d332b451d95351d14c951542942a816eea69c9e04b240"}
	for _, original := range originals {
		t.Run("Parses "+original, func(t *testing.T) {
			ref, e := FromOriginal(original)
			assert.Nil(t, e)
			assert.Equal(t, "docker.io/library/nginx", ref.Name())
			assert.Equal(t, "1.2", ref.Tag())

			assert.Equal(t, "sha256:d21b79794850b4b15d8d332b451d95351d14c951542942a816eea69c9e04b240", ref.DigestString())
			assert.Equal(t, "sha256:d21b79794850b4b15d8d332b451d95351d14c951542942a816eea69c9e04b240", string(ref.Digest()))
		})
	}
}

func TestDomainAndNameAndTagAndDigest(t *testing.T) {
	originals := []string{"nginx"}
	for _, original := range originals {
		t.Run("Parses "+original, func(t *testing.T) {
			ref, e := FromOriginal(original)
			assert.Nil(t, e)
			assert.Equal(t, "docker.io", ref.Domain())
			assert.Equal(t, "library/nginx", ref.Path())
		})
	}
	originals = []string{"my.com/nginx"}
	for _, original := range originals {
		t.Run("Parses "+original, func(t *testing.T) {
			ref, e := FromOriginal(original)
			assert.Nil(t, e)
			assert.Equal(t, "my.com", ref.Domain())
			assert.Equal(t, "nginx", ref.Path())
		})
	}
}

func TestDockref_Named(t *testing.T) {
	t.Run("Returns Named for named references", func(t *testing.T) {
		ref, e := FromOriginal("nginx")
		assert.Nil(t, e)
		assert.NotNil(t, ref.Named())
	})
	t.Run("Returns nil for unnamed references", func(t *testing.T) {
		ref, e := FromOriginal("d21b79794850b4b15d8d332b451d95351d14c951542942a816eea69c9e04b240")
		assert.Nil(t, e)
		assert.Nil(t, ref.Named())
	})
}