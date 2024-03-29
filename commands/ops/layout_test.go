package ops

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/statecrafthq/borg/utils"

	"github.com/statecrafthq/borg/geometry"

	"github.com/stretchr/testify/assert"
)

func TestLayout(t *testing.T) {
	poly := geometry.NewProjectedSimplePolygon([]geometry.PointGeo{
		{-73.998824, 40.716576},
		{-73.998862, 40.716515},
		{-73.998523, 40.716396},
		{-73.998485, 40.716457},
	})
	layout := LayoutRectangle(poly, 3.6576, 10.668)
	assert.True(t, layout.Analyzed)
	assert.True(t, layout.Fits)
	assert.True(t, layout.HasLocation)
	assert.InEpsilon(t, 2.0045, layout.Angle, 0.0001)
}

func TestConvex(t *testing.T) {
	poly := geometry.NewProjectedSimplePolygon([]geometry.PointGeo{
		{-73.999994, 40.72815},
		{-73.999787, 40.728396},
		{-74.000423, 40.728709},
		{-74.000626, 40.728467},
		{-74.000314, 40.728308},
	})
	layout := LayoutRectangle(poly, 3.6576, 10.668)
	assert.True(t, layout.Analyzed)
	assert.True(t, layout.Fits)
	assert.True(t, layout.HasLocation)
}

func TestNarrow(t *testing.T) {
	poly := geometry.NewProjectedSimplePolygon([]geometry.PointGeo{
		{-74.008608, 40.717205},
		{-74.008557, 40.71727},
		{-74.008904, 40.717317},
		{-74.008917, 40.717247},
		{-74.008763, 40.717226},
	})
	layout := LayoutRectangle(poly, 3.6576, 10.668)
	assert.True(t, layout.Analyzed)
	assert.True(t, layout.Fits)
	assert.True(t, layout.HasLocation)
}

func TestComplex(t *testing.T) {
	poly := geometry.NewProjectedSimplePolygon([]geometry.PointGeo{
		{-74.002324, 40.719039},
		{-74.002587, 40.719159},
		{-74.00265, 40.719208},
		{-74.002773, 40.719072},
		{-74.00243, 40.718915},
	})
	layout := LayoutRectangle(poly, 3.6576, 10.668)
	assert.True(t, layout.Analyzed)
	assert.True(t, layout.Fits)
	assert.True(t, layout.HasLocation)
}

func loadParcel(str string) geometry.Polygon2D {
	var coordinates []interface{}
	e := json.Unmarshal([]byte(str), &coordinates)
	if e != nil {
		log.Panic()
	}
	c := utils.ParseFloat4(coordinates)[0]
	// p := c[:len(c)-1]
	// fmt.Println(c)
	// fmt.Println(p)
	poly := geometry.NewGeoPolygon(c)
	proj := geometry.NewProjection(poly.Center())
	// fmt.Println(poly.Center())
	projected := poly.Project(proj)
	return projected
}

func testLayoutCase(t *testing.T, bbl string, str string) {
	testLayoutCaseF(t, bbl, str, true, true)
}

func testLayoutCaseOnly2(t *testing.T, bbl string, str string) {
	testLayoutCaseF(t, bbl, str, false, true)
}

func testLayoutCaseOnly1(t *testing.T, bbl string, str string) {
	testLayoutCaseF(t, bbl, str, true, false)
}

func testLayoutCaseNone(t *testing.T, bbl string, str string) {
	testLayoutCaseF(t, bbl, str, false, false)
}

