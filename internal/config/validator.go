package config

import (
    "fmt"
    "strings"
)

// ValidationError represents a configuration validation error
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error in field '%s': %s", e.Field, e.Message)
}

// ValidationErrors holds multiple validation errors
type ValidationErrors []ValidationError

func (errs ValidationErrors) Error() string {
    var messages []string
    for _, err := range errs {
        messages = append(messages, err.Error())
    }
    return strings.Join(messages, "; ")
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
    var errors ValidationErrors
    
    // Validate platforms
    if len(c.Platforms) == 0 {
        errors = append(errors, ValidationError{
            Field:   "platforms",
            Message: "at least one platform must be configured",
        })
    }
    
    for name, platform := range c.Platforms {
        if platform.Enabled && platform.Handle == "" {
            errors = append(errors, ValidationError{
                Field:   fmt.Sprintf("platforms.%s.handle", name),
                Message: "handle is required when platform is enabled",
            })
        }
    }
    
    // Validate repository config (if using GitHub API)
    if c.Repository.UseGitHub {
        if c.Repository.Owner == "" {
            errors = append(errors, ValidationError{
                Field:   "repository.owner",
                Message: "owner is required when using GitHub API",
            })
        }
        
        if c.Repository.Name == "" {
            errors = append(errors, ValidationError{
                Field:   "repository.name",
                Message: "repository name is required when using GitHub API",
            })
        }
        
        if c.Repository.Token == "" {
            errors = append(errors, ValidationError{
                Field:   "repository.token",
                Message: "token is required when using GitHub API",
            })
        }
        
        if c.Repository.Branch == "" {
            errors = append(errors, ValidationError{
                Field:   "repository.branch",
                Message: "branch cannot be empty",
            })
        }
    } else {
        // Using local git
        if c.Repository.LocalPath == "" {
            errors = append(errors, ValidationError{
                Field:   "repository.local_path",
                Message: "local_path is required when not using GitHub API",
            })
        }
    }
    
    // Validate sync config
    validStrategies := []string{"keep_latest", "keep_all", "ask"}
    isValidStrategy := false
    for _, s := range validStrategies {
        if c.Sync.ConflictStrategy == s {
            isValidStrategy = true
            break
        }
    }
    
    if !isValidStrategy {
        errors = append(errors, ValidationError{
            Field: "sync.conflict_strategy",
            Message: fmt.Sprintf("must be one of: %s",
                strings.Join(validStrategies, ", ")),
        })
    }
    
    if c.Sync.StateFile == "" {
        errors = append(errors, ValidationError{
            Field:   "sync.state_file",
            Message: "state_file cannot be empty",
        })
    }
    
    if c.Sync.MaxConcurrency < 1 {
        errors = append(errors, ValidationError{
            Field:   "sync.max_concurrency",
            Message: "must be at least 1",
        })
    }
    
    if c.Sync.MaxConcurrency > 20 {
        errors = append(errors, ValidationError{
            Field:   "sync.max_concurrency",
            Message: "should not exceed 20 to avoid rate limiting",
        })
    }
    
    if len(errors) > 0 {
        return errors
    }
    
    return nil
}