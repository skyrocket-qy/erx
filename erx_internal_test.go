package erx

import (
	"strings"
	"testing"
)

func TestGetCallStack(t *testing.T) {
	// We use a helper function to control the call stack depth.
	helper := func(skip int) []CallerInfo {
		// This function adds one level to the stack.
		// getCallStack is called from here.
		return getCallStack(skip)
	}

	// getCallStack has a default skip of 2.
	// Stack: runtime.Callers (0), getCallStack (1), helper (2), TestGetCallStack (3)
	// With a skip of 3, we should skip getCallStack, helper, and land on TestGetCallStack.
	callerInfos := helper(3)

	if len(callerInfos) == 0 {
		t.Fatal("expected caller info, got none")
	}

	// The first frame should be TestGetCallStack.
	if !strings.Contains(callerInfos[0].Function, "erx.TestGetCallStack") {
		t.Errorf("expected function name to contain 'erx.TestGetCallStack', got '%s'", callerInfos[0].Function)
	}
	// It should not be the helper function.
	if strings.Contains(callerInfos[0].Function, "erx.TestGetCallStack.func1") {
		t.Errorf("did not expect function name to contain 'erx.TestGetCallStack.func1', got '%s'", callerInfos[0].Function)
	}
}