func testLayoutCaseF(t *testing.T, bbl string, str string, e1 bool, e2 bool) {

	poly := loadParcel(str)

	fmt.Println(poly.DebugString())

	// Kassita-1: 12ft x 35ft (3.6576 x 10.668)
	layout := LayoutRectangle(poly, 3.6576, 10.668)
	assert.True(t, layout.Analyzed, "Element1 "+bbl+" should be analyzed")
	if layout.Analyzed {
		if e1 {
			assert.True(t, layout.Fits, "Element1 "+bbl+" should fit")
			if layout.Fits {
				assert.True(t, layout.HasLocation, "Element1 "+bbl+" should has location")

				// footprint := geometry.NewSimplePolygon([]geometry.Point2D{
				// 	{-3.6576 / 2, 10.668 / 2},
				// 	{3.6576 / 2, 10.668 / 2},
				// 	{3.6576 / 2, -10.668 / 2},
				// 	{-3.6576 / 2, -10.668 / 2},
				// }).Rotate(layout.Angle).Shift(layout.Center)
				// assert.True(t, poly.Contains(footprint), "Element1 "+bbl+" layout should be within poly")
			}
		} else {
			assert.False(t, layout.Fits, "Element1 "+bbl+" should NOT fit")
		}
	}

	// Kassita-2: 10ft x 35ft (3.048  x 12.192)
	if e2 {
		layout := LayoutRectangle(poly, 3.048, 12.192)
		assert.True(t, layout.Analyzed, "Element2 "+bbl+" should be analyzed")
		if layout.Analyzed {
			if e2 {
				assert.True(t, layout.Fits, "Element2 "+bbl+" should fit")
				if layout.Fits {
					assert.True(t, layout.HasLocation, "Element2 "+bbl+" should has location")

					// footprint := geometry.NewSimplePolygon([]geometry.Point2D{
					// 	{X: -3.048 / 2, Y: 12.192 / 2},
					// 	{X: 3.048 / 2, Y: 12.192 / 2},
					// 	{X: 3.048 / 2, Y: -12.192 / 2},
					// 	{X: -3.048 / 2, Y: -12.192 / 2},
					// }).Rotate(layout.Angle).Shift(layout.Center)
					// assert.True(t, poly.Contains(footprint), "Element2 "+bbl+" layout should be within poly")
				}
			} else {
				assert.False(t, layout.Fits, "Element2 "+bbl+" should NOT fit")
			}
		}
	}
}

func TestUrbynNegative(t *testing.T) {
	parcel := loadParcel("[[[[-74.012129,40.704325],[-74.012181,40.704322],[-74.012161,40.704304],[-74.012152,40.704088],[-74.012125,40.704087],[-74.012129,40.704325]]]]")
	layout := LayoutRectangle(parcel, 3.6576, 10.668)

	footprint := geometry.NewSimplePolygon(
		[]geometry.Point2D{
			{-3.6576 / 2, 10.668 / 2},
			{3.6576 / 2, 10.668 / 2},
			{3.6576 / 2, -10.668 / 2},
			{-3.6576 / 2, -10.668 / 2},
		}).
		Rotate(4.052654523130833).
		Shift(geometry.Point2D{-0.277551, 16.671435})

	fmt.Println(parcel.Contains(footprint))

	fmt.Println(footprint.DebugString())
	fmt.Println(layout.Center.DebugString())

	assert.False(t, parcel.Contains(footprint))
	// [[[[-74.012129,40.704325],[-74.012181,40.704322],[-74.012161,40.704304],[-74.012152,40.704088],[-74.012125,40.704087],[-74.012129,40.704325]]]]
}

