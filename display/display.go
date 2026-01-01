package display

type Options struct {
	ShowWords bool
	ShowLines bool
	ShowBytes bool
}

func (d *Options) ShouldShowWords() bool {
	if !d.ShowBytes && !d.ShowLines && !d.ShowWords {
		return true
	}
	return d.ShowWords
}

func (d *Options) ShouldShowLines() bool {
	if !d.ShowBytes && !d.ShowLines && !d.ShowWords {
		return true
	}
	return d.ShowLines
}

func (d *Options) ShouldShowBytes() bool {
	if !d.ShowBytes && !d.ShowLines && !d.ShowWords {
		return true
	}
	return d.ShowBytes
}
