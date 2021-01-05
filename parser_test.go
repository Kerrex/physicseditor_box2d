package physicseditor_box2d

import (
	"github.com/ByteArena/box2d"
	"testing"
)

const sampleXml = `
<?xml version="1.0" encoding="UTF-8"?>
<!-- created with http://www.physicseditor.de -->
<bodydef version="1.0">
	<bodies>

		<body name="Zrzut ekranu 2020-08-18 o 15">
            <anchorpoint>0.5000,0.5000</anchorpoint>
			<fixtures>

				<fixture>
					<density>2</density>
					<friction>0</friction>
					<restitution>0</restitution>
					<filter_categoryBits>1</filter_categoryBits>
					<filter_groupIndex>0</filter_groupIndex>
					<filter_maskBits>65535</filter_maskBits>
					<fixture_type>POLYGON</fixture_type>


					<polygons>

                        <polygon>  788.5000, -76.0000  ,  788.5000, 76.0000  ,  -788.5000, 76.0000  ,  -344.5000, 7.0000 </polygon>

					</polygons>

				</fixture>

				<fixture>
					<density>3</density>
					<friction>1</friction>
					<restitution>1</restitution>
					<filter_categoryBits>1</filter_categoryBits>
					<filter_groupIndex>0</filter_groupIndex>
					<filter_maskBits>65535</filter_maskBits>
					<fixture_type>CIRCLE</fixture_type>


                    <circle r="20.000" x="-687.500" y="-12.000"/>

				</fixture>

				<fixture>
					<density>2</density>
					<friction>0</friction>
					<restitution>0</restitution>
					<filter_categoryBits>1</filter_categoryBits>
					<filter_groupIndex>0</filter_groupIndex>
					<filter_maskBits>65535</filter_maskBits>
					<fixture_type>POLYGON</fixture_type>


					<polygons>

                        <polygon>  -356.5000, -14.0000  ,  -396.5000, -14.0000  ,  -396.5000, -54.0000 </polygon>

					</polygons>

				</fixture>

			</fixtures>
		</body>

		<body name="Zrzut ekranu 2020-09-2 o 11">
            <anchorpoint>0.5000,0.5040</anchorpoint>
			<fixtures>

			</fixtures>
		</body>

	</bodies>
	<metadata>
		<format>1</format>
		<ptm_ratio>32</ptm_ratio>
	</metadata>
</bodydef>
`

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}

func TestShouldParseBodyData(t *testing.T) {
	// when
	parsed, err := Parse([]byte(sampleXml))
	if err != nil {
		t.Fatalf("parse should be successful: %v", err)
	}

	// then
	assertEqual(t, len(parsed), 2)
	assertEqual(t, parsed[0].Name, "Zrzut ekranu 2020-08-18 o 15")
	assertEqual(t, len(parsed[0].Shapes), 3)

	// and
	assertEqual(t, len(parsed[1].Shapes), 0)
}

func TestShouldParseFixtureData(t *testing.T) {
	// when
	parsed, err := Parse([]byte(sampleXml))
	if err != nil {
		t.Fatalf("parse should be successful: %v", err)
	}

	// then
	firstShape := parsed[0].Shapes[0]
	assertEqual(t, firstShape.Restitution, 0.0)
	assertEqual(t, firstShape.Friction, 0.0)
	assertEqual(t, firstShape.Density, 2.0)
	assertEqual(t, firstShape.Shape.GetType(), box2d.B2Shape_Type.E_chain)

	secondShape := parsed[0].Shapes[1]
	assertEqual(t, secondShape.Restitution, 1.0)
	assertEqual(t, secondShape.Friction, 1.0)
	assertEqual(t, secondShape.Density, 3.0)
	assertEqual(t, secondShape.Shape.GetType(), box2d.B2Shape_Type.E_circle)
}

func TestShouldParsePolygonData(t *testing.T) {
	// when
	parsed, err := Parse([]byte(sampleXml))
	if err != nil {
		t.Fatalf("parse should be successful: %v", err)
	}

	// then
	firstShape := parsed[0].Shapes[0]
	chain := firstShape.Shape.(*box2d.B2ChainShape)
	assertEqual(t, chain.M_vertices[0].X, 788.5000)
	assertEqual(t, chain.M_vertices[0].Y, -76.0000)

	assertEqual(t, chain.M_vertices[1].X, 788.5000)
	assertEqual(t, chain.M_vertices[1].Y, 76.0000)

	assertEqual(t, chain.M_vertices[2].X, -788.5000)
	assertEqual(t, chain.M_vertices[2].Y, 76.0000)

	assertEqual(t, chain.M_vertices[3].X, -344.5000)
	assertEqual(t, chain.M_vertices[3].Y, 7.0000)
}

func TestShouldParseCircleData(t *testing.T) {
	// when
	parsed, err := Parse([]byte(sampleXml))
	if err != nil {
		t.Fatalf("parse should be successful: %v", err)
	}

	// then
	secondShape := parsed[0].Shapes[1]
	circle := secondShape.Shape.(*box2d.B2CircleShape)
	assertEqual(t, circle.GetRadius(), 20.0)
	assertEqual(t, circle.M_p.X, -687.500)
	assertEqual(t, circle.M_p.Y, -12.000)

	assertEqual(t, len(parsed[1].Shapes), 0)
}

func TestShouldScaleFixtures(t *testing.T) {
	// when
	parsed, err := ParseScaled([]byte(sampleXml), 2)
	if err != nil {
		t.Fatalf("parse should be successful: %v", err)
	}

	// then
	firstShape := parsed[0].Shapes[0]
	chain := firstShape.Shape.(*box2d.B2ChainShape)
	assertEqual(t, chain.M_vertices[0].X, 788.5000 * 2)
	assertEqual(t, chain.M_vertices[0].Y, -76.0000 * 2)

	// and
	secondShape := parsed[0].Shapes[1]
	circle := secondShape.Shape.(*box2d.B2CircleShape)
	assertEqual(t, circle.GetRadius(), 20.0 * 2)
	assertEqual(t, circle.M_p.X, -687.500 * 2)
	assertEqual(t, circle.M_p.Y, -12.000 * 2)
}

func TestShouldReturnErrorOnInvalidXml(t *testing.T) {
	invalidXml := "<xml> this in invalid xml </xml>"
	_, err := Parse([]byte(invalidXml))

	if err == nil {
		t.Fatalf("should return error")
	}
}