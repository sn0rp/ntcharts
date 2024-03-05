// Package graph contains functions to help draw runes on to a canvas
package graph

// https://en.wikipedia.org/wiki/Bresenham%27s_line_algorithm
// https://en.wikipedia.org/wiki/Midpoint_circle_algorithm

import (
	"math"
	"sort"

	"github.com/NimbleMarkets/bubbletea-charts/canvas"
	"github.com/NimbleMarkets/bubbletea-charts/canvas/runes"
	"github.com/charmbracelet/lipgloss"
)

// DrawVerticalLineUp draws a vertical line going up starting from (X,Y) coordinates.
// Applies given style to all runes.
// Coordinates (0,0) is top left of canvas.
func DrawVerticalLineUp(m *canvas.Model, p canvas.Point, s lipgloss.Style) {
	x := p.X
	r := canvas.NewCell(runes.LineVertical, s)
	for i := p.Y; i >= 0; i-- {
		m.SetCell(canvas.Point{x, i}, r)
	}
}

// DrawVerticalLineDown draws a vertical line going down starting from (X,Y) coordinates.
// Applies given style to all runes.
// Coordinates (0,0) is top left of canvas.
func DrawVerticalLineDown(m *canvas.Model, p canvas.Point, s lipgloss.Style) {
	x := p.X
	r := canvas.NewCell(runes.LineVertical, s)
	for i := p.Y; i < m.Height(); i++ {
		m.SetCell(canvas.Point{x, i}, r)
	}
}

// DrawHorizonalLineLeft draws a horizontal line going to the left starting from (X,Y) coordinates.
// Applies given style to all runes.
// Coordinates (0,0) is top left of canvas.
func DrawHorizonalLineLeft(m *canvas.Model, p canvas.Point, s lipgloss.Style) {
	y := p.Y
	r := canvas.NewCell(runes.LineHorizontal, s)
	for i := p.X; i >= 0; i-- {
		m.SetCell(canvas.Point{i, y}, r)
	}
}

// DrawHorizonalLineRight draws a horizontal line going to the right starting from (X,Y) coordinates.
// Applies given style to all runes.
// Coordinates (0,0) is top left of canvas.
func DrawHorizonalLineRight(m *canvas.Model, p canvas.Point, s lipgloss.Style) {
	y := p.Y
	r := canvas.NewCell(runes.LineHorizontal, s)
	for i := p.X; i < m.Width(); i++ {
		m.SetCell(canvas.Point{i, y}, r)
	}
}

// DrawXYAxis draws X and Y axes with origin at (X,Y cordinates) with given style.
// Y axis extends up, and X axis extends right.
// Coordinates (0,0) is top left of canvas.
func DrawXYAxis(m *canvas.Model, p canvas.Point, s lipgloss.Style) {
	m.SetCell(p, canvas.NewCell(runes.LineUpRight, s))
	DrawVerticalLineUp(m, canvas.Point{p.X, p.Y - 1}, s)
	DrawHorizonalLineRight(m, canvas.Point{p.X + 1, p.Y}, s)
}

// DrawXYAxisDown draws X and Y axes with origin at (X,Y cordinates) with given style.
// Y axis extends up and down, and X axis extends right.
// Coordinates (0,0) is top left of canvas.
func DrawXYAxisDown(m *canvas.Model, p canvas.Point, s lipgloss.Style) {
	m.SetCell(p, canvas.NewCell(runes.LineVerticalRight, s))
	DrawVerticalLineUp(m, canvas.Point{p.X, p.Y - 1}, s)
	DrawVerticalLineDown(m, canvas.Point{p.X, p.Y + 1}, s)
	DrawHorizonalLineRight(m, canvas.Point{p.X + 1, p.Y}, s)
}

// DrawXYAxisLeft draws X and Y axes with origin at (X,Y cordinates) with given style.
// Y axis extends up, and X axis extends left and right.
// Coordinates (0,0) is top left of canvas.
func DrawXYAxisLeft(m *canvas.Model, p canvas.Point, s lipgloss.Style) {
	m.SetCell(p, canvas.NewCell(runes.LineHorizontalUp, s))
	DrawVerticalLineUp(m, canvas.Point{p.X, p.Y - 1}, s)
	DrawHorizonalLineRight(m, canvas.Point{p.X + 1, p.Y}, s)
	DrawHorizonalLineLeft(m, canvas.Point{p.X - 1, p.Y}, s)
}

