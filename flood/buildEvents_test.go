package flood

import (
	"testing"
)

const (
	LONG_ACTING = "#la="
	SMBG        = "#bg="
)

func TestBuilderPatterns(t *testing.T) {

	patterns := make(map[Pattern]string)
	patterns[PatternBasal] = LONG_ACTING
	patterns[PatternSmbg] = SMBG

	bob := NewEventBuilder().AddPatterns(patterns)

	if bob.registerPatterns[PatternBasal] != LONG_ACTING {
		t.Fatalf("the pattern %s  should be %s ", PatternBasal, LONG_ACTING)
	}

	if bob.registerPatterns[PatternSmbg] != SMBG {
		t.Fatalf("the pattern %s  should be %s ", PatternSmbg, SMBG)
	}

	if bob.registerPatterns[PatternBolus] != "" {
		t.Fatalf("the pattern %s should be empty ", PatternBolus)
	}

}

func TestBuilderUnprocessedEvents(t *testing.T) {

	unprocessed := []UnprocessedEvent{UnprocessedEvent{Text: "#bg=9.9", Date: "1/1/2010", Device: "mine"}, UnprocessedEvent{Text: "#bg=5.5", Date: "1/2/2010", Device: "mine"}}

	bob := NewEventBuilder().AddUnprocessedEvents(unprocessed)

	if bob.unprocessed == nil {
		t.Fatal("there should be events to process")
	}

	if len(bob.unprocessed) != len(unprocessed) {
		t.Fatal("unprocessed length sould be the same")
	}

}

func TestBuilderProcess(t *testing.T) {

	unprocessed := []UnprocessedEvent{UnprocessedEvent{Text: "#bg=9.9", Date: "1/1/2010", Device: "mine"}, UnprocessedEvent{Text: "#bg=5.5", Date: "1/2/2010", Device: "mine"}}

	patterns := make(map[Pattern]string)
	patterns[PatternSmbg] = SMBG

	processed := NewEventBuilder().AddUnprocessedEvents(unprocessed).AddPatterns(patterns).Process()

	if processed == nil {
		t.Fatal("there should be event processed")
	}

	if len(processed) != 2 {
		t.Fatal("there should only be one event processed")
	}

}
