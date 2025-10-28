package svg

import (
	"fmt"
	"html/template"
	//"math"
)

const Radius = 100

type CircleDiagram struct {
	Name string
	Data []struct{}
}

func (c *CircleDiagram) Render() template.HTML {
	res :=
		fmt.Sprintf("<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"500\" height=\"500\" viewBox=\"-250 -250 500 500\" role=\"img\" aria-label=%s>\n", c.Name) +
			fmt.Sprintf("<circle r=\"40\" fill=\"none\" stroke=\"#ddd\" stroke-width=\"0.8\"/>\n")

	return template.HTML(res)
}