// DrawXYAxisAll draws X and Y axes with origin at (X,Y cordinates) with given style.
// Y axis extends up and down, and X axis extends left and right.
// Coordinates (0,0) is top left of canvas.
func DrawXYAxisAll(m *canvas.Model, p canvas.Point, s lipgloss.Style) {
	m.SetCell(p, canvas.NewCell(runes.LineHorizontalVertical, s))
	DrawVerticalLineUp(m, canvas.Point{p.X, p.Y - 1}, s)
	DrawVerticalLineDown(m, canvas.Point{p.X, p.Y + 1}, s)
	DrawHorizonalLineRight(m, canvas.Point{p.X + 1, p.Y}, s)
	DrawHorizonalLineLeft(m, canvas.Point{p.X - 1, p.Y}, s)
}

// DrawColumnRune draws a braille rune on to the canvas at given (X,Y) coordinates with given style.
// The function checks for existing braille runes already on the canvas and
// will draw a new braille pattern with the dot patterns of both the existing and given runes.
// Does nothing if given rune is Null or is not a braille rune.
func DrawBrailleRune(m *canvas.Model, p canvas.Point, r rune, s lipgloss.Style) {
	if (r == runes.Null) || !runes.IsBraillePattern(r) {
		return
	}
	cr := m.Cell(p).Rune
	if cr == 0 { // set rune if nothing exists on canvas
		m.SetCell(p, canvas.NewCell(r, s))
		return
	}
	m.SetCell(p, canvas.NewCell(runes.CombineBraillePatterns(m.Cell(p).Rune, r), s))
}

// DrawBraillePatterns draws braille runes from a [][]rune representing a 2D grid of
// Braille Pattern runes.  The runes will be drawn onto the canvas from starting from top
// left of the grid to the bottom right of the grid starting at the given canvas Point.
// Given style will be applied to all runes drawn.
// This function can be used with the output [][]rune from PatternDotsGrid.BraillePatterns().
func DrawBraillePatterns(m *canvas.Model, p canvas.Point, b [][]rune, s lipgloss.Style) {
	for y, row := range b {
		for x, r := range row {
			if r != runes.BrailleBlockOffset {
				DrawBrailleRune(m, p.Add(canvas.Point{X: x, Y: y}), r, s)
			}
		}
	}
}

// DrawLineSequence draws line runes on to the canvas starting
// from a given X coordinate and a sequence of Y coordinates.
// `startYAxis` should be true if `startX` is the Y axis.
// Sequential Y coordinates will increment X coordinates.
// Applies style to all line runes.
// Handles overlapping lines.
// Handles X and Y axes drawn using DrawXYAxis functions.
// Coordinates (0,0) is top left of canvas.
func DrawLineSequence(m *canvas.Model, startYAxis bool, startX int, seqY []int, ls runes.LineStyle, s lipgloss.Style) {
	var prevY int
	for i, y := range seqY {
		if i == 0 { // draw first point
			p := canvas.Point{startX, y}
			r := runes.LineHorizontal
			if startYAxis {
				c := m.Cell(p).Rune
				if c == runes.LineUpRight { // first point is origin
					m.SetCell(p, canvas.NewCell(runes.LineUpRight, s))
				} else if c == runes.LineVertical { // first point on Y axis
					m.SetCell(p, canvas.NewCell(runes.LineVerticalRight, s))
				} else if c == runes.LineVerticalRight { // first point on Y axis overlapping another line
					m.SetCell(p, canvas.NewCell(runes.LineVerticalRight, s))
				} else {
					DrawLineRune(m, p, r, ls, s)
				}
			} else {
				DrawLineRune(m, p, r, ls, s)
			}
		} else {
			DrawLineSequenceLeftToRight(m, canvas.Point{i + startX - 1, prevY}, canvas.Point{i + startX, y}, ls, s)
		}
		prevY = y
	}
}

