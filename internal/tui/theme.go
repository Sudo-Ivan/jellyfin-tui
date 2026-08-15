package tui

// Theme is the oxide palette. Names are materials, not UI roles, so screens
// pick meaning at the call site instead of a generic primary/secondary map.
type Theme struct {
	Void     Color
	Basin    Color
	Well     Color
	Ridge    Color
	Filament Color
	Dust     Color
	Ember    Color
	Arc      Color
	Signal   Color
	Fault    Color
	Bloom    Color
}

const (
	hexVoid     = 0x0B0D10
	hexBasin    = 0x13161C
	hexWell     = 0x1A1F28
	hexRidge    = 0x2A313C
	hexFilament = 0xD6CFC4
	hexDust     = 0x6E675C
	hexEmber    = 0xE08A5A
	hexArc      = 0x6FC3C9
	hexSignal   = 0xB7C86A
	hexFault    = 0xD05C5C
	hexBloom    = 0xC48AC4
)

func Oxide() Theme {
	return Theme{
		Void:     colorHex(hexVoid),
		Basin:    colorHex(hexBasin),
		Well:     colorHex(hexWell),
		Ridge:    colorHex(hexRidge),
		Filament: colorHex(hexFilament),
		Dust:     colorHex(hexDust),
		Ember:    colorHex(hexEmber),
		Arc:      colorHex(hexArc),
		Signal:   colorHex(hexSignal),
		Fault:    colorHex(hexFault),
		Bloom:    colorHex(hexBloom),
	}
}

func (t Theme) Body() Style     { return Style{Fg: t.Filament, Bg: t.Void} }
func (t Theme) Panel() Style    { return Style{Fg: t.Filament, Bg: t.Basin} }
func (t Theme) Muted() Style    { return Style{Fg: t.Dust, Bg: t.Void} }
func (t Theme) Accent() Style   { return Style{Fg: t.Ember, Bg: t.Void} }
func (t Theme) Focus() Style    { return Style{Fg: t.Void, Bg: t.Ember} }
func (t Theme) Cursor() Style   { return Style{Fg: t.Void, Bg: t.Arc} }
func (t Theme) Title() Style    { return Style{Fg: t.Ember, Bg: t.Void, Attr: AttrBold} }
func (t Theme) Chrome() Style   { return Style{Fg: t.Ridge, Bg: t.Void} }
func (t Theme) Playing() Style  { return Style{Fg: t.Signal, Bg: t.Void} }
func (t Theme) Danger() Style   { return Style{Fg: t.Fault, Bg: t.Void} }
func (t Theme) BarTrack() Style { return Style{Fg: t.Ridge, Bg: t.Well} }
func (t Theme) BarFill() Style  { return Style{Fg: t.Ember, Bg: t.Well} }
