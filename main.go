package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/fjukstad/luftkvalitet"
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"github.com/lucasb-eyer/go-colorful"
)

func main() {

	p, err := plot.New()
	if err != nil {
		fmt.Println(err)
		return
	}

	xticks := plot.TimeTicks{Format: "2006-01-02"}

	p.Title.Text = "Air quality in Tromsø"
	p.X.Label.Text = "Date"
	p.Y.Label.Text = "PM10 (µg/m³) "
	p.X.Tick.Marker = xticks

	stations, err := luftkvalitet.GetStations()
	if err != nil {
		fmt.Println(err)
		return
	}

	colors, err := generateColors(stations)
	if err != nil {
		fmt.Println(err)
		return
	}

	fromTime := time.Date(2016, time.September, 1, 0, 0, 0, 0, time.UTC)
	toTime := time.Date(2016, time.December, 4, 0, 0, 0, 0, time.UTC)

	f := luftkvalitet.Filter{
		Areas:      []string{"Tromsø"},
		Components: []string{"PM10"},
		FromTime:   fromTime,
		ToTime:     toTime,
	}

	res, err := luftkvalitet.GetHistorical(f)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, result := range res {

		pts := make(plotter.XYs, len(result.Measurements))

		for i, m := range result.Measurements {
			pts[i].X = float64(m.FromTime.Unix())
			pts[i].Y = m.Value
		}

		// Make a line plotter and set its style.
		l, err := plotter.NewLine(pts)
		if err != nil {
			panic(err)
		}

		l.LineStyle.Width = vg.Points(1)
		r, g, b := colors[result.Station.Station].RGB255()
		l.LineStyle.Color = color.RGBA{R: r, G: g, B: b}

		p.Add(l)
		p.Legend.Add(result.Station.Station, l)
		p.Legend.Top = true
	}

	p.Y.Min = 0
	p.Y.Max = 300

	if err := p.Save(16*vg.Inch, 16*vg.Inch, "plot.pdf"); err != nil {
		panic(err)
	}
}

func generateColors(stations []luftkvalitet.Station) (map[string]colorful.Color, error) {
	pal, err := colorful.HappyPalette(len(stations))
	if err != nil {
		return map[string]colorful.Color{}, err
	}
	colors := make(map[string]colorful.Color, len(stations))
	for i, s := range stations {
		colors[s.Station] = pal[i]
	}

	return colors, nil
}