// DrawLineSequenceLeftToRight draws line runes from point A to point B where B.X = A.X+1.
// Assumes point A has already been drawn and does not draw point A.
// Applies style to all line runes.
// Handles overlapping lines.
// Handles X and Y axes drawn using DrawXYAxis functions.
// Coordinates (0,0) is top left of canvas.
func DrawLineSequenceLeftToRight(m *canvas.Model, a canvas.Point, b canvas.Point, ls runes.LineStyle, s lipgloss.Style) {
	if a.X >= b.X {
		return
	}

	prevY := a.Y
	y := b.Y
	x := b.X
	r := runes.LineHorizontal // default: point A has same Y coordinate as point B

	// if not the same Y coordinates,
	// draw vertical lines from point A to point B
	if prevY > y { // drawing line up
		r = runes.ArcDownRight
		DrawLineRune(m, canvas.Point{x, prevY}, runes.ArcUpLeft, ls, s)
		for j := prevY - 1; j > y; j-- { // draw vertical lines
			DrawLineRune(m, canvas.Point{x, j}, runes.LineVertical, ls, s)
		}
	} else if prevY < y { // drawing line down
		r = runes.ArcUpRight
		DrawLineRune(m, canvas.Point{x, prevY}, runes.ArcDownLeft, ls, s)
		for j := prevY + 1; j < y; j++ { // draw vertical lines
			DrawLineRune(m, canvas.Point{x, j}, runes.LineVertical, ls, s)
		}
	}

	DrawLineRune(m, b, r, ls, s)
}

// DrawLinePoints draws line runes on to the canvas from a []canvas.Point.
// Each canvas Point is expected to be either adjacent or diagonal from each other.
// At least two Points are required to draw any runes on to the canvas.
// This function can be used with the []canvas.Point output from GetLinePoints().
func DrawLinePoints(m *canvas.Model, points []canvas.Point, ls runes.LineStyle, s lipgloss.Style) {
	if len(points) < 2 {
		return
	}
	extraPoints := []canvas.Point{}
	extraRunes := []rune{} // additional corner runes to draw
	dir := make([]runes.LineSegments, len(points), len(points))
	for i := 1; i < len(points); i++ {
		p := points[i]
		prev := points[i-1]

		if p.X > prev.X {
			if p.Y > prev.Y { // p down right of prev
				dir[i-1].Right = true
				dir[i].Up = true
				extraPoints = append(extraPoints, canvas.Point{X: p.X, Y: p.Y - 1})
				extraRunes = append(extraRunes, runes.ArcDownLeft)
			} else if p.Y < prev.Y { // p up right of prev
				dir[i-1].Right = true
				dir[i].Down = true
				extraPoints = append(extraPoints, canvas.Point{X: p.X, Y: p.Y + 1})
				extraRunes = append(extraRunes, runes.ArcUpLeft)
			} else { // p right of prev
				dir[i-1].Right = true
				dir[i].Left = true
			}
		} else if p.X < prev.X {
			if p.Y > prev.Y { // p down left of prev
				dir[i-1].Left = true
				dir[i].Up = true
				extraPoints = append(extraPoints, canvas.Point{X: p.X, Y: p.Y - 1})
				extraRunes = append(extraRunes, runes.ArcDownRight)
			} else if p.Y < prev.Y { // p up left of prev
				dir[i-1].Left = true
				dir[i].Down = true
				extraPoints = append(extraPoints, canvas.Point{X: p.X, Y: p.Y + 1})
				extraRunes = append(extraRunes, runes.ArcUpRight)
			} else { // p left of prev
				dir[i-1].Left = true
				dir[i].Right = true
			}
		} else {
			if p.Y > prev.Y { // p below prev
				dir[i-1].Down = true
				dir[i].Up = true
			} else if p.Y < prev.Y { // p above prev
				dir[i-1].Up = true
				dir[i].Down = true
			} else {
				// same point - do nothing
			}
		}
	}
	for i, l := range dir {
		DrawLineRune(m, points[i], runes.ArcLineFromLineSegments(l), ls, s)
	}
	for i, r := range extraRunes {
		DrawLineRune(m, extraPoints[i], r, ls, s)
	}
}

// DrawLineRune draws a line rune on to the canvas at given (X,Y) coordinates with given style.
// The given rune is used to check line directions, and the final output line rune
// depends on the given runes.LineStyle.
// The function checks for existing X,Y axis or line runes already on the canvas and draws runes
// such that the lines appear overlapping.
// Does nothing if given rune is empty or is not a line rune.
func DrawLineRune(m *canvas.Model, p canvas.Point, r rune, ls runes.LineStyle, s lipgloss.Style) {
	if (r == runes.Null) || !runes.IsLine(r) {
		return
	}
	m.SetCell(p, canvas.NewCell(runes.CombineLines(m.Cell(p).Rune, r, ls), s))
}

