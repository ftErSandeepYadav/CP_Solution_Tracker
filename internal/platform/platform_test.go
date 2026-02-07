package platform

import (
    "testing"
    "time"
)

func TestSubmission_Creation(t *testing.T) {
    sub := Submission{
        ID:             "12345",
        ProblemID:      "1000A",
        ProblemName:    "Theatre Square",
        ContestID:      "1000",
        Verdict:        "OK",
        Language:       "GNU C++17",
        Code:           "#include<iostream>",
        SubmissionTime: time.Now(),
    }
    
    // Test basic fields
    if sub.ID != "12345" {
        t.Errorf("Expected ID '12345', got '%s'", sub.ID)
    }
    
    if sub.ProblemID != "1000A" {
        t.Errorf("Expected ProblemID '1000A', got '%s'", sub.ProblemID)
    }
    
    if sub.Verdict != "OK" {
        t.Errorf("Expected Verdict 'OK', got '%s'", sub.Verdict)
    }
}

func TestProblem_Creation(t *testing.T) {
    prob := Problem{
        ID:         "1000A",
        Name:       "Theatre Square",
        ContestID:  "1000",
        Index:      "A",
        Tags:       []string{"math", "implementation"},
        Difficulty: 800,
        URL:        "https://codeforces.com/problemset/problem/1000/A",
    }
    
    if prob.ID != "1000A" {
        t.Errorf("Expected ID '1000A', got '%s'", prob.ID)
    }
    
    if len(prob.Tags) != 2 {
        t.Errorf("Expected 2 tags, got %d", len(prob.Tags))
    }
    
    if prob.Difficulty != 800 {
        t.Errorf("Expected difficulty 800, got %d", prob.Difficulty)
    }
}

func TestConfig_Creation(t *testing.T) {
    cfg := Config{
        Handle:  "testuser",
        Cookies: "session=abc123",
        Extra:   make(map[string]string),
    }
    
    cfg.Extra["custom_field"] = "custom_value"
    
    if cfg.Handle != "testuser" {
        t.Errorf("Expected handle 'testuser', got '%s'", cfg.Handle)
    }
    
    if cfg.Extra["custom_field"] != "custom_value" {
        t.Error("Extra field not set correctly")
    }
}

// Test that Platform interface can be implemented
type MockPlatform struct{}

func (m *MockPlatform) GetName() string {
    return "mock"
}

func (m *MockPlatform) GetAcceptedSubmissions(handle string) ([]Submission, error) {
    return []Submission{}, nil
}

func (m *MockPlatform) GetSubmissionCode(submissionID string) (string, error) {
    return "", nil
}

func (m *MockPlatform) GetProblemMetadata(problemID string) (*Problem, error) {
    return nil, nil
}

func TestPlatformInterface(t *testing.T) {
    var platform Platform = &MockPlatform{}
    
    name := platform.GetName()
    if name != "mock" {
        t.Errorf("Expected name 'mock', got '%s'", name)
    }
    
    // Should not panic
    _, _ = platform.GetAcceptedSubmissions("test")
    _, _ = platform.GetSubmissionCode("123")
    _, _ = platform.GetProblemMetadata("1000A")
}