package tagcloud

import "sort"

// TagCloud aggregates statistics about used tags
type TagCloud struct {
	Stat map[string]int
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
	return &TagCloud{Stat}
}

// AddTag should add a tag to the cloud if it wasn't present and increase tag occurrence count
// thread-safety is not needed
// TODO: You decide whether receiver should be a pointer or a value
func (tagCloud TagCloud) AddTag(tag string) {
	if _, hasValue := tagCloud.Stat[tag]; hasValue {
		tagCloud.Stat[tag] += 1
	} else {
		tagCloud.Stat[tag] = 1
	}
}

// TopN should return top N most frequent tags ordered in descending order by occurrence count
// if there are multiple tags with the same occurrence count then the order is defined by implementation
// if n is greater that TagCloud size then all elements should be returned
// thread-safety is not needed
// there are no restrictions on time complexity
// TODO: You decide whether receiver should be a pointer or a value
func (tagCloud TagCloud) TopN(n int) []TagStat {
	topN := make([]TagStat, 0, n)
	for tag, count := range tagCloud.Stat {
		node := TagStat{tag, count}
		topN = append(topN, node)
	}
	sort.SliceStable(topN, func(i, j int) bool {
		return topN[i].OccurrenceCount > topN[j].OccurrenceCount
	})
	if n >= len(topN) {
		return topN
	}
	return topN[:n]
}
