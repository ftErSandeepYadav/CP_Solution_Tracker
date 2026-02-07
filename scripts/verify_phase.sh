#!/bin/bash

# Verification script for TDD phases
# Usage: ./scripts/verify_phase.sh <phase_number>

set -e

PHASE=$1
PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"

if [ -z "$PHASE" ]; then
    echo "Usage: $0 <phase_number>"
    echo "Example: $0 1"
    exit 1
fi

cd "$PROJECT_ROOT"

case $PHASE in
    1)
        echo "=========================================="
        echo "Verifying Phase 1: Platform Layer"
        echo "=========================================="
        echo ""
        
        echo "Running platform tests..."
        cd internal/platform
        go test -v -cover
        
        echo ""
        echo "✓ Phase 1 verification complete!"
        echo "Expected: Coverage should be >80%"
        ;;
        
    2)
        echo "=========================================="
        echo "Verifying Phase 2: Configuration Layer"
        echo "=========================================="
        echo ""
        
        echo "Running config tests..."
        cd internal/config
        go test -v -cover
        
        echo ""
        echo "Creating test config file..."
        cat > /tmp/test-config.json << 'EOF'
{
  "platforms": {
    "codeforces": {"enabled": true, "handle": "test"}
  },
  "repository": {
    "owner": "test", "name": "repo", "branch": "main",
    "token": "token", "use_github": true
  },
  "sync": {
    "auto_sync": false, "conflict_strategy": "keep_latest",
    "state_file": ".state.json", "max_concurrency": 5,
    "organize_by_contest": true
  }
}
EOF
        
        echo "Testing config loading..."
        cd "$PROJECT_ROOT"
        go run cmd/test_config/main.go
        
        echo ""
        echo "✓ Phase 2 verification complete!"
        echo "Expected: All tests pass with >70% coverage"
        ;;
        
    *)
        echo "Error: Unknown phase '$PHASE'"
        echo "Available phases: 1, 2"
        exit 1
        ;;
esac
