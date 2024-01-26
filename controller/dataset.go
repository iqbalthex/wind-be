package controller

import (
  "encoding/json"
  "fmt"
  "github.com/gofiber/fiber/v2"
  "io"
  "net/http"
  "strconv"
  "strings"
)

type Dict map[string]any

type Dataset struct {
  DateRange, AltRange, LatRange, LonRange string
}

type ErddapResponse struct {
  Table ErddapTable `json:"table"`
}

type ErddapTable struct {
  Columns []string `json:"columnNames"`
  Rows    [][]any  `json:"rows"`
}

var (
  dataset    Dataset
  dataSource Dict
)

// constructor
func init() {
  dataset = Dataset{
    DateRange: "2023-06-01_2023-06-03",
    AltRange: "10_10",
    LatRange: "-6_-8",
    LonRange: "112_114",
  }

  dataSource = Dict{
    "NOAA": "https://coastwatch.pfeg.noaa.gov/erddap/griddap/erdQCwindproducts1day.json?",
  }
}

// factory
func NewDatasetController() *Dataset { return &Dataset{} }

// router
func (d *Dataset) BindRoute(router fiber.Router) {
  router.Get("/winds", d.getWind)
}

// methods
func (d *Dataset) getWind(c *fiber.Ctx) error {
  endpoint := dataSource["NOAA"].(string)

  dtime := c.Query("datetime")
  lats  := c.Query("lat")
  lons  := c.Query("lon")

  queries := []string{}
  params  := []string{"wind_direction", "wind_speed"}
  ranges  := fmt.Sprintf("[%s][%s][%s][%s]",
    parseRange(dtime),
    parseRange("10_10"),
    parseRange(lats),
    parseRange(lons),
  )

  // ["wind_direction[][][]", "wind_speed[][][]"]
  for _, param := range params {
    queries = append(queries, param + ranges)
  }

  // "...?" += "wind_direction[][][],wind_speed[][][]"
  endpoint += strings.Join(queries, ",")

  res, err := http.Get(endpoint)
  if err != nil {
    fmt.Println(err)
    return c.Status(500).JSON(&fiber.Map{ "error": err })
  }

  body, err := io.ReadAll(res.Body)
  res.Body.Close()

  if res.StatusCode >= 300 {
    fmt.Println(res.StatusCode, body)
    return c.Status(res.StatusCode).JSON(&fiber.Map{
      "error": fmt.Sprintf("Response failed: [%d] %s", res.StatusCode, body),
    })
  }

  if err != nil {
    fmt.Println(err)
    return c.Status(500).JSON(&fiber.Map{ "error": err })
  }

  var response ErddapResponse
	if err := json.Unmarshal(body, &response); err != nil {
    fmt.Println(err)
    return c.Status(500).JSON(&fiber.Map{ "error": err })
  }

  datasets := []Dict{}

  for _, row := range response.Table.Rows {
    if row[5] == nil { continue }

    // .toFixed()
    lat, _ := strconv.ParseFloat(fmt.Sprintf("%6f", row[2]), 6)
    lon, _ := strconv.ParseFloat(fmt.Sprintf("%6f", row[3]), 6)
    dir, _ := strconv.ParseFloat(fmt.Sprintf("%2f", row[4]), 2)
    spd, _ := strconv.ParseFloat(fmt.Sprintf("%2f", row[5]), 2)

    datasets = append(datasets, Dict{
      "time": row[0],
      "lat": lat,
      "lon": lon,
      "dir": dir,
      "spd": spd,
    })
  }

  return c.Status(200).JSON(&fiber.Map{
    "data_count": len(datasets),
    "datasets": datasets,
  })
}

func parseRange(ranges string) string {
  startEnd := strings.Split(ranges, "_")

  return fmt.Sprintf("(%s):1:(%s)", startEnd[0], startEnd[1])
}
