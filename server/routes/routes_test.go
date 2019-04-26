package routes

import (
	"os"
	"testing"

	"github.com/go-park-mail-ru/2019_1_three_in_a_boat/shared/test-utils"
)

// Launches all tests between the calls to test_utils.SetUp and
// test_utils.TearDown. Tests should be launched using this function, if used
// outside of the go test framework.
func TestMain(m *testing.M) {
	test_utils.SetUp()
	code := m.Run()
	test_utils.TearDown()
	os.Exit(code)
}
