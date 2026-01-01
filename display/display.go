package display

type DisplayOpts struct {
	ShowWords bool
	ShowLines bool
	ShowBytes bool
}

func (d *DisplayOpts) ShouldShowWords() bool {
	if !d.ShowBytes && !d.ShowLines && !d.ShowWords {
		return true
	}
	return d.ShowWords
}

func (d *DisplayOpts) ShouldShowLines() bool {
	if !d.ShowBytes && !d.ShowLines && !d.ShowWords {
		return true
	}
	return d.ShowLines
}

func (d *DisplayOpts) ShouldShowBytes() bool {
	if !d.ShowBytes && !d.ShowLines && !d.ShowWords {
		return true
	}
	return d.ShowBytes
}
