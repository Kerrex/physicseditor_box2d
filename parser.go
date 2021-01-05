package physicseditor_box2d

import (
	"encoding/xml"
	"errors"
	"github.com/ByteArena/box2d"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type ParsedBody struct {
	Name string
	Shapes []box2d.B2FixtureDef
}

type Circle struct {
	XMLName xml.Name `xml:"circle"`
	Radius  float64  `xml:"r,attr"`
	X       float64  `xml:"x,attr"`
	Y       float64  `xml:"y,attr"`
}

type Polygons struct {
	XMLName  xml.Name `xml:"polygons"`
	Polygons []string `xml:"polygon"`
}

type Fixture struct {
	XMLName     xml.Name `xml:"fixture"`
	Density     float64      `xml:"density"`
	Friction    float64      `xml:"friction"`
	Restitution float64      `xml:"restitution"`
	FixtureType string   `xml:"fixture_type"`
	Polygons    Polygons `xml:"polygons"`
	Circle      *Circle  `xml:"circle"`
}

type Fixtures struct {
	XMLName  xml.Name  `xml:"fixtures"`
	Fixtures []Fixture `xml:"fixture"`
}

type Body struct {
	XMLName     xml.Name `xml:"body"`
	Name        string   `xml:"name,attr"`
	AnchorPoint string   `xml:"anchorpoint"`
	Fixtures    Fixtures `xml:"fixtures"`
}

type Bodies struct {
	XMLName xml.Name `xml:"bodies"`
	Bodies  []Body   `xml:"body"`
}

type BodyDef struct {
	XMLName xml.Name `xml:"bodydef"`
	Bodies  Bodies   `xml:"bodies"`
}

var splitRegexp = regexp.MustCompile("\\s+,\\s+")

func Parse(file []byte) ([]ParsedBody, error) {
	return ParseScaled(file, 1)
}

func ParseScaled(file []byte, scale float64) ([]ParsedBody, error) {
	result := BodyDef{}
	err := xml.Unmarshal(file, &result)
	if err != nil {
		log.Println("Unable to parse XML:")
		return nil, err
	}

	bodies, err := parsedBodies(result, scale)
	if err != nil {
		return nil, err
	}

	return bodies, nil
}

func parsedBodies(result BodyDef, scale float64) ([]ParsedBody, error) {
	bodies := make([]ParsedBody, 0)
	for _, body := range result.Bodies.Bodies {
		shapes, err := parseFixtures(body, scale)
		if err != nil {
			return nil, err
		}

		bodyToAdd := ParsedBody{
			Name:   body.Name,
			Shapes: shapes,
		}

		bodies = append(bodies, bodyToAdd)
	}
	return bodies, nil
}

func parseFixtures(body Body, scale float64) ([]box2d.B2FixtureDef, error) {
	shapes := make([]box2d.B2FixtureDef, 0)
	for _, fixture := range body.Fixtures.Fixtures {
		fixtureToAdd := box2d.MakeB2FixtureDef()
		fixtureToAdd.Density = fixture.Density
		fixtureToAdd.Friction = fixture.Friction
		fixtureToAdd.Restitution = fixture.Restitution

		shape, err := parseShape(fixture, body.Name, scale)
		if err != nil {
			return nil, err
		}
		fixtureToAdd.Shape = shape
		shapes = append(shapes, fixtureToAdd)
	}
	return shapes, nil
}

func parseShape(fixture Fixture, bodyName string, scale float64) (box2d.B2ShapeInterface, error) {
	if strings.ToUpper(fixture.FixtureType) == "POLYGON" {
		return parsePolygon(fixture, scale)
	} else if strings.ToUpper(fixture.FixtureType) == "CIRCLE" {
		return parseCircle(fixture, bodyName, scale)
	} else {
		log.Printf("Invalid fixture type %s", fixture.FixtureType)
		return nil, errors.New("invalid fixture type")
	}
}

func parseCircle(fixture Fixture, bodyName string, scale float64) (box2d.B2ShapeInterface, error) {
	if fixture.Circle == nil {
		log.Printf("Circle cannot be null if fixture type is circle! Body %s", bodyName)
		return nil, errors.New("invalid fixture type")
	}

	circle := box2d.NewB2CircleShape()
	circle.SetRadius(fixture.Circle.Radius * scale)
	circle.M_p = box2d.MakeB2Vec2(fixture.Circle.X * scale, fixture.Circle.Y * scale)
	return circle, nil
}

func parsePolygon(fixture Fixture, scale float64) (box2d.B2ShapeInterface, error) {
	vectors := make([]box2d.B2Vec2, 0)
	for _, polygon := range fixture.Polygons.Polygons {
		vecs := splitRegexp.Split(polygon, -1)
		for _, vec := range vecs {
			trimmedVec := strings.TrimSpace(vec)
			xyArr := strings.Split(trimmedVec, ",")

			x, err := strconv.ParseFloat(strings.TrimSpace(xyArr[0]), 64)
			if err != nil {
				return nil, err
			}
			y, err := strconv.ParseFloat(strings.TrimSpace(xyArr[1]), 64)
			if err != nil {
				return nil, err
			}

			vectors = append(vectors, box2d.B2Vec2{
				X: x * scale,
				Y: y * scale,
			})
		}
	}
	shape := box2d.MakeB2ChainShape()
	shape.CreateChain(vectors, len(vectors))

	return &shape, nil
}

