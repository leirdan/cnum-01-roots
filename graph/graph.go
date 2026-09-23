package graph

import (
	"cnum/types"
	"fmt"
	"image/color"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func Plot(inst types.Problem, intervals types.RootIntervalList, dir string) error {
	full := fmt.Sprintf("%s/f%d.png", dir, inst.Id)
	if err := plotRange(inst, inst.Start, inst.End, intervals, full); err != nil {
		return err
	}

	if len(intervals) == 0 {
		return nil
	}

	xMin, xMax := intervals[0][0], intervals[0][len(intervals[0])-1]
	for _, interval := range intervals[1:] {
		if interval[0] < xMin {
			xMin = interval[0]
		}
		if last := interval[len(interval)-1]; last > xMax {
			xMax = last
		}
	}
	pad := (xMax - xMin) * 0.5
	if pad == 0 {
		pad = inst.H
	}

	zoom := fmt.Sprintf("%s/f%d_zoom.png", dir, inst.Id)
	return plotRange(inst, xMin-pad, xMax+pad, intervals, zoom)
}

func plotRange(inst types.Problem, xMin, xMax float64, intervals types.RootIntervalList, path string) error {
	p := plot.New()
	p.Title.Text = fmt.Sprintf("Função %d", inst.Id)
	p.X.Label.Text = "x"
	p.Y.Label.Text = "f(x)"
	p.Add(plotter.NewGrid())

	const samples = 500
	points := make(plotter.XYs, samples)
	step := (xMax - xMin) / float64(samples-1)
	for i := range points {
		x := xMin + float64(i)*step
		points[i] = plotter.XY{X: x, Y: inst.Function(x)}
	}

	fn, err := plotter.NewLine(points)
	if err != nil {
		return err
	}
	fn.Color = color.RGBA{B: 200, A: 255}
	p.Add(fn)
	p.Legend.Add("f(x)", fn)

	zero, err := plotter.NewLine(plotter.XYs{{X: xMin, Y: 0}, {X: xMax, Y: 0}})
	if err != nil {
		return err
	}
	zero.Color = color.Gray{Y: 150}
	zero.Dashes = []vg.Length{vg.Points(3), vg.Points(3)}
	p.Add(zero)

	zeroLabel, err := plotter.NewLabels(plotter.XYLabels{
		XYs:    plotter.XYs{{X: xMin, Y: 0}},
		Labels: []string{"y = 0"},
	})
	if err != nil {
		return err
	}
	zeroLabel.Offset = vg.Point{X: vg.Points(4), Y: vg.Points(4)}
	zeroLabel.TextStyle[0].Color = color.Gray{Y: 100}
	p.Add(zeroLabel)

	if len(intervals) > 0 {
		bounds := make(plotter.XYs, 0, 2*len(intervals))
		for _, interval := range intervals {
			bounds = append(bounds,
				plotter.XY{X: interval[0], Y: 0},
				plotter.XY{X: interval[len(interval)-1], Y: 0},
			)
		}
		marks, err := plotter.NewScatter(bounds)
		if err != nil {
			return err
		}
		marks.Color = color.RGBA{R: 220, A: 255}
		marks.Radius = vg.Points(3)
		p.Add(marks)
		p.Legend.Add("intervalo isolado", marks)
	}

	return p.Save(6*vg.Inch, 4*vg.Inch, path)
}
