package map_image

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	// this library is used in order to remove unnecessary complexity, witch arises when using similar libraries
	"github.com/paulmach/go.geojson"
)

type MapImage struct {
	host string
}

type MapImageOptions struct {
	Host string `default:"http://static-maps.yandex.ru/1.x/"`
}

func NewMapImageDefault() *MapImage {
	mi := &MapImage{
		host: "http://static-maps.yandex.ru/1.x/",
	}

	return mi
}

func NewMapImage(options MapImageOptions) *MapImage {
	mi := &MapImage{
		host: options.Host,
	}

	return mi
}

type GetImageOptions struct {
	/*
		max size is 650x450

		default is 650x450
	*/
	SizeX int
	SizeY int

	/*
		zoom is an int from 0 to 17
	*/
	Zoom int

	/*
		could be:
			map - schematic map
			sat - satellite image
			skl - landmarks
			trf - traffic
		or any combination of those delimited by coma

		default is map
	*/
	Maptype string

	/*
		for detailed description see tech.yandex.ru/maps/doc/staticapi/1.x/dg/concepts/markers-docpage/
	*/
	LabelType string

	/*
		line width in pixels
	*/
	LineThickness int

	/*
		colors are represented as a string of 4 concatinated 4-digit hex values in RGBA format
		for examlpe "EC473FFF"
	*/
	LineColor string
	FillColor string

	/*
		GeoJson should be a geojson in format [{"geometry":...},...] or {"geometry":...}

		!!! this parameter doesn't have a default value and is mandatorily passed by the user
	*/
	GeoJson []byte
}

type MapImageError struct {
	Status  string `xml:"status"`
	Message string `xml:"message"`
}

func NewMapImageError(data []byte) *MapImageError {
	obj := &MapImageError{}
	err := xml.Unmarshal(data, obj)
	if err != nil {
		log.Println(err)
	}
	return obj
}

func (e *MapImageError) Error() string {
	return "{ status: \"" + e.Status + "\", message: \"" + e.Message + "\" }"
}

func (mi *MapImage) GetImage(opt GetImageOptions) (io.ReadCloser, error) {
	if opt.GeoJson == nil {
		return nil, errors.New("GeoJson parameter is nil")
	}
	if opt.Maptype == "" {
		opt.Maptype = "map"
	}
	if opt.LineThickness == 0 {
		opt.LineThickness = 1
	}
	if opt.FillColor == "" {
		opt.FillColor = "00FF0020"
	}
	if opt.LineColor == "" {
		opt.LineColor = "ec473fFF"
	}
	if opt.LabelType == "" {
		opt.LabelType = "vkgrm"
	}

	geoopts, err := unmarshalGeoJson(opt.GeoJson)
	if err != nil {
		return nil, err
	}

	params := []string{
		"l=" + opt.Maptype,
		polygonParam(geoopts, opt.LineThickness, opt.LineColor, opt.FillColor),
		labelParam(geoopts, opt.LabelType),
	}

	if opt.SizeX != 0 && opt.SizeY != 0 {
		params = append(params, fmt.Sprintf("size=%d,%d", opt.SizeX, opt.SizeY))
	}
	if opt.Zoom != 0 {
		params = append(params, fmt.Sprintf("z=%d", opt.Zoom))
	}

	query := concatParams(params)
	reqURL := mi.host + query

	res, err := http.Get(reqURL)
	if err != nil {
		return nil, err
	}

	if len(res.Header["Content-Type"]) == 0 || res.Header["Content-Type"][0] == "text/xml" {
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		return nil, NewMapImageError(data)
	}

	return res.Body, nil
}

func unmarshalGeoJson(GeoJson []byte) ([]geojson.Geometry, error) {
	geoopts := []geojson.Geometry{}

	err := json.Unmarshal(GeoJson, &geoopts)
	if err != nil {
		geoopt := geojson.Geometry{}

		err = json.Unmarshal(GeoJson, &geoopt)
		if err != nil {
			return nil, err
		}

		geoopts = []geojson.Geometry{geoopt}
	}

	return geoopts, nil
}

func concatParams(params []string) string {
	str := strings.Builder{}
	flag := false

	str.WriteString("?")
	for i := 0; i < len(params); i++ {
		if params[i] != "" {
			if flag {
				str.WriteString("&")
			}
			flag = true

			str.WriteString(params[i])
		}
	}

	return str.String()
}

func labelParam(geoms []geojson.Geometry, labelType string) string {
	str := strings.Builder{}
	flag := false

	str.WriteString("pt=")
	for _, obj := range geoms {

		if obj.Type == "Point" && len(obj.Point) == 2 {

			if flag {
				str.WriteString("~")
			}
			flag = true

			str.WriteString(fmt.Sprintf("%f,%f,%s", obj.Point[0], obj.Point[1], labelType))
		}
	}

	return str.String()
}

func polygonParam(geoms []geojson.Geometry, lineThickness int, lineColor string, fillColor string) string {
	str := strings.Builder{}
	flag := false

	str.WriteString("pl=")
	for _, obj := range geoms {

		if obj.Type == "LineString" && len(obj.LineString) > 0 {

			if flag {
				str.WriteString("~")
			}
			flag = true

			innerFlag := false
			str.WriteString(fmt.Sprintf("c:%s,w:%d,", lineColor, lineThickness))
			for _, p := range obj.LineString {

				if innerFlag {
					str.WriteString(",")
				}
				innerFlag = true

				if len(p) == 2 {
					str.WriteString(fmt.Sprintf("%f,%f", p[0], p[1]))
				}
			}
		}

		if obj.Type == "Polygon" && len(obj.Polygon) > 0 {

			if flag {
				str.WriteString("~")
			}
			flag = true

			innerFlag := false

			str.WriteString(fmt.Sprintf("c:%s,f:%s,w:%d,", lineColor, fillColor, lineThickness))
			for _, poly := range obj.Polygon {

				if innerFlag {
					str.WriteString(";")
				}
				innerFlag = true

				for _, p := range poly {
					if len(p) == 2 {
						str.WriteString(fmt.Sprintf("%f,%f,", p[0], p[1]))
					}
				}
				if len(poly) > 0 && len(poly[0]) == 2 {
					str.WriteString(fmt.Sprintf("%f,%f", poly[0][0], poly[0][1]))
				}
			}
		}
	}

	return str.String()
}
