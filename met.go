package main

import (
	"fmt"

	"github.com/fjukstad/met"
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
)

func main() {

	p, err := plot.New()
	if err != nil {
		fmt.Println(err)
		return
	}

	xticks := plot.TimeTicks{Format: "2006-01-02"}

	p.Title.Text = "Precipitation in Tromsø"
	p.X.Label.Text = "Date"
	p.Y.Label.Text = "mm"
	p.X.Tick.Marker = xticks

	f := met.Filter{
		Sources:       []string{"SN90451"},
		ReferenceTime: "2016-01-01T00:00:00.000Z/2016-06-01T00:00:00.000Z",
		Elements:      []string{"sum(precipitation_amount T1H)"},
	}

	data, err := met.GetObservations(f)
	if err != nil {
		fmt.Println(err)
	}

	pts := make(plotter.XYs, len(data))

	for i, d := range data {
		precip := d.Observations[0].Value
		pts[i].X = float64(d.ReferenceTime.Unix())
		pts[i].Y = precip
	}

	l, err := plotter.NewLine(pts)
	if err != nil {
		panic(err)
	}
	p.Add(l)

	if err := p.Save(16*vg.Inch, 16*vg.Inch, "met.pdf"); err != nil {
		panic(err)
	}

}
