package router

import "regexp"

type segment interface {
	Match(string) (string, Param, bool)
	Is(segment) bool
	Params() uint8
}

func parseSegments(path string) []segment {
	var segments = make([]segment, 0)
	var i = 0

	for i < len(path) {
		if path[i] == ':' {
			var j = i
			for j < len(path) && path[j] != '/' {
				j++
			}
			segments = append(segments, &paramSegment{name: path[i+1 : j]})
			i = j - 1
		} else if path[i] == '*' {
			segments = append(segments, &catchAllSegment{name: path[i+1:]})
			i = len(path) - 1
		} else {
			segments = append(segments, &staticSegment{value: string(path[i])})
		}
		i++
	}

	return segments
}

type staticSegment struct {
	value string
}

func (s *staticSegment) Match(path string) (string, Param, bool) {
	//fmt.Println(path, s.value)
	return path[len(s.value):], Param{}, string(path[:len(s.value)]) == s.value
}

func (s *staticSegment) Is(t segment) bool {
	return s.value == t.(*staticSegment).value
}

func (s *staticSegment) Params() uint8 {
	return 0
}

type paramSegment struct {
	name string
}

func (s *paramSegment) Match(path string) (string, Param, bool) {
	var i = 0
	for i < len(path) && path[i] != '/' {
		i++
	}
	//fmt.Println(path, s.name, path[:i])
	return path[i:], Param{Name: s.name, Value: path[:i]}, true
}

func (s *paramSegment) Is(t segment) bool {
	return s.name == t.(*paramSegment).name
}

func (s *paramSegment) Params() uint8 {
	return 1
}

type catchAllSegment struct {
	name string
}

func (s *catchAllSegment) Match(path string) (string, Param, bool) {
	return "", Param{Name: s.name, Value: path}, true
}

func (s *catchAllSegment) Is(t segment) bool {
	return s.name == t.(*catchAllSegment).name
}

func (s *catchAllSegment) Params() uint8 {
	return 1
}

type regexpSegment struct {
	regex regexp.Regexp
}

func (s *regexpSegment) Match(path string) (string, Param, bool) {
	return path[1:], Param{}, true
}
