package dockref

import (
	_ "crypto/sha256" // side effect: register sha256
	"github.com/docker/distribution/reference"
	"github.com/opencontainers/go-digest"
)

func FromOriginal(original string) (ref Reference, e error) {
	r, e := reference.ParseAnyReference(original)
	if e != nil {
		return
	}

	var name string
	var domain string
	var path string
	var named reference.Named
	var ok bool
	if named, ok = r.(reference.Named); ok {
		name = named.Name()
		domain = reference.Domain(named)
		path = reference.Path(named)
	}

	var tag string
	if tagged, ok := r.(reference.Tagged); ok {
		tag = tagged.Tag()
	}

	var dig string
	if digested, ok := r.(reference.Digested); ok {
		dig = string(digested.Digest())
	}

	ref = dockref{
		original: original,
		domain: domain,
		name:     name,
		tag:      tag,
		digest:   dig,
		path: path,
		named: named,
	}
	return
}

type Reference interface {
	Name() string
	Tag() string
	DigestString() string
	Digest() digest.Digest
	Original() string
	Domain() string
	Path() string
	Named() reference.Named
}

var _ Reference = (*dockref)(nil)

type dockref struct {
	name     string
	original string
	tag      string
	digest   string
	domain   string
	path     string
	named    reference.Named
}

func (r dockref) Named() reference.Named {
	return r.named
}

func (r dockref) Name() string {
	return r.name
}

func (r dockref) Tag() string {
	return r.tag
}

func (r dockref) DigestString() string {
	return r.digest
}

func (r dockref) Digest() digest.Digest {
	return digest.Digest(r.digest)
}

func (r dockref) Original() string {
	return r.original
}

func (r dockref) Domain() string {
	return r.domain
}

func (r dockref) Path() string {
	return r.path
}
