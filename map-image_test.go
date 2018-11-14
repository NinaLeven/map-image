package map_image

import (
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestPolygon(t *testing.T) {
	mi := NewMapImageDefault()

	file, err := os.Open("testdata/geojson_polygon.json")
	defer file.Close()
	if err != nil {
		t.Fatal(err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	opt := GetImageOptions{
		SizeX:     650,
		SizeY:     450,
		Maptype:   "map",
		LabelType: "vkgrm",
		GeoJson:   data,
	}
	_, err = mi.GetImage(opt)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLineString(t *testing.T) {
	mi := NewMapImageDefault()

	file, err := os.Open("testdata/geojson_line_string.json")
	defer file.Close()
	if err != nil {
		t.Fatal(err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	opt := GetImageOptions{
		SizeX:     650,
		SizeY:     450,
		Maptype:   "map",
		LabelType: "vkgrm",
		GeoJson:   data,
	}
	_, err = mi.GetImage(opt)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLabels(t *testing.T) {
	mi := NewMapImageDefault()

	file, err := os.Open("testdata/geojson_labels.json")
	defer file.Close()
	if err != nil {
		t.Fatal(err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	opt := GetImageOptions{
		SizeX:     650,
		SizeY:     450,
		Maptype:   "map",
		LabelType: "vkgrm",
		GeoJson:   data,
	}
	_, err = mi.GetImage(opt)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSingleObject(t *testing.T) {
	mi := NewMapImageDefault()

	file, err := os.Open("testdata/geojson_single.json")
	defer file.Close()
	if err != nil {
		t.Fatal(err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	opt := GetImageOptions{
		SizeX:     650,
		SizeY:     450,
		Maptype:   "map",
		LabelType: "vkgrm",
		GeoJson:   data,
	}
	_, err = mi.GetImage(opt)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAll(t *testing.T) {
	mi := NewMapImageDefault()

	file, err := os.Open("testdata/geojson.json")
	defer file.Close()
	if err != nil {
		t.Fatal(err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	opt := GetImageOptions{
		SizeX:     650,
		SizeY:     450,
		Maptype:   "map",
		LabelType: "vkgrm",
		GeoJson:   data,
	}
	body, err := mi.GetImage(opt)
	if err != nil {
		t.Fatal(err)
	}

	file, err = os.Create("testdata/img.png")
	if err != nil {
		t.Fatal(err)
	}

	_, err = io.Copy(file, body)
	if err != nil {
		t.Fatal(err)
	}
}