func TestLayoutWorking(t *testing.T) {
	// Only element2
	testLayoutCaseOnly2(t, "4-10234-0304", "[[[[-73.783445,40.702794],[-73.783413,40.702807],[-73.783581,40.703056],[-73.783613,40.703043],[-73.783445,40.702794]]]]")
	testLayoutCaseOnly2(t, "4-01656-0028", "[[[[-73.865268,40.766762],[-73.865296,40.766781],[-73.865589,40.766618],[-73.865567,40.766596],[-73.865268,40.766762]]]]")
	testLayoutCaseOnly2(t, "4-10235-0367", "[[[[-73.783155,40.702416],[-73.783185,40.702404],[-73.783023,40.702156],[-73.782988,40.702169],[-73.783155,40.702416]]]]")

	// All
	testLayoutCase(t, "2-03916-0040", "[[[[-73.867656,40.837069],[-73.867742,40.837067],[-73.86773,40.836785],[-73.867644,40.836787],[-73.867653,40.836995],[-73.867656,40.837069]]]]")
	testLayoutCase(t, "2-03758-0027", "[[[[-73.861931,40.829939],[-73.861982,40.829925],[-73.861939,40.829736],[-73.861846,40.829748],[-73.861902,40.829874],[-73.861931,40.829939]]]]")
	testLayoutCase(t, "3-07598-0051", "[[[[-73.942799,40.628888],[-73.943237,40.62886],[-73.94307,40.628692],[-73.942774,40.628864],[-73.942799,40.628888]]]]")
	testLayoutCase(t, "2-05519-0001", "[[[[-73.808321,40.821649],[-73.808485,40.821638],[-73.80835,40.821369],[-73.808288,40.821372],[-73.808303,40.821502],[-73.808321,40.821649]]]]")
	testLayoutCase(t, "4-10420-0095", "[[[[-73.765794,40.698434],[-73.766437,40.697635],[-73.765736,40.698407],[-73.765794,40.698434]]]]")
	testLayoutCase(t, "4-11197-0109", "[[[[-73.735619,40.707746],[-73.735637,40.707751],[-73.735754,40.707507],[-73.735682,40.707487],[-73.735668,40.707546],[-73.735619,40.707746]]]]")
	testLayoutCase(t, "4-11303-0105", "[[[[-73.737919,40.69649],[-73.737851,40.69647],[-73.737712,40.696737],[-73.737779,40.696757],[-73.737852,40.696617],[-73.737919,40.69649]]]]")
	testLayoutCase(t, "4-05258-0080", "[[[[-73.792611,40.76645],[-73.792618,40.766398],[-73.792319,40.766374],[-73.792312,40.766426],[-73.792611,40.76645]]]]")
	testLayoutCase(t, "4-13733-0033", "[[[[-73.747541,40.65496],[-73.747703,40.654708],[-73.747539,40.654657],[-73.747398,40.654915],[-73.747541,40.65496]]]]")
	testLayoutCase(t, "4-10167-0054", "[[[[-73.7904,40.697587],[-73.790365,40.69754],[-73.790051,40.697673],[-73.790086,40.697721],[-73.7904,40.697587]]]]")
	testLayoutCase(t, "5-00520-0013", "[[[[-74.083906,40.627185],[-74.083958,40.627045],[-74.083696,40.626951],[-74.083647,40.627013],[-74.083906,40.627185]]]]")
	testLayoutCase(t, "5-00822-0022", "[[[[-74.1325,40.602759],[-74.132858,40.602729],[-74.132851,40.602674],[-74.132492,40.602704],[-74.1325,40.602759]]]]")
	testLayoutCase(t, "5-03753-0012", "[[[[-74.083599,40.579272],[-74.083805,40.579039],[-74.083685,40.578979],[-74.083479,40.579213],[-74.083599,40.579272]]]]")
	testLayoutCase(t, "5-03753-0038", "[[[[-74.083268,40.578772],[-74.083463,40.578552],[-74.083281,40.578461],[-74.083086,40.578682],[-74.083145,40.578711],[-74.083268,40.578772]]]]")
	testLayoutCase(t, "5-03753-0041", "[[[[-74.083381,40.578828],[-74.083576,40.578608],[-74.083463,40.578552],[-74.083268,40.578772],[-74.083381,40.578828]]]]")
	testLayoutCase(t, "5-03753-0045", "[[[[-74.08362,40.578947],[-74.083815,40.578727],[-74.083695,40.578667],[-74.0835,40.578888],[-74.083562,40.578918],[-74.08362,40.578947]]]]")
	testLayoutCase(t, "5-05460-0147", "[[[[-74.156634,40.549031],[-74.156673,40.549011],[-74.156514,40.548826],[-74.156474,40.548845],[-74.156634,40.549031]]]]")
	testLayoutCase(t, "5-06740-0045", "[[[[-74.203615,40.516675],[-74.203976,40.516558],[-74.203944,40.516506],[-74.203898,40.516497],[-74.203615,40.516675]]]]")

	testLayoutCase(t, "5-03284-0047", "[[[[-74.069917,40.594958],[-74.070241,40.594853],[-74.07022,40.594803],[-74.069896,40.594908],[-74.069917,40.594958]]]]")
	testLayoutCase(t, "5-03284-0060", "[[[[-74.069917,40.594958],[-74.069896,40.594908],[-74.069613,40.595],[-74.069638,40.595049],[-74.069823,40.594988],[-74.069917,40.594958]]]]")
	testLayoutCase(t, "5-03671-0011", "[[[[-74.097985,40.579971],[-74.098162,40.580091],[-74.098225,40.580037],[-74.098047,40.579917],[-74.097985,40.579971]]]]")
	testLayoutCase(t, "5-03684-0016", "[[[[-74.096939,40.578282],[-74.097121,40.578079],[-74.097034,40.578035],[-74.096854,40.578236],[-74.096939,40.578282]]]]")
	testLayoutCase(t, "5-00569-0118", "[[[[-74.081099,40.636566],[-74.081389,40.636386],[-74.081272,40.636282],[-74.080982,40.636462],[-74.081099,40.636566]]]]")
	testLayoutCase(t, "4-12618-0032", "[[[[-73.754382,40.693197],[-73.754722,40.693102],[-73.754705,40.693066],[-73.754364,40.693161],[-73.754382,40.693197]]]]")
	testLayoutCase(t, "5-06559-0091", "[[[[-74.186316,40.531277],[-74.186096,40.530995],[-74.185485,40.53125],[-74.185615,40.53142],[-74.185673,40.531416],[-74.185728,40.531432],[-74.185774,40.531469],[-74.18579,40.531507],[-74.186316,40.531277]]]]")
	testLayoutCase(t, "4-13845-0033", "[[[[-73.744317,40.645254],[-73.744424,40.645207],[-73.744296,40.644987],[-73.744185,40.645025],[-73.744317,40.645254]]]]")
	testLayoutCase(t, "5-00018-0142", "[[[[-74.079726,40.640203],[-74.079625,40.640218],[-74.079658,40.640352],[-74.079759,40.640338],[-74.079741,40.640265],[-74.079726,40.640203]]]]")
	testLayoutCase(t, "5-01715-0100", "[[[[-74.169478,40.624249],[-74.170375,40.623017],[-74.168935,40.622592],[-74.168102,40.622411],[-74.167904,40.622757],[-74.167806,40.622999],[-74.167573,40.623906],[-74.169478,40.624249]]]]")
	testLayoutCase(t, "4-10509-0318", "[[[[-73.767245,40.718895],[-73.767138,40.718784],[-73.766978,40.718876],[-73.767004,40.718901],[-73.767148,40.718887],[-73.767245,40.718895]]]]")

	testLayoutCase(t, "2-03275-0075", "[[[[-73.891408,40.862843],[-73.891984,40.862254],[-73.891959,40.862238],[-73.891315,40.862781],[-73.891394,40.862834],[-73.891408,40.862843]]]]")
	testLayoutCase(t, "2-04944-0083", "[[[[-73.831393,40.886978],[-73.83154,40.886762],[-73.83137,40.886745],[-73.831393,40.886978]]]]")

	testLayoutCase(t, "5-00611-0095", "[[[[-74.084985,40.6262],[-74.084349,40.625993],[-74.084318,40.626057],[-74.084949,40.626262],[-74.084985,40.6262]]]]")
	testLayoutCase(t, "5-01392-0115", "[[[[-74.180289,40.628745],[-74.180366,40.628849],[-74.180334,40.629001],[-74.180143,40.62922],[-74.180243,40.629387],[-74.180316,40.629303],[-74.180755,40.629102],[-74.180533,40.628932],[-74.180289,40.628745]]]]")
	testLayoutCase(t, "5-02231-0135", "[[[[-74.169483,40.611129],[-74.169694,40.611163],[-74.169851,40.61049],[-74.169745,40.610474],[-74.169704,40.610652],[-74.169597,40.610636],[-74.16953,40.610924],[-74.169483,40.611129]]]]")
	testLayoutCase(t, "5-03760-0001", "[[[[-74.088245,40.578561],[-74.088046,40.578771],[-74.088295,40.578901],[-74.08849,40.578695],[-74.088268,40.578574],[-74.088245,40.578561]]]]")
	testLayoutCase(t, "4-12217-0052", "[[[[-73.789947,40.676434],[-73.790016,40.676423],[-73.789888,40.676162],[-73.789871,40.676165],[-73.789947,40.676434]]]]")
	testLayoutCase(t, "5-03753-0016", "[[[[-74.083356,40.579152],[-74.083562,40.578918],[-74.083268,40.578772],[-74.083062,40.579005],[-74.083356,40.579152]]]]")
	testLayoutCase(t, "4-01961-0001", "[[[[-73.853571,40.740931],[-73.853499,40.740789],[-73.853486,40.740785],[-73.85343,40.740886],[-73.853571,40.740931]]]]")
	testLayoutCase(t, "4-11214-0049", "[[[[-73.73954,40.704588],[-73.739674,40.704202],[-73.739589,40.7042],[-73.739573,40.704484],[-73.739499,40.704482],[-73.739493,40.704587],[-73.73954,40.704588]]]]")
	testLayoutCase(t, "3-01317-0141", "[[[[-73.945989,40.663811],[-73.946223,40.663797],[-73.946215,40.663737],[-73.945982,40.663751],[-73.945989,40.663811]]]]")
	testLayoutCase(t, "4-09948-0031", "[[[[-73.784436,40.71958],[-73.78474,40.719538],[-73.784427,40.719513],[-73.784436,40.71958]]]]")
	testLayoutCase(t, "4-11212-0015", "[[[[-73.729003,40.705698],[-73.729256,40.705762],[-73.72926,40.705701],[-73.729003,40.705698]]]]")

	//
	// We had misaligned center here
	//

	testLayoutCase(t, "1-00846-0021", "[[[[-73.990348,40.737378],[-73.990333,40.737446],[-73.990657,40.737526],[-73.990667,40.737513],[-73.990348,40.737378]]]]")
	// Complex shape and we have
	testLayoutCaseOnly1(t, "1-00817-0036", "[[[[-73.993561,40.737168],[-73.993495,40.73714],[-73.993447,40.737207],[-73.993468,40.737255],[-73.99352,40.737276],[-73.993708,40.737021],[-73.993678,40.737008],[-73.993561,40.737168]]]]")
}

