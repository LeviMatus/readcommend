package era

import (
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	exitCode := m.Run()

	if exitCode == 0 && testing.CoverMode() != "" {
		coverage := testing.Coverage()
		if coverage < 0.9 {
			log.Printf("Tests passed but only %.2f%% of tests lines are covered - require at least 90%%\n", coverage*100)
			exitCode = 1
		}
	}
	os.Exit(exitCode)
}
