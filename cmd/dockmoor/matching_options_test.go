package main

import (
	"fmt"
	"github.com/MeneDev/dockmoor/dockproc"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestEmptyPredicates(t *testing.T) {
	fo := &MatchingOptions{}
	err := verifyMatchOptions(fo)
	assert.Nil(t, err)
}

func applyPredicatesByName(fo *MatchingOptions, names ...string) {

	for _, name := range names {
		switch {
		case equalsAnyString("outdated", name):
			fo.TagPredicates.Outdated = true
		case equalsAnyString("unpinned", name):
			fo.DigestPredicates.Unpinned = true
		case equalsAnyString("latest", name):
			fo.TagPredicates.Latest = true
		case equalsAnyString("domains", name):
			fo.DomainPredicates.Domains = []string{"a", "b"}
		case equalsAnyString("names", name):
			fo.NamePredicates.Names = []string{"a", "b"}
		case equalsAnyString("untagged", name):
			fo.TagPredicates.Untagged = true
		case equalsAnyString("tags", name):
			fo.TagPredicates.Tags = []string{"a", "b"}
		case equalsAnyString("unpinned", name):
			fo.DigestPredicates.Unpinned = true
		case equalsAnyString("digests", name):
			fo.DigestPredicates.Digests = []string{"a", "b"}
		default:
			panic(fmt.Sprintf("Unknown predicate name '%s'", names))
		}
	}
}

func TestSinglePredicatesIsValid(t *testing.T) {
	for _, a := range predicateNames {
		t.Run(a, func(t *testing.T) {

			fo := &MatchingOptions{}
			applyPredicatesByName(fo, a)

			err := verifyMatchOptions(fo)
			assert.Nil(t, err)
		})
	}
}

func TestOnePredicateCanBeCombinedWithOneFromOtherGroup(t *testing.T) {
	for iGroupA, groupA := range predicateGroups {
		for iGroupB, groupB := range predicateGroups {
			if iGroupA == iGroupB {
				continue
			}
			for _, a := range groupA {
				for _, b := range groupB {

					t.Run(fmt.Sprintf("%s and %s", a, b), func(t *testing.T) {
						fo := &MatchingOptions{}
						applyPredicatesByName(fo, a, b)

						err := verifyMatchOptions(fo)
						assert.Nil(t, err)
					})
				}
			}
		}
	}
}

func TestMultipleFromSameGroupFail(t *testing.T) {

	for groupName, group := range predicateGroups {
		if len(group) <= 1 {
			continue
		}

		for _, a := range group {
			for _, b := range group {
				if a == b {
					continue
				}

				t.Run(a+" and "+b, func(t *testing.T) {
					fo := &MatchingOptions{}
					applyPredicatesByName(fo, a, b)

					err := verifyMatchOptions(fo)
					assert.Error(t, err)
					assert.Equal(t, ErrAtMostOnePredicate[groupName], err)
				})
			}
		}
	}
}

func TestNonExclusivePredicatesCanBeCombined(t *testing.T) {
	fo := &MatchingOptions{}

	for _, domain := range [][]string{nil, {"a", "b"}} {
		fo.DomainPredicates.Domains = domain
		for _, name := range [][]string{nil, {"a", "b"}} {
			fo.NamePredicates.Names = name
			for _, tag := range [][]string{nil, {"a", "b"}} {
				fo.TagPredicates.Tags = tag
				for _, digest := range []string{"unpinned"} {
					unpinned := digest == "unpinned"
					fo.DigestPredicates.Unpinned = unpinned

					testCase := fmt.Sprintf("predicates can be combined: %s/%s:%s@%t", domain, name, tag, unpinned)
					t.Run(testCase, func(t *testing.T) {
						err := verifyMatchOptions(fo)
						assert.Nil(t, err)
					})
				}
			}
		}
	}
}

func TestTagPredicateCombinations(t *testing.T) {
	fo := &MatchingOptions{}

	for _, untagged := range []bool{true, false} {
		for _, latest := range []bool{true, false} {
			for _, outdated := range []bool{true, false} {
				for _, tags := range [][]string{nil, {"a", "b"}} {
					set := 0
					if untagged {
						set++
					}
					if latest {
						set++
					}
					if outdated {
						set++
					}
					if tags != nil {
						set++
					}

					fo.TagPredicates.Untagged = untagged
					fo.TagPredicates.Latest = latest
					fo.TagPredicates.Outdated = outdated
					fo.TagPredicates.Tags = tags

					testCase := fmt.Sprintf("domain tag combination: untagged:%t-latest:%t-outdated%t-tags:%s", untagged, latest, outdated, tags)
					t.Run(testCase, func(t *testing.T) {
						err := verifyMatchOptions(fo)
						switch {
						case set > 1:
							assert.Error(t, err)
						case set <= 1:
							assert.Nil(t, err)
						}
					})
				}
			}
		}
	}
}

func TestDigestPredicateCombinations(t *testing.T) {
	fo := &MatchingOptions{}

	for _, unpinned := range []bool{true, false} {
		for _, digests := range [][]string{nil, {"a", "b"}} {
			set := 0
			if unpinned {
				set++
			}
			if digests != nil {
				set++
			}

			fo.DigestPredicates.Unpinned = unpinned
			fo.DigestPredicates.Digests = digests

			testCase := fmt.Sprintf("domain tag combination: unpinned:%t-digests:%s", unpinned, digests)
			t.Run(testCase, func(t *testing.T) {
				err := verifyMatchOptions(fo)
				switch {
				case set > 1:
					assert.Error(t, err)
				case set == 1:
					assert.Nil(t, err)
				}
			})
		}
	}
}


func TestAnyPredicateWhenNoFlagWithContains(t *testing.T) {
	fo := &MatchingOptions{}

	predicate := fo.getPredicate()

	assert.IsType(t, dockproc.AnyPredicateNew(), predicate)
}

func TestDomainsPredicateWhenDomainsSet(t *testing.T) {
	fo := &MatchingOptions{}
	fo.DomainPredicates.Domains = []string {"a", "b"}

	predicate := fo.getPredicate()

	assert.IsType(t, dockproc.DomainsPredicateNew([]string {"a", "b"}), predicate)
}

func TestNamesPredicateWhenNamesSet(t *testing.T) {
	fo := &MatchingOptions{}
	fo.NamePredicates.Names = []string {"a", "b"}

	predicate := fo.getPredicate()

	assert.IsType(t, dockproc.NamesPredicateNew([]string {"a", "b"}), predicate)
}

func TestTagsPredicateWhenTagsSet(t *testing.T) {
	fo := &MatchingOptions{}
	fo.TagPredicates.Tags = []string {"a", "b"}

	predicate := fo.getPredicate()

	assert.IsType(t, dockproc.TagsPredicateNew([]string {"a", "b"}), predicate)
}

func TestUntaggedPredicateWhenUntaggedSet(t *testing.T) {
	fo := &MatchingOptions{}
	fo.TagPredicates.Untagged = true

	predicate := fo.getPredicate()

	assert.IsType(t, dockproc.UntaggedPredicateNew(), predicate)
}

//func TestOutdatedPredicateWhenOutdatedSet(t *testing.T) {
//	fo := &MatchingOptions{}
//	fo.TagPredicates.Outdated = true
//
//	predicate := fo.getPredicate()
//
//	assert.IsType(t, dockproc.OutdatedPredicateNew(), predicate)
//}

func TestLatestPredicateWhenLatestSet(t *testing.T) {
	fo := &MatchingOptions{}
	fo.TagPredicates.Latest = true

	predicate := fo.getPredicate()

	assert.IsType(t, dockproc.LatestPredicateNew(), predicate)
}

func TestDigestsPredicateWhenDomainSet(t *testing.T) {
	fo := &MatchingOptions{}
	fo.DigestPredicates.Digests = []string {"a", "b"}

	predicate := fo.getPredicate()

	assert.IsType(t, dockproc.DigestsPredicateNew([]string {"a", "b"}), predicate)
}

func TestUnpinnedPredicateWhenUnpinnedSet(t *testing.T) {
	fo := &MatchingOptions{}
	fo.DigestPredicates.Unpinned = true

	predicate := fo.getPredicate()

	assert.IsType(t, dockproc.UnpinnedPredicateNew(), predicate)
}

func TestAndPredicateWhenUnpinnedAndLatestSet(t *testing.T) {
	fo := &MatchingOptions{}
	fo.DigestPredicates.Unpinned = true
	fo.TagPredicates.Latest = true

	predicate := fo.getPredicate()

	assert.IsType(t, dockproc.AndPredicateNew(nil), predicate)

	andPredicate, _ := predicate.(dockproc.AndPredicate)
	andPredicate.Predicates()

	assert.IsType(t, dockproc.AndPredicateNew(nil), predicate)

	matches := 0
	for _, p := range andPredicate.Predicates() {
		t := reflect.TypeOf(p)
		switch t {
		case reflect.TypeOf(dockproc.UnpinnedPredicateNew()):
			fallthrough
		case reflect.TypeOf(dockproc.LatestPredicateNew()):
			matches++
		}
	}

	assert.Equal(t, 2, matches)
}