func TestLayoutNegative(t *testing.T) {
	testLayoutCaseNone(t, "2-04544-0027", "[[[[-73.868414,40.870473],[-73.868342,40.870474],[-73.868343,40.870529],[-73.868417,40.870528],[-73.868414,40.870473]]]]")
	testLayoutCaseNone(t, "4-10234-0244", "[[[[-73.783613,40.703043],[-73.783581,40.703056],[-73.783746,40.703301],[-73.783778,40.703289],[-73.783613,40.703043]]]]")
	testLayoutCaseNone(t, "4-06192-0054", "[[[[-73.761959,40.769605],[-73.761947,40.769583],[-73.761448,40.769715],[-73.761959,40.769605]]]]")
	testLayoutCaseNone(t, "4-11069-0082", "[[[[-73.754656,40.694181],[-73.754681,40.694175],[-73.754565,40.693917],[-73.75454,40.693924],[-73.754656,40.694181]]]]")
	testLayoutCaseNone(t, "4-10067-0070", "[[[[-73.807296,40.687173],[-73.807327,40.687162],[-73.807211,40.686955],[-73.807179,40.686965],[-73.807296,40.687173]]]]")

	testLayoutCaseNone(t, "4-13363-0036", "[[[[-73.760702,40.660659],[-73.760759,40.660729],[-73.76078,40.660667],[-73.760702,40.660659]]]]")
	testLayoutCaseNone(t, "4-15801-0010", "[[[[-73.759378,40.595632],[-73.759526,40.595683],[-73.759532,40.595617],[-73.759378,40.595632]]]]")
	testLayoutCaseNone(t, "2-05084-0140", "[[[[-73.851362,40.901518],[-73.85131,40.901505],[-73.851284,40.901567],[-73.851335,40.901582],[-73.851362,40.901518]]]]")
	testLayoutCaseNone(t, "4-12580-0154", "[[[[-73.766198,40.672551],[-73.766201,40.672572],[-73.766543,40.672556],[-73.76654,40.672528],[-73.766198,40.672551]]]]")

	// Too small triangle
	testLayoutCaseNone(t, "4-15801-0010", "[[[[-73.759378,40.595632],[-73.759526,40.595683],[-73.759532,40.595617],[-73.759378,40.595632]]]]")
}

func TestCornerCases(t *testing.T) {
	// Narrow triangle
	// testLayoutCase(t, "1-01065-0132", "[[[[-73.98691,40.767126],[-73.986612,40.766955],[-73.986586,40.76699],[-73.98691,40.767126]]]]")
}
