package calc

import "testing"

func TestSipBumpShortensYears(t *testing.T) {
	in := DefaultFIRE()
	base := FIRE(in)
	moves := YearMoves(in)
	if len(moves) != 3 {
		t.Fatalf("moves=%d", len(moves))
	}
	sip := moves[0]
	if sip.Kind != "sip" || sip.Amount != 5_000 {
		t.Fatalf("sip move: %+v", sip)
	}
	if !sip.Reaches || sip.Years >= base.Years {
		t.Fatalf("extra SIP should be faster: base=%v sip=%v", base.Years, sip.Years)
	}
	if sip.DeltaYears >= 0 {
		t.Fatalf("delta should be negative (earlier): %v", sip.DeltaYears)
	}
}

func TestLowerReturnLengthensYears(t *testing.T) {
	in := DefaultFIRE()
	base := FIRE(in)
	ret := YearMoves(in)[1]
	if ret.Kind != "return" {
		t.Fatalf("kind=%s", ret.Kind)
	}
	if !ret.Reaches || ret.Years <= base.Years {
		t.Fatalf("1%% lower return should take longer: base=%v ret=%v", base.Years, ret.Years)
	}
}

func TestHouseMoveChangesYears(t *testing.T) {
	in := DefaultFIRE()
	in.Housing = "rent"
	base := FIRE(in)
	house := YearMoves(in)[2]
	if house.Kind != "house" {
		t.Fatalf("kind=%s", house.Kind)
	}
	if !house.Reaches || house.Years <= base.Years {
		t.Fatalf("buying should take longer than renting: rent=%v buy=%v", base.Years, house.Years)
	}
}

func TestPayCutLengthensYears(t *testing.T) {
	in := DefaultFIRE()
	base := FIRE(in)
	opts := YearOptions(in)
	if len(opts) != 3 {
		t.Fatalf("options=%d", len(opts))
	}
	cut := opts[0]
	if cut.Kind != "paycut" || !cut.Reaches || cut.Years <= base.Years {
		t.Fatalf("40%% pay cut should take longer: base=%v cut=%v", base.Years, cut.Years)
	}
	pause := opts[2]
	if pause.Kind != "pause" || !pause.Reaches || pause.Years <= base.Years {
		t.Fatalf("two years off should take longer: base=%v pause=%v", base.Years, pause.Years)
	}
}