// DrawColumns draws columns going upwards on to canvas
// starting from a given (X,Y) coordinate and a sequence of column lengths.
// Columns will be drawn from left to right and
// sequential column lengths will increment X coordinates for drawing.
// Handles overlapping columns of diferent rune heights.
// If there exists an existing column at given Point with same height as new column,
// then the existing column will be replaced.
// Applies style to all block runes.
// Coordinates (0,0) is top left of canvas.
func DrawColumns(m *canvas.Model, p canvas.Point, seqLen []float64, s lipgloss.Style) {
	y := p.Y
	x := p.X
	for i, f := range seqLen {
		DrawColumnBottomToTop(m, canvas.Point{x + i, y}, f, s)
	}
}

// DrawColumnBottomToTop draws block element runes up from given point.
// The value of float64 is the number of characters to draw up.
// A fractional value is used since there are 1/8th lower block elements and
// fractional values will map to the nearest 1/8th block for the last rune drawn.
// Handles overlapping columns of diferent rune heights.
// If there exists an existing column at given Point with same height as new column,
// then the existing column will be replaced.
// Applies style to all block runes.
// Coordinates (0,0) is top left of canvas.
func DrawColumnBottomToTop(m *canvas.Model, p canvas.Point, v float64, s lipgloss.Style) {
	x := p.X
	y := p.Y

	h := getColumnHeight(m, p) // height of existing column on canvas
	n := math.Floor(v)         // number of full blocks to show
	nh := int(n)               // height of new column to draw on canvas

	r := runes.LowerBlockElementFromFloat64(v - n)
	if r != runes.Null {
		nh++
	}

	fb := canvas.NewCell(runes.FullBlock, s)
	if (h == 0) || (nh == h) { // replace entire column if same height or no existing column
		// set full block columns
		end := int(n)
		for i := 0; i < end; i++ {
			m.SetCell(canvas.Point{x, y - i}, fb)
		}
		// set column top rune
		m.SetCell(canvas.Point{x, y - end}, canvas.NewCell(r, s))
	} else if nh < h { // new column shorter than old column
		// replace existing full blocks with new full blocks
		end := int(n)
		for i := 0; i < end; i++ {
			m.SetCell(canvas.Point{x, y - i}, fb)
		}
		// overlap new column top rune on top of old full block
		DrawColumnRune(m, canvas.Point{x, y - end}, r, s)
	} else if nh > h { // new column taller than old column
		oc := (h - 1) // index of existing column top
		if oc <= 0 {
			oc = 0
		}
		// overlap existing column top rune on top of new full block
		DrawColumnRune(m, canvas.Point{x, y - oc}, runes.FullBlock, s)
		// draw new full blocks above existing columns
		end := int(n)
		for i := h; i < end; i++ {
			m.SetCell(canvas.Point{x, y - i}, fb)
		}
		// set new column top rune
		m.SetCell(canvas.Point{x, y - end}, canvas.NewCell(r, s))
	}
}

// DrawColumnRune draws a column rune on to the canvas at given (X,Y) coordinates with given style.
// The function checks for existing column runes already on the canvas and attempts to
// draws runes such that the runes appear overlapping.
// Overlapping runes can only occur if either one of the runes is a full block element rune,
// and the other rune is not a full block element rune.
// If the runes cannot overlap, then it will the existing rune will be replaced.
// Does nothing if given rune is Null or is not a columne rune.
func DrawColumnRune(m *canvas.Model, p canvas.Point, r rune, s lipgloss.Style) {
	if (r == runes.Null) || !runes.IsLowerBlockElement(r) {
		return
	}
	rs := s.Copy()
	c := m.Cell(p)
	if runes.IsLowerBlockElement(c.Rune) {
		if (r == runes.FullBlock) && (c.Rune != runes.FullBlock) { // existing rune on top of new full block
			r = c.Rune
			rs.Background(s.GetForeground()).Foreground(c.Style.GetForeground())
		} else if (c.Rune == runes.FullBlock) && r != runes.FullBlock { // new rune on top of existing full block
			rs.Background(c.Style.GetForeground()).Foreground(s.GetForeground())
		}
	}
	m.SetCell(p, canvas.NewCell(r, rs))
}

// getColumnHeight obtains number of runes drawn
// by the DrawColumnBottomToTop function at given Point.
func getColumnHeight(m *canvas.Model, p canvas.Point) int {
	x := p.X
	y := p.Y
	i := 0
	c := m.Cell(canvas.Point{x, y})
	for runes.IsLowerBlockElement(c.Rune) {
		i++
		c = m.Cell(canvas.Point{x, y - i})
	}
	return i
}

