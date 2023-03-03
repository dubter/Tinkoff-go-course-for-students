package tagcloud

import "sort"

// TagCloud aggregates statistics about used tags
type TagCloud struct {
	IndexesInSlice map[string]int
	SortedTags     []TagStat
}

// TagStat represents statistics regarding single tag
type TagStat struct {
	Tag             string
	OccurrenceCount int
}

// New should create a valid TagCloud instance
// TODO: You decide whether this function should return a pointer or a value
func New() *TagCloud {
	Stat := map[string]int{}
	var SortedTags []TagStat
	return &TagCloud{Stat, SortedTags}
}

// AddTag should add a tag to the cloud if it wasn't present and increase tag occurrence count
// thread-safety is not needed
// TODO: You decide whether receiver should be a pointer or a value
func (tagCloud *TagCloud) AddTag(tag string) {
	if _, hasValue := tagCloud.IndexesInSlice[tag]; hasValue {
		tagCloud.SortedTags[tagCloud.IndexesInSlice[tag]].OccurrenceCount += 1
		sort.SliceStable(tagCloud.SortedTags, func(i, j int) bool {
			return tagCloud.SortedTags[i].OccurrenceCount > tagCloud.SortedTags[j].OccurrenceCount
		})
	} else {
		tagCloud.IndexesInSlice[tag] = len(tagCloud.SortedTags)
		tagCloud.SortedTags = append(tagCloud.SortedTags, TagStat{tag, 1})
	}
}

// TopN should return top N most frequent tags ordered in descending order by occurrence count
// if there are multiple tags with the same occurrence count then the order is defined by implementation
// if n is greater that TagCloud size then all elements should be returned
// thread-safety is not needed
// there are no restrictions on time complexity
// TODO: You decide whether receiver should be a pointer or a value
func (tagCloud TagCloud) TopN(n int) []TagStat {
	if len(tagCloud.SortedTags) <= n {
		return tagCloud.SortedTags
	}
	return tagCloud.SortedTags[:n]
}
