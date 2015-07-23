package router

import "testing"

func TestStaticSegment(t *testing.T) {
	var s segment = &staticSegment{
		value: "sample/segment",
	}
	var s_equal segment = &staticSegment{
		value: "sample/segment",
	}

	if !s.Is(s_equal) {
		t.Errorf("Segments %v and %v should be equals by segment.Is\n", s, s_equal)
	}

	if s.Params() != 0 {
		t.Errorf("These segments shouldn't have any parameters\n")
	}

	expected_path := "/abc"
	remaining_path, _, matched := s.Match("sample/segment/abc")
	if remaining_path != expected_path {
		t.Errorf("Expected '%s' from segment.Match and '%s' returned.\n", remaining_path, expected_path)
	}
	if !matched {
		t.Errorf("segment.Match with value '%s' does not match '%s'.\n", s.(*staticSegment).value, expected_path)
	}
}

func TestParamSegment(t *testing.T) {
	var s segment = &paramSegment{
		name: "test",
	}
	var s_equal segment = &paramSegment{
		name: "test",
	}

	if s.Params() != 1 {
		t.Errorf("These segments should have exactly one parameter\n")
	}

	if !s.Is(s_equal) {
		t.Errorf("Segments %v and %v should be equals by segment.Is\n", s, s_equal)
	}

	expected_path := "/abc"
	expected_param := param{name: "test", value: "segment"}
	remaining_path, param, matched := s.Match("segment/abc")
	if remaining_path != expected_path {
		t.Errorf("Expected '%s' from segment.Match and '%s' returned.\n", remaining_path, expected_path)
	}
	if param.name != expected_param.name || param.value != expected_param.value {
		t.Errorf("Expected param '%v' from segment.Match and '%v' returned.\n", expected_param, param)
	}
	if !matched {
		t.Errorf("segment.Match with value '%s' does not match '%s'.\n", s.(*staticSegment).value, expected_path)
	}
}