// abs returns absolute value of given integer.
func abs(i int) int {
	if i < 0 {
		return i * -1
	}
	return i
}

// GetFullCirclePoints returns a []canvas.Point containing points
// that approximates a filled circle of radius r for center Point c.
func GetFullCirclePoints(c canvas.Point, r int) (p []canvas.Point) {
	// sort points
	cPoints := GetCirclePoints(c, r)
	sort.Slice(cPoints, func(i, j int) bool {
		a := cPoints[i]
		b := cPoints[j]
		if a.Y == b.Y {
			return a.X < b.X
		}
		return a.Y < b.Y
	})
	// set all cells between first and last point of a row
	f := cPoints[0]
	l := cPoints[0]
	for _, v := range cPoints {
		// if new row, draw line between previous row first and last points
		if v.Y != l.Y {
			for i := f.X; i < l.X; i++ {
				p = append(p, canvas.Point{X: i, Y: l.Y})
			}
			f = v
		}
		l = v
		p = append(p, v)
	}
	return
}

// GetCirclePoints returns a []canvas.Point containing points
// that approximates a circle of radius r for center Point c.
func GetCirclePoints(c canvas.Point, r int) (p []canvas.Point) {
	t1 := r / 16
	t2 := 0
	x := r
	y := 0
	for x >= y {
		p = append(p, c.Add(canvas.Point{X: x, Y: y}))
		p = append(p, c.Add(canvas.Point{X: x, Y: -y}))
		p = append(p, c.Add(canvas.Point{X: -x, Y: y}))
		p = append(p, c.Add(canvas.Point{X: -x, Y: -y}))
		p = append(p, c.Add(canvas.Point{X: y, Y: x}))
		p = append(p, c.Add(canvas.Point{X: y, Y: -x}))
		p = append(p, c.Add(canvas.Point{X: -y, Y: x}))
		p = append(p, c.Add(canvas.Point{X: -y, Y: -x}))
		y++
		t1 += y
		t2 = t1 - x
		if t2 >= 0 {
			t1 = t2
			x--
		}
	}
	return
}

// GetLinePoints returns a []canvas.Point containing points
// that approximates a line between points p1 and p2.
func GetLinePoints(p1 canvas.Point, p2 canvas.Point) []canvas.Point {
	if abs(p2.Y-p1.Y) < abs(p2.X-p1.X) {
		if p1.X > p2.X {
			return getLinePointsLow(p2, p1)
		} else {
			return getLinePointsLow(p1, p2)
		}
	} else {
		if p1.Y > p2.Y {
			return getLinePointsHigh(p2, p1)
		} else {
			return getLinePointsHigh(p1, p2)
		}
	}
}

// getLinePointsLow returns a []canvas.Point containing points
// that approximates a line between points p1 and p2 for
// slight line slopes between -1 and 1.
func getLinePointsLow(p1 canvas.Point, p2 canvas.Point) (r []canvas.Point) {
	dx := (p2.X - p1.X)
	dy := (p2.Y - p1.Y)
	yi := 1
	if dy < 0 {
		yi = -1
		dy = -dy
	}
	D := (2 * dy) - dx
	y := p1.Y

	start := p1.X
	end := p2.X
	if start > end {
		start = p2.X
		end = p1.X
	}
	for x := start; x <= end; x++ {
		r = append(r, canvas.Point{X: x, Y: y})
		if D > 0 {
			y += yi
			D += (2 * (dy - dx))
		} else {
			D += 2 * dy
		}
	}
	return
}

// getLinePointsHigh returns a []canvas.Point containing points
// that approximates a line between points p1 and p2 for
// steep line slopes <= -1 or >= 1.
func getLinePointsHigh(p1 canvas.Point, p2 canvas.Point) (r []canvas.Point) {
	dx := (p2.X - p1.X)
	dy := (p2.Y - p1.Y)
	xi := 1
	if dx < 0 {
		xi = -1
		dx = -dx
	}
	D := (2 * dx) - dy
	x := p1.X

	start := p1.Y
	end := p2.Y
	if start > end {
		start = p2.Y
		end = p1.Y
	}
	for y := start; y <= end; y++ {
		r = append(r, canvas.Point{X: x, Y: y})
		if D > 0 {
			x += xi
			D += (2 * (dx - dy))
		} else {
			D += 2 * dx
		}
	}
	return
}
