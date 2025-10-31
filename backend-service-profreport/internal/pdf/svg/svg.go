package svg

import (
	"fmt"
	"html/template"
	"log/slog"
	"math"
	"strconv"
)

const Radius = 100

type Point struct {
	Label string
	Value int
}

type CircleDiagram struct {
	Name   string
	Points []Point
}

func (c *CircleDiagram) Render() template.HTML {
	res :=
		fmt.Sprintf("<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"500\" height=\"500\" viewBox=\"-250 -250 500 500\" role=\"img\" aria-label=%s>\n", c.Name) +
			fmt.Sprintf("<circle r=\"%s\" fill=\"none\" stroke=\"#ddd\" stroke-width=\"0.8\"/>\n", strconv.Itoa(Radius*0.25)) +
			fmt.Sprintf("<circle r=\"%s\" fill=\"none\" stroke=\"#ddd\" stroke-width=\"0.8\"/>\n", strconv.Itoa(Radius*0.5)) +
			fmt.Sprintf("<circle r=\"%s\" fill=\"none\" stroke=\"#ddd\" stroke-width=\"0.8\"/>\n", strconv.Itoa(Radius*0.75)) +
			fmt.Sprintf("<circle r=\"%s\" fill=\"none\" stroke=\"#ddd\" stroke-width=\"0.8\"/>\n", strconv.Itoa(Radius)) +
			c.renderAxis()

	return template.HTML(res)
}

func (c *CircleDiagram) renderAxis() string {

	n := len(c.Points)
	if n < 3 {
		return ""
	}

	maxValue := c.maxValue()
	k := float64(Radius) / float64(maxValue)

	phi0 := -math.Pi / 2

	log := slog.With("func", "renderAxis")

	log.Info(
		"render info",
		slog.Int("n", n),
	)

	type xy struct {
		x, y int
	}

	xyArray := make([]xy, n)

	res := "<g stroke=\"#999\" stroke-width=\"1\">\n"

	for i := range n {
		i := i
		log.Info(
			"render info",
			slog.Int("i", i),
		)

		rad := phi0 + ((float64(i*(360/n)) * math.Pi) / 180)

		res += fmt.Sprintf(
			"<line x1=\"0\" y1=\"0\" x2=\"%s\" y2=\"%s\"/>\n",
			strconv.Itoa(int(math.Round((Radius * math.Cos(rad))))),
			strconv.Itoa(int(math.Round((Radius * math.Sin(rad))))),
		)

		xyArray[i] = xy{
			x: int(math.Round((k * float64(c.Points[i].Value) * math.Cos(rad)))),
			y: int(math.Round((k * float64(c.Points[i].Value) * math.Sin(rad)))),
		}
	}

	res += "</g>\n"
	res += "<g font-family=\"Inter, sans-serif\" font-size=\"12\" fill=\"#333\" text-anchor=\"middle\">\n"

	for i := range n {

		rad := phi0 + ((float64(i*(360/n)) * math.Pi) / 180)

		res += fmt.Sprintf(
			"<text x=\"%s\" y=\"%s\">%s</text>\n",
			strconv.Itoa(int(math.Round((1.1 * Radius * math.Cos(rad))))),
			strconv.Itoa(int(math.Round((1.1 * Radius * math.Sin(rad))))),
			c.Points[i].Label,
		)
	}

	res += "</g>\n"
	res += "<polygon points=\"\n"

	for i := range n {
		res += fmt.Sprintf(
			"%s,%s ",
			strconv.Itoa(xyArray[i].x),
			strconv.Itoa(xyArray[i].y),
		)
	}

	res += "\" fill=\"#1f77b4\" fill-opacity=\"0.10\" stroke=\"#1f77b4\" stroke-width=\"2\"/>\n"
	res += "</svg>"
	return res
}

func (c *CircleDiagram) maxValue() int {
	if len(c.Points) == 0 {
		return 0
	}

	max := c.Points[0].Value
	for _, p := range c.Points[1:] {
		if p.Value > max {
			max = p.Value
		}
	}
	return max
}
